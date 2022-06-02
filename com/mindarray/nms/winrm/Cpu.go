package winrm

import (
	exception "ModataPlugins/com/mindarray/nms/exceptionHandler"
	"encoding/json"
	"fmt"
	"github.com/masterzen/winrm"
	"regexp"
	"strings"
)

func Cpu(credentials map[string]interface{}) {
	defer exception.ErrorHandle(credentials)
	host := (credentials["ip"]).(string)
	port := int(credentials["port"].(float64))
	username := credentials["username"].(string)
	password := credentials["password"].(string)

	endpoint := winrm.NewEndpoint(host, port, false, false, nil, nil, nil, 0)
	result := make(map[string]interface{})
	var errors []string
	client, err := winrm.NewClient(endpoint, username, password)
	if err != nil {
		errors = append(errors, err.Error())
	}
	clients, er := client.CreateShell()
	defer clients.Close()
	if er != nil {
		errors = append(errors, er.Error())
		result["status"] = "fail"
		result["error"] = errors
		data, _ := json.Marshal(result)
		fmt.Print(string(data))
	} else {
		output := ""
		a := "aa"
		ac := "(Get-Counter '\\Processor(*)\\% Idle Time','\\Processor(*)\\% Processor Time','\\Processor(*)\\% user time' -ErrorAction SilentlyContinue).countersamples | Format-List  -Property Path,Cookedvalue;"
		output, _, _, err = client.RunPSWithString(ac, a)
		re := regexp.MustCompile("Path\\s*\\:\\s\\\\+[\\w\\-#]+\\\\(\\w*\\([\\w\\-#]+\\))\\\\%?\\s*(\\w*\\s*\\w*)\\s*\\w*\\s*:\\s*([\\d\\.]+)")
		value := re.FindAllStringSubmatch(output, -1)
		var counters = 3
		var cores []map[string]interface{}
		size := len(value) / counters

		for counterIndex := 0; counterIndex < len(value)/counters; counterIndex++ {
			count := 0
			core := make(map[string]interface{})
			res := strings.Split(value[counterIndex][1], "(")
			if strings.Split(res[1], ")")[0] == "_total" {
				result["system.cpu.idle.percent"] = value[counterIndex][3]
				result["system.cpu.process.percent"] = value[count+size][3]
				result["system.cpu.user.percent"] = value[count+size+size][3]
			} else {
				core["core.name"] = value[counterIndex][2]
				core["core.idle.percent"] = value[counterIndex][3]
				core["core.process.percent"] = value[count+size][3]
				core["core.user.percent"] = value[count+size+size][3]

			}
			cores = append(cores, core)
		}
		result["cores"] = cores
		result["ip"] = credentials["ip"]
		result["metric.group"] = credentials["metric.group"]
		result["status"] = "success"
		data, _ := json.Marshal(result)
		fmt.Print(string(data))

	}
}

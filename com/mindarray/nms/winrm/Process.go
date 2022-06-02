package winrm

import (
	exception "ModataPlugins/com/mindarray/nms/exceptionHandler"
	"encoding/json"
	"fmt"
	"github.com/masterzen/winrm"
	"regexp"
)

func Process(credentials map[string]interface{}) {
	defer exception.ErrorHandle(credentials)
	host := (credentials["ip"]).(string)
	port := int(credentials["port"].(float64))
	username := credentials["username"].(string)
	password := credentials["password"].(string)
	endpoint := winrm.NewEndpoint(host, port, false, false, nil, nil, nil, 0)
	client, err := winrm.NewClient(endpoint, username, password)
	var errors []string
	result := make(map[string]interface{})
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
		a := "aa"
		output := ""
		ac := "(Get-Counter '\\Process(*)\\ID Process','\\Process(*)\\% Processor Time','\\Process(*)\\Thread Count' -ErrorAction SilentlyContinue).countersamples | Format-List  -Property Path,Cookedvalue;"
		output, _, _, err = client.RunPSWithString(ac, a)
		re := regexp.MustCompile("Path\\s*\\:\\s*\\\\+[\\w\\-#.]+\\\\\\w*\\(([\\w\\-#.]+)\\)\\\\%?\\s*(\\w*\\s*\\w*)\\s*\\w*\\s*:\\s*([\\d\\.]+)")
		value := re.FindAllStringSubmatch(output, -1)
		var processes []map[string]interface{}
		processes = append(processes, result)
		var count int
		for i := 0; i < len(value); i++ {
			temp := make(map[string]interface{})
			temp1 := make(map[string]interface{})
			processName := value[i][1]
			for j := 0; j < len(processes); j++ {
				temp = processes[j]
				if temp[processName] != nil {
					count = 1
					break
				} else {
					count = 0
				}
			}
			if count == 0 {
				temp1["process.name"] = processName
				if (value[i][2]) == "id process\r" {
					temp1["process.id"] = value[i][3]
				} else if value[i][2] == "% processor time\r" {
					temp1["process.processor.time.percent"] = value[i][3]
				} else if value[i][2] == "thread count\r" {
					temp1["process.thread.count"] = value[i][3]
				}
				processes = append(processes, temp1)

			} else {
				if (value[i][2]) == "id process\r" {
					temp["process.id"] = value[i][3]
				} else if value[i][2] == "% processor time\r" {
					temp["process.processor.time.percent"] = value[i][3]
				} else if value[i][2] == "thread count\r" {
					temp["process.thread.count"] = value[i][3]
				}
				count = 1
				processes = append(processes, temp)
			}
		}
		processes = processes[1:len(processes)]
		size := (len(processes)) / 3
		var Values []map[string]interface{}
		for k := 0; k < len(processes)/3; k = k + 1 {
			temp2 := make(map[string]interface{})
			temp2 = processes[k]
			temp2["process.processor.time.percent"] = value[k+size][3]
			temp2["process.thread.count"] = value[k+size+size][3]
			Values = append(Values, temp2)
		}
		result["process"] = Values
		result["ip"] = credentials["ip"]
		result["metric.group"] = credentials["metric.group"]
		result["status"] = "success"
		data, _ := json.Marshal(result)
		fmt.Print(string(data))
	}
}

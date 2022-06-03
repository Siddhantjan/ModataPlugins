package winrm

import (
	exception "ModataPlugins/com/mindarray/nms/exceptionHandler"
	"encoding/json"
	"fmt"
	"github.com/masterzen/winrm"
	"regexp"
	"strings"
)

func Process(credentials map[string]interface{}) {
	defer exception.ErrorHandle(credentials)
	host := (credentials["ip"]).(string)
	port := int(credentials["port"].(float64))
	username := credentials["username"].(string)
	password := credentials["password"].(string)
	endpoint := winrm.NewEndpoint(host, port, false, false, nil, nil, nil, 0)
	client, err := winrm.NewClient(endpoint, username, password)
	result := make(map[string]interface{})
	if err != nil {
		result["error"] = err.Error()
		result["status"] = "fail"
		data, _ := json.Marshal(result)
		fmt.Print(string(data))

	} else {
		clients, er := client.CreateShell()
		defer clients.Close()
		if er != nil {
			userReadableError := strings.Contains(er.Error(), "connection refused")
			if userReadableError {
				result["error"] = "wrong ip or port ( connection refused )"
			} else {
				result["error"] = "wrong username or password ( unable to authenticate )"
			}
			result["status"] = "fail"
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
			for index := 0; index < len(value); index++ {
				temp := make(map[string]interface{})
				temp1 := make(map[string]interface{})
				processName := value[index][1]
				for subIndex := 0; subIndex < len(processes); subIndex++ {
					temp = processes[subIndex]
					if temp[processName] != nil {
						count = 1
						break
					} else {
						count = 0
					}
				}
				if count == 0 {
					temp1["process.name"] = processName
					if (value[index][2]) == "id process\r" {
						temp1["process.id"] = value[index][3]
					} else if value[index][2] == "% processor time\r" {
						temp1["process.processor.time.percent"] = value[index][3]
					} else if value[index][2] == "thread count\r" {
						temp1["process.thread.count"] = value[index][3]
					}
					processes = append(processes, temp1)

				} else {
					if (value[index][2]) == "id process\r" {
						temp["process.id"] = value[index][3]
					} else if value[index][2] == "% processor time\r" {
						temp["process.processor.time.percent"] = value[index][3]
					} else if value[index][2] == "thread count\r" {
						temp["process.thread.count"] = value[index][3]
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
}

package winrm

import (
	exception "ModataPlugins/com/mindarray/nms/exceptionHandler"
	"encoding/json"
	"fmt"
	"github.com/masterzen/winrm"
	"strings"
)

func System(credentials map[string]interface{}) {

	defer exception.ErrorHandle(credentials)

	host := (credentials["ip"]).(string)
	port := int(credentials["port"].(float64))
	username := credentials["username"].(string)
	password := credentials["password"].(string)

	endpoint := winrm.NewEndpoint(host, port, false, false, nil, nil, nil, 0)

	result := make(map[string]interface{})

	client, err := winrm.NewClient(endpoint, username, password)

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

			ac := "(Get-WmiObject win32_operatingsystem).name;(Get-WMIObject win32_operatingsystem).version;whoami;(Get-WMIObject win32_operatingsystem).LastBootUpTime;"

			output, _, _, err = client.RunPSWithString(ac, a)

			res1 := strings.Split(output, "\n")

			result["system.os.name"] = strings.Replace(strings.Split(res1[0], "\r")[0], "\\", ": ", 9)
			result["system.os.version"] = strings.Split(res1[1], "\r")[0]
			result["system.user.name"] = strings.Replace(strings.Split(res1[2], "\r")[0], "\\", ": ", 2)
			result["system.up.time"] = strings.Split(res1[3], "\r")[0]
			result["metric.group"] = credentials["metric.group"]
			result["ip"] = credentials["ip"]
			result["status"] = "success"

			data, _ := json.Marshal(result)
			fmt.Print(string(data))
		}
	}
}

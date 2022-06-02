package winrm

import (
	exception "ModataPlugins/com/mindarray/nms/exceptionHandler"
	"encoding/json"
	"fmt"
	"github.com/masterzen/winrm"
	"strings"
)

func Discovery(credentials map[string]interface{}) {
	defer exception.ErrorHandle(credentials)
	var errorOccurred []string
	defer exception.ErrorHandle(credentials)
	result := make(map[string]interface{})
	host := (credentials["ip"]).(string)
	port := int(credentials["port"].(float64))
	username := credentials["username"].(string)
	password := credentials["password"].(string)
	endpoint := winrm.NewEndpoint(host, port, false, false, nil, nil, nil, 0)
	client, err := winrm.NewClient(endpoint, username, password)
	if err != nil {
		result["status"] = "fail"

		errorOccurred = append(errorOccurred, err.Error())
	}
	_, err2 := client.CreateShell()
	if err2 != nil {
		result["status"] = "fail"
		errorOccurred = append(errorOccurred, err2.Error())
	}
	if len(errorOccurred) == 0 {
		result["status"] = "success"
		a := "aa"
		output := ""
		cmd := "hostname"
		output, _, _, err = client.RunPSWithString(cmd, a)
		if err != nil {
			result["status"] = "fail"
			result["error"] = err.Error()

		} else {
			result["host"] = strings.Split(output, "\r\n")[0]
		}
	} else {
		result["status"] = "fail"
		result["error"] = errorOccurred
	}
	data, _ := json.Marshal(result)
	fmt.Print(string(data))
}

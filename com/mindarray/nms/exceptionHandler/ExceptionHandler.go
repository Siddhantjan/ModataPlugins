package exceptionHandler

import (
	"encoding/json"
	"fmt"
)

func ErrorHandle(credentials map[string]interface{}) {

	var data = make(map[string]interface{})
	data["ip"] = credentials["ip"]
	error := recover()
	if error != nil {
		data["Panic"] = "Yes"
		data["error"] = error
		data["status"] = "fail"
		data, _ := json.Marshal(data)
		fmt.Println(string(data))
	}
}

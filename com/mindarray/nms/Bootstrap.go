package main

import (
	"ModataPlugins/com/mindarray/nms/snmp"
	"ModataPlugins/com/mindarray/nms/ssh"
	"ModataPlugins/com/mindarray/nms/winrm"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	argument, _ := base64.StdEncoding.DecodeString(os.Args[1])
	credentials := make(map[string]interface{})

	var error = json.Unmarshal([]byte(string(argument)), &credentials)

	result := make(map[string]interface{})

	if error != nil {
		result["status"] = "fail"
		result["error"] = "yes"
		result["Cause"] = error

	}

	if credentials["category"] == "discovery" {

		if credentials["type"] == "linux" {
			ssh.Discovery(credentials)

		} else if credentials["type"] == "windows" {
			winrm.Discovery(credentials)

		} else if credentials["type"] == "network" {
			snmp.Discovery(credentials)
		}

	} else if credentials["category"] == "polling" {

		if credentials["type"] == "linux" {
			switch credentials["metric.group"] {

			case "system":
				ssh.System(credentials)
				break

			case "disk":
				ssh.Disk(credentials)
				break

			case "memory":
				ssh.Memory(credentials)
				break

			case "process":
				ssh.Process(credentials)
				break

			case "cpu":
				ssh.Cpu(credentials)

			default:
				result["status"] = "fail"
				result["error"] = "yes"
				result["Cause"] = "Wrong metric group selected for metric type linux"

			}

		} else if credentials["type"] == "windows" {

			switch credentials["metric.group"] {

			case "system":
				winrm.System(credentials)
				break

			case "disk":
				winrm.Disk(credentials)
				break

			case "memory":
				winrm.Memory(credentials)
				break

			case "process":
				winrm.Process(credentials)
				break

			case "cpu":
				winrm.Cpu(credentials)

			default:
				result["status"] = "fail"
				result["error"] = "yes"
				result["Cause"] = "Wrong metric group selected for metric type Windows"

			}

		} else if credentials["type"] == "network" {

			switch credentials["metric.group"] {

			case "system":
				snmp.System(credentials)
				break

			case "interface":
				snmp.InterfaceData(credentials)
				break

			default:
				result["status"] = "fail"
				result["error"] = "yes"
				result["Cause"] = "Wrong metric group selected for metric type Network Devices"
			}
		}
	} else {

		result["status"] = "fail"
		result["error"] = "yes"
		result["Cause"] = "Wrong Category Given"

	}

	data, _ := json.Marshal(result)

	if string(data) != "{}" {
		fmt.Print(string(data))
	}

}

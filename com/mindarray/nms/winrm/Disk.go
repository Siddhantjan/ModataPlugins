package winrm

import (
	exception "ModataPlugins/com/mindarray/nms/exceptionHandler"
	"encoding/json"
	"fmt"
	"github.com/masterzen/winrm"
	"strconv"
	"strings"
)

func Disk(credentials map[string]interface{}) {
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
		a := "aa"
		output := ""
		ac := "Get-WmiObject win32_logicaldisk |Foreach-Object {$_.DeviceId,$_.Freespace,$_.Size}"
		output, _, _, err = client.RunPSWithString(ac, a)
		res := strings.Split(output, "\r\n")
		var disks []map[string]interface{}
		var usedBytes int64
		var totalBytes int64

		for diskCounter := 0; diskCounter < len(res); diskCounter = diskCounter + 3 {
			disk := make(map[string]interface{})
			disk["Disk.Name"] = strings.Split(res[diskCounter], ":")[0]
			if (diskCounter+1) > len(res) || res[diskCounter+1] == "" {
				disk["disk.free.bytes"] = 0
				disk["disk.total.bytes"] = 0
				disk["disk.available.bytes"] = 0
				disk["disk.used.percent"] = 0
				disk["disk.free.percent"] = 0
				disks = append(disks, disk)
				break
			}
			bytes, _ := strconv.ParseInt(res[diskCounter+1], 10, 64)
			usedBytes = usedBytes + bytes
			disk["disk.available.bytes"], _ = strconv.ParseInt(res[diskCounter+1], 10, 64)
			bytes, _ = strconv.ParseInt(res[diskCounter+2], 10, 64)
			totalBytes = totalBytes + bytes
			disk["disk.total.bytes"] = bytes
			disk["disk.used.bytes"] = (disk["disk.total.bytes"]).(int64) - (disk["disk.available.bytes"]).(int64)
			disk["disk.used.percent"] = ((float64((disk["disk.total.bytes"]).(int64)) - float64((disk["disk.used.bytes"]).(int64))) / float64((disk["disk.total.bytes"].(int64)))) * 100
			disk["disk.free.percent"] = 100 - disk["disk.used.percent"].(float64)
			disks = append(disks, disk)
		}
		result["disk.total.bytes"] = totalBytes
		result["disk.used.byes"] = usedBytes
		result["disk.available.bytes"] = totalBytes - usedBytes
		result["disk.used.percent"] = ((float64(totalBytes) - float64(usedBytes)) / float64(totalBytes)) * 100
		result["disk.available.percent"] = 100.00 - (result["disk.used.percent"]).(float64)
		result["disks"] = disks
		result["ip"] = credentials["ip"]
		result["metric.group"] = credentials["metric.group"]
		result["status"] = "success"
		data, _ := json.Marshal(result)
		fmt.Print(string(data))
	}
}

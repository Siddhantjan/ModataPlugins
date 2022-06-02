package ssh

import (
	exception "ModataPlugins/com/mindarray/nms/exceptionHandler"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strconv"
	"strings"
	"time"
)

func Memory(credentials map[string]interface{}) {
	defer exception.ErrorHandle(credentials)
	sshHost := credentials["ip"].(string)
	sshPort := int(credentials["port"].(float64))
	sshUser := credentials["username"].(string)
	sshPassword := credentials["password"].(string)

	config := &ssh.ClientConfig{
		Timeout:         10 * time.Second,
		User:            sshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Config: ssh.Config{Ciphers: []string{
			"aes128-ctr", "aes192-ctr", "aes256-ctr",
		}},
	}
	config.Auth = []ssh.AuthMethod{ssh.Password(sshPassword)}
	addr := fmt.Sprintf("%s:%d", sshHost, sshPort)
	sshClient, er := ssh.Dial("tcp", addr, config)
	result := make(map[string]interface{})
	session, err := sshClient.NewSession()
	if err != nil {

		result["error"] = "yes"
		result["Cause"] = er
		result["status"] = "fail"
		data, _ := json.Marshal(result)
		fmt.Print(string(data))

	} else {

		cmd := "free -b | awk  '{if ($1 != \"total\") print $1 \" \" $2 \" \" $3 \" \" $4 \" \"$7}'"
		combo, _ := session.CombinedOutput(cmd)
		output := string(combo)
		res := strings.Split(output, "\n")

		memoryValue := strings.Split(res[0], " ")
		totalBytes, _ := strconv.ParseInt(memoryValue[1], 10, 64)
		result["memory.total.bytes"] = totalBytes

		usedBytes, _ := strconv.ParseInt(memoryValue[2], 10, 64)
		result["memory.used.bytes"] = usedBytes
		result["memory.free.bytes"], _ = strconv.ParseInt(memoryValue[3], 10, 64)
		result["memory.available.bytes"], _ = strconv.ParseInt(memoryValue[4], 10, 64)

		swapValue := strings.Split(res[1], " ")
		result["memory.swap.total.bytes"], _ = strconv.ParseInt(swapValue[1], 10, 64)
		result["memory.swap.used.bytes"], _ = strconv.ParseInt(swapValue[2], 10, 64)
		result["memory.swap.free.bytes"], _ = strconv.ParseInt(swapValue[3], 10, 64)
		usedPercent := float64(totalBytes-usedBytes) / float64(totalBytes)

		result["memory.used.percent"] = usedPercent
		result["memory.available.percent"] = 100 - usedPercent

		result["ip"] = credentials["ip"]
		result["metric.group"] = credentials["metric.group"]

		data, _ := json.Marshal(result)
		fmt.Print(string(data))
	}
}

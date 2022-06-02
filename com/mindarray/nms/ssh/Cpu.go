package ssh

import (
	exception "ModataPlugins/com/mindarray/nms/exceptionHandler"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
	"time"
)

func Cpu(credentials map[string]interface{}) {
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
		result["Cause"] = er.Error()
		result["status"] = "fail"
		data, _ := json.Marshal(result)
		fmt.Print(string(data))
	} else {
		cmd := "mpstat -P ALL |awk  '{print $4 \" \" $5 \" \" $7 \" \" $14}'"
		combo, _ := session.CombinedOutput(cmd)
		output := string(combo)
		res := strings.Split(output, "\n")
		system := strings.Split(res[3], " ")
		result["system.cpu.user.percent"] = system[1]
		result["system.cpu.system.percent"] = system[2]
		result["system.cpu.idle.percent"] = system[3]
		//fmt.Println(res)
		var cores []map[string]interface{}

		for outputIterator := 4; outputIterator < len(res)-1; outputIterator++ {
			core := make(map[string]interface{})
			value := strings.Split(res[outputIterator], " ")
			core["core.name"] = value[0]
			core["core.user.percent"] = value[1]
			core["core.system.percent"] = value[2]
			core["core.idle.percent"] = value[3]
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

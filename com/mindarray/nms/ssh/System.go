package ssh

import (
	exception "ModataPlugins/com/mindarray/nms/exceptionHandler"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
	"time"
)

func System(credentials map[string]interface{}) {
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

		terminalCommand := "uname -a | awk  '{ print $1 \" \" $2  \" \" $4 \" \"$6 \" \" $7 \" \" $8 \" \"$9 }'"
		combo, er := session.CombinedOutput(terminalCommand)
		output := string(combo)

		res := strings.Split(output, "\n")
		systemValue := strings.Split(res[0], " ")

		result["system.os.name"] = systemValue[0]
		result["system.user.name"] = systemValue[1]
		result["system.os.version"] = systemValue[2]
		result["system.up.time"] = systemValue[3] + " " + systemValue[4] + " " + systemValue[5] + " " + systemValue[6]

		session.Close()

		session, err = sshClient.NewSession()

		if err != nil {

			result["error"] = "yes"
			result["Cause"] = er
			result["status"] = "fail"

		} else {

			result["error"] = "no"
			result["status"] = "success"

		}
		runningProcess := " vmstat | awk '{print $1 \" \" $2 \" \"  $12}'"

		combo, er = session.CombinedOutput(runningProcess)

		output = string(combo)

		res = strings.Split(output, "\n")

		processValue := strings.Split(res[2], " ")

		result["system.running.process"] = processValue[0]
		result["system.blocking.process"] = processValue[1]
		result["system.context.switching"] = processValue[2]
		result["ip"] = credentials["ip"]
		result["metric.group"] = credentials["metric.group"]

		data, _ := json.Marshal(result)
		fmt.Print(string(data))
	}
}

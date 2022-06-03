package ssh

import (
	exception "ModataPlugins/com/mindarray/nms/exceptionHandler"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
	"time"
)

func Discovery(credentials map[string]interface{}) {
	defer exception.ErrorHandle(credentials)
	sshHost := credentials["ip"].(string)
	sshPort := int(credentials["port"].(float64))
	sshUser := credentials["username"].(string)
	sshPassword := credentials["password"].(string)

	config := &ssh.ClientConfig{
		Timeout:         6 * time.Second,
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
		session, err := sshClient.NewSession()
		if err != nil {
			result["status"] = "fail"
			result["error"] = err.Error()
		} else {
			result["status"] = "success"
		}
		cmd := "hostname"
		combo, err := session.CombinedOutput(cmd)
		output := string(combo)
		if err != nil {
			result["status"] = "fail"
			result["error"] = er.Error()
		} else {
			result["status"] = "success"
			result["host"] = strings.Split(output, "\n")[0]
		}
		data, _ := json.Marshal(result)
		fmt.Print(string(data))

	}
}

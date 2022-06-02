package snmp

import (
	exception "ModataPlugins/com/mindarray/nms/exceptionHandler"
	"encoding/json"
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"time"
)

func Discovery(credentials map[string]interface{}) {
	defer exception.ErrorHandle(credentials)
	result := make(map[string]interface{})
	var version = g.Version1
	switch credentials["version"] {
	case "version1":
		version = g.Version1
		break
	case "version2":
		version = g.Version2c
		break
	case "version3":
		version = g.Version3
		break
	}

	params := &g.GoSNMP{
		Target:    credentials["ip"].(string),
		Port:      uint16(int(credentials["port"].(float64))),
		Community: credentials["community"].(string),
		Version:   version,
		Timeout:   time.Duration(2) * time.Second,
	}
	err := params.Connect()
	if err != nil {
		result["error"] = "yes"
		result["Cause"] = err.Error()
		data, _ := json.Marshal(result)
		fmt.Print(string(data))
	} else {
		result["error"] = "no"
	}
	sysName, err := params.Get([]string{".1.3.6.1.2.1.1.5.0"})

	if err != nil {
		result["error"] = err.Error()
		result["status"] = "fail"
		data, _ := json.Marshal(result)
		fmt.Print(string(data))
	} else {
		result["error"] = "no"
		result["status"] = "success"
		for _, value := range sysName.Variables {
			result["host"] = string(value.Value.([]byte))
		}
	}
	walkOid := "1.3.6.1.2.1.2.2.1.1"
	index := "1.3.6.1.2.1.2.2.1.1."
	description := "1.3.6.1.2.1.2.2.1.2."
	name := "1.3.6.1.2.1.31.1.1.1.1."
	operationalStatus := "1.3.6.1.2.1.2.2.1.8."
	alias := "1.3.6.1.2.1.31.1.1.1.18."

	var walkOidArray []string
	walk := params.Walk(walkOid, func(pdu g.SnmpPDU) error {
		switch pdu.Type {
		case g.IPAddress:
			result := pdu.Value
			walkOidArray = append(walkOidArray, result.(string))
			break
		case g.Integer:
			result := g.ToBigInt(pdu.Value)
			walkOidArray = append(walkOidArray, result.String())
			break
		case g.OctetString:
			result := pdu.Value.([]byte)
			walkOidArray = append(walkOidArray, string(result))
			break
		case g.Gauge32:
			result := pdu.Value
			walkOidArray = append(walkOidArray, result.(string))
		default:
			result := pdu.Value
			walkOidArray = append(walkOidArray, result.(string))
		}
		return nil
	},
	)
	if walk != nil {
		result["error"] = walk.Error()
	}

	var oids []string
	for i := 0; i < len(walkOidArray); i++ {
		oids = append(oids, index+walkOidArray[i])
		oids = append(oids, description+walkOidArray[i])
		oids = append(oids, name+walkOidArray[i])
		oids = append(oids, operationalStatus+walkOidArray[i])
		oids = append(oids, alias+walkOidArray[i])
	}
	var startIndex = 0
	var endIndex = 50
	var resultArray []interface{}
	for {
		if len(resultArray) == len(oids) {
			break
		}
		output, error := params.Get(oids[startIndex:endIndex])
		if error != nil {
			result["error"] = error.Error()

			return
		}
		for _, variable := range output.Variables {

			resultArray = append(resultArray, ConvertData(variable))
		}
		startIndex = endIndex
		endIndex = endIndex + 50

		if endIndex > len(oids) {
			endIndex = len(oids)
		}
	}
	var interfaces []map[string]interface{}
	for i := 0; i < len(resultArray); i = i + 5 {
		interfaceValue := make(map[string]interface{})
		interfaceValue["interface.index"] = resultArray[i].(int)
		interfaceValue["interface.description"] = resultArray[i+1]
		interfaceValue["interface.name"] = resultArray[i+2]
		if resultArray[i+3] == 1 {
			interfaceValue["interface.operational.status"] = "up"
		} else {
			interfaceValue["interface.operational.status"] = "down"
		}
		if resultArray[i+4] == "" {
			interfaceValue["interface.alias.name"] = "empty"
		} else {
			interfaceValue["interface.alias.name"] = resultArray[i+5]
		}
		interfaces = append(interfaces, interfaceValue)
	}
	fmt.Println(interfaces)

	result["interfaces"] = interfaces
	result["ip"] = credentials["ip"]
	if result["error"] == "no" {
		result["status"] = "success"
	} else {
		result["status"] = "fail"
	}
	data, _ := json.Marshal(result)
	fmt.Print(string(data))
}

func ConvertData(pdu g.SnmpPDU) interface{} {
	if pdu.Value == " " {
		return pdu.Value
	}
	switch pdu.Type {
	case g.IPAddress:
		return pdu.Value
	case g.Integer:
		return pdu.Value
	case g.OctetString:
		return string(pdu.Value.([]byte))
	default:
		return pdu.Value
	}

}

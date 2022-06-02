package snmp

import (
	exception "ModataPlugins/com/mindarray/nms/exceptionHandler"
	"encoding/json"
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"time"
)

func InterfaceData(credentials map[string]interface{}) {
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
	var errors []string
	params := &g.GoSNMP{
		Target:    credentials["ip"].(string),
		Port:      uint16(int(credentials["port"].(float64))),
		Community: credentials["community"].(string),
		Version:   version,
		Timeout:   time.Duration(1) * time.Second,
	}
	err := params.Connect()
	if err != nil {

		result["error"] = err.Error()
		result["status"] = "fail"
		data, _ := json.Marshal(result)
		fmt.Print(string(data))
	} else {
		walkOid := "1.3.6.1.2.1.2.2.1.1"
		snmpIndex := "1.3.6.1.2.1.2.2.1.1."
		description := "1.3.6.1.2.1.2.2.1.2."
		name := "1.3.6.1.2.1.31.1.1.1.1."
		operationalStatus := "1.3.6.1.2.1.2.2.1.8."
		adminStatus := "1.3.6.1.2.1.2.2.1.7."
		alias := "1.3.6.1.2.1.31.1.1.1.18."
		sentError := "1.3.6.1.2.1.2.2.1.20."
		receiveError := "1.3.6.1.2.1.2.2.1.14."
		sentOctets := "1.3.6.1.2.1.2.2.1.16."
		receiveOctets := "1.3.6.1.2.1.2.2.1.10."
		ifSpeed := "1.3.6.1.2.1.2.2.1.5."

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
			errors = append(errors, walk.Error())
		}

		var oids []string
		for index := 0; index < len(walkOidArray); index++ {
			oids = append(oids, snmpIndex+walkOidArray[index])
			oids = append(oids, description+walkOidArray[index])
			oids = append(oids, name+walkOidArray[index])
			oids = append(oids, operationalStatus+walkOidArray[index])
			oids = append(oids, adminStatus+walkOidArray[index])
			oids = append(oids, alias+walkOidArray[index])
			oids = append(oids, sentError+walkOidArray[index])
			oids = append(oids, receiveError+walkOidArray[index])
			oids = append(oids, sentOctets+walkOidArray[index])
			oids = append(oids, receiveOctets+walkOidArray[index])
			oids = append(oids, ifSpeed+walkOidArray[index])
		}
		var startIndex = 0
		var endIndex = 60

		var resultArray []interface{}

		for {
			if len(resultArray) == len(oids) {
				break
			}
			output, error := params.Get(oids[startIndex:endIndex])
			if error != nil {
				errors = append(errors, error.Error())
				return
			}
			for _, variable := range output.Variables {
				resultArray = append(resultArray, ConvertData(variable))
			}
			startIndex = endIndex
			endIndex = endIndex + 60
			if endIndex > len(oids) {
				endIndex = len(oids)
			}

		}
		var interfaces []map[string]interface{}
		for index := 0; index < len(resultArray); index = index + 11 {
			interfaceValue := make(map[string]interface{})
			interfaceValue["interface.index"] = resultArray[index].(int)
			interfaceValue["interface.description"] = resultArray[index+1]
			interfaceValue["interface.name"] = resultArray[index+2]
			if resultArray[index+3] == 1 {
				interfaceValue["interface.operational.status"] = "up"
			} else {
				interfaceValue["interface.operational.status"] = "down"
			}
			if resultArray[index+4] == 1 {
				interfaceValue["interface.admin.status"] = "up"
			} else {
				interfaceValue["interface.admin.status"] = "down"
			}
			if resultArray[index+5] == "" {
				interfaceValue["interface.alias.name"] = "empty"
			} else {
				interfaceValue["interface.alias.name"] = resultArray[index+5]
			}
			interfaceValue["interface.sent.errors"] = resultArray[index+6]
			interfaceValue["interface.receive.errors"] = resultArray[index+7]
			interfaceValue["interface.sent.octets"] = resultArray[index+8]
			interfaceValue["interface.receive.octets"] = resultArray[index+8]
			interfaceValue["interface.speed"] = resultArray[index+9]
			interfaces = append(interfaces, interfaceValue)
		}

		result["interfaces"] = interfaces
		result["ip"] = credentials["ip"]
		result["metric.group"] = credentials["metric.group"]
		if len(errors) == 0 {
			result["status"] = "success"
		} else {
			result["status"] = "fail"
			result["error"] = errors
		}
		data, _ := json.Marshal(result)
		fmt.Print(string(data))
	}
}

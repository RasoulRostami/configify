package powerdns

import (
	"configify/databases"
	"configify/services"
	"fmt"
	"log"
)

type TxtRecordSplit struct {
	next services.Process
}

func NewTxtRecordSplit(next services.Process) *TxtRecordSplit {
	return &TxtRecordSplit{next: next}
}

// splite long txt content
func (t *TxtRecordSplit) Update(message *databases.Message) bool {
	dns_records := message.Value["dns_records"].([]interface{})
	for i := 0; i < len(dns_records); i++ {
		item := dns_records[i].(map[string]interface{})
		if item["type"] == "TXT" {
			dns_records[i].(map[string]interface{})["content"] = t.chunckString(item["content"].(string), 255)
		}
	}
	log.Printf("DEBUG PowerDNS TXT Record Split process (%s) \n", message.Key)
	return t.next.Update(message)
}

func (t *TxtRecordSplit) Reverse(message *databases.Message) bool {
	return t.next.Reverse(message)
}

func (t *TxtRecordSplit) chunckString(longString string, maxLength int) string {
	var result string
	for len(longString) > 0 {
		if len(longString) <= maxLength {
			result += fmt.Sprintf("\"%s\"", longString) + " "
			break
		}
		result += fmt.Sprintf("\"%s\"", longString[:maxLength]) + " "
		longString = longString[maxLength:]
	}
	return result
}

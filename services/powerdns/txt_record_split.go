package powerdns

import (
	"configify/databases"
	"configify/services"
	"fmt"
	"log"
)

type TxtRecordSplit struct {
	Next services.Process
}

func NewTxtRecordSplit(next services.Process) *TxtRecordSplit {
	return &TxtRecordSplit{Next: next}
}

func (t *TxtRecordSplit) Update(message *databases.Message) {
	// separete long TXT content
	dns_records := message.Value["dns_records"].([]interface{})
	for i := 0; i < len(dns_records); i++ {
		item := dns_records[i].(map[string]interface{})
		if item["type"] == "TXT" {
			dns_records[i].(map[string]interface{})["content"] = t.chunckString(item["content"].(string), 255)
		}
	}
	log.Printf("DEBUG PowerDNS TXT Record Split process (%s) \n", message.Key)
	// call next step
	t.Next.Update(message)
}

func (t *TxtRecordSplit) Reverse(message *databases.Message) {
	t.Next.Reverse(message)
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

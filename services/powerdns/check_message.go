package powerdns

import (
	"configify/databases"
	"configify/services"
	"log"
	"regexp"
)

type CheckRelatedMessage struct {
	Next         services.Process
	Prefix       string
	prefix_regex *regexp.Regexp
}

func NewCheckRelatedmessage(prefix string, next services.Process) *CheckRelatedMessage {
	regex := regexp.MustCompile(prefix)
	return &CheckRelatedMessage{Prefix: prefix, Next: next, prefix_regex: regex}
}

func (c *CheckRelatedMessage) isRelatedMessage(key string) bool {
	return c.prefix_regex.MatchString(key)
}

// Check message key and if it is related go to next step
func (c *CheckRelatedMessage) Update(message *databases.Message) {
	if c.isRelatedMessage(message.Key) {
		log.Printf("DEBUG (%s) was accepted in PowerDNS services \n", message.Key)
		c.Next.Update(message)
	}
	log.Printf("DEBUG (%s) was rejected in PowerDNS services \n", message.Key)
}

// Check message key and if it is related go to next step
func (c *CheckRelatedMessage) Reverse(message *databases.Message) {
	if c.isRelatedMessage(message.Key) {
		c.Next.Reverse(message)
	}
}

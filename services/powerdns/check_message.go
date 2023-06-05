package powerdns

import (
	"configify/databases"
	"configify/services"
	"log"
	"regexp"
)

type CheckMessage struct {
	next        services.Process
	prefix      string
	prefixRegex *regexp.Regexp
}

func NewCheckMessage(prefix string, next services.Process) *CheckMessage {
	regex := regexp.MustCompile(prefix)
	return &CheckMessage{prefix: prefix, next: next, prefixRegex: regex}
}

func (c *CheckMessage) isRelatedMessage(key string) bool {
	return c.prefixRegex.MatchString(key)
}

// Check message schema and message key, if it is ok go to next process
func (c *CheckMessage) Update(message *databases.Message) bool {
	if c.isRelatedMessage(message.Key) {
		log.Printf("DEBUG PowerDNS Check Message accepte (%s)  \n", message.Key)
		return c.next.Update(message)
	} else {
		return false
	}
}

// Check message key
func (c *CheckMessage) Reverse(message *databases.Message) bool {
	if c.isRelatedMessage(message.Key) {
		log.Printf("DEBUG PowerDNS Check Message accepte (%s)  \n", message.Key)
		return c.next.Reverse(message)
	}
	return false
}

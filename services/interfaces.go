package services

import "configify/databases"

// Service interface
type Service interface {
	Update(message *databases.Message)
	Reverse(message *databases.Message)
	Reload()
}

// Each service have got some process
type Process interface {
	Update(message *databases.Message) bool
	Reverse(message *databases.Message) bool
}

package services

import "configify/databases"

type Process interface {
	Execute(message databases.Message)
}

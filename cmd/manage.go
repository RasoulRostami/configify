package cmd

import (
	"configify/databases"
	"fmt"
	"sync"
)

// Get new message and call related service
func Performer(messages chan databases.Message, wg *sync.WaitGroup,
) {
	// TODO For loop and call related service
	message := <-messages
	fmt.Println(message.Key, message.Value)
	wg.Done()
}

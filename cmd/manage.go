package cmd

import (
	"configify/databases"
	"fmt"
	"sync"
)

func Performer(messages chan databases.Message, wg *sync.WaitGroup,
) {
	message := <-messages
	fmt.Println(message.Key, message.Value)
	wg.Done()
}

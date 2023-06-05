package cmd

import (
	"configify/databases"
	"log"
	"sync"
)

// Listen to channel, get new message and call services
func Performer(messages chan databases.Message, wg *sync.WaitGroup,
) {
	for {
		message := <-messages
		log.Printf("DEBUG Performer recived (%s) \n", message.Key)
		if message.Type == databases.Set {
			for _, service := range ActiveServices {
				service.Update(&message)
				service.Reload()
			}
		} else {
			for _, service := range ActiveServices {
				service.Reverse(&message)
				service.Reload()
			}
		}
		wg.Done()
	}
}

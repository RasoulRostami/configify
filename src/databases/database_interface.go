package databases

import "sync"

type Database interface {
	Get(key string) (map[string]any, error)
	Keys(prefix string, messages chan Message, wg *sync.WaitGroup)
	Stream(message chan Message, wg *sync.WaitGroup)
}

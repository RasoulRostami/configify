package databases

type Database interface {
	Get(key string) (map[string]any, error)
	Keys(prefix string, messages chan map[string]any)
	Stream()
}

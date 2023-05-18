package publishers

type Publisher interface {
	Get(key string) (map[string]any, error)
	Keys(prefix string)
}

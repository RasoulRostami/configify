package databases

type Type string

const (
	Set    Type = "set"
	Remove Type = "remove"
)

type Message struct {
	Key   string
	Value map[string]interface{}
	Type  Type
}

package databases

type Type string

const (
	Set    Type = "Set"
	Remove Type = "Remove"
)

type Message struct {
	Key   string
	Value map[string]any
	Type  Type
}

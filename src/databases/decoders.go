package databases

import "encoding/json"

type Decoder interface {
	Decode(data string) (map[string]any, error)
}

type JsonDecoder struct {
}

func (j *JsonDecoder) Decode(data string) (map[string]any, error) {
	var result map[string]any
	error := json.Unmarshal([]byte(data), &result)

	if error != nil {
		return nil, error
	}

	return result, nil
}

// TODO and more decoder

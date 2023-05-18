package publishers

import (
	"github.com/spf13/cast"
)

func DecoderFactory(decoderType string) Decoder {
	if decoderType == "json" {
		return &JsonDecoder{}
	} else {
		panic("Decoder is invalid.")
	}
}

func PublisherFactory(publisherType string, config map[string]any, decoder Decoder) Publisher {
	if publisherType == "redis" {
		return NewRedisPublisher(
			cast.ToString(config["address"]),
			cast.ToString(config["password"]),
			cast.ToInt(config["db"]),
			decoder,
		)
	} else {
		panic("Publisher is invalid.")
	}
}

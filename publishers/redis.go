package publishers

import (
	"fmt"

	"github.com/go-redis/redis"
)

type RedisPublisher struct {
	client  *redis.Client
	decoder Decoder
}

func (r *RedisPublisher) Get(key string) (map[string]any, error) {
	data, err := r.client.Get(key).Result()
	if err != nil {
		return nil, err
	}
	result, err := r.decoder.Decode(data)
	if err != nil {
		return nil, err
	}
	return result, nil

}

func (r *RedisPublisher) Keys(prefix string) {
	if prefix == "" {
		prefix = "*"
	}

	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = r.client.Scan(cursor, prefix, 0).Result()

		if err != nil {
			panic(err)
		}

		for _, key := range keys {
			// TODO call channels
			fmt.Println("key", key)
		}
		if cursor == 0 { // no more keys
			break
		}
	}

}

func NewRedisPublisher(address string, password string, db int, decoder Decoder) *RedisPublisher {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	return &RedisPublisher{client: client, decoder: decoder}
}

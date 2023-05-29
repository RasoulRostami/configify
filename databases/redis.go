package databases

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type RedisDB struct {
	client  *redis.Client
	pubsub  *redis.PubSub
	decoder Decoder
}

func (r *RedisDB) Get(key string) (map[string]any, error) {
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

func (r *RedisDB) Keys(prefix string, messages chan Message, wg *sync.WaitGroup) {
	if prefix == "" {
		prefix = "*"
	}

	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = r.client.Scan(cursor, prefix, 100).Result()

		if err != nil {
			panic(err)
		}

		for _, key := range keys {
			data, err := r.Get(key)
			if err != nil {
				log.Printf("Error to get (%s) {%s}", key, err)
			} else {
				wg.Add(1)
				messages <- Message{Key: key, Value: data, Type: Set}
			}
		}
		// no more keys
		if cursor == 0 {
			break
		}
	}

}

func (r *RedisDB) Stream() {
	// Wait for confirmation that the subscription is created
	_, err := r.pubsub.Receive()
	if err != nil {
		panic(err)
	}

	// Start listening for events in a separate goroutine
	go func() {
		ch := r.pubsub.Channel()
		for msg := range ch {
			fmt.Println(msg.Channel, msg.Payload)
		}
	}()
	// Keep the main goroutine alive
	for {
		time.Sleep(1 * time.Second)
		fmt.Println(time.Now())
	}
}

func NewRedisDB(address string, password string, db int, decoder Decoder) *RedisDB {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	redis := RedisDB{client: client, decoder: decoder}
	pubsub := redis.client.Subscribe(
		fmt.Sprintf("__keyevent@%d__:set", db),
		fmt.Sprintf("__keyevent@%d__:del", db),
	)
	return &RedisDB{client: client, decoder: decoder, pubsub: pubsub}
}

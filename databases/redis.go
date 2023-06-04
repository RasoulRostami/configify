package databases

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type RedisDB struct {
	client     *redis.Client
	pubsub     *redis.PubSub
	setCahnnel string
	delChannel string
	decoder    Decoder
}

func NewRedisDB(address string, password string, db int, decoder Decoder) *RedisDB {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	redis := RedisDB{client: client, decoder: decoder}
	setCahnnel := fmt.Sprintf("__keyevent@%d__:set", db)
	delCahnnel := fmt.Sprintf("__keyevent@%d__:del", db)
	pubsub := redis.client.Subscribe(setCahnnel, delCahnnel)
	return &RedisDB{
		client:     client,
		decoder:    decoder,
		pubsub:     pubsub,
		setCahnnel: setCahnnel,
		delChannel: delCahnnel,
	}
}

func (r *RedisDB) Get(key string) (map[string]interface{}, error) {
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
				log.Printf("WARNING Redis DB, Keys, (%s) {%s} \n", key, err)
			} else {
				log.Printf("DEBUG Redis DB, Keys: Sent (%s) to channel", key)
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

// DB streaming, listen to pubsub events,
func (r *RedisDB) startStream(messages chan Message, wg *sync.WaitGroup) {
	ch := r.pubsub.Channel()
	for msg := range ch {
		if msg.Channel == r.setCahnnel {
			data, err := r.Get(msg.Payload)
			if err != nil {
				log.Printf("WARNING Redis DB, Stream, (%s) {%s} \n", msg.Payload, err)
			} else {
				log.Printf("DEBUG Redis DB, Stream: Sent (%s) to channel", msg.Payload)
				wg.Add(1)
				messages <- Message{Key: msg.Payload, Value: data, Type: Set}
			}
		} else if msg.Channel == r.delChannel {
			wg.Add(1)
			messages <- Message{Key: msg.Payload, Type: Remove}
		}
	}
}

// Start stream and keep the main goroutine alive
func (r *RedisDB) Stream(messages chan Message, wg *sync.WaitGroup) {
	// Wait for confirmation that the subscription is created
	_, err := r.pubsub.Receive()
	if err != nil {
		panic(err)
	}
	// Start listening for events in a separate goroutine
	go r.startStream(messages, wg)

	// Keep the main goroutine alive
	for {
		time.Sleep(1 * time.Second)
		fmt.Println(time.Now())
	}
}

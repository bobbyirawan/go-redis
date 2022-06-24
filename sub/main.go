package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"test-pubsub/redis"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

var ctx = context.Background()

func main() {
	subscriber := redis.RedisClient.Subscribe(ctx, "send-user-data")

	user := User{}

	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}

		if err := json.Unmarshal([]byte(msg.Payload), &user); err != nil {
			panic(err)
		}

		log.Println(msg.Payload)

		fmt.Println("Received message from " + msg.Channel + " channel.")
		fmt.Printf("%+v\n", user)
	}
}

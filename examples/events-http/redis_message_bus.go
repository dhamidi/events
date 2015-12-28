package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dhamidi/events"
	"gopkg.in/redis.v3"
)

type UnknownMessageError struct {
	BusName string
	Message events.Message
}

func NewUnknownMessageError(busName string, message events.Message) *UnknownMessageError {
	return &UnknownMessageError{
		BusName: busName,
		Message: message,
	}
}

func (self *UnknownMessageError) Error() string {
	return fmt.Sprintf("%s: unknown message: %s", self.BusName, self.Message.RoutingKey())
}

type RedisMessageBus struct {
	channel  string
	client   *redis.Client
	handlers map[string]events.MessageHandler
}

func NewRedisMessageBus(channel string, client *redis.Client) *RedisMessageBus {
	return &RedisMessageBus{
		channel:  channel,
		client:   client,
		handlers: map[string]events.MessageHandler{},
	}
}

func (self *RedisMessageBus) HandleMessage(routingKey string, handler events.MessageHandler) events.MessageBus {
	self.handlers[routingKey] = handler
	return self
}

func (self *RedisMessageBus) Listen() error {
	pubsub, err := self.client.Subscribe(self.channel)
	if err != nil {
		return err
	}
	defer pubsub.Close()

	for {
		redisMessage, err := pubsub.ReceiveMessage()
		if err != nil {
			return err
		}
		message := NewRedisMessage(self.client)
		if err := json.Unmarshal([]byte(redisMessage.Payload), message); err != nil {
			message.Reject(err)
			continue
		}
		fmt.Fprintf(os.Stderr, "DEBUG:\n%s\n", (func() []byte { data, _ := json.Marshal(message); return data })())
		handler, found := self.handlers[message.RoutingKey()]
		if !found {
			handler = self.handlers["/"]
		}
		if handler == nil {
			message.Reject(NewUnknownMessageError("redis", message))
			continue
		}

		handler.HandleMessage(message)
	}
}

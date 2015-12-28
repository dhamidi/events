package main

import (
	"github.com/dhamidi/events"
	"gopkg.in/redis.v3"
)

type RedisMessage struct {
	ReplyTo string
	Headers map[string]string
	Message string

	client *redis.Client
}

func NewRedisMessage(client *redis.Client) *RedisMessage {
	return &RedisMessage{
		client: client,
	}
}

func (self *RedisMessage) RoutingKey() string {
	return self.Headers["Routing-Key"]
}

func (self *RedisMessage) ContentType() string {
	return self.Headers["Content-Type"]
}

func (self *RedisMessage) Body() []byte {
	return []byte(self.Message)
}

func (self *RedisMessage) Header(name string) string {
	return self.Headers[name]
}

func (self *RedisMessage) Acknowledge(event events.Event) error {
	return self.client.Publish(self.ReplyTo, event.EventName()).Err()
}

func (self *RedisMessage) Reject(err error) error {
	channel := self.ReplyTo
	if channel == "" {
		channel = "results"
	}
	return self.client.Publish(channel, err.Error()).Err()
}

package main

import (
	"errors"
	"log"

	"github.com/dhamidi/events"
	"github.com/dhamidi/events/examples/post"
	"github.com/dhamidi/events/examples/sessions"
	"github.com/dhamidi/events/examples/user"
	"gopkg.in/redis.v3"
)

var (
	ErrLoginRequired = errors.New("login required")
)

func main() {
	app := events.NewApplication()
	makeSignUpCommand := func() events.Command {
		return user.NewSignUp()
	}
	makeLogInCommand := func() events.Command {
		return &user.LogIn{}
	}

	sessionStore := sessions.NewInMemory()
	app.EventStore.Subscribe(sessionStore)

	app.RegisterCommand("/posts/publish", func() events.Command { return post.NewPublish() })
	app.RegisterCommand("/posts/draft", func() events.Command { return post.NewDraft() })
	app.RegisterCommand("/users/sign-up", makeSignUpCommand)
	app.RegisterCommand("/users/log-in", makeLogInCommand)

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	transport := NewRedisMessageBus("commands", redisClient)
	handleCommandWithLogin := func(msg events.Message) error {
		if !requireLogin(msg, sessionStore) {
			return ErrLoginRequired
		}
		return app.HandleCommand(msg)
	}
	events.HandleMessageFunc(transport, "/posts/draft", handleCommandWithLogin)
	events.HandleMessageFunc(transport, "/posts/publish", handleCommandWithLogin)
	events.HandleMessageFunc(transport, "/", app.HandleCommand)

	log.Fatal(transport.Listen())
}

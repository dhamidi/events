package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/dhamidi/events"
	"github.com/dhamidi/events/examples/post"
	"github.com/dhamidi/events/examples/sessions"
	"github.com/dhamidi/events/examples/user"
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

	transport := NewHTTPMessageBus()
	handleCommandWithLogin := func(msg events.Message) error {
		if !requireLogin(msg, sessionStore) {
			return ErrLoginRequired
		}
		return app.HandleCommand(msg)
	}
	events.HandleMessageFunc(transport, "/posts/draft", handleCommandWithLogin)
	events.HandleMessageFunc(transport, "/posts/publish", handleCommandWithLogin)
	events.HandleMessageFunc(transport, "/", app.HandleCommand)

	log.Fatal(http.ListenAndServe(":8080", transport))
}

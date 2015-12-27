package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

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

	handleCommandWithLogin := func(w http.ResponseWriter, req *http.Request) {
		msg := NewHTTPMessage(w, req)
		if !requireLogin(msg, sessionStore) {
			return
		}
		app.HandleCommand(msg)
	}
	handleCommand := func(w http.ResponseWriter, req *http.Request) {
		msg := NewHTTPMessage(w, req)
		app.HandleCommand(msg)
		fmt.Fprintf(os.Stderr, "DEBUG:\n%s\n", (func() []byte { data, _ := json.Marshal(sessionStore); return data })())
	}

	http.HandleFunc("/posts/publish", handleCommandWithLogin)
	http.HandleFunc("/posts/draft", handleCommandWithLogin)
	http.HandleFunc("/", handleCommand)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

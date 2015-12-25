package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

	app.RegisterCommand("/posts/draft", func() events.Command { return post.NewDraft() })
	app.RegisterCommand("/users/sign-up", makeSignUpCommand)
	app.RegisterCommand("/users/log-in", makeLogInCommand)
	http.HandleFunc("/posts/draft", func(w http.ResponseWriter, req *http.Request) {
		msg := NewHTTPMessage(w, req)
		if !requireLogin(msg, sessionStore) {
			return
		}
		app.HandleCommand(msg)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		msg := NewHTTPMessage(w, req)
		app.HandleCommand(msg)
		fmt.Fprintf(os.Stderr, "DEBUG:\n%s\n", (func() []byte { data, _ := json.Marshal(sessionStore); return data })())
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

type HTTPMessage struct {
	Request  *http.Request
	Response http.ResponseWriter
	body     []byte
}

func NewHTTPMessage(w http.ResponseWriter, req *http.Request) *HTTPMessage {
	return &HTTPMessage{
		Request:  req,
		Response: w,
	}
}

func (self *HTTPMessage) RoutingKey() string {
	return self.Request.URL.Path
}

func (self *HTTPMessage) ContentType() string {
	return self.Request.Header.Get("Content-Type")
}

func (self *HTTPMessage) Body() []byte {
	if self.body == nil {
		body, err := ioutil.ReadAll(self.Request.Body)
		if err != nil {
			self.body = []byte{}
		} else {
			self.body = body
		}
	}

	return self.body
}

func (self *HTTPMessage) String() string {
	return fmt.Sprintf("%s %s [%s]\n%s\n", self.Request.Method, self.RoutingKey(), self.ContentType(), self.Body())
}

func (self *HTTPMessage) Reject(err error) error {
	self.Response.Header().Set("Content-Type", "application/json")
	self.Response.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(self.Response).Encode(map[string]string{
		"error": err.Error(),
	})
	return nil
}

func (self *HTTPMessage) Acknowledge(event events.Event) error {
	switch e := event.(type) {
	case *user.LoggedIn:
		http.SetCookie(self.Response, &http.Cookie{
			Name:  "session_id",
			Value: e.SessionId,
		})
	}
	return nil
}

func (self *HTTPMessage) Cookie(name string) string {
	cookie, err := self.Request.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}

type SessionStore interface {
	IsLoggedIn(sessionId string) bool
}

func requireLogin(msg events.Message, sessionStore SessionStore) bool {
	if !sessionStore.IsLoggedIn(msg.Cookie("session_id")) {
		msg.Reject(ErrLoginRequired)
		return false
	}
	return true
}

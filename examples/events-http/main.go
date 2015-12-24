package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dhamidi/events"
	"github.com/dhamidi/events/examples/user"
)

func main() {
	app := events.NewApplication()
	makeSignUpCommand := func() events.Command {
		return user.NewSignUp()
	}
	app.RegisterCommand("/users/sign-up", makeSignUpCommand)
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		msg := NewHTTPMessage(req)
		err := app.HandleCommand(msg)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

type HTTPMessage struct {
	Request *http.Request
	body    []byte
}

func NewHTTPMessage(req *http.Request) *HTTPMessage {
	return &HTTPMessage{
		Request: req,
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

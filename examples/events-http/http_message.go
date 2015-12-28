package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dhamidi/events"
	"github.com/dhamidi/events/examples/user"
)

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
			Path:  "/",
		})
	}
	return nil
}

func (self *HTTPMessage) Header(name string) string {
	headers, ok := self.Request.Header[name]
	if ok {
		if len(headers) > 0 {
			return headers[0]
		}
	}

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
	if !sessionStore.IsLoggedIn(msg.Header("session_id")) {
		msg.Reject(ErrLoginRequired)
		return false
	}
	return true
}

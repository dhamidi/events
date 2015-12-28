package main

import (
	"net/http"

	"github.com/dhamidi/events"
)

type HTTPMessageBus struct {
	mux          *http.ServeMux
	errorHandler func(w http.ResponseWriter, req *http.Request, err error)
}

func NewHTTPMessageBus() *HTTPMessageBus {
	return &HTTPMessageBus{
		mux: http.NewServeMux(),
		errorHandler: func(w http.ResponseWriter, req *http.Request, err error) {
		},
	}
}

func (self *HTTPMessageBus) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	self.mux.ServeHTTP(w, req)
}

func (self *HTTPMessageBus) HandleMessage(routingKey string, handler events.MessageHandler) events.MessageBus {
	self.mux.HandleFunc(routingKey, self.httpHandlerFor(handler))
	return self
}

func (self *HTTPMessageBus) OnError(fn func(w http.ResponseWriter, req *http.Request, err error)) *HTTPMessageBus {
	self.errorHandler = fn
	return self
}

func (self *HTTPMessageBus) httpHandlerFor(handler events.MessageHandler) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		msg := NewHTTPMessage(w, req)
		if err := handler.HandleMessage(msg); err != nil {
			self.errorHandler(w, req, err)
		}
	}
}

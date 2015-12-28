# Description

This repository contains experiments with finding common abstractions
in message-based applications using Event Sourcing.

An example application is provided in [./examples/events-http/main.go].

# Components

## Commands

Commands are messages that request a change in the application state,
e.g. "Log In User".  In the framework presented here, a command is
always addressed to one specific domain object instance, which is
tasked with validating the state change.

## Events

Events document state changes that have taken place.  They are the
resulted of successfully processed commands.

In the implementation presented here, there is a one-to-one
correspondence between events and commands, i.e. every commands
results in either one event or one error.

## Aggregates

Aggregates process commands and maintain the necessary state to
enforce any business rules.  The current state of an aggregate can be
derived from its event history.

## Projections

Projections are optimized data structures for serving requests that
only read data, analoguous to a relational database-level view.  Since
projections are derived from the event history, they can be discarded
and rebuilt at any time.

## Event Store

The event store functions as an append-only log for events and
supports efficient retrieval of an aggregate's history.  Subscribers
can be attached to the event store in order to update projections.

## Application

The application accepts commands, restores the aggregate which handles
the command in its current state and lets the aggregate handle the
command.  A resulting event is appended to the event store and a
resulting error is reported to the command sender.

## Messages

Messages are used for communication between different components.
Semantically, a message is similar to an HTTP request, in that it is
self describing.  A message is addressed to one specific recipient, as
determined by its routing key.  The body of a message is an opaque
blob of bytes, with message headers describing the format of the
content.

Deviating a bit from HTTP's request-response semantics, messages can
be acknowledged (ACK) or rejected (NACK).  Acknowleding a message
means that is has been processed successfully and need not be
redelivered.  Rejecting a message means that the message could not be
processed at the time.  Redelivery is possible, depending on the
reason why the message has been rejected.

These ACK/NACK semantics can be easily mapped onto HTTP responses with
the respective response status codes.

Messages hide the details of the underlying transport mechanism, such
as HTTP, AMQP or Redis' Pub/Sub.

## Message Buses

Message buses listen to messages and route them to registered
handlers.

# Missing bits and pieces

Currently commands are sent to the application via a message bus.  The
messaging model should also be adopted for communication between the
event store and event subscribers.

It has to be possible to attach a new projection to a running system
and forward it to the current state.  A projection should be
forwardable to "now" from any point in time, e.g. if a projection
missed yesterday's events, restarting it should result in the
projection "catching-up" with all events until today.

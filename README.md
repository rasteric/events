# Events
*-- a simple event manager for Go*

[![GoDoc](https://godoc.org/github.com/rasteric/minidb/go?status.svg)](https://godoc.org/github.com/rasteric/events)
[![Go Report Card](https://goreportcard.com/badge/github.com/rasteric/minidb)](https://goreportcard.com/report/github.com/rasteric/events)
[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)

This package provides simple event management for Go. Events have a integer ID and a number of arguments. They can be published. Subscribers can subscribe to all events and then dispatch on their ID in a switch statement, or you may register event callbacks by event ID. The package stores all information internally for simplicity; therefore, it ought not be used in cases when different event systems are to be used within the same application. There is just one event system, the one you import when importing this package, and its state is initialized in an internal `init()` function.

Publishing and subscribing to events is thread-safe, but be aware that it's very easy to introduce race conditions and make wrong assumptions when using an internal event system. In case of doubt, your publisher and subscribers should use callbacks, because the publisher can guarantee that callbacks are processed at the right time (for example, within a function protected by a mutex).

## Usage

`go get github.com/rasteric/events`

`Pub(id int, args ...any)` publish event with ID `id` and any number of arguments.

`PubEvent(evt Event)` publish an event `evt` which may hold a number of arguments.

`PubLimited(id int, args ...any)` publish event with ID `id` with a built-in rate limitation. This only publishes the event if at least 200 milliseconds have passed since the last event of that type has been published. This may be used for cases when too frequent publications of miscellaneous events would have a performance impact and is not really needed.

`PubLimitedEvent(evt Event)` like `PubLimited` but for events.

`Sub(func(evt Event))` subscribe a consumer function. A subscriber will usually dispatch on the event ID in a switch statement to process the interesting cases.

`SubCb(id int, func(evt Event))` subscribe the given function to all events with ID `id`. This dispatches on the event ID internally, which is of course not faster than when you do it. However, callbacks are sometimes needed to avoid race conditions.

Events are implemented in a very straightforward way:

```
type Event struct {
	ID   int
	Args []any
}
```

You may use `Arg(idx int)` and `Count()` to get the arguments of an event, or use `Args` directly.

## Dependencies

None.
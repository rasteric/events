package events

import (
	"time"

	"github.com/mattn/go-pubsub"
)

var ps *pubsub.PubSub
var lastTime time.Time
var lastEvent Event

func init() {
	ps = pubsub.New()
	lastTime = time.Now()
	lastEvent = Event{}
}

// Event describes an event with an ID and an arbitrary datum associated with the ID.
// The subscriber is responsible to correctly type dispatch on the Datum, based on the event ID.
type Event struct {
	ID   int
	Args []any
}

// New creates a new event with given ID and optional arguments.
func New(id int, args ...any) Event {
	return Event{ID: id, Args: args}
}

// Arg returns the nth argument and true, nil and false if no argument with that index exists.
func (e Event) Arg(idx int) (any, bool) {
	if e.Args == nil {
		return nil, false
	}
	if idx < 0 || idx > len(e.Args)-1 {
		return nil, false
	}
	return e.Args[idx], true
}

// Count returns the event's number of arguments.
func (e Event) Count() int {
	if e.Args == nil {
		return 0
	}
	return len(e.Args)
}

// PubEvent publishes an event.
func PubEvent(evt Event) {
	ps.Pub(evt)
}

// Pub publishes an event given by the event ID and optional arguments.
func Pub(id int, args ...any) {
	PubEvent(New(id, args...))
}

// PubLimitedEvent publishes an event but drops the same event was published no longer than 200 milliseconds
// ago. This is used for unimportant, ephemeral events to rate limit them.
func PubLimitedEvent(evt Event) {
	now := time.Now()
	if evt.ID == lastEvent.ID && now.Sub(lastTime) < 200*time.Millisecond {
		return
	}
	lastTime = now
	lastEvent = evt
	PubEvent(evt)
}

// PubLimited publishes an event given by the ID and optional arguments.
func PubLimited(id int, args ...any) {
	PubLimitedEvent(New(id, args...))
}

// SubEvent subscribes to an event.
func SubEvent(fn func(evt Event)) {
	ps.Sub(fn)
}

// Sub subscribes to an event with a function that receives an event id and optional many arguments.
func Sub(fn func(id int, args ...any)) {

}

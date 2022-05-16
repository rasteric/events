package events

import (
	"sync"
	"time"
)

// TODO implement pubsub here myself
var lastTime time.Time
var lastEvent Event
var subscriptions []func(evt Event)
var mutex sync.RWMutex
var callbacks map[int][]func(evt Event)

func init() {
	subscriptions = make([]func(evt Event), 0)
	callbacks = make(map[int][]func(evt Event))
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
	mutex.RLock()
	defer mutex.RUnlock()
	cbs, ok := callbacks[evt.ID]
	if ok {
		for _, cb := range cbs {
			cb(evt)
		}
	}
	for _, fn := range subscriptions {
		fn(evt)
	}
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

// Sub subscribes to an event.
func Sub(fn func(evt Event)) {
	mutex.Lock()
	defer mutex.Unlock()
	subscriptions = append(subscriptions, fn)
}

// SubCb subscribes a callback function to an event ID. If an event with the given ID type
// occurs, the callback function is executed with the event as argument. This can be used for
// certain cases when Sub alone would lead to race conditions.
func SubCb(eventID int, fn func(evt Event)) {
	mutex.Lock()
	defer mutex.Unlock()
	cbs, ok := callbacks[eventID]
	if ok {
		cbs = append(cbs, fn)
		return
	}
	queue := make([]func(evt Event), 1)
	queue[0] = fn
	callbacks[eventID] = queue
}

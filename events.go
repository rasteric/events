package events

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

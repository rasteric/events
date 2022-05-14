package events


// Event describes an event with an ID and an arbitrary datum associated with the ID.
// The subscriber is responsible to correctly type dispatch on the Datum, based on the event type.
type Event struct {
	ID   int
	Arg1 any
	Arg2 any
	Arg3 any
	Arg4 any
	Arg5 any
	Arg6 any
	Err  error
}


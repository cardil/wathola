package event

// RegisterStep contains methods that register step event type
type RegisterStep interface {
	RegisterStep(step Step)
}

// RegisterFinished registers a finished event type
type RegisterFinished interface {
	RegisterFinished(finished Finished)
}

// Type says a type of an event
type Type interface {
	Type() string
}

package event

const (
	// StepType is a string type representation of step event
	StepType     = "com.github.cardil.wathola.step"
	// FinishedType os a string type representation of finished event
	FinishedType = "com.github.cardil.wathola.finished"
)

// Step is a event call at each step of verification
type Step struct {
	Number int
}

// Finished is step call after verification finishes
type Finished struct {
	Count int
}

// Type returns a type of a event
func (s Step) Type() string {
	return StepType
}

// Type returns a type of a event
func (f Finished) Type() string {
	return FinishedType
}

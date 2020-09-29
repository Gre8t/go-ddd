package aggregate

type eventRecorder struct {
	events []interface{}
}

func (e *eventRecorder) Record(event interface{}) {
	e.events = append(e.events, event)
}

func (e *eventRecorder) Events() []interface{} { return e.events }

func (e *eventRecorder) Clear() {
	e.events = []interface{}{}
}

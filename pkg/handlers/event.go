package handlers

import (
	"fmt"
	"time"
)

type Event struct {
	Key          string
	EventType    string
	ResourceType string

	EventTime time.Time

	//optional fields
	New, Old interface{}
}

func (e Event) String() string {
	return fmt.Sprintf("Key=%s, EventType=%s, ResourceType=%s, EventTime='%s'", e.Key, e.EventType, e.ResourceType, e.EventTime.String())
}

package connector

import (
	"fmt"
	"log"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/google/uuid"
)

type CEvent struct {
	*event.Event
}

var source string = "https://github.com/ec2ainun/poc-nats/business/connector"

func NewEvent(subject, eventType string, payload interface{}) *CEvent {
	cEvent := cloudevents.NewEvent()
	cEvent.SetID(uuid.New().String())
	cEvent.SetSource(source)
	cEvent.SetSpecVersion("1.0")
	cEvent.SetSubject(subject)
	cEvent.SetType(fmt.Sprintf("org.company.stream.%s", eventType))
	cEvent.SetTime(time.Now())
	err := cEvent.SetData("application/json", payload)
	if err != nil {
		log.Fatalf("err gen cloudevents: %s", err.Error())
	}
	return &CEvent{
		&cEvent,
	}
}

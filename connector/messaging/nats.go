package messaging

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type Stream struct {
	svc          *nats.Conn
	event        nats.JetStreamContext
	streamName   string
	consumerName string
	queueName    string
	pubType      string
	subType      string
}

type Handler func(msg *nats.Msg)

func NewNATSProvider(server *nats.Conn) *Stream {
	return &Stream{
		svc: server,
	}
}

func NewNATSQueueProvider(server *nats.Conn, queueName string) *Stream {
	return &Stream{
		svc:       server,
		queueName: queueName,
	}
}

func (s *Stream) Publish(subject string, payload []byte) error {
	if s.event == nil {
		if err := s.svc.Publish(subject, payload); err != nil {
			return fmt.Errorf("err publishing: %s", err.Error())
		}
	} else {
		if s.pubType == "async" {
			futureAck, err := s.event.PublishAsync(subject, payload)
			if err != nil {
				return fmt.Errorf("err publishing: %s", err.Error())
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			select {
			case ack := <-futureAck.Ok():
				log.Printf("Published in %s, on stream: %s, with sequence number: %d\n", subject, ack.Stream, ack.Sequence)
			case err = <-futureAck.Err():
				return fmt.Errorf("err publishing: %s", err.Error())
			case <-ctx.Done():
				return fmt.Errorf("err unable to finish in time:  %s", ctx.Err().Error())
			}
		} else {
			if _, err := s.event.Publish(subject, payload); err != nil {
				return fmt.Errorf("err publishing: %s", err.Error())
			}
		}
	}
	return nil
}

func (s *Stream) Subscribe(subject string, callback Handler) (*nats.Subscription, error) {
	var sub *nats.Subscription
	var err error

	if s.event == nil {
		if s.queueName == "" {
			sub, err = s.svc.Subscribe(subject, func(msg *nats.Msg) {
				callback(msg)
			})
		} else {
			sub, err = s.svc.QueueSubscribe(subject, s.queueName, func(msg *nats.Msg) {
				callback(msg)
			})
		}
	} else {
		if s.subType == "pull" {
			// pull
			sub, err = s.event.PullSubscribe(subject, s.consumerName, nats.Bind(s.streamName, s.consumerName))
		} else {
			// push
			if s.consumerName == "" {
				// create ephermal consumer
				sub, err = s.event.Subscribe(subject, func(msg *nats.Msg) {
					callback(msg)
				}, nats.BindStream(s.streamName))
			} else {
				// subscribe using an existing durable consumer
				if s.queueName == "" {
					sub, err = s.event.Subscribe(subject, func(msg *nats.Msg) {
						callback(msg)
					}, nats.Durable(s.consumerName))
				} else {
					sub, err = s.event.QueueSubscribe(subject, s.queueName, func(msg *nats.Msg) {
						callback(msg)
					}, nats.Durable(s.consumerName))
				}
			}
		}
	}
	s.svc.Flush()
	return sub, err
}

func (s *Stream) Request(subject string, payload []byte, callback Handler) error {
	msg, err := s.svc.Request(subject, payload, 2*time.Second)
	if err != nil {
		if s.svc.LastError() != nil {
			return fmt.Errorf("err requesting to %s:%v", subject, s.svc.LastError())
		}
		return fmt.Errorf("err requesting to %s:%v", subject, err)
	}
	callback(msg)
	return nil
}

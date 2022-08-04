package messaging

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/nats-io/nats.go"
)

type Stream struct {
	svc          *nats.Conn
	eventStream  nats.JetStreamContext
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

func (s *Stream) Publish(subject string, data string) {
	if s.eventStream == nil {
		if err := s.svc.Publish(subject, []byte(data)); err != nil {
			log.Fatalf("err publishing: %s", err.Error())
		}
	} else {
		if s.pubType == "async" {
			futureAck, err := s.eventStream.PublishAsync(subject, []byte(data))
			if err != nil {
				log.Fatalf("err publishing: %s", err.Error())
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			select {
			case ack := <-futureAck.Ok():
				log.Printf("Published on %s, on stream: %s, with sequence number: %d\n", subject, ack.Stream, ack.Sequence)
			case err = <-futureAck.Err():
				log.Fatalf("err publishing message: %s", err.Error())
			case <-ctx.Done():
				log.Fatalf("err unable to finish in time: %s", ctx.Err().Error())
			}
		} else {
			if _, err := s.eventStream.Publish(subject, []byte(data)); err != nil {
				log.Fatalf("err publishing: %s", err.Error())
			}
		}
	}
}

func (s *Stream) Subscribe(subject string, callback Handler) (*nats.Subscription, error) {
	var sub *nats.Subscription
	var err error

	if s.eventStream == nil {
		sub, err = s.svc.Subscribe(subject, func(msg *nats.Msg) {
			callback(msg)
		})
		s.svc.Flush()
	} else {
		if s.subType == "pull" {
			// pull
			sub, err = s.eventStream.PullSubscribe(subject, s.consumerName, nats.Bind(s.streamName, s.consumerName))
		} else {
			// push
			if s.consumerName == "" {
				// create ephermal consumer
				sub, err = s.eventStream.Subscribe(subject, func(msg *nats.Msg) {
					callback(msg)
				}, nats.BindStream(s.streamName))
			} else {
				// subscribe using an existing durable consumer
				if s.queueName == "" {
					sub, err = s.eventStream.Subscribe(subject, func(msg *nats.Msg) {
						callback(msg)
					}, nats.Durable(s.consumerName))
				} else {
					sub, err = s.eventStream.QueueSubscribe(subject, s.queueName, func(msg *nats.Msg) {
						callback(msg)
					}, nats.Durable(s.consumerName))
				}
			}

			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
			<-c
			if err := sub.Unsubscribe(); err != nil {
				log.Fatalf("err :%v\n", err)
			}
			s.svc.Drain()
		}
	}
	return sub, err
}

func (s *Stream) Request(subject string, data string, callback Handler) {
	msg, err := s.svc.Request(subject, []byte(data), 2*time.Second)
	if err != nil {
		if s.svc.LastError() != nil {
			log.Fatalf("err requesting to %s:%v\n", subject, s.svc.LastError())
		}
		log.Fatalf("err requesting to %s:%v\n", subject, err)
	}
	callback(msg)
}

func (s *Stream) Queue(subject, queue string, callback Handler) {
	_, err := s.svc.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		callback(msg)
	})
	s.svc.Flush()

	if err != nil {
		if s.svc.LastError() != nil {
			log.Fatalf("err requesting to %s,%s:%v\n", subject, queue, s.svc.LastError())
		}
		log.Fatalf("err subscribing queue to %s,%s:%v\n", subject, queue, err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	s.svc.Drain()
}

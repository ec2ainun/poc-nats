package messaging

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/nats-io/nats.go"
)

type Stream struct {
	svc *nats.Conn
}

type Handler func(msg *nats.Msg)

func NewNATSProvider(server *nats.Conn) *Stream {
	return &Stream{
		svc: server,
	}
}

func (s *Stream) Publish(subject string, data string) {
	if err := s.svc.Publish(subject, []byte(data)); err != nil {
		log.Fatal("error publishing:", err)
	}
}

func (s *Stream) Subscribe(subject string, callback Handler) {
	_, err := s.svc.Subscribe(subject, func(msg *nats.Msg) {
		callback(msg)
	})
	s.svc.Flush()
	if err != nil {
		log.Fatalf("err subscribing to %s:%v\n", subject, err)
	}
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
	log.Println()
	log.Printf("Draining...")
	s.svc.Drain()
	log.Printf("Exiting")
}

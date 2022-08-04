package messaging

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

func NewJetStreamProvider(server *nats.Conn, stream nats.JetStreamContext) *Stream {
	return &Stream{
		svc:         server,
		eventStream: stream,
	}
}

func NewJetStreamProducerProvider(server *nats.Conn, stream nats.JetStreamContext, pubType string) *Stream {
	return &Stream{
		svc:         server,
		eventStream: stream,
		pubType:     pubType,
	}
}

func NewJetStreamConsumerProvider(server *nats.Conn, stream nats.JetStreamContext, subType, streamName, consumerName, queueName string) *Stream {
	return &Stream{
		svc:          server,
		eventStream:  stream,
		subType:      subType,
		streamName:   streamName,
		consumerName: consumerName,
		queueName:    queueName,
	}
}

func (s *Stream) Create(name, subject string) error {
	subs := strings.Split(subject, ",")
	_, err := s.eventStream.AddStream(&nats.StreamConfig{
		Name:     name,
		Subjects: subs,
		MaxBytes: 256 << 20,
	})
	if err != nil {
		return err
	}
	log.Printf("created stream %s: consumed subject %s\n", name, subject)
	return nil
}

func (s *Stream) Update(name, subject string) error {
	subs := strings.Split(subject, ",")
	_, err := s.eventStream.UpdateStream(&nats.StreamConfig{
		Name:     name,
		Subjects: subs,
		MaxBytes: 256 << 20,
	})
	if err != nil {
		return err
	}
	log.Printf("updated stream %s: consumed subject %s\n", name, subject)
	return nil
}

func (s *Stream) Delete(name string) error {
	err := s.eventStream.DeleteStream(name)
	if err != nil {
		return err
	}
	log.Printf("deleted stream %s", name)
	return nil
}

func (s *Stream) CreateConsumer(_type, streamName, consumerName, subject, queue string) error {
	_, err := s.eventStream.StreamInfo(streamName)
	if err != nil && !errors.Is(err, nats.ErrStreamNotFound) {
		return err
	}
	if errors.Is(err, nats.ErrStreamNotFound) {
		s.Create(streamName, subject)
	}
	var consumerSetup *nats.ConsumerConfig
	if _type == "push" {
		consumerSetup = &nats.ConsumerConfig{
			Durable:        consumerName,
			DeliverSubject: fmt.Sprintf("deliver.%s", subject),
			AckPolicy:      nats.AckExplicitPolicy,
			FilterSubject:  subject,
		}
		if queue != "" {
			consumerSetup.DeliverGroup = queue
		}
	} else {
		// pull
		consumerSetup = &nats.ConsumerConfig{
			Durable:       consumerName,
			AckPolicy:     nats.AckExplicitPolicy,
			FilterSubject: subject,
		}
		if queue != "" {
			consumerSetup.DeliverGroup = queue
		}
	}
	_, err = s.eventStream.AddConsumer(streamName, consumerSetup)
	if err != nil {
		return err
	}
	log.Printf("created consumer %s", consumerName)
	return nil
}

func (s *Stream) DeleteConsumer(streamName, consumerName string) error {
	err := s.eventStream.DeleteConsumer(streamName, consumerName)
	if err != nil {
		return err
	}
	log.Printf("deleted consumer %s", consumerName)
	return nil
}

func (s *Stream) DelayPublish(subject, data string, delay int) {
	customMsg := nats.NewMsg(subject)
	customMsg.Data = []byte(data)
	customMsg.Header.Add("AI-Delayed-Time", fmt.Sprintf("%d", delay))
	if err := s.svc.PublishMsg(customMsg); err != nil {
		log.Fatalf("err publishing: %s", err.Error())
	}
}

func (s *Stream) IsDelayedMsg(msg *nats.Msg) bool {
	isDelayMsg := msg.Header.Get("AI-Delayed-Time")
	return isDelayMsg != ""
}

func (s *Stream) GetDelayValue(msg *nats.Msg) time.Duration {
	delayedValue := msg.Header.Get("AI-Delayed-Time")
	if delayedValue != "" {
		delay, err := strconv.ParseInt(delayedValue, 10, 64)
		if err != nil {
			log.Fatalf("err convert int: %s", err.Error())
			return 0
		}
		return time.Duration(delay) * time.Second
	}
	return 0
}

func (s *Stream) IsFirstMsg(msg *nats.Msg) bool {
	metadata, err := msg.Metadata()
	if err != nil {
		log.Fatalf("err read metadata: %s", err.Error())
		return false
	}
	return metadata.NumDelivered == 1
}

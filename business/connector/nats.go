package connector

import (
	"encoding/json"
	"fmt"
	"log"

	model "github.com/ec2ainun/poc-nats/business/dto"
	messaging "github.com/ec2ainun/poc-nats/connector/messaging"
	"github.com/nats-io/nats.go"
)

type streamConnectorImpl struct {
	svc *messaging.Stream
}

type StreamConnector interface {
	MonitorProfit(subject string)
	SendProfit(subject string, data model.ProfitInvestment) string
	ProcessProfit(subject string)
	QueueProcessProfit(subject, queue string)
	RequestProfit(subject string, data model.ProfitInvestment) string
	RespondProfit(subject, queue string)
}

func NewStreamConnector(_stream *messaging.Stream) StreamConnector {
	return &streamConnectorImpl{
		svc: _stream,
	}
}

func (s *streamConnectorImpl) MonitorProfit(subject string) {
	s.svc.Subscribe(subject, func(msg *nats.Msg) {
		// get only string
		log.Printf("Received from %s:%s\n", msg.Subject, string(msg.Data))
	})
}

func (s *streamConnectorImpl) SendProfit(subject string, data model.ProfitInvestment) string {
	// sent only string
	// s.svc.Publish(subject, data.String())

	// sent JSON string
	out, err := json.Marshal(data)
	if err != nil {
		log.Fatal("err marshal:", err)
	}
	s.svc.Publish(subject, string(out))
	return data.String()
}

func (s *streamConnectorImpl) ProcessProfit(subject string) {
	s.svc.Subscribe(subject, func(msg *nats.Msg) {
		// get only string
		// log.Printf("Received from %s:%s\n", msg.Subject, string(msg.Data))

		// get JSON string
		data := model.ProfitInvestment{}
		err := json.Unmarshal([]byte(msg.Data), &data)
		if err != nil {
			log.Fatal("err unmarshal:", err)
		}
		log.Printf("Received from %s:%s\n", msg.Subject, data.String())
	})
}

func (s *streamConnectorImpl) QueueProcessProfit(subject, queue string) {
	i := 0
	s.svc.Queue(subject, queue, func(msg *nats.Msg) {
		i++
		// get only string
		// log.Printf("[#%d] Received on [%s]: '%s'\n", i, msg.Subject, string(msg.Data))

		// get JSON string
		data := model.ProfitInvestment{}
		err := json.Unmarshal([]byte(msg.Data), &data)
		if err != nil {
			log.Fatal("err unmarshal:", err)
		}
		log.Printf("[#%d] Received on [%s]: '%s'\n", i, msg.Subject, data.String())
	})
}

func (s *streamConnectorImpl) RequestProfit(subject string, data model.ProfitInvestment) string {
	// sent only string
	// s.svc.Request(subject, data.String(), processingReply)

	// sent JSON string
	out, err := json.Marshal(data)
	if err != nil {
		log.Fatal("err marshal:", err)
	}
	s.svc.Request(subject, string(out), processingReply)
	return data.String()
}

func processingReply(msg *nats.Msg) {
	log.Printf("Received data from %s: %s\n", msg.Subject, string(msg.Data))
}

func (s *streamConnectorImpl) RespondProfit(subject, queue string) {
	i := 0
	s.svc.Queue(subject, queue, func(msg *nats.Msg) {
		i++

		// get only string
		// log.Printf("[#%d] Received on [%s]: '%s'\n", i, msg.Subject, string(msg.Data))
		// resp := 10.0
		// sent respond back
		// msg.Respond([]byte(fmt.Sprintf("%.2f", resp)))
		// log.Printf("[#%d] Replied to [%s]: '%s'\n", i, msg.Subject, fmt.Sprintf("%.2f", resp))

		// get JSON string
		data := model.ProfitInvestment{}
		err := json.Unmarshal([]byte(msg.Data), &data)
		if err != nil {
			log.Fatal("err unmarshal:", err)
		}
		log.Printf("[#%d] Received on [%s]: '%s'\n", i, msg.Subject, data.String())

		// sent respond back
		resp := data.ProfitAmount - 10
		msg.Respond([]byte(fmt.Sprintf("%.2f", resp)))
		log.Printf("[#%d] Replied to [%s]: '%s'\n", i, msg.Subject, fmt.Sprintf("%.2f", resp))
	})
}

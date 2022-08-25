package connector

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	model "github.com/ec2ainun/poc-nats/business/dto"
	messaging "github.com/ec2ainun/poc-nats/connector/messaging"
	"github.com/nats-io/nats.go"
)

type streamConnectorImpl struct {
	svc *messaging.Stream
}

type StreamConnector interface {
	MonitorAll(subject string)
	SendProfit(subject string, data model.ProfitInvestment) string
	ProcessProfit(subject string)
	QueueProcessProfit(subject, queue string)
	RequestProfit(subject string, data model.ProfitInvestment) string
	RespondProfit(subject, queue string)
	BatchProcessProfit(subject string, batch int)
	PushProcessProfit(subject string)
	DelayedProfit(subject string, data model.ProfitInvestment, delay int) string
}

func NewStreamConnector(_stream *messaging.Stream) StreamConnector {
	return &streamConnectorImpl{
		svc: _stream,
	}
}

func (s *streamConnectorImpl) MonitorAll(subject string) {
	_, err := s.svc.Subscribe(subject, func(msg *nats.Msg) {
		log.Printf("Received from %s: %s\n", msg.Subject, string(msg.Data))
	})
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (s *streamConnectorImpl) SendProfit(subject string, data model.ProfitInvestment) string {
	out, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("err marshal: %s", err.Error())
	}
	err = s.svc.Publish(subject, string(out))
	if err != nil {
		log.Fatalf(err.Error())
	}
	return data.String()
}

func (s *streamConnectorImpl) ProcessProfit(subject string) {
	_, err := s.svc.Subscribe(subject, func(msg *nats.Msg) {
		data := model.ProfitInvestment{}
		if err := json.Unmarshal([]byte(msg.Data), &data); err != nil {
			log.Fatalf("err unmarshal: %s", err.Error())
		}
		log.Printf("Received from %s: %s\n", msg.Subject, data.String())
	})
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (s *streamConnectorImpl) QueueProcessProfit(subject, queue string) {
	i := 0
	err := s.svc.Queue(subject, queue, func(msg *nats.Msg) {
		i++
		data := model.ProfitInvestment{}
		if err := json.Unmarshal([]byte(msg.Data), &data); err != nil {
			log.Fatalf("err unmarshal: %s", err.Error())
		}
		log.Printf("[#%d] Received on [%s, %s]: '%s'\n", i, msg.Subject, queue, data.String())
	})
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (s *streamConnectorImpl) RequestProfit(subject string, data model.ProfitInvestment) string {
	out, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("err marshal: %s", err.Error())
	}
	err = s.svc.Request(subject, string(out), processingReply)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return data.String()
}

func processingReply(msg *nats.Msg) {
	log.Printf("Received data from %s: %s\n", msg.Subject, string(msg.Data))
}

func (s *streamConnectorImpl) RespondProfit(subject, queue string) {
	i := 0
	err := s.svc.Queue(subject, queue, func(msg *nats.Msg) {
		i++
		data := model.ProfitInvestment{}
		if err := json.Unmarshal([]byte(msg.Data), &data); err != nil {
			log.Fatalf("err unmarshal: %s", err.Error())
		}
		log.Printf("[#%d] Received on [%s, %s]: '%s'\n", i, msg.Subject, queue, data.String())

		// sent respond back
		resp := data.ProfitAmount - 10
		msg.Respond([]byte(fmt.Sprintf("%.2f", resp)))
		log.Printf("[#%d] Replied to [%s]: '%s'\n", i, msg.Subject, fmt.Sprintf("%.2f", resp))
	})
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (s *streamConnectorImpl) BatchProcessProfit(subject string, batch int) {
	sub, err := s.svc.Subscribe(subject, nil)
	if err != nil {
		log.Fatalf(err.Error())
	}
	i := 0
	for {
		msgs, err := sub.Fetch(batch, nats.MaxWait(2*time.Second))
		if err != nil && errors.Is(err, nats.ErrConnectionClosed) {
			return
		}
		if err != nil && errors.Is(err, nats.ErrTimeout) {
			log.Printf("all message processed")
			os.Exit(0)
		}
		if err != nil {
			log.Fatalf("err fetch : %s", err.Error())
		}

		for _, msg := range msgs {
			data := model.ProfitInvestment{}
			if err := json.Unmarshal([]byte(msg.Data), &data); err != nil {
				log.Fatalf("err unmarshal : %s", err.Error())
			}
			log.Printf("Received from %s: %s\n", msg.Subject, data.String())
			if err := msg.AckSync(); err != nil {
				log.Fatalf("err ack : %s", err.Error())
			}
			i++
			log.Printf("Processed message %d", i)
		}

		// Poll every 5 second
		time.Sleep(5 * time.Second)
	}
}

func (s *streamConnectorImpl) PushProcessProfit(subject string) {
	i := 0
	_, err := s.svc.Subscribe(subject, func(msg *nats.Msg) {
		// check delayed msg
		if s.svc.IsDelayedMsg(msg) {
			delayedVal := s.svc.GetDelayValue(msg)
			if s.svc.IsFirstMsg(msg) {
				// delay msg
				if err := msg.NakWithDelay(delayedVal); err != nil {
					log.Fatalf("err nack :%s", err.Error())
				}
				data := model.ProfitInvestment{}
				if err := json.Unmarshal([]byte(msg.Data), &data); err != nil {
					log.Fatalf("err unmarshal: %s", err.Error())
				}
				log.Printf("msg %d for %s, delayed for %s", data.Id, data.Receiver, delayedVal)
			} else {
				i++
				data := model.ProfitInvestment{}
				if err := json.Unmarshal([]byte(msg.Data), &data); err != nil {
					log.Fatalf("err unmarshal: %s", err.Error())
				}
				log.Printf("Received delayed msg from [%s, %s]: %s\n", msg.Sub.Subject, msg.Sub.Queue, data.String())
				if err := msg.AckSync(); err != nil {
					log.Fatalf("err ack :%s", err.Error())
				}
				log.Printf("Processed message %d", i)
			}
		} else {
			i++
			data := model.ProfitInvestment{}
			if err := json.Unmarshal([]byte(msg.Data), &data); err != nil {
				log.Fatalf("err unmarshal: %s", err.Error())
			}
			log.Printf("Received from [%s, %s]: %s\n", msg.Sub.Subject, msg.Sub.Queue, data.String())
			if err := msg.Ack(); err != nil {
				log.Fatalf("err ack :%s", err.Error())
			}
			log.Printf("Processed message %d", i)
		}
	})
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (s *streamConnectorImpl) DelayedProfit(subject string, data model.ProfitInvestment, delay int) string {
	out, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("err marshal: %s", err.Error())
	}
	err = s.svc.DelayPublish(subject, string(out), delay)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return data.String()
}

package connector

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	model "github.com/ec2ainun/poc-nats/business/dto"
	messaging "github.com/ec2ainun/poc-nats/connector/messaging"
	"github.com/nats-io/nats.go"
)

// concurrency control, adapted based on certain boundary
// such as db connection, max ack pending etc
var CNCR int = 1000

type streamConnectorImpl struct {
	svc *messaging.Stream
}

type StreamConnector interface {
	MonitorAll(subject string)
	SendProfit(subject string, data model.ProfitInvestment) string
	ProcessProfit(subject string)
	ChanProcessProfit(subject string)
	QueueProcessProfit(subject string)
	ChanQueueProcessProfit(subject string)
	RequestProfit(subject string, data model.ProfitInvestment) string
	RespondProfit(subject string)
	ChanRespondProfit(subject string)
	BatchProcessProfit(subject string, batch int)
	PushProcessProfit(subject string)
	ChanPushProcessProfit(subject string)
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

// https://github.com/cloudevents/spec/blob/main/cloudevents/formats/cloudevents.json
func (s *streamConnectorImpl) SendProfit(subject string, data model.ProfitInvestment) string {
	cEvent := NewEvent(subject, "SendProfit", data)
	out, err := cEvent.MarshalJSON()
	if err != nil {
		log.Fatalf("err marshal: %s", err.Error())
	}
	err = s.svc.Publish(subject, out)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return data.String()
}

func (s *streamConnectorImpl) ProcessProfit(subject string) {
	_, err := s.svc.Subscribe(subject, func(msg *nats.Msg) {
		cEvent := cloudevents.Event{}
		data := model.ProfitInvestment{}
		if err := json.Unmarshal([]byte(msg.Data), &cEvent); err != nil {
			log.Fatalf("err unmarshal: %s", err.Error())
		}
		if err := cEvent.DataAs(&data); err != nil {
			fmt.Printf("err got cloudevents: %s\n", err.Error())
		}
		log.Printf("Received from %s: %s\n", msg.Subject, data.String())
	})
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (s *streamConnectorImpl) ChanProcessProfit(subject string) {
	msgs := make(chan *nats.Msg, CNCR)
	_, err := s.svc.ChanSubscribe(subject, msgs)
	if err != nil {
		log.Fatalf(err.Error())
	}
	go func() {
		for msg := range msgs {
			cEvent := cloudevents.Event{}
			data := model.ProfitInvestment{}
			if err := json.Unmarshal([]byte(msg.Data), &cEvent); err != nil {
				log.Fatalf("err unmarshal: %s", err.Error())
			}
			if err := cEvent.DataAs(&data); err != nil {
				fmt.Printf("err got cloudevents: %s\n", err.Error())
			}
			log.Printf("Received from %s: %s\n", msg.Subject, data.String())
		}
	}()
}

func (s *streamConnectorImpl) QueueProcessProfit(subject string) {
	_, err := s.svc.Subscribe(subject, func(msg *nats.Msg) {
		cEvent := cloudevents.Event{}
		data := model.ProfitInvestment{}
		if err := json.Unmarshal([]byte(msg.Data), &cEvent); err != nil {
			log.Fatalf("err unmarshal: %s", err.Error())
		}
		if err := cEvent.DataAs(&data); err != nil {
			fmt.Printf("err got cloudevents: %s\n", err.Error())
		}
		log.Printf("Received from [%s, %s]: %s\n", msg.Sub.Subject, msg.Sub.Queue, data.String())
	})
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (s *streamConnectorImpl) ChanQueueProcessProfit(subject string) {
	msgs := make(chan *nats.Msg, CNCR)
	_, err := s.svc.ChanSubscribe(subject, msgs)
	if err != nil {
		log.Fatalf(err.Error())
	}
	go func() {
		for msg := range msgs {
			cEvent := cloudevents.Event{}
			data := model.ProfitInvestment{}
			if err := json.Unmarshal([]byte(msg.Data), &cEvent); err != nil {
				log.Fatalf("err unmarshal: %s", err.Error())
			}
			if err := cEvent.DataAs(&data); err != nil {
				fmt.Printf("err got cloudevents: %s\n", err.Error())
			}
			log.Printf("Received from [%s, %s]: %s\n", msg.Sub.Subject, msg.Sub.Queue, data.String())
		}
	}()
}

func (s *streamConnectorImpl) RequestProfit(subject string, data model.ProfitInvestment) string {
	cEvent := NewEvent(subject, "RequestProfit", data)
	out, err := cEvent.MarshalJSON()
	if err != nil {
		log.Fatalf("err marshal: %s", err.Error())
	}
	err = s.svc.Request(subject, out, processingReply)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return data.String()
}

func processingReply(msg *nats.Msg) {
	log.Printf("Received data from %s: %s\n", msg.Subject, string(msg.Data))
}

func (s *streamConnectorImpl) RespondProfit(subject string) {
	_, err := s.svc.Subscribe(subject, func(msg *nats.Msg) {
		cEvent := cloudevents.Event{}
		data := model.ProfitInvestment{}
		if err := json.Unmarshal([]byte(msg.Data), &cEvent); err != nil {
			log.Fatalf("err unmarshal: %s", err.Error())
		}
		if err := cEvent.DataAs(&data); err != nil {
			fmt.Printf("err got cloudevents: %s\n", err.Error())
		}
		log.Printf("Received from [%s, %s]: %s\n", msg.Sub.Subject, msg.Sub.Queue, data.String())

		// sent respond back
		resp := data.ProfitAmount - 10
		msg.Respond([]byte(fmt.Sprintf("%.2f", resp)))
		log.Printf("Replied to [%s]: %s\n", msg.Subject, fmt.Sprintf("%.2f", resp))
	})
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (s *streamConnectorImpl) ChanRespondProfit(subject string) {
	msgs := make(chan *nats.Msg, CNCR)
	_, err := s.svc.ChanSubscribe(subject, msgs)
	if err != nil {
		log.Fatalf(err.Error())
	}
	go func() {
		for msg := range msgs {
			cEvent := cloudevents.Event{}
			data := model.ProfitInvestment{}
			if err := json.Unmarshal([]byte(msg.Data), &cEvent); err != nil {
				log.Fatalf("err unmarshal: %s", err.Error())
			}
			if err := cEvent.DataAs(&data); err != nil {
				fmt.Printf("err got cloudevents: %s\n", err.Error())
			}
			log.Printf("Received from [%s, %s]: %s\n", msg.Sub.Subject, msg.Sub.Queue, data.String())

			// sent respond back
			resp := data.ProfitAmount - 10
			msg.Respond([]byte(fmt.Sprintf("%.2f", resp)))
			log.Printf("Replied to [%s]: %s\n", msg.Subject, fmt.Sprintf("%.2f", resp))
		}
	}()
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
			cEvent := cloudevents.Event{}
			data := model.ProfitInvestment{}
			if err := json.Unmarshal([]byte(msg.Data), &cEvent); err != nil {
				log.Fatalf("err unmarshal: %s", err.Error())
			}
			if err := cEvent.DataAs(&data); err != nil {
				fmt.Printf("err got cloudevents: %s\n", err.Error())
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
	_, err := s.svc.Subscribe(subject, func(msg *nats.Msg) {
		// check delayed msg
		if s.svc.IsDelayedMsg(msg) {
			delayedVal := s.svc.GetDelayValue(msg)
			if s.svc.IsFirstMsg(msg) {
				// delay msg
				if err := msg.NakWithDelay(delayedVal); err != nil {
					log.Fatalf("err nack :%s", err.Error())
				}
				cEvent := cloudevents.Event{}
				data := model.ProfitInvestment{}
				if err := json.Unmarshal([]byte(msg.Data), &cEvent); err != nil {
					log.Fatalf("err unmarshal: %s", err.Error())
				}
				if err := cEvent.DataAs(&data); err != nil {
					fmt.Printf("err got cloudevents: %s\n", err.Error())
				}
				log.Printf("msg %d for %s, delayed for %s", data.Id, data.Receiver, delayedVal)
			} else {
				cEvent := cloudevents.Event{}
				data := model.ProfitInvestment{}
				if err := json.Unmarshal([]byte(msg.Data), &cEvent); err != nil {
					log.Fatalf("err unmarshal: %s", err.Error())
				}
				if err := cEvent.DataAs(&data); err != nil {
					fmt.Printf("err got cloudevents: %s\n", err.Error())
				}
				log.Printf("Received delayed msg from [%s, %s]: %s\n", msg.Sub.Subject, msg.Sub.Queue, data.String())
				if err := msg.AckSync(); err != nil {
					log.Fatalf("err ack :%s", err.Error())
				}
				log.Printf("Processed message")
			}
		} else {
			cEvent := cloudevents.Event{}
			data := model.ProfitInvestment{}
			if err := json.Unmarshal([]byte(msg.Data), &cEvent); err != nil {
				log.Fatalf("err unmarshal: %s", err.Error())
			}
			if err := cEvent.DataAs(&data); err != nil {
				fmt.Printf("err got cloudevents: %s\n", err.Error())
			}
			log.Printf("Received from [%s, %s]: %s\n", msg.Sub.Subject, msg.Sub.Queue, data.String())
			if err := msg.Ack(); err != nil {
				log.Fatalf("err ack :%s", err.Error())
			}
			log.Printf("Processed message")
		}
	})
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (s *streamConnectorImpl) ChanPushProcessProfit(subject string) {
	msgs := make(chan *nats.Msg, CNCR)
	_, err := s.svc.ChanSubscribe(subject, msgs)
	if err != nil {
		log.Fatalf(err.Error())
	}
	go func() {
		for msg := range msgs {
			// check delayed msg
			if s.svc.IsDelayedMsg(msg) {
				delayedVal := s.svc.GetDelayValue(msg)
				if s.svc.IsFirstMsg(msg) {
					// delay msg
					if err := msg.NakWithDelay(delayedVal); err != nil {
						log.Fatalf("err nack :%s", err.Error())
					}
					cEvent := cloudevents.Event{}
					data := model.ProfitInvestment{}
					if err := json.Unmarshal([]byte(msg.Data), &cEvent); err != nil {
						log.Fatalf("err unmarshal: %s", err.Error())
					}
					if err := cEvent.DataAs(&data); err != nil {
						fmt.Printf("err got cloudevents: %s\n", err.Error())
					}
					log.Printf("msg %d for %s, delayed for %s", data.Id, data.Receiver, delayedVal)
				} else {
					cEvent := cloudevents.Event{}
					data := model.ProfitInvestment{}
					if err := json.Unmarshal([]byte(msg.Data), &cEvent); err != nil {
						log.Fatalf("err unmarshal: %s", err.Error())
					}
					if err := cEvent.DataAs(&data); err != nil {
						fmt.Printf("err got cloudevents: %s\n", err.Error())
					}
					log.Printf("Received delayed msg from [%s, %s]: %s\n", msg.Sub.Subject, msg.Sub.Queue, data.String())
					if err := msg.AckSync(); err != nil {
						log.Fatalf("err ack :%s", err.Error())
					}
					log.Printf("Processed message")
				}
			} else {
				cEvent := cloudevents.Event{}
				data := model.ProfitInvestment{}
				if err := json.Unmarshal([]byte(msg.Data), &cEvent); err != nil {
					log.Fatalf("err unmarshal: %s", err.Error())
				}
				if err := cEvent.DataAs(&data); err != nil {
					fmt.Printf("err got cloudevents: %s\n", err.Error())
				}
				log.Printf("Received from [%s, %s]: %s\n", msg.Sub.Subject, msg.Sub.Queue, data.String())
				if err := msg.Ack(); err != nil {
					log.Fatalf("err ack :%s", err.Error())
				}
				log.Printf("Processed message")
			}
		}
	}()
}

func (s *streamConnectorImpl) DelayedProfit(subject string, data model.ProfitInvestment, delay int) string {
	cEvent := NewEvent(subject, "DelayedProfit", data)
	out, err := cEvent.MarshalJSON()
	if err != nil {
		log.Fatalf("err marshal: %s", err.Error())
	}
	err = s.svc.DelayPublish(subject, out, delay)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return data.String()
}

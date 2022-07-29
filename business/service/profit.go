package service

import (
	connector "github.com/ec2ainun/poc-nats/business/connector"
)

type profitServiceImpl struct {
	Stream connector.StreamConnector
}

type ProfitService interface {
	Monitor(subject string)
	PublishProfit(subject string, id int) string
	SubscribeProfit(subject string)
	QueueSubscribeProfit(subject, queue string)
	RequestProfit(subject string, id int) string
	ReplyProfit(subject, queue string)
}

func NewProfitService(_stream connector.StreamConnector) ProfitService {
	return &profitServiceImpl{
		Stream: _stream,
	}
}

func (p *profitServiceImpl) Monitor(subject string) {
	p.Stream.MonitorProfit(subject)
}

func (p *profitServiceImpl) PublishProfit(subject string, id int) string {
	data := GetRandomProfitDistribute(id)
	dataString := p.Stream.SendProfit(subject, data)
	return dataString
}

func (p *profitServiceImpl) SubscribeProfit(subject string) {
	p.Stream.ProcessProfit(subject)
}

func (p *profitServiceImpl) QueueSubscribeProfit(subject, queue string) {
	p.Stream.QueueProcessProfit(subject, queue)
}

func (p *profitServiceImpl) RequestProfit(subject string, id int) string {
	data := GetRandomProfitDistribute(id)
	dataString := p.Stream.RequestProfit(subject, data)
	return dataString
}

func (p *profitServiceImpl) ReplyProfit(subject, queue string) {
	p.Stream.RespondProfit(subject, queue)
}

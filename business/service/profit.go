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
	QueueSubscribeProfit(subject string)
	RequestProfit(subject string, id int) string
	ReplyProfit(subject string)
	BatchSubscribeProfit(subject string, batch int)
	PushSubscribeProfit(subject string)
	PublishDelayedProfit(subject string, i, delay int) string
}

func NewProfitService(_stream connector.StreamConnector) ProfitService {
	return &profitServiceImpl{
		Stream: _stream,
	}
}

func (p *profitServiceImpl) Monitor(subject string) {
	p.Stream.MonitorAll(subject)
}

func (p *profitServiceImpl) PublishProfit(subject string, id int) string {
	data := GetRandomProfitDistribute(id)
	dataString := p.Stream.SendProfit(subject, data)
	return dataString
}

func (p *profitServiceImpl) SubscribeProfit(subject string) {
	p.Stream.ProcessProfit(subject)
}

func (p *profitServiceImpl) QueueSubscribeProfit(subject string) {
	p.Stream.QueueProcessProfit(subject)
}

func (p *profitServiceImpl) RequestProfit(subject string, id int) string {
	data := GetRandomProfitDistribute(id)
	dataString := p.Stream.RequestProfit(subject, data)
	return dataString
}

func (p *profitServiceImpl) ReplyProfit(subject string) {
	p.Stream.RespondProfit(subject)
}

func (p *profitServiceImpl) BatchSubscribeProfit(subject string, batch int) {
	p.Stream.BatchProcessProfit(subject, batch)
}

func (p *profitServiceImpl) PushSubscribeProfit(subject string) {
	p.Stream.PushProcessProfit(subject)
}

func (p *profitServiceImpl) PublishDelayedProfit(subject string, i, delay int) string {
	data := GetRandomProfitDistribute(i)
	dataString := p.Stream.DelayedProfit(subject, data, delay)
	return dataString
}

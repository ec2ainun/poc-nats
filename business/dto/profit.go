package dto

import (
	"fmt"
	"math"
)

type ProfitInvestment struct {
	Id           int
	ProfitAmount float64
	ProductName  string
	Period       string
	Currency     string
	Sender       string
	Receiver     string
}

func (p ProfitInvestment) String() string {
	return fmt.Sprintf("[%d]: %s sent a payment of %0.2f %s from portofolio product %s to %s at period %s", p.Id, p.Sender, p.ProfitAmount, p.Currency, p.ProductName, p.Receiver, p.Period)
}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

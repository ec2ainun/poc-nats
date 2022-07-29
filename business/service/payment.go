package service

import (
	"math/rand"
	"time"

	model "github.com/ec2ainun/poc-nats/business/dto"
	"github.com/google/uuid"
)

const maxRandomAmount float64 = 500
const currency string = "USD"

var product = [4]string{
	"Moon",
	"Star",
	"Supernova",
	"Galaxy",
}

func GetRandomProfitDistribute(id int) model.ProfitInvestment {
	profit := model.RoundFloat(rand.Float64()*maxRandomAmount, 2)
	r := rand.Intn(4)
	pI := model.ProfitInvestment{
		Id:           id,
		ProfitAmount: profit,
		ProductName:  product[r],
		Period:       time.Now().Format("2006-01-02 15:04:05"),
		Currency:     currency,
		Sender:       "system",
		Receiver:     uuid.New().String(),
	}

	return pI
}

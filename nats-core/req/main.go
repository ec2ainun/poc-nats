package main

import (
	"flag"
	"log"
	"runtime"
	"time"

	connector "github.com/ec2ainun/poc-nats/business/connector"
	service "github.com/ec2ainun/poc-nats/business/service"
	config "github.com/ec2ainun/poc-nats/business/utils"
	stream "github.com/ec2ainun/poc-nats/connector/messaging"
	messaging "github.com/ec2ainun/poc-nats/foundation/messaging"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

const genMessages = 500

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
	// Exit the main goroutine but allow subscribe to continue running
	runtime.Goexit()
}

func run() error {
	goal := "POC NATS Core requestor"
	log.Println(goal)

	flag.Parse()
	args := flag.Args()
	subject := args[0]

	conf, err := config.Read("config/config.yaml")
	if err != nil {
		return errors.Wrap(err, "err open config")
	}

	opts := []nats.Option{nats.Name(goal)}
	opts = messaging.SetupConnOptions(opts)

	nc, err := messaging.Open(false, conf, opts)
	if err != nil {
		return err
	}
	defer nc.Close()

	streamProvider := stream.NewNATSProvider(nc)
	streamSvc := connector.NewStreamConnector(streamProvider)
	profitSvc := service.NewProfitService(streamSvc)

	for i := 0; i < genMessages; i++ {
		p := profitSvc.RequestProfit(subject, i)
		log.Printf("Requested on %s:%s\n", subject, p)
		time.Sleep(2 * time.Second)
	}

	return nil
}

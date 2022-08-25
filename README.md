# NATS

stands for Neural Autonomic Transport System

the rationale why choose NATS over other messaging system comes from [this](https://docs.google.com/document/d/1N9-QybgQbp2Om2WLtFE-d43jJHxlW17XIRGnx-hCRtU/edit?usp=sharing)

## Getting Started

### Installation

- nats server

```sh
brew install nats-server
```

[another way to install nats-server](https://docs.nats.io/running-a-nats-service/introduction/installation)

- nats cli

```sh
brew tap nats-io/nats-tools
brew install nats-io/nats-tools/nats
```

[another way to install nats cli](https://docs.nats.io/using-nats/nats-tools/nats_cli)

- nats monitoring

```sh
go install github.com/nats-io/nats-top@latest
```

[another way to monitor](https://docs.nats.io/reference/faq#how-can-i-monitor-my-nats-cluster)

## Core NATS

- start nats server first

```sh
nats-server --config nats-core/nats.conf
```

### Subscribe to all event

- subscribe

```sh
go run nats-core/monitor/main.go ">"
```

### Publish Subscribe

- publish

```sh
go run nats-core/pub/main.go "profit"
```

- subscribe

```sh
go run nats-core/sub/main.go "profit"
```

### Fan In

- publish

```sh
go run nats-core/pub/main.go "profit.sg"
```

```sh
go run nats-core/pub/main.go "profit.id"
```

- subscribe

```sh
go run nats-core/sub/main.go "profit.*"
```

### Fan Out

- publish

```sh
go run nats-core/pub/main.go "profit.id"
```

- subscribe

```sh
go run nats-core/sub/main.go "profit.id"
```

```sh
go run nats-core/sub/main.go "profit.id"
```

### Queue Load Balanced

- publish

```sh
go run nats-core/pub/main.go "profit.id"
```

- subscribe using queue

```sh
go run nats-core/subq/main.go -q poc "profit.id" 
```

```sh
go run nats-core/subq/main.go -q poc "profit.id" 
```

### Request Reply

- start replier first

```sh
go run nats-core/rep/main.go -q poc "profit.get" 
```

- request

```sh
go run nats-core/req/main.go "profit.get"
```

## Jetstream

- start nats server first

```sh
nats-server --config nats-jetstream/js.conf
```

### Stream

what kind event data to persist

- check list of stream

```sh
nats str ls --server 0.0.0.0:5222
```

- check report of stream

```sh
nats stream report --server 0.0.0.0:5222
```

- create stream `profit` in which interest to persist event data from subject `profit.*`

```sh
go run nats-jetstream/manager/main.go --cmd add -n profit -s "profit.*"
```

- check stream detail

```sh
nats stream info profit --server 0.0.0.0:5222
```

- update stream `profit` in which interest to persist event data from subject `profit.sg profit.id`

```sh
go run nats-jetstream/manager/main.go --cmd update -n profit -s "profit.sg,profit.id"
```

- delete stream

```sh
go run nats-jetstream/manager/main.go --cmd delete -n profit 
```

### Publish

- sync publish

```sh
go run nats-jetstream/producer/sync/main.go -msgs 1000 "profit.id"
```

- async publish

```sh
go run nats-jetstream/producer/async/main.go -msgs 1000 "profit.id"
```

### Consumer

how application getting a messages from stream

- check list of consumer

```sh
nats consumer ls --server 0.0.0.0:5222
```

- push consumer

```sh
go run nats-jetstream/consumer/main.go --cmd add-push -sn profit -cn one-push-cprofit_id -s profit.id
```

- push consumer with shared load balanced queue

```sh
go run nats-jetstream/consumer/main.go --cmd add-push -sn profit -cn queue-push-cprofit_id -qn poc -s profit.id
```

- pull consumer is automatically shared if there are multiple subscriber

```sh
go run nats-jetstream/consumer/main.go --cmd add-pull -sn profit -cn pull-cprofit_id -s profit.id
```

- check consumer detail

```sh
nats consumer info --server 0.0.0.0:5222
```

- delete consumer

```sh
go run nats-jetstream/consumer/main.go --cmd add-push -sn profit -cn test-delete -qn poc -s profit.id
```

```sh
go run nats-jetstream/consumer/main.go --cmd delete -sn profit -cn test-delete
```

### Ephermeral (direct) consumer

- publish

```sh
go run nats-jetstream/producer/sync/main.go -msgs 1000 -d 1 "profit.id"
```

- subscribe

```sh
go run nats-jetstream/consumer/push/main.go -sn profit -s "profit.id"
```

### Durable consumer

(save last state if disconnected)

#### push one subscriber only

- publish

```sh
go run nats-jetstream/producer/sync/main.go -msgs 1000 -d 1 "profit.id"
```

- subscribe

```sh
go run nats-jetstream/consumer/push/main.go -sn profit -cn one-push-cprofit_id -s "profit.id"
```

#### push load balanced queue consumer

- publish

```sh
go run nats-jetstream/producer/sync/main.go -msgs 1000 -d 1 "profit.id"
```

- subscribe

```sh
go run nats-jetstream/consumer/push/main.go -sn profit -cn queue-push-cprofit_id -qn poc -s "profit.id"
```

```sh
go run nats-jetstream/consumer/push/main.go -sn profit -cn queue-push-cprofit_id -qn poc -s "profit.id"
```

#### Shared pull consumer with batch processing

- publish

```sh
go run nats-jetstream/producer/sync/main.go -msgs 1000 -d 1 "profit.id"
```

- subscribe

```sh
go run nats-jetstream/consumer/pull/main.go -sn profit -cn pull-cprofit_id -s "profit.id"
```

```sh
go run nats-jetstream/consumer/pull/main.go -sn profit -cn pull-cprofit_id -s "profit.id"
```

#### Delayed messaging

- create stream `scheduled-message` that interest to persist event data from subject `delayed.profit`

```sh
go run nats-jetstream/manager/main.go --cmd add -n scheduled-message -s "delayed.profit"
```

- create durable push consumer to get data from stream `scheduled-message`

```sh
go run nats-jetstream/consumer/main.go --cmd add-push -sn scheduled-message -cn queue-push-delayed_profit -qn poc -s "delayed.profit"
```

- subscribe

```sh
go run nats-jetstream/consumer/push/main.go -sn scheduled-message -cn queue-push-delayed_profit -qn poc -s "delayed.profit"
```

- publish

```sh
go run nats-jetstream/producer/delayed/main.go -msgs 5 -d 10 "delayed.profit"
```

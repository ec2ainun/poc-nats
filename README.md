# NATS
stands for Neural Autonomic Transport System

the rationale why choose NATS over other messaging system comes from [this](https://docs.google.com/document/d/1N9-QybgQbp2Om2WLtFE-d43jJHxlW17XIRGnx-hCRtU/edit?usp=sharing)

## Getting Started
### Installation
- nats server
```
brew install nats-server
```
[another way to install nats-server](https://docs.nats.io/running-a-nats-service/introduction/installation)
- nats cli
```
brew tap nats-io/nats-tools
brew install nats-io/nats-tools/nats
```
[another way to install nats cli](https://docs.nats.io/using-nats/nats-tools/nats_cli)
- nats monitoring
```
go install github.com/nats-io/nats-top@latest
```
[another way to monitor](https://docs.nats.io/reference/faq#how-can-i-monitor-my-nats-cluster)

## Core NATS
- start nats server first
```
nats-server --config nats-core/nats.conf
```
### Subscribe to all event
- subscribe
```
go run nats-core/monitor/main.go ">"
```

### Publish Subscribe
- publish
```
go run nats-core/pub/main.go "profit"
```
- subscribe
```
go run nats-core/sub/main.go "profit"
```

### Fan In
- publish
```
go run nats-core/pub/main.go "profit.sg"
```
```
go run nats-core/pub/main.go "profit.id"
```
- subscribe
```
go run nats-core/sub/main.go "profit.*"
```

### Fan Out
- publish
```
go run nats-core/pub/main.go "profit.id"
```
- subscribe
```
go run nats-core/sub/main.go "profit.id"
```
```
go run nats-core/sub/main.go "profit.id"
```

### Queue Load Balanced
- publish
```
go run nats-core/pub/main.go "profit.id"
```
- subscribe using queue
```
go run nats-core/subq/main.go -q poc "profit.id" 
```
```
go run nats-core/subq/main.go -q poc "profit.id" 
```

### Request Reply
- start replier first
```
go run nats-core/rep/main.go -q poc "profit.get" 
```
- request
```
go run nats-core/req/main.go "profit.get"
```

## Jetstream
- start nats server first
```
nats-server --config nats-jetstream/js.conf
```

### Stream
what kind event data to persist
- check list of stream
```
nats str ls --server 0.0.0.0:5222
```
- check report of stream
```
nats stream report --server 0.0.0.0:5222
```
- create stream `profit` in which interest to persist event data from subject `profit.*`
```
go run nats-jetstream/manager/main.go --cmd add -n profit -s "profit.*"
```
- check stream detail
```
nats stream info profit --server 0.0.0.0:5222
```
- update stream `profit` in which interest to persist event data from subject `profit.sg profit.id`
```
go run nats-jetstream/manager/main.go --cmd update -n profit -s "profit.sg,profit.id"
```
- delete stream
```
go run nats-jetstream/manager/main.go --cmd delete -n profit 
```

### Publish
- sync publish
```
go run nats-jetstream/producer/sync/main.go -msgs 1000 "profit.id"
```
- async publish
```
go run nats-jetstream/producer/async/main.go -msgs 1000 "profit.id"
```

### Consumer
how application getting a messages from stream
- check list of consumer
```
nats consumer ls --server 0.0.0.0:5222
```
- push consumer
```
go run nats-jetstream/consumer/main.go --cmd add-push -sn profit -cn push-cprofit_id -s profit.id
```
- push consumer with shared load balanced queue 
```
go run nats-jetstream/consumer/main.go --cmd add-push -sn profit -cn queue-push-cprofit_id -qn poc -s profit.id
```

- pull consumer is automatically shared if there are multiple subscriber
```
go run nats-jetstream/consumer/main.go --cmd add-pull -sn profit -cn pull-cprofit_id -s profit.id
```

- check consumer detail
```
nats consumer info --server 0.0.0.0:5222
```
- delete consumer
```
go run nats-jetstream/consumer/main.go --cmd add-push -sn profit -cn test-delete -qn poc -s profit.id
```
```
go run nats-jetstream/consumer/main.go --cmd delete -sn profit -cn test-delete
```
### Ephermeral (direct) consumer 
- publish
```
go run nats-jetstream/producer/sync/main.go -msgs 1000 -d 1 "profit.id"
```
- subscribe
```
go run nats-jetstream/consumer/push/main.go -sn profit -s "profit.id"
```
### Durable consumer 
(save last state if disconnected)
#### push one subscriber only
- publish
```
go run nats-jetstream/producer/sync/main.go -msgs 1000 -d 1 "profit.id"
```
- subscribe
```
go run nats-jetstream/consumer/push/main.go -sn profit -cn push-cprofit_id -s "profit.id"
```
#### push load balanced queue consumer
- publish
```
go run nats-jetstream/producer/sync/main.go -msgs 1000 -d 1 "profit.id"
```
- subscribe
```
go run nats-jetstream/consumer/push/main.go -sn profit -cn queue-push-cprofit_id -qn poc -s "profit.id"
```
```
go run nats-jetstream/consumer/push/main.go -sn profit -cn queue-push-cprofit_id -qn poc -s "profit.id"
```

#### Shared pull consumer with batch processing
- publish
```
go run nats-jetstream/producer/sync/main.go -msgs 1000 -d 1 "profit.id"
```
- subscribe
```
go run nats-jetstream/consumer/pull/main.go -sn profit -cn pull-cprofit_id -s "profit.id"
```
```
go run nats-jetstream/consumer/pull/main.go -sn profit -cn pull-cprofit_id -s "profit.id"
```

#### Delayed messaging
- create stream `scheduled-message` that interest to persist event data from subject `delayed.profit`
```
go run nats-jetstream/manager/main.go --cmd add -n scheduled-message -s "delayed.profit"
```
- create durable push consumer to get data from stream `scheduled-message`
```
go run nats-jetstream/consumer/main.go --cmd add-push -sn scheduled-message -cn queue-push-delayed_profit -qn poc -s "delayed.profit"
```
- subscribe
```
go run nats-jetstream/consumer/push/main.go -sn scheduled-message -cn queue-push-delayed_profit -qn poc -s "delayed.profit"
```
- publish
```
go run nats-jetstream/producer/delayed/main.go -msgs 5 -d 10 "delayed.profit"
```

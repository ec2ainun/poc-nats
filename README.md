# NATS
stands for Neural Autonomic Transport System

the rationale why choose NATS over other messaging system comes from [this](https://docs.google.com/document/d/1N9-QybgQbp2Om2WLtFE-d43jJHxlW17XIRGnx-hCRtU/edit?usp=sharing)

## Getting Started
### installation
On Mac OS
```
brew install nats-server
```
[another way to install nats-server](https://docs.nats.io/running-a-nats-service/introduction/installation)

### setup
```
nats-server
```
then open new terminal and follow instruction below

## Core
### Monitor Event
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
- publish
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
- subscribe
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
go run nats-core/subq/main.go "profit.id" -q poc
```
- subscribe using queue
```
go run nats-core/subq/main.go "profit.id" -q poc
```

### Request Reply
- start replier first
```
go run nats-core/rep/main.go "profit.get" -q poc
```
- request
```
go run nats-core/req/main.go "profit.get"
```

## Jetstream

# Composing example

if you curios what kind of command is executed, you can refer to [this](/Makefile)

## Core NATS

### Publish and Subscribe

```sh
make run-pub-sub
```

- clean up

```sh
make stop-pub-sub
```

### Fan In

```sh
make run-fan-in
```

- clean up

```sh
make stop-fan-in
```

### Fan Out

```sh
make run-fan-out
```

- clean up

```sh
make stop-fan-out
```

### Queue Load Balanced

```sh
make run-qlb
```

- clean up

```sh
make stop-qlb
```

### Request Replay

```sh
make run-req-rep
```

- clean up

```sh
make stop-req-rep
```

## JetStream NATS

### Create data-stream Network

```sh
make data-network
```

### Stream Setup

```sh
make run-stream-manager
```

### Ephermeral

```sh
make run-ephermeral
```

- clean up

```sh
make stop-ephermeral
```

### Durable Push one sub

```sh
make run-durable-push-one
```

- clean up

```sh
make stop-durable-push-one
```

### Durable Push Load Balanced

```sh
make run-durable-push-lb
```

- clean up

```sh
make stop-durable-push-lb
```

### Durable Batch Processing

```sh
make run-durable-pull-batch
```

- clean up

```sh
make stop-durable-pull-batch
```

### Durable Delayed Messaging

```sh
make run-durable-push-delay
```

- clean up

```sh
make stop-durable-push-delay
```

### Clean up Jetstream

```sh
make stop-stream-manager
```

```sh
make rm-data-network
```

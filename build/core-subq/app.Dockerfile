FROM golang:1.18-alpine AS builder

ARG GIT_SHA
ARG GIT_TAG

WORKDIR /go/poc
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./
RUN CGO_ENABLED=0 go build -ldflags "-X main.gitSHA=${GIT_SHA} -X main.version=${GIT_TAG}" -o /poc nats-core/subq/main.go

FROM alpine

WORKDIR /
COPY --from=builder /poc poc
COPY --from=builder /go/poc/config/container.yaml config/config.yaml
# USER nonroot:nonroot

ENTRYPOINT ["/poc"]
CMD ["profit"]
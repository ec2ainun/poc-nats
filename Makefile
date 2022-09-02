
GIT_SHA=$(shell git rev-parse --short HEAD)
GIT_TAG=$(shell git describe --tags --abbrev=0)

build-pub:
	docker build -t ec2ainun/core-pub --build-arg GIT_SHA=$(GIT_SHA) --build-arg GIT_TAG=$(GIT_TAG) -f build/core-pub/app.Dockerfile .

build-sub:
	docker build -t ec2ainun/core-sub --build-arg GIT_SHA=$(GIT_SHA) --build-arg GIT_TAG=$(GIT_TAG) -f build/core-sub/app.Dockerfile .

build-subq:
	docker build -t ec2ainun/core-subq --build-arg GIT_SHA=$(GIT_SHA) --build-arg GIT_TAG=$(GIT_TAG) -f build/core-subq/app.Dockerfile .

build-req:
	docker build -t ec2ainun/core-req --build-arg GIT_SHA=$(GIT_SHA) --build-arg GIT_TAG=$(GIT_TAG) -f build/core-req/app.Dockerfile .

build-rep:
	docker build -t ec2ainun/core-rep --build-arg GIT_SHA=$(GIT_SHA) --build-arg GIT_TAG=$(GIT_TAG) -f build/core-rep/app.Dockerfile .

build-stream-manager:
	docker build -t ec2ainun/stream-manager --build-arg GIT_SHA=$(GIT_SHA) --build-arg GIT_TAG=$(GIT_TAG) -f build/js-stream/app.Dockerfile .

build-consumer-manager:
	docker build -t ec2ainun/consumer-manager --build-arg GIT_SHA=$(GIT_SHA) --build-arg GIT_TAG=$(GIT_TAG) -f build/js-stream/consumer.Dockerfile .

build-producer-sync:
	docker build -t ec2ainun/producer-sync --build-arg GIT_SHA=$(GIT_SHA) --build-arg GIT_TAG=$(GIT_TAG) -f build/producer-sync/app.Dockerfile .

build-producer-async:
	docker build -t ec2ainun/producer-async --build-arg GIT_SHA=$(GIT_SHA) --build-arg GIT_TAG=$(GIT_TAG) -f build/producer-async/app.Dockerfile .

build-producer-delay:
	docker build -t ec2ainun/producer-delay --build-arg GIT_SHA=$(GIT_SHA) --build-arg GIT_TAG=$(GIT_TAG) -f build/producer-delay/app.Dockerfile .

build-consumer-pull:
	docker build -t ec2ainun/consumer-pull --build-arg GIT_SHA=$(GIT_SHA) --build-arg GIT_TAG=$(GIT_TAG) -f build/consumer-pull/app.Dockerfile .

build-consumer-push:
	docker build -t ec2ainun/consumer-push --build-arg GIT_SHA=$(GIT_SHA) --build-arg GIT_TAG=$(GIT_TAG) -f build/consumer-push/app.Dockerfile .

docker-build:
	make build-pub
	make build-sub
	make build-subq
	make build-req
	make build-rep
	make build-stream-manager
	make build-consumer-manager
	make build-producer-sync
	make build-producer-async
	make build-producer-delay
	make build-consumer-pull
	make build-consumer-push

docker-push:
	docker push ec2ainun/core-pub
	docker push ec2ainun/core-sub
	docker push ec2ainun/core-subq
	docker push ec2ainun/core-req
	docker push ec2ainun/core-rep
	docker push ec2ainun/stream-manager
	docker push ec2ainun/consumer-manager
	docker push ec2ainun/producer-sync
	docker push ec2ainun/producer-async
	docker push ec2ainun/producer-delay
	docker push ec2ainun/consumer-pull
	docker push ec2ainun/consumer-push

run-pub-sub:
	docker compose -f deploy/compose/pub-sub/mnfst.yaml up -d

stop-pub-sub:
	docker compose -f deploy/compose/pub-sub/mnfst.yaml down

run-fan-in:
	docker compose -f deploy/compose/fan-in/mnfst.yaml up -d

stop-fan-in:
	docker compose -f deploy/compose/fan-in/mnfst.yaml down

run-fan-out:
	docker compose -f deploy/compose/fan-out/mnfst.yaml up -d

stop-fan-out:
	docker compose -f deploy/compose/fan-out/mnfst.yaml down

run-qlb:
	docker compose -f deploy/compose/qlb/mnfst.yaml up -d

stop-qlb:
	docker compose -f deploy/compose/qlb/mnfst.yaml down

run-req-rep:
	docker compose -f deploy/compose/req-rep/mnfst.yaml up -d

stop-req-rep:
	docker compose -f deploy/compose/req-rep/mnfst.yaml down

data-network:
	docker network create data-stream

rm-data-network:
	docker network rm data-stream

run-stream-manager:
	docker compose -f deploy/compose/stream/mnfst.yaml up -d

stop-stream-manager:
	docker compose -f deploy/compose/stream/mnfst.yaml down

run-ephermeral:
	docker compose -f deploy/compose/ephermeral/mnfst.yaml up -d

stop-ephermeral:
	docker compose -f deploy/compose/ephermeral/mnfst.yaml down

run-durable-push-one:
	docker compose -f deploy/compose/durable-push-one/mnfst.yaml up -d

stop-durable-push-one:
	docker compose -f deploy/compose/durable-push-one/mnfst.yaml down

run-durable-push-lb:
	docker compose -f deploy/compose/durable-push-lb/mnfst.yaml up -d

stop-durable-push-lb:
	docker compose -f deploy/compose/durable-push-lb/mnfst.yaml down

run-durable-pull-batch:
	docker compose -f deploy/compose/durable-pull-batch/mnfst.yaml up -d

stop-durable-pull-batch:
	docker compose -f deploy/compose/durable-pull-batch/mnfst.yaml down
	
run-durable-push-delay:
	docker compose -f deploy/compose/durable-push-delay/mnfst.yaml up -d

stop-durable-push-delay:
	docker compose -f deploy/compose/durable-push-delay/mnfst.yaml down
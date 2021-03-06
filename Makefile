.PHONY: test
test:
	go test -race -coverprofile=cover.out -covermode=atomic ./...

.PHONY: build
build:
	- docker build --rm --file=./build/docker/starter.dockerfile --tag jasonsoft/starter:latest .

run:
	- docker-compose up --scale worker=1

build-run: build run


lint:
	docker run --rm -v ${LOCAL_WORKSPACE_FOLDER}:/app -w /app golangci/golangci-lint:v1.41-alpine golangci-lint run ./... -v

infra:
	- docker-compose -f docker-compose-infra.yml up -d

infra-down:
	- docker-compose -f docker-compose-infra.yml down

proto:
	- protoc --proto_path=. --go_out=plugins=grpc:. ./pkg/event/proto/*.proto
	- protoc --proto_path=. --go_out=plugins=grpc:. ./pkg/wallet/proto/*.proto
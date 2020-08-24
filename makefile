build:
	COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose build
run:
	COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose up --build
infra:
	COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose -f docker-compose-infra.yml up -d

proto:
	protoc --proto_path=. --go_out=plugins=grpc:. ./pkg/event/proto/*.proto
	protoc --proto_path=. --go_out=plugins=grpc:. ./pkg/wallet/proto/*.proto
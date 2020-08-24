run:
	docker-compose up --build

infra-run:
	docker-compose -f docker-compose-infra.yml up -d

proto:
	protoc --proto_path=. --go_out=plugins=grpc:. ./pkg/event/proto/*.proto
	protoc --proto_path=. --go_out=plugins=grpc:. ./pkg/wallet/proto/*.proto
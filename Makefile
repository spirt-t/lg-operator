export GO111MODULE=on
export GOSUMDB=off

.PHONY: install
install:
	go mod tidy && go mod download

.PHONY: test
test:
	go test -v -cover -race ./...

.PHONY: cover
cover:
	go test -v -coverprofile=coverage.out ./...  && go tool cover -html=coverage.out

.PHONY: run
run:
	go run cmd/lg-operator/main.go

.PHONY: go-generate
go-generate:
	go generate ./...

.PHONY: generate
generate:
	rm -f pkg/lg-operator/*.go
	protoc --proto_path=api --go_out=pkg --go_opt=paths=source_relative \
	--go-grpc_out=pkg --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pkg --grpc-gateway_opt=paths=source_relative \
	--swagger_out=pkg \
	api/lg-operator/*.proto



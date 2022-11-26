# Makefile for local and docker development
ifeq ($(MAKEFILE_ENV), local)
CMD = ""
VENDOR_PATH = "$(GOPATH)/src"
else
CMD = "docker exec -it tamago"
VENDOR_PATH = "/go/src"
endif

PROTO_SRC="./protobuf"
PROTO_DST="./internal/api"

GOLANGCI_FILE = "./build/golangci.yml"
DOCKER_COMPOSE_FILE = "build/docker-compose.local.yml"


gen-proto:
	@eval $(CMD) echo "building protobuf"
	@eval $(CMD) mkdir -p $(PROTO_DST)
	@eval $(CMD) protoc -I=$(VENDOR_PATH) -I=$(PROTO_SRC) --gogo_out=$(PROTO_DST) $(PROTO_SRC)/api.proto
	@eval $(CMD) echo "built protobuf"

gen-mocks:
	@eval $(CMD) echo "building mock files"
	@eval $(CMD) mkdir -p ./internal/mocks
	@eval $(CMD) mockery -name=DataAccessIface -dir=./internal/dao -output=./internal/mocks
	@eval $(CMD) mockery -name=SessionIface -dir=./internal/sessions -output=./internal/mocks
	@eval $(CMD) echo "built mock files"

gen: gen-proto gen-mocks vendor

lint: gen 
	@eval $(CMD) golangci-lint run -c $(GOLANGCI_FILE) -v ./...
	@eval $(CMD) buf check lint --input-config build/buf.yaml --log-level debug

docker-build:
	@docker network create bento || true 
	@docker-compose -f $(DOCKER_COMPOSE_FILE) build

up:
	@docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

down:
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down

logs:
	@docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

build: gen
	@eval $(CMD) go build -o ./tmp/main ./cmd/tamago/main.go

run: build 
	@eval $(CMD) ./tmp/main

run-ci:
	@circleci local execute --job build

db-init:
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down
	@docker-compose -f $(DOCKER_COMPOSE_FILE) run -d -p 5432:5432 --name postgres bento_postgres || true
	@docker exec -it postgres createdb -U postgres bento || true
	@go run ./cmd/migrate/main.go
	@docker stop postgres

db-drop:
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down
	@docker-compose -f $(DOCKER_COMPOSE_FILE) run -d -p 5432:5432 --name postgres bento_postgres || true
	@docker exec -it postgres dropdb -U postgres bento
	@docker stop postgres

test: gen
	@eval $(CMD) go test ./... -v

testsum: gen
	@eval $(CMD) gotestsum -f short-with-failures

coverage: gen
	@eval $(CMD) go test ./... -covermode=count -coverpkg=./... -coverprofile=/tmp/tamago_cover.out fmt
	@eval $(CMD) go tool cover -html=/tmp/tamago_cover.out

download:
	@eval $(CMD) go mod download

vendor:
	@eval $(CMD) go mod vendor

install-base:
	GO111MODULE=off eval $(CMD) go get -u -d github.com/gogo/protobuf/gogoproto
	@eval $(CMD) go get github.com/gogo/protobuf/protoc-gen-gogo
	@eval $(CMD) go get github.com/vektra/mockery/.../
	@eval $(CMD) go get gotest.tools/gotestsum
	@eval $(CMD) go get -u github.com/cosmtrek/air

install-mac: install-base
	brew tap bufbuild/buf
	brew install buf
	brew install golangci/tap/golangci-lint

# TODO: complete this
install: install-base

SRC := $(shell find . -name "*.go")

all: bin/chat

bin/chat: cmd/main.go $(SRC) pb
	go build -o bin/chat cmd/main.go

run: bin/chat
	./bin/chat

test:
	go test -v ./tests -run Test_Generate
	go test -v ./tests -run Test_CreateTable

test_with_env:
	docker-compose -f docker-compose.test.yaml build
	docker-compose -f docker-compose.test.yaml up -d fdb
	docker-compose -f docker-compose.test.yaml run fdb-service
	docker-compose -f docker-compose.test.yaml down

run_with_env:
	docker-compose up
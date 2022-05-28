test:
	go test -v ./tests -run Test_Generate
	go test -v ./tests -run Test_CreateTable

test_with_env:
	docker-compose -f docker-compose.test.yaml build
	docker-compose -f docker-compose.test.yaml up -d fdb
	docker-compose -f docker-compose.test.yaml run fdb-service
	docker-compose -f docker-compose.test.yaml down

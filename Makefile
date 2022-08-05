test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

test-iterator-coverage:
	cd structures/iterator && go test . -coverprofile=coverage.out
	cd structures/iterator && go tool cover -html=coverage.out -o coverage.html
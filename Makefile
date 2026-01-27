.PHONY: test
# run go tests with coverage
test:
	go test ./... -coverprofile=coverage.out -covermode=atomic

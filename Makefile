.PHONY: test 
test: 
	go test -v --cover --race --count=1 ./...

.PHONY: lint 
lint: 
	golangci-lint --enable gosec,misspell run ./...

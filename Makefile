run:
	go run main.go

run-race:
	go run -race main.go

tidy:
	go get -u && go mod tidy && go get -u ./...

test:
	CGO_ENABLED=0 go test -v ./...

lint:
	golangci-lint run ./...



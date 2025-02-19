run:
	go run main.go

tidy:
	go mod tidy

test:
	CGO_ENABLED=0 go test -v ./...

build:
  go build

test:
  go test ./...

lint:
  go fmt ./... && golangci-lint run

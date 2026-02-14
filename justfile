default: build

watch:
  go build -o ssh-list *.go || exit 1 && ./ssh-list

build:
  go build -o ssh-list *.go

test:
  go test ./...

lint:
  go fmt ./... && golangci-lint run

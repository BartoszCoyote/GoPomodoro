build: fmt lint test
	mkdir -p bin && GO111MODULE=on go build -o ./bin/gopom ./cmd/gopom

run: build
	./bin/gopom "$(command)"

test:
	GO111MODULE=on go test ./...

fmt:
	GO111MODULE=on go fmt ./...

lint:
	GO111MODULE=on golint ./...
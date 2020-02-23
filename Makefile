ensure-deps:
	./ensure-deps.sh

build: fmt test
	mkdir -p bin && GO111MODULE=on go build -o ./bin/gopom ./cmd/gopom

run: build
	./bin/gopom "$(command)"

test: ensure-deps
	GO111MODULE=on go test ./...

fmt:
	GO111MODULE=on go fmt ./...

lint:
	GO111MODULE=on golint ./...
build:
	mkdir -p bin
	GO111MODULE=on go build -o ./bin/gopom ./cmd/gopom

run: build
	./bin/gopom "$(command)"

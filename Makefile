build: fmt test
	./build.sh

run: build
	./bin/gopom "$(command)"

test:
	GO111MODULE=on go test ./...

fmt:
	GO111MODULE=on go fmt ./...

lint:
	GO111MODULE=on golint ./...
.PHONY: build test clean lint

GO_SRC := .engine/scripts/src
GO_OUT := scripts/spacecraft

build:
	cd $(GO_SRC) && go build -o ../../../$(GO_OUT) .

test:
	cd $(GO_SRC) && go test ./...

clean:
	rm -f $(GO_OUT)

lint:
	cd $(GO_SRC) && go vet ./...

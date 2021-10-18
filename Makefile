TARGETS := dpl

.PHONY: all clean format test ${TARGETS}

all: ${TARGETS}

clean:
	go clean

format:
	find -name '*.go' | xargs gofmt -w

dpl:
	go build -o ${@}

test:
	go test -race ./...

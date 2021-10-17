TARGETS := dpl

.PHONY: all clean format ${TARGETS}

all: ${TARGETS}

clean:
	go clean

format:
	find -name '*.go' | xargs gofmt -w

dpl:
	go build -o ${@}

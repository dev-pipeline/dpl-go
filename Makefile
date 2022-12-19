TARGETS := dpl

.PHONY: \
	all \
	clean \
	coverage \
	format \
	install \
	lint \
	test \
	${TARGETS}

all: ${TARGETS}

clean:
	go clean

format:
	find -name '*.go' | xargs gofmt -w

dpl:
	go build -o ${@}

test:
	go test -race ./...

lint:
	go vet ./...
	staticcheck ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

install: dpl
	if [ -z "${DESTDIR}" ]; then \
		cp dpl "\${DESTDIR}/bin/"; \
	else \
		cp dpl "$(shell go env GOPATH)/bin"; \
	fi

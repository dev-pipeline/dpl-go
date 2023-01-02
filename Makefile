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
	find -name '*.go' | xargs gofmt -w -s

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
	if [ -n "${DESTDIR}" ]; then \
		mkdir -p "${DESTDIR}/bin/"; \
		cp dpl "${DESTDIR}/bin/"; \
	else \
		mkdir -p "$(shell go env GOPATH)/bin"; \
		cp dpl "$(shell go env GOPATH)/bin"; \
	fi

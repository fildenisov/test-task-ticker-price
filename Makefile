.PHONY: mod
mod:
	go mod tidy

.PHONY: build
build:
	go build \
			--trimpath \
			-o bin/app \
			cmd/http/main.go

.PHONY: check
check:
	golangci-lint run -v --config .golangci.yml

.PHONY: test
test:
	go test -v ./...

.PHONY: race
race:
	go test -v -race ./...

.PHONY: cover 
cover:
	./scripts/coverage.sh html

.PHONY: run
run: build
	./bin/app --config=configs/local/config.yaml


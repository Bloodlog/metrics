GOLANGCI_LINT_CACHE?=/tmp/praktikum-golangci-lint-cache

.PHONY: golangci-lint-run
golangci-lint-run: _golangci-lint-rm-unformatted-report

.PHONY: _golangci-lint-reports-mkdir
_golangci-lint-reports-mkdir:
	mkdir -p ./golangci-lint

.PHONY: _golangci-lint-run
_golangci-lint-run: _golangci-lint-reports-mkdir
	-docker run --rm \
    -v $(shell pwd):/app \
    -v $(GOLANGCI_LINT_CACHE):/root/.cache \
    -w /app \
    golangci/golangci-lint:v1.57.2 \
        golangci-lint run \
            -c .golangci.yml \
	> ./golangci-lint/report-unformatted.json

.PHONY: _golangci-lint-format-report
_golangci-lint-format-report: _golangci-lint-run
	cat ./golangci-lint/report-unformatted.json | jq > ./golangci-lint/report.json

.PHONY: _golangci-lint-rm-unformatted-report
_golangci-lint-rm-unformatted-report: _golangci-lint-format-report
	rm ./golangci-lint/report-unformatted.json

.PHONY: golangci-lint-clean
golangci-lint-clean:
	sudo rm -rf ./golangci-lint

style:
	go mod tidy
	go fmt ./...
	go vet ./...
	goimports -w .

docker-up:
	docker-compose up --build

go-test:
	go test ./...

go-test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

go-doc:
	godoc -http=:8081

swag:
	swag init -g cmd/server/main.go --output ./swagger/

multichecker:
	go build -o cmd/staticlint/multichecker cmd/staticlint/main.go
	cmd/staticlint/multichecker internal

build:
	go build -ldflags "-X main.Version=v1.0.1 -X 'main.BuildTime=$(date +'%Y/%m/%d %H:%M:%S')'" -o cmd/agent/agent cmd/agent/main.go
	go build -ldflags "-X main.Version=v1.0.1 -X 'main.BuildTime=$(date +'%Y/%m/%d %H:%M:%S')'" -o cmd/server/server cmd/server/main.go

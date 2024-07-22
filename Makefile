-include .env
export $(shell sed 's/=.*//' .env)

VERSION=$(shell git describe --tags --abbrev=0)

.PHONY: install
# install development dependencies
install:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.56.2
	go install github.com/google/wire/cmd/wire@v0.6.0
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	go install github.com/vektra/mockery/v2@v2.43.0
	go install github.com/cosmtrek/air@v1.45.0

.PHONY: wire
# generate wire
wire:
	@cd cmd/service && wire gen ./...

.PHONY: api
# generate server code from openapi spec
api:
	@mkdir -p gen/api/profiles
	@oapi-codegen --config=.oapi-codegen.yaml -include-tags Profiles -package profiles -o gen/api/profiles/api.go internal/openapi/openapi.yaml

.PHONY: mocks
# generate mocks
mocks:
	@mockery

.PHONY: test
# run tests
test:
	@go test -timeout 30s -cover -v ./internal/...

.PHONY: clean
# clean
clean:
	@rm -rf cmd/service/wire_gen.go
	@rm -rf bin/
	@rm -rf gen/
	@rm -rf vendor/

.PHONY: lint
# run linter
lint:
	@golangci-lint run -v

.PHONY: lint-fix
# run linter and fix
lint-fix:
	@golangci-lint run -v --fix

.PHONY: build
# build
build:
	@mkdir -p bin/ && CGO_ENABLED=0 GOOS=linux go build -o ./bin/service ./cmd/service

.PHONY: run
# run server with air
run:
	@air \
		--build.include_dir gen/api,cmd/service,internal,pkg \
		--build.args_bin "-conf ./configs/config.yaml" \
		--build.cmd "go build -buildvcs=false -o ./bin/service ./cmd/service" \
		--build.bin ./bin/service \
		--build.send_interrupt true \
		--build.poll true \
		.

.PHONY: compose
# run docker compose
compose:
	@docker compose -f docker/compose.yaml up --build --remove-orphans

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
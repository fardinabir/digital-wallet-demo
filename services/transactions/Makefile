COMMON_PACKAGE=github.com/fardinabir/digital-wallet-demo/services/transactions/internal/common
CURRENT_DIR=$(shell pwd)
DIST_DIR=${CURRENT_DIR}/dist
CLI_NAME=wallet-cli

HOST_OS:=$(shell go env GOOS)
HOST_ARCH:=$(shell go env GOARCH)

VERSION=$(shell cat ${CURRENT_DIR}/VERSION)
BUILD_DATE:=$(if $(BUILD_DATE),$(BUILD_DATE),$(shell date -u +'%Y-%m-%dT%H:%M:%SZ'))
GIT_COMMIT:=$(if $(GIT_COMMIT),$(GIT_COMMIT),$(shell git rev-parse HEAD))

ifeq (${COVERAGE_ENABLED}, true)
# We use this in the cli-local target to enable code coverage for e2e tests.
COVERAGE_FLAG=-cover
else
COVERAGE_FLAG=
endif

override LDFLAGS += \
  -X ${COMMON_PACKAGE}.version=${VERSION} \
  -X ${COMMON_PACKAGE}.buildDate=${BUILD_DATE} \
  -X ${COMMON_PACKAGE}.gitCommit=${GIT_COMMIT} \

.PHONY: migrate
migrate:
	go run main.go migrate --config config.yaml

.PHONY: migrate-test
migrate-test:
	go run main.go migrate --config config.test.yaml

.PHONY: reset-db
reset-db:
	PGPASSWORD=postgres psql -h localhost -U postgres -d postgres -c "DROP DATABASE IF EXISTS wallet;"
	PGPASSWORD=postgres psql -h localhost -U postgres -d postgres -c "CREATE DATABASE wallet WITH TEMPLATE = template0 OWNER = postgres ENCODING = 'UTF8';"
	make migrate

.PHONY: reset-test-db
reset-test-db:
	PGPASSWORD=postgres psql -h localhost -U postgres -d postgres -c "DROP DATABASE IF EXISTS transaction_test;"
	PGPASSWORD=postgres psql -h localhost -U postgres -d postgres -c "CREATE DATABASE transaction_test WITH TEMPLATE = template0 OWNER = postgres ENCODING = 'UTF8';"
	make migrate-test

.PHONY: start
start:
	docker-compose up -d

.PHONY: stop
stop:
	docker-compose down

.PHONY: clear
clear:
	docker-compose down -v

.PHONY: serve-backend
serve:
	go run main.go server --config config.yaml

.PHONY: lint
lint:
	$(shell go env GOPATH)/bin/golangci-lint run

.PHONY: fmt
fmt:
	go mod tidy
	$(shell go env GOPATH)/bin/golangci-lint run --fix
	swag fmt

.PHONY: test
test: reset-test-db
	gotestsum --format=testname --rerun-fails

.PHONY: test-ci
test-ci: reset-test-db
	gotestsum --format=testname -- -cover -coverprofile=coverage.out ./...

# Go files to check during build
SWAG_GO_FILES:=$(shell find internal/controller -type f -name '*.go' -print)

docs/swagger.yaml: main.go $(SWAG_GO_FILES)
	swag init

docs/swagger.html: docs/swagger.yaml
	npx @redocly/cli@1.25.3 build-docs -o docs/swagger.html docs/swagger.yaml

.PHONY: swagger
swagger: docs/swagger.html
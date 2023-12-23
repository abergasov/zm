PROJECT_NAME:=zm
FILE_HASH := $(shell git rev-parse HEAD)
GOLANGCI_LINT := $(shell command -v golangci-lint 2> /dev/null)

# test coverage threshold
COVERAGE_THRESHOLD:=70
COVERAGE_TOTAL := $(shell go tool cover -func=cover.out | grep total | grep -Eo '[0-9]+\.[0-9]+')
COVERAGE_PASS_THRESHOLD := $(shell echo "$(COVERAGE_TOTAL) $(COVERAGE_THRESHOLD)" | awk '{print ($$1 >= $$2)}')

init_repo: ## create necessary configs
	cp configs/sample.common.env configs/common.env
	cp configs/sample.app_conf.yml configs/app_conf.yml
	cp configs/sample.app_conf_docker.yml configs/app_conf_docker.yml

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

install-lint: ## Installs golangci-lint tool which a go linter
ifndef GOLANGCI_LINT
	${info golangci-lint not found, installing golangci-lint@latest}
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
endif

gogen: ## generate code
	${info generate code...}
	go generate ./internal...

test: ## Runs tests
	${info Running tests...}
	go test -v -race ./... -cover -coverprofile cover.out
	go tool cover -func cover.out | grep total

bench: ## Runs benchmarks
	${info Running benchmarks...}
	go test -bench=. -benchmem ./... -run=^#

vulcheck: ## Runs vulnerability check
	${info Running vulnerability check...}
	govulncheck ./...

lint: install-lint ## Runs linters
	@echo "-- linter running"
	golangci-lint run -c .golangci.yaml ./internal...
	golangci-lint run -c .golangci.yaml ./cmd...

lint_d:
	docker run --rm -v ${PWD}:/app -w /app golangci/golangci-lint:v1.50 golangci-lint run --timeout 5m ./...

stop: init_repo ## Stops the local environment
	${info Stopping containers...}
	docker container ls -q --filter name=${PROJECT_NAME} ; true
	${info Dropping containers...}
	docker rm -f -v $(shell docker container ls -q --filter name=${PROJECT_NAME}) ; true

dev_up: stop ## Runs local environment
	${info Running docker-compose up...}
	GIT_HASH=${FILE_HASH} docker compose -p ${PROJECT_NAME} up --build dbPostgres

build: ## Builds binary
	@echo "-- building binary"
	go build -o ./bin/binary ./cmd

build_in_docker: ## Builds binary in docker
	@echo "-- building docker image"
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/binary ./cmd

run: init_repo ## Runs binary local with environment in docker
	${info Run app containered}
	GIT_HASH=${FILE_HASH} docker compose -p ${PROJECT_NAME} up --build

migrate_new: ## Create new migration
	migrate create -ext sql -dir migrations -seq data

coverage: ## Check test coverage is enough
	@echo "Threshold:                ${COVERAGE_THRESHOLD}%"
	@echo "Current test coverage is: ${COVERAGE_TOTAL}%"
	@if [ "${COVERAGE_PASS_THRESHOLD}" -eq "0" ] ; then \
		echo "Test coverage is lower than threshold"; \
		exit 1; \
	fi

generate_files: clean ## Generate files for tests
	mkdir -p data_folder
	for i in $$(seq 1 10); do uuid=$$(uuidgen); echo "$$uuid" > data_folder/file_$$i.txt; done

clean: ## Clean generated files
	rm -rf data_folder

upload: # generate_files ## Uploads generated files to server
	go run cmd/uploader/main.go

loop:
	@for i in {1..100}; do \
		echo "Iteration $$i"; \
		make test || { echo "Exiting loop due to test failure"; exit 1; }; \
	done

.PHONY: help install-lint test gogen lint stop dev_up build run init_repo migrate_new vulcheck coverage build_in_docker generate_files clean upload
.DEFAULT_GOAL := help
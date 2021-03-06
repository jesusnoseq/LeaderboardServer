UNAME := $(shell uname)
CODE_DIR = functions
LINTER_ARGS = run -c .golangci.yml --timeout 5m
CGO_CFLAGS = ""
CMD_FILE=$(CURDIR)/$(CODE_DIR)/api/cmd/*.go
AWS_REGION=eu-central-1
ENVIRONMENT=qa


.PHONY: help
help:	## Show a list of available commands
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

.PHONY: make-debug
make-debug:	## Debug Makefile itself
	@echo $(UNAME)

.PHONY: build
build:
	set GOOS=linux
	set GOARCH=amd64
	set CGO_ENABLED=0
	sam build
	echo "rebuilding to shrink binary size"
	cd $(CODE_DIR) go build -o ../.aws-sam/build/EntryApiFunction/EntryApi -ldflags="-s -w" lambda/cmd/main.go

.PHONY: fmt
fmt:	## Format code
	cd $(CODE_DIR) && gofmt -w -s ./$(CODE_DIR)
	cd $(CODE_DIR) && goimports -w -l ./$(CODE_DIR)

.PHONY: tidy
tidy:	## Prune any no-longer-needed dependencies from go.mod and add any dependencies needed
	cd $(CODE_DIR) && go mod tidy -v

.PHONY: test
test:	## Run unitary test
	cd $(CODE_DIR) go test -p 1 -cover -v ./... -timeout 5m

lint:	## Run static linting of source files. See .golangci.yml for options
	cd $(CODE_DIR) && golangci-lint $(LINTER_ARGS)

.PHONY: download-tools
download-tools:  ## Download all required tools to generate documentation, code analysis...
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.44.2
	go install golang.org/x/tools/cmd/goimports@v0.1.9
	go install github.com/go-swagger/go-swagger/cmd/swagger@v0.29.0

.PHONY: run-debug
run-debug:	## Debug application with CLI
	dlv debug $(CMD_FILE) --headless=false 

.PHONY: run-remote-debug
run-remote-debug:	## Debug remote application
	dlv debug $(CMD_FILE) --headless=true --listen=:$(DEBUG_PORT) --api-version=2 --log

.PHONY: run
run-api:	## Run API
	cd $(CODE_DIR) && CGO_CFLAGS=${CGO_CFLAGS} go run $(CMD_FILE)

.PHONY: validate-template
validate-template:
	sam validate

.PHONY: deploy
deploy: ## Deploy in default env
	sam deploy --config-env $(ENVIRONMENT)

remove:
	aws cloudformation delete-stack --stack-name $(ENVIRONMENT)-LeaderboardServer --region $(AWS_REGION)
USER_API_DIR = pkg/api/user
USER_API_PROTO_SOURCE = $(wildcard ${USER_API_DIR}/*.proto)
USER_API_PROTO_OUT = ${USER_API_DIR}/api.swagger.json $(wildcard pkg/api/users/*.go)

INVOICES_API_DIR = pkg/api/invoice
INVOICES_API_PROTO_SOURCE = $(wildcard ${INVOICES_API_DIR}/*.proto)
INVOICES_API_PROTO_OUT = ${INVOICES_API_DIR}/api.swagger.json $(wildcard pkg/api/invoices/*.go)

.PHONY: all
all: help

$(INVOICES_API_PROTO_OUT): $(INVOICES_API_PROTO_SOURCE)
	protoc -I ./proto/googleapis -I $(INVOICES_API_DIR) \
		--go_out=$(INVOICES_API_DIR)  \
        --go_opt=paths=source_relative \
        --go-grpc_out=$(INVOICES_API_DIR) \
        --go-grpc_opt=paths=source_relative \
        --grpc-gateway_out=$(INVOICES_API_DIR) \
        --grpc-gateway_opt paths=source_relative \
        --openapiv2_out $(INVOICES_API_DIR) \
        --openapiv2_opt logtostderr=true \
        $(INVOICES_API_PROTO_SOURCE)

$(USER_API_PROTO_OUT): $(USER_API_PROTO_SOURCE)
	protoc -I ./proto/googleapis -I $(USER_API_DIR) \
		--go_out=$(USER_API_DIR)  \
        --go_opt=paths=source_relative \
        --go-grpc_out=$(USER_API_DIR) \
        --go-grpc_opt=paths=source_relative \
        --grpc-gateway_out=$(USER_API_DIR) \
        --grpc-gateway_opt paths=source_relative \
        --openapiv2_out $(USER_API_DIR) \
        --openapiv2_opt logtostderr=true \
        $(USER_API_PROTO_SOURCE)

## codegen:
.PHONY: codegen
codegen: $(USER_API_PROTO_OUT) $(INVOICES_API_PROTO_OUT) ## Generate GRPC stubs

## build:
.PHONY: build
build: codegen ## Build application and generate grpc stubs
	go build -o invoice ./cmd/invoice

## Help:
.PHONY: help
help: ## Show this help
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)

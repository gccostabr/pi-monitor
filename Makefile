BUILDER_IMAGE := pi-monitor-builder

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build-docker: ## Build the Docker builder image, that will be used to compile the binary 
	@docker build -t ${BUILDER_IMAGE} .
	# Builder image created.

.PHONY: build 
build: build-docker ## Build the application binary
	@docker run -it -v $(shell pwd):/app ${BUILDER_IMAGE} go build -o pi-monitor .
	# Binary created.
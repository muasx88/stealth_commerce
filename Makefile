export GOBIN := $(PWD)/bin
# Exporting bin folder to the path for makefile
export PATH   := $(PWD)/bin:$(PATH)
# Default Shell
export SHELL  := bash
# Type of OS: Linux or Darwin.
export OSTYPE := $(shell uname -s)

BINARY=stealth-commerce

include ./misc/tools.Makefile


install-deps: air  ## Install Development Dependencies (localy).
deps: $(AIR) ## Checks for Global Development Dependencies.
deps:
	@echo "Required Tools Are Available"


docker-up: ## Bootstrap Environment (with a Docker-Compose help).
	@ docker-compose up --build -d

dev: $(AIR) ## Starts AIR ( Continuous Development app).
	air

build:
	@ printf "Building aplication... "

	@ go build \
 		-trimpath \
 		-o target/${BINARY} \
 		./app/

	@ echo "done"

build-race: ## Builds binary (with -race flag)
	@ printf "Building aplication with race flag... "

	@ go build \
		-trimpath  \
		-race      \
		-o target/${BINARY} \
		./app/

	@ echo "done"
VERSION ?= latest

# Create a variable that contains the current date in UTC
# Different flow if this script is running on Darwin or Linux machines.
ifeq (Darwin,$(shell uname))
	NOW_DATE = $(shell date -u +%d-%m-%Y)
else
	NOW_DATE = $(shell date -u -I)
endif

all: test

.PHONY: test
test:
	@go test ./... -race -cover

.PHONY: coverage
coverage:
	@go test ./... -race -cover -coverprofile cover.out

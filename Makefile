GOCMD := GOPRIVATE="" go

UNAME := $(strip $(shell uname -m))
COMMA := $(shell echo ",")
SPACE := $(shell echo " ")
BUILD_TAGS ?=
TAG_LIST = $(subst $(COMMA),$(SPACE),$(BUILD_TAGS))

ifeq ($(UNAME),$(filter $(UNAME),arm64 aarch64))
  TAG_LIST += dynamic
endif
TAGS := $(subst $(SPACE),$(COMMA),$(strip $(TAG_LIST)))
ifneq ($(strip $(TAGS)),)
  TAGS := --tags $(TAGS)
endif

generate:
	$(GOCMD) generate $(TAGS) ./...
.PHONY: generate

test:
	$(GOCMD) test $(TAGS) -buildmode=pie -v -coverprofile=coverage.out  ./...
.PHONY: test

test/%:
	$(GOCMD) test $(TAGS) -buildmode=pie -v -count=1 ./$*/...

coverage: test
	$(GOCMD) tool cover -func=coverage.out
.PHONY: coverage

htmlcoverage: test
	$(GOCMD) tool cover -html=coverage.out
.PHONY: htmlcoverage

lint:
	$(GOCMD) vet $(TAGS) ./...
.PHONY: lint

download:
	$(GOCMD) mod download
.PHONY: download

changelog:
	git-chglog --next-tag=$(shell cat VERSION) -o CHANGELOG.md
.PHONY: changelog

mocks:
 # Generating mocks for each directory separately so we don't generate mocks for the
 # integration tests nor the workflows directory
	  rm -rf mocks \
  		&& mockery --dir internal --all --keeptree --with-expecter --output internal/mocks
.PHONY: mocks

rsa:
 # Generate rsa keys for testing
	rm -f "./id_rsa"
	ssh-keygen -t rsa -b 4096 -C "`hostname`" -f "./id_rsa" -P "" -q
	ssh-keygen -f "./id_rsa.pub" -e -m pem > "./id_rsa.pem"
	ssh-keygen -f "./id_rsa" -p -N "" -e -m pem > /dev/null;
	rm -f "./id_rsa.pub"
.PHONY: rsa

monolith: download
	$(GOCMD) build $(TAGS) -o dist/monolith ./cmd/monolith/main.go
.PHONY: monolith

two_microservices_session: download
	$(GOCMD) build $(TAGS) -o dist/session ./cmd/two_microservices/session/main.go
.PHONY: two_microservices_session

two_microservices_user_auth: download
	$(GOCMD) build $(TAGS) -o dist/user_auth ./cmd/two_microservices/user_auth/main.go
.PHONY: two_microservices_user_auth

three_microservices_user: download
	$(GOCMD) build $(TAGS) -o dist/user ./cmd/three_microservices/user/main.go
.PHONY: three_microservices_user

three_microservices_session: download
	$(GOCMD) build $(TAGS) -o dist/session ./cmd/three_microservices/session/main.go
.PHONY: three_microservices_session

three_microservices_auth: download
	$(GOCMD) build $(TAGS) -o dist/auth ./cmd/three_microservices/auth/main.go
.PHONY: three_microservices_auth

three_microservices_crypto_user: download
	$(GOCMD) build $(TAGS) -o dist/user ./cmd/three_microservices_crypto/user/main.go
.PHONY: three_microservices_crypto_user

three_microservices_crypto_session: download
	$(GOCMD) build $(TAGS) -o dist/session ./cmd/three_microservices_crypto/session/main.go
.PHONY: three_microservices_crypto_session

three_microservices_crypto_auth: download
	$(GOCMD) build $(TAGS) -o dist/auth ./cmd/three_microservices_crypto/auth/main.go
.PHONY: three_microservices_crypto_auth

binaries: monolith
.PHONY: binaries
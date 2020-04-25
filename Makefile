SHELL = /bin/bash -o pipefail

ENV_TYPE ?= desktop

# Code style formatting
fmt:
	go fmt ./...

# Quality assurance
qa: lint test-cover

# Getting deps with required to build the tests for the specified packages
test-deps:
	@echo "Getting deps with required to build the tests for the specified packages"
	go get -t -d ./...

# Testing with check cover
test-cover: test-deps
	@echo "Testing with check cover"
	bash script/go-test.sh -v -covermode count -coverprofile=profile-$(ENV_TYPE).cov ./...

# Testing with check race
test-race: test-deps
	go test -v -race ./...

# Store current bechmarck
bench-store:
	bash script/go-test.sh -benchmem -benchtime=2s -bench=. -run='^$$' ./... | tee testdata/current.bench.txt

# Take current bechmarck
bench-now:
	bash script/go-test.sh -benchmem -benchtime=2s -bench=. -run='^$$' ./... | tee testdata/new.bench.txt

# Make compares benchmarks
bench: bench-now
	@echo "Make compares benchmarks"
	$(shell go env GOPATH)/bin/benchstat testdata/current.bench.txt testdata/new.bench.txt

# Linting and code analyzes
lint: lint-go

say:
	@echo "$(ENV_TYPE)-$(shell go env GOOS)"

# Linting only go files
lint-go:
ifeq ("$(ENV_TYPE)", "destop")
	@echo "Checking for analyzes to identify unnecessary type conversions"
	$(shell go env GOPATH)/bin/unconvert -v ./...
endif

# Linting only shell scripts
lint-shell:
	@echo "Checking shell scripts for analyzes to identify unnecessary type conversions"
	shellcheck script/*

# Send cover report to coveralls.io
coveralls: profile-$(ENV_TYPE).cov
	@echo "Sending cover report to coveralls.io"
	$(shell go env GOPATH)/bin/goveralls -coverprofile=profile-$(ENV_TYPE).cov -service=github

# Verify demo, tutorial, etc
verify: verify-demo verify-tutorial

# Verify demos
verify-demo:
	bash script/verify-demo.sh

# Verify tutorials
verify-tutorial:
	bash script/verify-tutorial.sh

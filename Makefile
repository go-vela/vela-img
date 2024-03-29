# Copyright (c) 2022 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

# The `clean` target is intended to clean the workspace
# and prepare the local changes for submission.
#
# Usage: `make clean`
.PHONY: clean
clean: tidy vet fmt fix

# The `run` target is intended to build and
# execute the Docker image for the plugin.
#
# Usage: `make run`
.PHONY: run
run: build docker-build docker-run

# The `tidy` target is intended to clean up
# the Go module files (go.mod & go.sum).
#
# Usage: `make tidy`
.PHONY: tidy
tidy:
	@echo
	@echo "### Tidying Go module"
	@go mod tidy

# The `vet` target is intended to inspect the
# Go source code for potential issues.
#
# Usage: `make vet`
.PHONY: vet
vet:
	@echo
	@echo "### Vetting Go code"
	@go vet ./...

# The `fmt` target is intended to format the
# Go source code to meet the language standards.
#
# Usage: `make fmt`
.PHONY: fmt
fmt:
	@echo
	@echo "### Formatting Go Code"
	@go fmt ./...

# The `fix` target is intended to rewrite the
# Go source code using old APIs.
#
# Usage: `make fix`
.PHONY: fix
fix:
	@echo
	@echo "### Fixing Go Code"
	@go fix ./...

# The `test` target is intended to run
# the tests for the Go source code.
#
# Usage: `make test`
.PHONY: test
test:
	@echo
	@echo "### Testing Go Code"
	@go test -race ./...

# The `test-cover` target is intended to run
# the tests for the Go source code and then
# open the test coverage report.
#
# Usage: `make test-cover`
.PHONY: test-cover
test-cover:
	@echo
	@echo "### Creating test coverage report"
	@go test -race -covermode=atomic -coverprofile=coverage.out ./...
	@echo
	@echo "### Opening test coverage report"
	@go tool cover -html=coverage.out

# The `build` target is intended to compile
# the Go source code into a binary.
#
# Usage: `make build`
.PHONY: build
build:
	@echo
	@echo "### Building release/vela-img binary"
	GOOS=linux CGO_ENABLED=0 \
		go build -a \
		-o release/vela-img \
		github.com/go-vela/vela-img/cmd/vela-img

# The `build-static` target is intended to compile
# the Go source code into a statically linked binary.
#
# Usage: `make build-static`
.PHONY: build-static
build-static:
	@echo
	@echo "### Building static release/vela-img binary"
	GOOS=linux CGO_ENABLED=0 \
		go build -a \
		-ldflags '-s -w -extldflags "-static"' \
		-o release/vela-img \
		github.com/go-vela/vela-img/cmd/vela-img

# The `check` target is intended to output all
# dependencies from the Go module that need updates.
#
# Usage: `make check`
.PHONY: check
check: check-install
	@echo
	@echo "### Checking dependencies for updates"
	@go list -u -m -json all | go-mod-outdated -update

# The `check-direct` target is intended to output direct
# dependencies from the Go module that need updates.
#
# Usage: `make check-direct`
.PHONY: check-direct
check-direct: check-install
	@echo
	@echo "### Checking direct dependencies for updates"
	@go list -u -m -json all | go-mod-outdated -direct

# The `check-full` target is intended to output
# all dependencies from the Go module.
#
# Usage: `make check-full`
.PHONY: check-full
check-full: check-install
	@echo
	@echo "### Checking all dependencies for updates"
	@go list -u -m -json all | go-mod-outdated

# The `check-install` target is intended to download
# the tool used to check dependencies from the Go module.
#
# Usage: `make check-install`
.PHONY: check-install
check-install:
	@echo
	@echo "### Installing psampaz/go-mod-outdated"
	@go get -u github.com/psampaz/go-mod-outdated

# The `bump-deps` target is intended to upgrade
# non-test dependencies for the Go module.
#
# Usage: `make bump-deps`
.PHONY: bump-deps
bump-deps: check
	@echo
	@echo "### Upgrading dependencies"
	@go get -u ./...

# The `bump-deps-full` target is intended to upgrade
# all dependencies for the Go module.
#
# Usage: `make bump-deps-full`
.PHONY: bump-deps-full
bump-deps-full: check
	@echo
	@echo "### Upgrading all dependencies"
	@go get -t -u ./...

# The `docker-build` target is intended to build
# the Docker image for the plugin.
#
# Usage: `make docker-build`
.PHONY: docker-build
docker-build:
	@echo
	@echo "### Building vela-img:local image"
	@docker build --no-cache -t vela-img:local .

# The `docker-test` target is intended to execute
# the Docker image for the plugin with test variables.
#
# Usage: `make docker-test`
.PHONY: docker-test
docker-test:
	@echo
	@echo "### Testing vela-img:local image"
	@docker run --rm \
		-e DOCKER_PASSWORD=secretPassword \	
		-e PARAMETER_REGISTRY=index.docker.io \		
		-e DOCKER_USERNAME=moby \		
		-e PARAMETER_CONTEXT=/workspace/ \
		-e PARAMETER_FILE=Dockerfile.example \
		-e PARAMETER_TAGS=index.docker.io/target/vela-img:latest \
		-v $(shell pwd):/home/user/src:ro \
		--workdir /home/user/src \
		--security-opt seccomp=unconfined --security-opt apparmor=unconfined \
		vela-img:local

# The `docker-run` 	target is intended to execute
# the Docker image for the plugin.
#
# Usage: `make docker-run`
.PHONY: docker-run
docker-run:
	@echo
	@echo "### Executing vela-img:local image"
	@docker run --rm \
		-e PARAMETER_LOG_LEVEL \
		-e DOCKER_PASSWORD \
		-e PARAMETER_REGISTRY \
		-e DOCKER_USERNAME \
		-e PARAMETER_BUILD_ARGS \
		-e PARAMETER_CACHE_FROM \
		-e PARAMETER_FILE \
		-e PARAMETER_LABELS \
		-e PARAMETER_NO_CACHE \
		-e PARAMETER_NO_CONSOLE \
		-e PARAMETER_OUTPUT \
		-e PARAMETER_PLATFORMS \
		-e PARAMETER_TAGS \
		-e PARAMETER_TARGET \
		-v $(shell pwd):/home/user/src:ro \
		--workdir /home/user/src \
		--security-opt seccomp=unconfined --security-opt apparmor=unconfined \
		vela-img:local

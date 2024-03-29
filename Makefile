PACKAGES=$(shell go list ./... | grep -v /vendor/)
RACE := $(shell test $$(go env GOARCH) != "amd64" || (echo "-race"))
VERSION := $(shell cat VERSION)


help:
	@echo 'Available commands:'
	@echo
	@echo 'Usage:'
	@echo '    make deps     		Install go deps.'
	@echo '    make build    		Compile the App.'
	@echo '    make build/docker	Build Docker image.'
	@echo '    make run				Run Locally'
	@echo '    make run/docker		Build Docker image.'
	@echo '    make restore  		Restore all dependencies.'
	@echo '    make clean    		Clean the directory tree.'
	@echo


run:
	@go run -ldflags "-X main.version=`cat VERSION`" `ls *.go | grep -v _test.go` -i=$(i) -o=$(o)

test: ## run tests, except integration tests
	@go test ${RACE} ${PACKAGES}

test/benchmark: ## run benchmark tests
	@go test --bench=^Benchmark.* ${PACKAGES}

deps:
	@go get -u github.com/mitchellh/gox

build:
	@echo "Compiling..."
	@mkdir -p ./bin
	@gox -output "bin/{{.Dir}}_${VERSION}_{{.OS}}_{{.Arch}}" -os="linux" -os="darwin" -arch="386" -arch="amd64" ./
	@go build -i -o ./bin/simscale
	@echo "All done! The binaries are in ./bin Check it out!"

build/docker: build
	@docker build -t simscale:latest .

run/docker: build/docker
	@docker run --name simscale\
	 	--volume ${DOCKER_MOUNT_DIR}:${DOCKER_MOUNT_DIR}\
		--env INPUT_FILE="${DOCKER_MOUNT_DIR}/${DOCKER_INPUT_FILE}" \
		--env OUTPUT_FILE="${DOCKER_MOUNT_DIR}/${DOCKER_OUTPUT_FILE}" \
		 simscale:latest

vet: ## run go vet
	@test -z "$$(go vet ${PACKAGES} 2>&1 | grep -v '*composite literal uses unkeyed fields|exit status 0)' | tee /dev/stderr)"



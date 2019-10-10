MOCK_DIR = "test/mocks"

all: dep test server runner

test:
	go test go-brunel/internal...

integration-test:
	go test go-brunel/test...

server:
	go build cmd/server.go

runner:
	go build cmd/runner.go

run-server:
	go run cmd/server.go --config-file ./server.yaml

run-runner-local:
	go run cmd/runner.go --config-file ./runner-kube-local.yaml

cover:
	go test -coverprofile cover.out go-brunel/internal...
	go tool cover -html=cover.out

certs:
	go run cmd/cert.go

lint:
	wget -O - -q https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s v1.16.0
	${GOPATH}/bin/golangci-lint run ./internal/...

dep:
	dep ensure

mocks:
	go get github.com/golang/mock/gomock
	go install github.com/golang/mock/mockgen

	# Make runtime environment mocks
	rm -rf $(MOCK_DIR)
	mkdir -p $(MOCK_DIR)/mock_docker
	mkdir -p $(MOCK_DIR)/go-brunel/pkg/runner/remote
	mkdir -p $(MOCK_DIR)/go-brunel/pkg/runner/vcs

	${GOPATH}/bin/mockgen -package client github.com/docker/docker/client CommonAPIClient > $(MOCK_DIR)/mock_docker/client.go

	${GOPATH}/bin/mockgen -package vcs go-brunel/internal/pkg/runner/vcs VCS > $(MOCK_DIR)/go-brunel/pkg/runner/vcs/vcs.go
	${GOPATH}/bin/mockgen -package remote go-brunel/internal/pkg/runner/remote Remote > $(MOCK_DIR)/go-brunel/pkg/runner/remote/remote.go

.PHONY: test cover mocks
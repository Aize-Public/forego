test: vet test-long lint test-race

test-race:
	go test --short --race --trimpath ./...

test-short:
	go test --short --trimpath ./...

test-long:
	go test --trimpath ./...

test-ci-cover:
	go test -cover --short --trimpath ./... -coverprofile=profile.cov -covermode=count 2>&1
	go get github.com/boumenot/gocover-cobertura
	go run github.com/boumenot/gocover-cobertura < profile.cov > coverage.xml

vet:
	go vet ./...

fmt:
	go fmt ./...

lint:
	golangci-lint run --timeout 5m --out-format tab --skip-dirs '(example)' ./...

lint-setup:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.2
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

mod:
	go mod tidy -v

GOFMT_FILES?=$$(find . -name '*.go')

build: terraform-provider-lavinmq

tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest

terraform-provider-lavinmq:
	go build -o terraform-provider-lavinmq

install:
	go install .

fmt:
	gofmt -s -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

lint:
	golangci-lint run ./...

test:
	TF_ACC=1 go test ./lavinmq -v -count 1

coverage:
	TF_ACC=1 go test ./lavinmq -v -coverprofile=coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -f terraform-provider-lavinmq coverage.out coverage.html

docs: clean build
	go generate ./...

.PHONY: clean install fmt fmtcheck lint tools test coverage docs

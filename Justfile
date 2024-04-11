# List the available justfile recipes.
@default:
  just --list

# Format, vet, and test Go code.
check:
	go fmt ./...
	go vet ./...
	GOEXPERIMENT=loopvar go test ./... -cover

# Verbosely format, vet, and test Go code.
checkv:
	go fmt ./...
	go vet ./...
	GOEXPERIMENT=loopvar go test -v ./... -cover

# Lint code using staticcheck.
lint:
	staticcheck -f stylish ./...

# Test and provide HTML coverage report.
cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

# List the outdated go modules.
outdated:
  go list -u -m all

# Build and run the LXI Keysight 33220A example application.
k33220 ip:
  #!/usr/bin/env bash
  echo '# IVI LXI Keysight 33220A Example Application'
  cd {{justfile_directory()}}/examples/key33220
  env go build -o key33220
  ./key33220 -ip={{ip}}

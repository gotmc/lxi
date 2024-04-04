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

# Build and run Agilent 33220A example application.
a33220 ip:
  env go build -o ./examples/agilent33220/agilent33220 ./examples/agilent33220/
  ./examples/agilent33220/agilent33220 -ip=TCPIP0::{{ip}}::5025::SOCKET

# List the outdated go modules.
outdated:
  go list -u -m all

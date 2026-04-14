# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go package implementing the LAN eXtensions for Instrumentation (LXI) standard for sending SCPI commands to test equipment over TCP/IP. Part of the gotmc ecosystem: [visa](https://github.com/gotmc/visa) defines a common instrument interface across transports (GPIB, USB, TCP/IP, serial), [asrl](https://github.com/gotmc/asrl) provides serial transport, and [ivi](https://github.com/gotmc/ivi) builds on visa for standardized instrument-class APIs per IVI Foundation specs.

## Build & Test Commands

Requires [just](https://just.systems/man/en/) task runner.

```bash
# Format and vet
just check

# Run unit tests (includes check, race detector, coverage)
just unit

# Verbose unit tests
just unit -v

# Lint (requires golangci-lint v2, configured via .golangci.yaml)
just lint

# Test coverage (opens HTML report)
just cover

# Run single test
go test -run TestParsingVisaResourceString ./...

# Dependency management
just tidy          # go mod tidy + verify
just outdated      # list outdated direct deps (requires go-mod-outdated)
just update <mod>  # update a specific module

# Run example against real hardware
just k33220 192.168.1.101
```

## Architecture

The package has two core types:

- **`Device`** (`lxi.go`): Wraps a TCP connection to an LXI instrument. Implements `io.Reader`, `io.Writer`, `io.Closer`, and `io.StringWriter`. High-level SCPI methods: `Command(ctx, ...)` auto-appends the EndMark character, and `Query(ctx, cmd)` sends a command then reads the response (stripping the trailing EndMark). Context-aware binary I/O: `ReadBinary` and `WriteBinary`. Non-context `Read`, `Write`, and `WriteString` delegate to their context-aware counterparts with `context.Background()`.
- **`VisaResource`** (`visa.go`): Parses VISA resource strings (format: `TCPIP<board>::<host>::<port>::SOCKET`) using a package-level compiled regex. Only TCPIP/SOCKET interface type is supported. Input is case-insensitive; output is always uppercase canonical form.

`NewDevice(ctx, address)` parses the address via `NewVisaResource()`, then dials a TCP connection using the context for timeout/cancellation.

The internal `applyContext` method is the key context-handling mechanism: it sets connection deadlines from context deadlines, and for cancelable contexts without deadlines, spawns a goroutine that watches for cancellation and forces an immediate deadline to unblock pending I/O. The returned cleanup function stops the goroutine and resets the deadline.

## Conventions

- **Zero external dependencies** — the library itself uses only the Go standard library (Go 1.21+). The `gotmc/query` dependency in `go.mod` is used only by the example application in `examples/`.
- **Error prefixes** — errors use the format `"visa: ..."` or `"lxi: ..."` to identify their origin. Sentinel errors are defined in `visa.go` and use `%w` wrapping for `errors.Is()` checking.
- **Context handling** — `Command`, `Query`, `ReadBinary`, and `WriteBinary` apply context deadlines/cancellation to the TCP connection via `applyContext`. When context cancellation causes an I/O error, the context error (`context.Canceled` / `context.DeadlineExceeded`) is returned instead of the raw network timeout, following the `net.Dialer.DialContext` pattern.
- **Test pattern** — table-driven tests with structs defining inputs and expected outputs. Tests use `net.Pipe` for in-process TCP simulation via `newTestDevice` helper in `lxi_test.go`.

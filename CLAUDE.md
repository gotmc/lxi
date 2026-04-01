# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go package implementing the LAN eXtensions for Instrumentation (LXI) standard for sending SCPI commands to test equipment over TCP/IP. Part of the gotmc ecosystem alongside [ivi](https://github.com/gotmc/ivi) and [visa](https://github.com/gotmc/visa) packages.

## Build & Test Commands

Requires [just](https://just.systems/man/en/) task runner.

```bash
# Format and vet
just check

# Run unit tests (includes check, race detector, coverage)
just unit

# Verbose unit tests
just unit -v

# Lint (requires golangci-lint, configured via .golangci.yaml)
just lint

# Test coverage (opens HTML report)
just cover

# Run single test
go test -run TestParsingVisaResourceString ./...
```

## Architecture

The package has two core types:

- **`Device`** (`lxi.go`): Wraps a TCP connection to an LXI instrument. Implements `io.Reader`, `io.Writer`, `io.Closer`, and `io.StringWriter`. Provides `Command(ctx, ...)` (auto-appends EndMark) and `Query(ctx, cmd)` (sends command, reads response) for SCPI communication. Also provides `ReadContext` and `WriteContext` for context-aware raw I/O.
- **`VisaResource`** (`visa.go`): Parses VISA resource strings (format: `TCPIP<board>::<host>::<port>::SOCKET`) using a package-level compiled regex. Only TCPIP/SOCKET interface type is supported.

`NewDevice()` takes a VISA address string, parses it via `NewVisaResource()`, then dials a TCP connection.

The internal `applyContext` method is the key context-handling mechanism: it sets connection deadlines from context deadlines, and for cancelable contexts without deadlines, spawns a goroutine that watches for cancellation and forces an immediate deadline to unblock pending I/O. The returned cleanup function stops the goroutine and resets the deadline.

## Conventions

- **Zero external dependencies** — only the Go standard library is used.
- **Error prefixes** — errors use the format `"visa: ..."` or `"lxi: ..."` to identify their origin.
- **Context handling** — `Command`, `Query`, `ReadContext`, and `WriteContext` apply context deadlines/cancellation to the TCP connection via `applyContext`.
- **Test pattern** — table-driven tests with structs defining inputs and expected outputs. Tests use `net.Pipe` for in-process TCP simulation (see `newTestDevice` in `lxi_test.go`).

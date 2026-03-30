# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go package implementing the LAN eXtensions for Instrumentation (LXI) standard for sending SCPI commands to test equipment over TCP/IP. Part of the gotmc ecosystem alongside [ivi](https://github.com/gotmc/ivi) and [visa](https://github.com/gotmc/visa) packages.

## Build & Test Commands

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

- **`Device`** (`lxi.go`): Wraps a TCP connection to an LXI instrument. Implements `io.Reader`, `io.Writer`, and `io.Closer`. Provides `Command(ctx, ...)` (auto-appends EndMark) and `Query(ctx, cmd)` (sends command, reads response) for SCPI communication. Both accept `context.Context` and apply its deadline to the underlying connection.
- **`VisaResource`** (`visa.go`): Parses VISA resource strings (format: `TCPIP<board>::<host>::<port>::SOCKET`) using regex. Only TCPIP/SOCKET interface type is supported.

`NewDevice()` takes a VISA address string, parses it via `NewVisaResource()`, then dials a TCP connection.

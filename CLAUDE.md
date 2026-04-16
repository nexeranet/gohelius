# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Go client library for the [Helius API](https://docs.helius.dev/) — a Solana blockchain data provider. Package name: `helius`, module: `github.com/nexeranet/gohelius`.

## Commands

- **Run all tests:** `go test ./...`
- **Run tests (skip live API):** `go test -short ./...`
- **Run a single test:** `go test -run TestGetTransactions ./...`
- **Build:** `go build ./...`

## Architecture

- `client.go` — HTTP client with built-in rate limiting (`golang.org/x/time/rate`, ~5 req/s). All API calls go through `Client.call()` which handles auth (api-key query param), rate limiting, and JSON deserialization.
- `types.go` — Response types (`Transaction`, `TokenTransfer`, `NativeTransfer`).
- `client_test.go` — Live integration tests against the real Helius API. Requires `HELIUS_API_KEY` env var; skipped in `-short` mode or when the key is absent.

## Key Details

- Tests hit the real Helius API — there are no mocks. Set `HELIUS_API_KEY` to run them.
- The client targets Helius mainnet by default (`https://api-mainnet.helius-rpc.com`). Override via `Client.BaseURL`.
- Rate limiter is exposed as `Client.Limiter` for caller customization.

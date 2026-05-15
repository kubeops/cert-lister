# AGENTS.md

This file provides guidance to coding agents (e.g. Claude Code, claude.ai/code) when working with code in this repository.

## Repository purpose

Go module `kubeops.dev/cert-lister` — a small CLI that enumerates X.509 certificates and signing material that live in a Kubernetes cluster's Secrets, ConfigMaps, and APIService `caBundle`s, and prints their key name, serial number, and age in a tab-aligned table. Useful for auditing TLS material across a cluster without manually decoding every PEM.

The produced binary is `cert-lister`. Single-file utility: `main.go` is the entire program. No long-running process; you just run it against a kubeconfig.

## Architecture

- `main.go` — the entire program. Connects via `client-go`, walks `Secrets`, `ConfigMaps`, and `APIService` objects, decodes the PEM blocks using `gomodules.xyz/cert`, and tabulates output. Includes the cert-manager scheme so cert-manager-issued Secrets show through.
- `hack/`, `Makefile` — AppsCode build harness (everything runs inside `ghcr.io/appscode/golang-dev`). Binary builds for **5 platforms**: linux amd64/arm/arm64 plus `windows/amd64`, `darwin/amd64`, `darwin/arm64` (operators run this from their workstations).
- `vendor/` — checked-in deps.

There is no Docker image; this is a host CLI.

## Common commands

All Make targets run inside `ghcr.io/appscode/golang-dev` — Docker must be running.

- `make ci` — CI pipeline.
- `make build` — build for the host OS/ARCH into `bin/<os>_<arch>/cert-lister`.
- `make all-build` — build for every `BIN_PLATFORMS` (linux amd64/arm/arm64 + windows/amd64 + darwin/amd64 + darwin/arm64).
- `make fmt`, `make lint`, `make unit-tests` / `make test` — standard.
- `make verify` — `verify-gen verify-modules`; `go mod tidy && go mod vendor` must leave the tree clean.
- `make add-license` / `make check-license` — manage Apache-2.0 license headers.

Run a single Go test (requires a local Go toolchain):

```
go test ./... -run TestName -v
```

To run against a cluster:

```
./bin/<os>_<arch>/cert-lister --kubeconfig ~/.kube/config
```

## Conventions

- Module path is `kubeops.dev/cert-lister` (vanity URL). Imports must use that.
- License: Apache-2.0 (`LICENSE`). New code needs the standard "Copyright AppsCode Inc. and Contributors" header (`make add-license`).
- Sign off commits (`git commit -s`).
- Vendor directory is checked in — `go mod tidy && go mod vendor` must leave the tree clean (enforced by `verify-modules`).
- Keep this binary small. It's intentionally a one-file utility — if a new feature needs heavy structure (controller-runtime, codegen, etc.), it probably belongs in a different repo.
- Builds linux/windows/darwin host binaries; do not pull in linux-only or cgo deps.

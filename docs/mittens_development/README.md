# Mittens Development

This section provides information about Mittens development processes and how
to contribute.

## Mittens Overview

Mittens is a fork of the original [kubetap](https://github.com/soluble-ai/kubetap) project by [Soluble](https://www.soluble.ai/), maintained as an [Apache v2 Licensed](https://github.com/Lappihuan/mittens/blob/master/LICENSE) open source project.

Mittens focuses on providing a modern, mitmproxy-centric approach to intercepting and inspecting Kubernetes Service traffic, emphasizing interactive debugging and security testing workflows.

### Building Mittens

Building mittens requires the following dependencies:

| dependency | purpose               | notes                                         |
| ---        | ---                   | ---                                           |
| `kubectl`  | ...                   | mandatory for integration tests               |
| `docker`   | Build containers      | not needed to build `kubectl-mittens` binary  |
| `go`       | Build `kubectl-mittens`| minimum Go version 1.25                       |
| `zsh`      | Build scripts         | scripting is nicer than `bash` or `sh`        |

### Script-managed Dependencies

Installed using `go get` and `ci.mod` or `ig-tests.mod`:

| dependency      | purpose | notes                                                      |
| ---             | ---     | ---                                                        |
| `golangci-lint` | Linting | (`ci.mod`) used as Go code linter                          |
| `gotestsum`     | Testing | (`ci.mod`) used to make test output prettier               |
| `kind`          | Testing | (`ig-tests.mod`) required for integration tests            |
| `helm`          | Testing | (`ig-tests.mod`) used to deploy test apps to  kind cluster |

## Hacking on Mittens

Assuming you have [built mittens from source](../getting_started/installation.md),
you're ready to hack on mittens.

```sh
$ cd ${GOPATH}/src/github.com/Lappihuan/mittens

$ go generate .

$ go build ./cmd/kubectl-mittens
```

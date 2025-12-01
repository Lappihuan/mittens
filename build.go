//go:generate ./scripts/deps.sh
//go:generate go clean -i ./...
//go:generate rm -f ./cmd/kubectl-mittens/kubectl-mittens
//go:generate go mod download
//go:generate gotestsum --format=short-verbose --no-summary=skipped --junitfile=coverage.xml -- -count=1 -race -coverprofile=coverage.txt -covermode=atomic ./...
package main

func main() {}

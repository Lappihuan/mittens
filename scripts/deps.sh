#!/bin/sh

GOLANGCI_LINT_VERSION=v1.35.2
GOTESTSUM_VERSION=v0.6.0
KIND_VERSION=v0.9.0
HELM_VERSION=v3.5.0
GOFUMPT_VERSION=v0.3.0
GOFUMPORTS_VERSION=v0.3.0

cd

if ! [ -x "$(command -v kubectl)" ]; then
  echo "kubectl is not installed"
  exit 1
fi


if ! [ -x "$(command -v golangci-lint)" ]; then
  # Use 'go install' with a version tag (works outside modules)
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}
fi

if ! [ -x "$(command -v gotestsum)" ]; then
  # Try to install via 'go install'; some module versions expose replace
  # directives that can cause 'go install' to fail. Fall back to the
  # released binary if go install fails.
  if ! go install gotest.tools/gotestsum@${GOTESTSUM_VERSION}; then
    echo "go install gotestsum failed, downloading release binary..."
    GOTESTSUM_BIN="gotestsum_linux_amd64"
    GOTESTSUM_URL="https://github.com/gotestyourself/gotestsum/releases/download/${GOTESTSUM_VERSION}/${GOTESTSUM_BIN}"
    mkdir -p "${GOBIN:-$(go env GOPATH)/bin}"
    curl -fsSL -o "${GOBIN:-$(go env GOPATH)/bin}/gotestsum" "$GOTESTSUM_URL" && chmod +x "${GOBIN:-$(go env GOPATH)/bin}/gotestsum"
  fi
fi

if ! [ -x "$(command -v helm)" ]; then
  # official module path for Helm v3
  go install helm.sh/helm/v3/cmd/helm@${HELM_VERSION}
fi


if ! [ -x "$(command -v kind)" ]; then
  go install sigs.k8s.io/kind@${KIND_VERSION}
fi

if ! [ -x "$(command -v gofumpt)" ]; then
  go install mvdan.cc/gofumpt@${GOFUMPT_VERSION}
fi

if ! [ -x "$(command -v gofumports)" ]; then
  # Some module versions do not contain the gofumports subpackage in the
  # top-level module version; try the pinned version first and fall back
  # to @latest if needed.
  if ! go install mvdan.cc/gofumpt/gofumports@${GOFUMPORTS_VERSION}; then
    echo "Pinned gofumports install failed, trying @latest..."
    go install mvdan.cc/gofumpt/gofumports@latest || true
  fi
fi

cd -

#!/usr/bin/env zsh

script_dir=${0:A:h}
source ${script_dir}/_pre.zsh

#
# Build and install mittens plugin
#
echo "Building and installing mittens plugin..."
go install -v -trimpath -ldflags="-s -w" ./cmd/kubectl-mittens
if [[ $? -ne 0 ]]; then
  echo "Failed to build mittens plugin"
  return 1
fi

# Ensure Go bin is in PATH
export PATH="$(go env GOPATH)/bin:${PATH}"
echo "PATH set to include Go bin: $(go env GOPATH)/bin"

#
# Prep env
#

KIND_VERSION=v0.30.0
HELM_VERSION=v4.0.1

# modfile hack to avoid package collisions and work around azure go-autorest bug
print "module mittens-ig-tests

go 1.13

require (
)

replace (
  github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.0.0+incompatible
)" >! ig-tests.mod

# return error if kubectl is not available
if [[ =kubectl == '' ]]; then
  echo "kubectl not installed"
  return 1
fi

# install helm if not available
if [[ =helm == '' ]]; then
  GO111MODULE=on go get -modfile=ig-tests.mod helm.sh/helm/v4/cmd/helm@${HELM_VERSION}
fi
helm repo add grafana https://grafana.github.io/helm-charts --force-update
helm repo update

# we use kind to establish a local testing cluster
GO111MODULE=on go get -modfile=ig-tests.mod sigs.k8s.io/kind@${KIND_VERSION}

#
# Establish a local testing cluster
#

# remove stale mittens clusters if they exist
_mittens_kind_clusters=$(kind get clusters 2>&1)
if [[ ${_mittens_kind_clusters} == *'mittens'* ]]; then
  kind delete cluster --name mittens
fi
unset _mittens_kind_clusters

# catch sigints and exits to delete the cluster, keeping the last exit code
trap '{ e=${?}; sleep 1; kind delete cluster --name mittens ; exit ${e} }' SIGINT SIGTERM EXIT
kind create cluster --name mittens


#
# Test mittens using helm ${chart}
#
_mittens_helm_charts=('grafana/grafana' 'oci://registry-1.docker.io/bitnamicharts/nginx')
_mittens_helm_services=('grafana' 'nginx')
_mittens_helm_svc_port=('80' '80')

typeset -i _mittens_iter
for chart in ${_mittens_helm_charts[@]}; do
  ((_mittens_iter+=1))
  _mittens_helm=${chart:t}
  _mittens_port=${_mittens_helm_svc_port[${_mittens_iter}]}
  _mittens_service=${_mittens_helm_services[${_mittens_iter}]}

  helm install --kube-context kind-mittens ${_mittens_helm} ${chart}
  
  # We need to run mittens in a non-interactive way for CI environments
  # Use timeout and pipe input to avoid waiting for interactive prompt
  # This will cause mittens to fail after pod is ready (expected in CI)
  _mittens_output=$(timeout 45 bash -c "echo '' | kubectl mittens ${_mittens_service} -p${_mittens_port} --context kind-mittens" 2>&1)
  _mittens_exit_code=$?
  
  # Check if mittens successfully set up the proxy (pod should be ready)
  # We expect exit code 1 in CI since kubectl exec -it will fail without a real terminal
  if [[ ${_mittens_output} == *"Pod ready!"* ]]; then
    echo "✓ Mittens proxy setup successful for ${_mittens_service}"
  else
    echo "✗ Mittens failed to set up proxy for ${_mittens_service}"
    echo "Output: ${_mittens_output}"
    return 1
  fi
done
unset _mittens_helm_charts _mittens_helm_services _mittens_helm_svc_port _mittens_iter

#source ${script_dir}/_post.zsh

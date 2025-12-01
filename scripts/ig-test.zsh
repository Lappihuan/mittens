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
helm repo add podinfo https://stefanprodan.github.io/podinfo --force-update
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
_mittens_helm_charts=('podinfo/podinfo')
_mittens_helm_services=('podinfo')
_mittens_helm_svc_port=('9898')

typeset -i _mittens_iter
for chart in ${_mittens_helm_charts[@]}; do
  ((_mittens_iter+=1))
  _mittens_helm=${chart:t}
  _mittens_port=${_mittens_helm_svc_port[${_mittens_iter}]}
  _mittens_service=${_mittens_helm_services[${_mittens_iter}]}

  helm install --kube-context kind-mittens ${_mittens_helm} ${chart}
  
  # Wait for the service to be running
  sleep 5
  
  # Run mittens in background and check for sidecar injection while it's running
  # Start mittens in background (will fail on kubectl exec but that's ok)
  timeout 45 bash -c "echo '' | kubectl mittens ${_mittens_service} -p${_mittens_port} --context kind-mittens" >/dev/null 2>&1 &
  _mittens_pid=$!
  
  # Wait for pod with mittens sidecar to become ready (check every second, max 40 seconds)
  _ready=0
  for ((i=0; i<40; i++)); do
    # Get pods for this specific service and check if they have mittens container and are ready
    # Use app.kubernetes.io/name label which podinfo sets
    _pod_status=$(kubectl get pods --context kind-mittens -l "app.kubernetes.io/name=${_mittens_service}" -o jsonpath="{.items[*].status.conditions[?(@.type=='Ready')].status}" 2>/dev/null)
    _mittens_pods=$(kubectl get pods --context kind-mittens -l "app.kubernetes.io/name=${_mittens_service}" -o jsonpath="{.items[*].spec.containers[*].name}" 2>/dev/null)
    
    if [[ ${_pod_status} == *"True"* ]] && [[ ${_mittens_pods} == *"mittens"* ]]; then
      _ready=1
      break
    fi
    sleep 1
  done
  
  # Check if pod became ready with mittens container
  if [[ $_ready -eq 0 ]]; then
    echo "✗ Pod did not become ready or mittens sidecar not injected for ${_mittens_service}"
    kill $_mittens_pid 2>/dev/null
    wait $_mittens_pid 2>/dev/null
    return 1
  fi
  
  # Kill mittens process if still running (ignore errors if already dead)
  kill $_mittens_pid 2>/dev/null || true
  wait $_mittens_pid 2>/dev/null || true
  
  echo "✓ Mittens proxy setup successful for ${_mittens_service}"
done
unset _mittens_helm_charts _mittens_helm_services _mittens_helm_svc_port _mittens_iter

#source ${script_dir}/_post.zsh

#!/usr/bin/env zsh

script_dir=${0:A:h}
#source ${script_dir}/_pre.zsh

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
  # Start mittens in background; it will auto-cleanup when interrupted
  kubectl mittens ${_mittens_service} -p${_mittens_port} --context kind-mittens &
  _mittens_mittens_pid=${!}
  sleep 20

  _mittens_ready_state=""
  for i in {0..20}; do
    sleep 6
    _mittens_pod=($(kubectl --context kind-mittens get pods -ojsonpath='{.items[*].metadata.name}'))
    if (( ${#_mittens_pod} != 1 )); then
      continue
    fi
    _mittens_ready_state=$(kubectl --context kind-mittens get pod ${_mittens_pod} -ojsonpath='{.status.containerStatuses[*].ready}')
    if [[ ${_mittens_ready_state} == 'true true' ]]; then
      break
    fi
  done
  if [[ ${_mittens_ready_state} != 'true true' ]]; then
    echo "container did not come up within 90 seconds"
    echo ""
    echo "=== DEBUG INFO ==="
    echo "Pod status:"
    kubectl --context kind-mittens get pods -o wide
    echo ""
    echo "Pod events:"
    kubectl --context kind-mittens describe pod ${_mittens_pod}
    echo ""
    echo "Pod logs (all containers):"
    for _container in $(kubectl --context kind-mittens get pod ${_mittens_pod} -o jsonpath='{.spec.containers[*].name}'); do
      echo "--- Container: ${_container} ---"
      kubectl --context kind-mittens logs ${_mittens_pod} -c ${_container} 2>&1 || echo "No logs available"
    done
    echo "=================="
    echo ""
    return 1
  fi
  unset _mittens_pod _mittens_ready_state i

  sleep 1
  kubectl port-forward svc/${_mittens_service} -n default 4000:${_mittens_port} &
  _mittens_pf_pid=${!}
  sleep 5

  # check that we can reach the application port
  curl -v http://127.0.0.1:4000 || return 1
  kill ${_mittens_pf_pid}
  unset _mittens_pf_pid

  # cleanup - kill mittens process which triggers auto-cleanup
  kill ${_mittens_mittens_pid}
  sleep 2
  helm delete --kube-context kind-mittens ${_mittens_helm}

  unset _mittens_helm _mittens_port _mittens_service
done
unset _mittens_helm_charts _mittens_helm_services _mittens_helm_svc_port _mittens_iter

#source ${script_dir}/_post.zsh

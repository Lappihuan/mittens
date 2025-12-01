#!/usr/bin/env zsh

script_dir=${0:A:h}
source ${script_dir}/_pre.zsh

go clean -i ./...
rm -f ./cmd/kubectl-mittens/kubectl-mittens
rm -f ./kubectl-mittens
rm -rf ./site/

source ${script_dir}/_post.zsh

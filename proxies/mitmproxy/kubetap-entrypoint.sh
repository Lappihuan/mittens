#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

# HACK: this fixes permission issues
# Ensure the .mitmproxy directory exists and is writable
# Note: We skip this if we don't have permissions, as the directory should 
# already exist from the Dockerfile
if [ -d /home/mitmproxy/.mitmproxy ] && [ -w /home/mitmproxy/.mitmproxy ]; then
  # Only copy the config file if it exists and we have read access
  if [ -f /home/mitmproxy/config/config.yaml ] && [ -r /home/mitmproxy/config/config.yaml ]; then
    cp /home/mitmproxy/config/config.yaml /home/mitmproxy/.mitmproxy/config.yaml
  fi
fi

prog="${1}"
if [[ "${1}" == 'mitmdump' || "${1}" == 'mitmproxy' || "${1}" == 'mitmweb' ]]; then
  MITMPROXY_PATH='/home/mitmproxy/.mitmproxy'
  
  # For mitmproxy interactive terminal mode, use tmux to handle TTY requirements
  if [[ "${1}" == 'mitmproxy' ]]; then
    # Start a tmux session with mitmproxy to allow interactive access without requiring a TTY
    # Users can 'kubectl exec -it <pod> -- tmux attach-session -t mitmproxy' to interact
    # Use -c to specify a new window command, avoiding terminal requirement at session creation
    tmux new-session -d -s mitmproxy -c /home/mitmproxy \
      mitmproxy --set "confdir=${MITMPROXY_PATH}" "${@:2}"
    
    # Keep the container running - tail a log or sleep indefinitely
    # This allows users to attach via: kubectl exec -it <pod> -- tmux attach-session -t mitmproxy
    sleep infinity
  else
    # For mitmdump or mitmweb (or other commands), use direct execution
    exec "${@}" --set "confdir=${MITMPROXY_PATH}"
  fi
else
  exec "${@}"
fi

#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

# Ensure the .mitmproxy directory exists with proper permissions
mkdir -p /home/mitmproxy/.mitmproxy
chmod 777 /home/mitmproxy/.mitmproxy

# Copy the config file if it exists and is readable
if [ -f /home/mitmproxy/config/config.yaml ] && [ -r /home/mitmproxy/config/config.yaml ]; then
  cp /home/mitmproxy/config/config.yaml /home/mitmproxy/.mitmproxy/config.yaml
  chmod 666 /home/mitmproxy/.mitmproxy/config.yaml
  echo "Config file copied to /home/mitmproxy/.mitmproxy/config.yaml" >&2
else
  echo "Warning: Config file not found or not readable at /home/mitmproxy/config/config.yaml" >&2
fi

prog="${1}"
if [[ "${1}" == 'mitmdump' || "${1}" == 'mitmproxy' || "${1}" == 'mitmweb' ]]; then
  MITMPROXY_PATH='/home/mitmproxy/.mitmproxy'
  
  # For mitmproxy interactive terminal mode, use tmux to handle TTY requirements
  if [[ "${1}" == 'mitmproxy' ]]; then
    # Start a tmux session with mitmproxy to allow interactive access without requiring a TTY
    # Users can 'kubectl exec -it <pod> -- tmux attach-session -t mitmproxy' to interact
    echo "Starting mitmproxy in tmux session with confdir=${MITMPROXY_PATH}" >&2
    tmux new-session -d -s mitmproxy -c /home/mitmproxy \
      mitmproxy --set "confdir=${MITMPROXY_PATH}" "${@:2}"
    
    # Keep the container running - sleep indefinitely
    # This allows users to attach via: kubectl exec -it <pod> -- tmux attach-session -t mitmproxy
    sleep infinity
  else
    # For mitmdump or mitmweb (or other commands), use direct execution
    echo "Starting ${prog} with confdir=${MITMPROXY_PATH}" >&2
    exec "${@}" --set "confdir=${MITMPROXY_PATH}"
  fi
else
  echo "Running command: ${@}" >&2
  exec "${@}"
fi

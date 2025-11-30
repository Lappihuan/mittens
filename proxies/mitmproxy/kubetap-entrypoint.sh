#!/bin/sh

set -o errexit
set -o pipefail
set -o nounset

# HACK: this fixes permission issues
# Ensure the .mitmproxy directory exists and is writable
# Note: We skip this if we don't have permissions, as the directory should 
# already exist from the Dockerfile
if [ -d /home/mitmproxy/.mitmproxy ] && [ -w /home/mitmproxy/.mitmproxy ]; then
  # Only copy the config file if it exists and we have write access
  if [ -f /home/mitmproxy/config/config.yaml ] && [ -r /home/mitmproxy/config/config.yaml ]; then
    cp /home/mitmproxy/config/config.yaml /home/mitmproxy/.mitmproxy/config.yaml
  fi
fi

prog=${1}
if [[ ${1} == 'mitmdump' || ${1} == 'mitmproxy' || ${1} == 'mitmweb' ]]; then
  MITMPROXY_PATH='/home/mitmproxy/.mitmproxy'
  
  # For mitmproxy interactive terminal mode, use tmux to handle TTY requirements
  if [[ ${1} == 'mitmproxy' ]]; then
    # Start a tmux session with mitmproxy to allow interactive access without requiring a TTY
    # Users can 'kubectl exec -it <pod> -- tmux attach-session -t mitmproxy' to interact
    tmux new-session -d -s mitmproxy \; \
      send-keys -t mitmproxy "exec mitmproxy --set confdir=${MITMPROXY_PATH} ${@:2}" Enter \; \
      capture-pane -t mitmproxy -p
    # Keep the container running by attaching to the session (this allows logs to flow)
    exec tmux attach-session -t mitmproxy
  else
    # For mitmproxy (or other commands), use direct execution
    exec ${@} --set "confdir=${MITMPROXY_PATH}"
  fi
else
  exec ${@}
fi

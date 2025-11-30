#!/bin/bash

set -o pipefail

# Copy the config file if it exists and is readable
if [ -f /home/mitmproxy/config/config.yaml ] && [ -r /home/mitmproxy/config/config.yaml ]; then
  cp /home/mitmproxy/config/config.yaml /home/mitmproxy/.mitmproxy/config.yaml
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
    
    # Create a tmux session and capture any startup errors
    tmux new-session -d -s mitmproxy -c /home/mitmproxy -x 200 -y 50 \
      "mitmproxy --set confdir=${MITMPROXY_PATH} ${@:2}; bash"
    
    echo "Mitmproxy tmux session created. Waiting for it to start..." >&2
    sleep 3
    
    # Check if the session is still running
    if tmux has-session -t mitmproxy 2>/dev/null; then
      echo "Mitmproxy session is running" >&2
      echo "=== Current tmux pane content ===" >&2
      tmux capture-pane -t mitmproxy -p >&2
      echo "=== End tmux pane content ===" >&2
    else
      echo "ERROR: Mitmproxy session exited immediately" >&2
    fi
    
    # Create a bash wrapper script in /usr/local/bin that auto-attaches to tmux
    cat > /usr/local/bin/shell-wrapper.sh << 'WRAPPER_EOF'
#!/bin/bash
# Auto-attach to mitmproxy tmux session if it exists and we're in an interactive shell
if [[ $- == *i* ]] && [ -z "$TMUX" ]; then
  if tmux has-session -t mitmproxy 2>/dev/null; then
    exec tmux attach-session -t mitmproxy
  fi
fi
# Fall through to normal bash if no tmux session or if already in tmux
exec /bin/bash "$@"
WRAPPER_EOF
    chmod +x /usr/local/bin/shell-wrapper.sh
    
    # Keep the container running - sleep indefinitely
    # This allows users to attach via: kubectl exec -it <pod> -- tmux attach-session -t mitmproxy
    # or just: kubectl exec -it <pod> -- bash
    echo "Container keeping alive with sleep infinity" >&2
    sleep infinity
  else
    # For mitmdump or mitmweb (or other commands), use direct execution
    echo "Starting ${prog} with confdir=${MITMPROXY_PATH}" >&2
    exec "${@}" --set "confdir=${MITMPROXY_PATH}"
  fi
else
  # If no command specified, default to bash which will auto-attach to tmux
  if [ -z "$1" ]; then
    if tmux has-session -t mitmproxy 2>/dev/null; then
      echo "Attaching to mitmproxy session..." >&2
      exec tmux attach-session -t mitmproxy
    else
      exec /bin/bash
    fi
  fi
  
  # For other shells/commands, try to auto-attach if it's a shell command
  if [[ "$1" == "bash" || "$1" == "/bin/bash" || "$1" == "sh" || "$1" == "/bin/sh" ]]; then
    if tmux has-session -t mitmproxy 2>/dev/null; then
      echo "Attaching to mitmproxy session..." >&2
      exec tmux attach-session -t mitmproxy
    fi
  fi
  echo "Running command: ${@}" >&2
  exec "${@}"
fi

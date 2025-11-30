FROM mitmproxy/mitmproxy:12.2.1

# Update and install tmux in the Debian-based image
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y --no-install-recommends tmux && \
    rm -rf /var/lib/apt/lists/*

# Ensure the directory exists and is world-writable BEFORE we switch users
RUN mkdir -p /home/mitmproxy/.mitmproxy && \
    chmod 777 /home/mitmproxy && \
    chmod 777 /home/mitmproxy/.mitmproxy

# Hijack the mitmproxy entrypoint (docker-entrypoint.sh) so that
# configuration can be built from within the container using the
# kubetap binary.

COPY kubetap-entrypoint.sh /usr/local/bin/
ENTRYPOINT ["kubetap-entrypoint.sh"]
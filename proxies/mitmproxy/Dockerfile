FROM mitmproxy/mitmproxy:12.2.1

# Update and install tmux in the Debian-based image
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y --no-install-recommends tmux && \
    rm -rf /var/lib/apt/lists/*

# HACK: the security context of the injected pod could be run as any user, therefore
# all users must be able to write to the directory.
RUN chmod -R 777 /home/mitmproxy/.mitmproxy/

# Hijack the mitmproxy entrypoint (docker-entrypoint.sh) so that
# configuration can be built from within the container using the
# kubetap binary.

COPY kubetap-entrypoint.sh /usr/local/bin/
ENTRYPOINT ["kubetap-entrypoint.sh"]

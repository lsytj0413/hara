FROM mcr.microsoft.com/devcontainers/go:1.23-bullseye

RUN sed -i 's|http://deb.debian.org/|http://mirrors.aliyun.com/|g' /etc/apt/sources.list \
    && apt-get update \
    && PREFIX="/usr/local" \
    && VERSION="1.51.0" \
    && curl -sSL "https://github.com/bufbuild/buf/releases/download/v${VERSION}/buf-$(uname -s)-$(uname -m).tar.gz" | \
    tar -xvzf - -C "${PREFIX}" --strip-components 1 \
    && rm -rf /var/lib/apt/lists/*
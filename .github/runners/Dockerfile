FROM ubuntu:24.04
ENV DEBIAN_FRONTEND=noninteractive RUNNER_ALLOW_RUNASROOT=1 AGENT_TOOLSDIRECTORY=/opt/hostedtoolcache
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates curl git python3 python3-venv build-essential cmake pkg-config libssl-dev libicu74 zlib1g liblttng-ust1 libkrb5-3 unzip xz-utils zstd procps sudo && rm -rf /var/lib/apt/lists/*
WORKDIR /runner
RUN curl -fsSL https://github.com/actions/runner/releases/download/v2.338.0/actions-runner-linux-arm64-2.338.0.tar.gz -o runner.tar.gz && echo '628b4a7258487b80c1d3c221095a7ded349f5ae0175fd3dfab708d07428f041b  runner.tar.gz' | sha256sum -c - && tar xzf runner.tar.gz && rm runner.tar.gz
RUN mkdir -p /opt/hostedtoolcache
RUN curl --proto =https --tlsv1.2 -sSf https://sh.rustup.rs -o /tmp/rustup.sh && sh /tmp/rustup.sh -y --profile minimal --default-toolchain 1.99.0 && rm /tmp/rustup.sh
ENV PATH=/root/.cargo/bin:$PATH
CMD ["./run.sh"]

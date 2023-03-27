FROM debian:buster-slim

WORKDIR /opt/k8s-faas-builder

ENV FUNCTION_NAME="lambda-function" \
    TARGET_HOST="12.0.0.1" \
    TARGET_PORT="22" \
    TARGET_USER="k8s-faas-builder" \
    TARGET_PATH="/opt/fass/lambda"

RUN apt-get update -y && \
    apt-get install -y curl ssh build-essential && \
    curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | bash -s -- --default-toolchain none -y && \
    /root/.cargo/bin/rustup toolchain install nightly --allow-downgrade --profile minimal --component clippy && \
    /root/.cargo/bin/rustup target add wasm32-wasi
VOLUME /opt/k8s-faas-builder/lambda
COPY . .
RUN chmod u+x ./build.sh
CMD ["./build.sh"]

FROM golang:1.17

WORKDIR /opt/k8s-faas-builder

ENV FUNCTION_NAME="lambda-function" \
    TARGET_HOST="12.0.0.1" \
    TARGET_PORT="22" \
    TARGET_USER="k8s-faas-builder" \
    TARGET_PATH="/opt/fass/lambda"

VOLUME /opt/k8s-faas-builder/lambda

COPY ./build.sh .

RUN chmod u+x ./build.sh

CMD ["./build.sh"]

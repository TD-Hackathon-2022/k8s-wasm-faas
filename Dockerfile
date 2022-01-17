FROM golang:1.17

WORKDIR /opt/k8s-faas-builder-controller

COPY ./out/k8s-faas-builder-controller-linux .

CMD ["k8s-faas-builder-controller"]

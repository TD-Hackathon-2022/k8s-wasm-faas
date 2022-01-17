FROM golang:1.17

WORKDIR /opt/k8s-faas-builder-controller

COPY k8s-faas-builder-controller .

CMD ["k8s-faas-builder-controller"]

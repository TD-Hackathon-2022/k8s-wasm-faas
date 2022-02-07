# k8s WASM FaaS Kubectl Plugin

## 使用

```text
wasm runtime is provided on k8s clusters, which may make it easier to run faas on k8s clusters

Usage:
  faas [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  list        list all faas function on k8s clusters
  run         run wasm lambda function by name
  submit      Store function in ConfigMap

Flags:
  -h, --help   help for faas

Use "faas [command] --help" for more information about a command.
```

## 编译

```shell
GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o out/kubectl-faas
```

然后将 **out/kubectl-faas** 放在 ***PATH*** 路径下
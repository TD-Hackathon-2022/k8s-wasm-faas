# k8s + wam = faas

## 概述

实现 WASM + WASI 来 k8s 的 runtime, 试图解决常规 FaaS 中冷启动延时的问题.

### 整体架构

整体分为两大部分:

1. runtime 实现在 [TD-Hackathon-2022/cri-wasm-runtime](https://github.com/TD-Hackathon-2022/cri-wasm-runtime)

2. FaaS 组件相关实现在 [TD-Hackathon-2022/k8s-wasm-faas](https://github.com/TD-Hackathon-2022/k8s-wasm-faas)

FaaS 组件分为三部分:

1. kubectl plugin 交互入口，可以提交本地的 FaaS 函数, 查看集群中已经存储的 FaaS 函数, 执行 FaaS 函数
2. builder controller 编译控制器, 监听存储在 configmaps 中的 FaaS 函数, 然后下发编译任务
3. builder 将用户编写的 FaaS 脚本编译为 WASM 可执行文件（目前只支持 Rust 语言的脚本编译)

### FaaS 组件

```text
./k8s-wasm-faas
├── README.md
├── builder            # 编译 rust faas script
├── builder-controller # 监听 configmap 中存储的 faas script, 并下发编译任务
└── plugin             # kubectl plugin
```
-- DUIZHANG UPDATE
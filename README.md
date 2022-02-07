# k8s + wam = faas

## 概述

实现 WASM + WASI 来 k8s 的 runtime, 试图解决常规 Faas 中冷启动延时的问题.

### 整体架构

整体分为两大部分: runtime + Faas 支撑组件 (cli plugin, builder controller, buider job)

runtime 实现在 [TD-Hackathon-2022/cri-wasm-runtime](https://github.com/TD-Hackathon-2022/cri-wasm-runtime)

Faas 组件在本仓库中 ([TD-Hackathon-2022/k8s-wasm-faas](https://github.com/TD-Hackathon-2022/k8s-wasm-faas))

### Faas 组件

```text
./k8s-wasm-faas
├── README.md
├── builder            # 负责编译 rust faas script
├── builder-controller # 监听 configmap 中存储的 faas script, 并下发编译任务
├── k8s-wasm-faas.iml
└── plugin
```

## 编译与运行



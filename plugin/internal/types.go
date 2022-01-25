package internal

var (
	faasLabels = map[string]string{
		"runtime": "wasm",
		"type":    "faas-wasm",
	}

	wasmPodLabel = map[string]string{
		"cri-runtime": "wasm",
	}
)

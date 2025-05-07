wasm:
	$(eval wasi_go_sdk := $(shell go list -mod=readonly -m -f '{{.Dir}}' github.com/rajatjindal/wasi-go-sdk))
	tinygo build -target=wasip2 --wit-package $(wasi_go_sdk)/wit --wit-world sdk -o main.wasm main.go
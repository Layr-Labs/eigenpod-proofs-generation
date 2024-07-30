all:
	cd core && $(MAKE) all
	go build

wasm:
	GOOS=js GOARCH=wasm go build -o proofs.wasm


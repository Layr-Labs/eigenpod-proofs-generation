TMP_DIR :=$(shell mktemp -d)

all: codegen

clone: 
	@echo "Cloning eigenlayer-contracts..."
	git clone --depth=1 git@github.com:Layr-Labs/eigenlayer-contracts.git $(TMP_DIR)
	
abi: clone
	@echo "Generating abi..."
	cd $(TMP_DIR) && make bindings
	./fixAbi.sh $(TMP_DIR)/out/EigenPod.sol/EigenPod.json

codegen: abi
	@echo "Generating golang..."
	abigen --abi $(TMP_DIR)/out/EigenPod.sol/EigenPod.json --pkg onchain --type EigenPod --out onchain/EigenPod.go

deps:
	go install github.com/ethereum/go-ethereum/cmd/abigen@latest
	brew install jq



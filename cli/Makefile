TMP_DIR :=/var/folders/r2/msprmz3n42q3qd9x3pc6gd0m0000gn/T/tmp.CQcLHh9hto
#$(shell mktemp -d)

clone: 
	@echo "Cloning eigenlayer-contracts..."
	git clone --depth=1 git@github.com:Layr-Labs/eigenlayer-contracts.git $(TMP_DIR)
	
abi: clone
	@echo "Generating abi..."
	cd $(TMP_DIR) && make bindings
	./fixAbi.sh $(TMP_DIR)/out/EigenPod.sol/EigenPod.json

codegen:
	@echo "Generating golang..."
	abigen --abi $(TMP_DIR)/out/EigenPod.sol/EigenPod.json --pkg onchain --type EigenPod --out onchain/EigenPod.go

all: codegen

deps:
	go install github.com/ethereum/go-ethereum/cmd/abigen@latest
	brew install jq



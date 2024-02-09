package onchain_tests

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"testing"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	beacon "github.com/Layr-Labs/eigenpod-proofs-generation/beacon"
	contractBeaconChainProofs "github.com/Layr-Labs/eigenpod-proofs-generation/bindings"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/assert"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	chainClient                    *eigenpodproofs.ChainClient
	ctx                            context.Context
	contractAddress                common.Address
	beaconChainProofs              *contractBeaconChainProofs.BeaconChainProofs
	oracleState                    deneb.BeaconState
	oracleBlockHeader              phase0.BeaconBlockHeader
	blockHeader                    phase0.BeaconBlockHeader
	blockHeaderIndex               uint64
	block                          deneb.BeaconBlock
	validatorIndex                 phase0.ValidatorIndex
	beaconBlockHeaderToVerifyIndex uint64
	executionPayload               deneb.ExecutionPayload
	epp                            *eigenpodproofs.EigenPodProofs
	executionPayloadFieldRoots     []phase0.Root
)

const GOERLI_CHAIN_ID = uint64(5)
const VALIDATOR_INDEX = uint64(61336)

func TestMain(m *testing.M) {
	// Setup
	log.Println("Setting up suite")
	setupSuite()

	// Run tests
	code := m.Run()

	// Teardown
	log.Println("Tearing down suite")
	teardownSuite()

	// Exit with test result code
	os.Exit(code)
}

func setupSuite() {
	RPC := "https://rpc.ankr.com/eth_goerli"
	PrivateKey := os.Getenv("PRIVATE_KEY")

	ethClient, err := ethclient.Dial(RPC)
	if err != nil {
		log.Panicf("failed to connect to the Ethereum client: %s", err)
	}

	chainClient, err = eigenpodproofs.NewChainClient(ethClient, PrivateKey)
	if err != nil {
		log.Panicf("failed to create chain client: %s", err)
	}
	ctx = context.Background()
	contractAddress = common.HexToAddress("0xd42a10709f0cc83855Af9B9fFeAa40dcE56D8fF6")
	beaconChainProofs, err = contractBeaconChainProofs.NewBeaconChainProofs(contractAddress, chainClient)
	if err != nil {
		log.Panicf("failed to create contract instance: %s", err)
	}

	log.Println("Setting up suite")
	stateFile := "../data/deneb_goerli_slot_7413760.json"
	oracleHeaderFile := "../data/deneb_goerli_block_header_7413760.json"
	headerFile := "../data/deneb_goerli_block_header_7426113.json"
	bodyFile := "../data/deneb_goerli_block_7426113.json"

	stateJSON, err := eigenpodproofs.ParseJSONFileDeneb(stateFile)
	if err != nil {
		fmt.Println("error with JSON parsing beacon state")
	}
	eigenpodproofs.ParseDenebBeaconStateFromJSON(*stateJSON, &oracleState)

	blockHeader, err = eigenpodproofs.ExtractBlockHeader(headerFile)
	if err != nil {
		fmt.Println("error with block header", err)
	}

	oracleBlockHeader, err = eigenpodproofs.ExtractBlockHeader(oracleHeaderFile)
	if err != nil {
		fmt.Println("error with oracle block header", err)
	}

	block, err = eigenpodproofs.ExtractBlock(bodyFile)
	if err != nil {
		fmt.Println("error with block body", err)
	}

	executionPayload = *block.Body.ExecutionPayload

	blockHeaderIndex = uint64(blockHeader.Slot) % beacon.SlotsPerHistoricalRoot

	epp, err = eigenpodproofs.NewEigenPodProofs(GOERLI_CHAIN_ID, 1000)
	if err != nil {
		fmt.Println("error in NewEigenPodProofs", err)
	}

	executionPayloadFieldRoots, _ = beacon.ComputeExecutionPayloadFieldRootsDeneb(block.Body.ExecutionPayload)
}

func teardownSuite() {

}

func TestValidatorContainersProofOnChain(t *testing.T) {

	versionedOracleState, err := beacon.CreateVersionedState(&oracleState)
	if err != nil {
		fmt.Println("error", err)
		return
	}

	verifyValidatorFieldsCallParams, err := epp.ProveValidatorContainers(&oracleBlockHeader, &versionedOracleState, []uint64{VALIDATOR_INDEX})
	if err != nil {
		fmt.Println("error", err)
	}

	validatorFieldsProof := verifyValidatorFieldsCallParams.ValidatorFieldsProofs[0].ToByteSlice()
	validatorIndex := new(big.Int).SetUint64(verifyValidatorFieldsCallParams.ValidatorIndices[0])
	versionedOracleStateRoot, err := versionedOracleState.Deneb.HashTreeRoot()
	var validatorFields [][32]byte

	for _, field := range verifyValidatorFieldsCallParams.ValidatorFields[0] {
		validatorFields = append(validatorFields, field)
	}

	err = beaconChainProofs.VerifyValidatorFields(
		&bind.CallOpts{},
		versionedOracleStateRoot,
		validatorFields,
		validatorFieldsProof,
		validatorIndex,
	)
	if err != nil {
		fmt.Println("error", err)
	}
	assert.Nil(t, err)
}

// func generateValidatorFieldsProofTransaction()

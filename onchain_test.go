package eigenpodproofs_test

import (
	"context"
	"log"
	"math/big"
	"os"
	"testing"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	chainClient     *eigenpodproofs.ChainClient
	ctx             context.Context
	contractAddress common.Address
)

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
	PrivateKey := "c5871389c9221e91d776f355c852f374156bf7799f3f63a361e12d0cb075a479"

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

}

func teardownSuite() {

}

func TestValidatorContainersProofOnChain(t *testing.T) {
	var transaction *types.Transaction

	txData := &types.DynamicFeeTx{
		ChainID:   big.NewInt(5),
		Nonce:     uint64(32),
		GasTipCap: big.NewInt(2e9),
		GasFeeCap: big.NewInt(100e9),
		Gas:       uint64(3000000),
		To:        &contractAddress, // The address of the contract
		Value:     big.NewInt(0),    // Value sent with the transaction (0 for a call)
		Data:      data,
	}

	chainClient.EstimateGasPriceAndLimitAndSendTx(ctx, transacton, "prove validator fields")
}

func generateValidatorFieldsProofTransaction()

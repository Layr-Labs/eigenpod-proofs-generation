package eigenpodproofs_test

import (
	"log"
	"os"
	"testing"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	chainClient *eigenpodproofs.ChainClient
	ctx 	   context.Context
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

	chainClient, err := eigenpodproofs.NewChainClient(ethClient, PrivateKey)
	if err != nil {
		log.Panicf("failed to create chain client: %s", err)
	}
	ctx = context.Background()

}

func teardownSuite() {

}

func TestValidatorContainersProofOnChain(t *testing.T) {
	

	chainClient.EstimateGasPriceAndLimitAndSendTx(ctx, 
}


func generateValidatorFieldsProofTransaction()

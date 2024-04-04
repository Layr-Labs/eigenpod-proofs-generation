package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"

	contractBeaconChainProofs "github.com/Layr-Labs/eigenpod-proofs-generation/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
)

type RestakeProofsResponse struct {
	VerifyWithdrawalCredentialsCallParams *eigenpodproofs.VerifyWithdrawalCredentialsCallParams `json:"verifyWithdrawalCredentialsCallParams"`
	VerifyAndProcessWithdrawalCallParams  *eigenpodproofs.VerifyAndProcessWithdrawalCallParams  `json:"verifyAndProcessWithdrawalCallParams"`
}

var DENEB_FORK_TIMESTAMP_HOLESKY = uint64(1707305664)

func main() {
	log.Println("Setting up suite")
	rpc := "https://ethereum-holesky-rpc.publicnode.com"
	privateKey := os.Getenv("PRIVATE_KEY")
	ethClient, err := ethclient.Dial(rpc)
	if err != nil {
		log.Panicf("failed to connect to the Ethereum client: %s", err)
	}

	chainClient, err := eigenpodproofs.NewChainClient(ethClient, privateKey)
	if err != nil {
		log.Panicf("failed to create chain client: %s", err)
	}

	//BeaconChainProofs.sol deployment: https://goerli.etherscan.io/address/0xd132dD701d3980bb5d66A21e2340f263765e4a19#code
	contractAddress := common.HexToAddress("0xd132dD701d3980bb5d66A21e2340f263765e4a19")
	beaconChainProofs, err := contractBeaconChainProofs.NewBeaconChainProofsTest(contractAddress, chainClient)
	if err != nil {
		log.Panicf("failed to create beacon chain proofs contract: %s", err)
	}

	data, err := ioutil.ReadFile("object.json")
	if err != nil {
		panic(err)
	}

	var restakeResponse RestakeProofsResponse
	err = json.Unmarshal(data, &restakeResponse)
	if err != nil {
		panic(err)
	}

	//withdrawal credential proof

	verifyValidatorFieldsCallParams := restakeResponse.VerifyWithdrawalCredentialsCallParams
	for i, _ := range verifyValidatorFieldsCallParams.ValidatorIndices {
		validatorFieldsProof := verifyValidatorFieldsCallParams.ValidatorFieldsProofs[i].ToByteSlice()
		validatorIndex := new(big.Int).SetUint64(verifyValidatorFieldsCallParams.ValidatorIndices[i])
		// oracleRoot, err := "INSERT HERE".ToBytes32()
		// if err != nil {
		// 	fmt.Println("error", err)
		// }

		// err = beaconChainProofs.VerifyStateRootAgainstLatestBlockRoot(
		// 	&bind.CallOpts{},
		// 	oracleRoot,
		// 	verifyValidatorFieldsCallParams.StateRootProof.BeaconStateRoot,
		// 	verifyValidatorFieldsCallParams.StateRootProof.StateRootProof.ToByteSlice(),
		// )
		// if err != nil {
		// 	fmt.Println("error", err)
		// }

		var validatorFields [][32]byte
		for _, field := range verifyValidatorFieldsCallParams.ValidatorFields[0] {
			validatorFields = append(validatorFields, field)
		}

		err = beaconChainProofs.VerifyValidatorFields(
			&bind.CallOpts{},
			verifyValidatorFieldsCallParams.StateRootProof.BeaconStateRoot,
			validatorFields,
			validatorFieldsProof,
			validatorIndex,
		)
		if err != nil {
			fmt.Println("error", err)
		}
	}
	////////////////////////

	verifyAndProcessWithdrawalCallParams := restakeResponse.VerifyAndProcessWithdrawalCallParams

	for i, _ := range verifyAndProcessWithdrawalCallParams.WithdrawalProofs {
		var withdrawalFields [][32]byte
		for _, field := range verifyAndProcessWithdrawalCallParams.WithdrawalFields[i] {
			withdrawalFields = append(withdrawalFields, field)
		}

		withdrawalProof := contractBeaconChainProofs.BeaconChainProofsContractWithdrawalProof{
			WithdrawalProof:                 verifyAndProcessWithdrawalCallParams.WithdrawalProofs[i].WithdrawalProof.ToByteSlice(),
			SlotProof:                       verifyAndProcessWithdrawalCallParams.WithdrawalProofs[i].SlotProof.ToByteSlice(),
			ExecutionPayloadProof:           verifyAndProcessWithdrawalCallParams.WithdrawalProofs[i].ExecutionPayloadProof.ToByteSlice(),
			TimestampProof:                  verifyAndProcessWithdrawalCallParams.WithdrawalProofs[i].TimestampProof.ToByteSlice(),
			HistoricalSummaryBlockRootProof: verifyAndProcessWithdrawalCallParams.WithdrawalProofs[i].HistoricalSummaryBlockRootProof.ToByteSlice(),
			BlockRootIndex:                  verifyAndProcessWithdrawalCallParams.WithdrawalProofs[i].BlockRootIndex,
			HistoricalSummaryIndex:          verifyAndProcessWithdrawalCallParams.WithdrawalProofs[i].HistoricalSummaryIndex,
			WithdrawalIndex:                 verifyAndProcessWithdrawalCallParams.WithdrawalProofs[i].WithdrawalIndex,
			BlockRoot:                       verifyAndProcessWithdrawalCallParams.WithdrawalProofs[i].BlockRoot,
			SlotRoot:                        verifyAndProcessWithdrawalCallParams.WithdrawalProofs[i].SlotRoot,
			TimestampRoot:                   verifyAndProcessWithdrawalCallParams.WithdrawalProofs[i].TimestampRoot,
			ExecutionPayloadRoot:            verifyAndProcessWithdrawalCallParams.WithdrawalProofs[i].ExecutionPayloadRoot,
		}

		fmt.Println("historicalSummaryndex ", verifyAndProcessWithdrawalCallParams.WithdrawalProofs[i].HistoricalSummaryIndex)
		fmt.Println("blockRootIndex ", verifyAndProcessWithdrawalCallParams.WithdrawalProofs[i].BlockRootIndex)
		fmt.Println("beacon state root: ", hex.EncodeToString(verifyAndProcessWithdrawalCallParams.StateRootProof.BeaconStateRoot[:]))

		err = beaconChainProofs.VerifyWithdrawal(
			&bind.CallOpts{},
			verifyAndProcessWithdrawalCallParams.StateRootProof.BeaconStateRoot,
			withdrawalFields,
			withdrawalProof,
			DENEB_FORK_TIMESTAMP_HOLESKY,
		)

		if err != nil {
			fmt.Println("error", err)
		}
	}
}

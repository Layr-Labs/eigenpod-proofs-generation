package main

import (
	"bytes"
	"crypto/sha256"
	"flag"
	"log"

	"context"

	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/onchain"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ValidatorWithIndex = struct {
	Validator *phase0.Validator
	Index     uint64
}

func main() {
	eigenpodAddress := flag.String("eigenpodAddress", "", "[required] The onchain address of your eigenpod contract (0x123123123123)")
	beacon := flag.String("beacon", "", "[required] URI to a functioning beacon node RPC (https://)")
	node := flag.String("node", "", "[required] URI to a functioning execution-layer RPC")
	out := flag.String("output", "", "Output path for the proof. (defaults to stdout)")
	command := flag.String("prove", "validators", "one of 'checkpoint' or 'validators'.\n\tIf checkpoint, produces a proof which can be submitted via EigenPod.VerifyCheckpointProofs().\n\tIf validators, generates a proof which can be submitted via EigenPod.VerifyWithdrawalCredentials().")
	help := flag.Bool("help", false, "Prints the help message and exits.")

	if !((*command == "validators") || (*command == "checkpoint")) {
		flag.Usage()
		log.Fatal("Invalid argument passed to --prove.")
	}

	flag.Parse()

	if help != nil && *help {
		flag.Usage()
		log.Fatal("Showing help.")
	}

	if *eigenpodAddress == "" || *beacon == "" || *node == "" {
		flag.Usage()
		log.Fatal("Must specify: --eigenpod, --beacon, and --node.")
	}

	ctx := context.Background()

	execute(ctx, *eigenpodAddress, *beacon, *node, *command, out)
}

func getBeaconClient(beaconUri string) (BeaconClient, error) {
	beaconClient, _, err := NewBeaconClient(beaconUri)
	return beaconClient, err
}

func lastCheckpointedForEigenpod(eigenpodAddress string, client *ethclient.Client) uint64 {
	eigenPod, err := onchain.NewEigenPod(common.HexToAddress(eigenpodAddress), client)
	PanicOnError("failed to locate eigenpod. is your address correct?", err)

	timestamp, err := eigenPod.CurrentCheckpointTimestamp(nil)
	PanicOnError("failed to locate eigenpod. Is your address correct?", err)

	return timestamp
}

// search through beacon state for validators whose withdrawal address is set to eigenpod.
func findAllValidatorsForEigenpod(eigenpodAddress string, beaconState *spec.VersionedBeaconState) []ValidatorWithIndex {
	allValidators, err := beaconState.Validators()
	PanicOnError("failed to fetch beacon state", err)

	eigenpodAddressBytes := common.FromHex(eigenpodAddress)

	var outputValidators []ValidatorWithIndex = []ValidatorWithIndex{}
	var i uint64 = 0
	maxValidators := uint64(len(allValidators))
	for i = 0; i < maxValidators; i++ {
		validator := allValidators[i]
		if validator == nil {
			continue
		}
		// we check that the last 20 bytes of expectedCredentials matches validatorCredentials.
		if bytes.Equal(
			eigenpodAddressBytes[:],
			validator.WithdrawalCredentials[12:], // first 12 bytes are not the pubKeyHash, see (https://github.com/Layr-Labs/eigenlayer-contracts/blob/d148952a2942a97a218a2ab70f9b9f1792796081/src/contracts/pods/EigenPod.sol#L663)
		) {
			outputValidators = append(outputValidators, ValidatorWithIndex{
				Validator: validator,
				Index:     i,
			})
		}
	}
	return outputValidators
}

func getOnchainValidatorInfo(client *ethclient.Client, eigenpodAddress string, allValidators []ValidatorWithIndex) []onchain.IEigenPodValidatorInfo {
	eigenPod, err := onchain.NewEigenPod(common.HexToAddress(eigenpodAddress), client)
	PanicOnError("failed to locate Eigenpod. Is your address correct?", err)

	var validatorInfo []onchain.IEigenPodValidatorInfo = []onchain.IEigenPodValidatorInfo{}

	// TODO: batch/multicall
	zeroes := [16]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	for i := 0; i < len(allValidators); i++ {
		// ssz requires values to be 32-byte aligned, which requires 16 bytes of 0's to be added
		// prior to hashing.
		pubKeyHash := sha256.Sum256(
			append(
				(allValidators[i]).Validator.PublicKey[:],
				zeroes[:]...,
			),
		)
		info, err := eigenPod.ValidatorPubkeyHashToInfo(nil, pubKeyHash)
		PanicOnError("failed to fetch validator eigeninfo.", err)
		validatorInfo = append(validatorInfo, info)
	}

	return validatorInfo
}

func getCurrentCheckpointBlockRoot(eigenpodAddress string, eth *ethclient.Client) (*[32]byte, error) {
	eigenPod, err := onchain.NewEigenPod(common.HexToAddress(eigenpodAddress), eth)
	PanicOnError("failed to locate Eigenpod. Is your address correct?", err)

	checkpoint, err := eigenPod.CurrentCheckpoint(nil)
	PanicOnError("failed to reach eigenpod.", err)

	return &checkpoint.BeaconBlockRoot, nil
}

func execute(ctx context.Context, eigenpodAddress, beacon_node_uri, node, command string, out *string) {
	eth, err := ethclient.Dial(node)
	PanicOnError("failed to reach eth --node.", err)

	chainId, err := eth.ChainID(ctx)
	PanicOnError("failed to fetch chain id", err)

	beaconClient, err := getBeaconClient(beacon_node_uri)
	PanicOnError("failed to reach beacon chain.", err)

	if command == "checkpoint" {
		RunCheckpointProof(ctx, eigenpodAddress, eth, chainId, beaconClient, out)
	} else {
		RunValidatorProof(ctx, eigenpodAddress, eth, chainId, beaconClient, out)
	}

}

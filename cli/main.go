package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"context"

	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/onchain"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
)

type ValidatorWithIndex = struct {
	Validator *phase0.Validator
	Index     uint64
}

func main() {
	eigenpod := flag.String("eigenpod", "", "[required] The onchain address of your eigenpod contract (0x123123123123)")
	beacon := flag.String("beacon", "", "[required] URI to a functioning beacon node RPC (https://)")
	node := flag.String("node", "", "[required] URI to a functioning execution-layer RPC")
	out := flag.String("output", "", "Output path for the proof. (defaults to stdout)")
	help := flag.Bool("help", false, "Prints the help message and exits.")

	flag.Parse()

	if help != nil && *help {
		flag.Usage()
		log.Fatal("Showing help.")
	}

	if *eigenpod == "" || *beacon == "" || *node == "" {
		flag.Usage()
		log.Fatal("Must specify: -eigenpod, -beacon, and -node.")
	}

	ctx := context.Background()

	execute(ctx, *eigenpod, *beacon, *node, out)
}

func getBeaconClient(beacon_uri string) (BeaconClient, error) {
	beaconClient, _, err := NewBeaconClient(beacon_uri)
	return beaconClient, err
}

func lastCheckpointedForEigenpod(eigenpodAddress string, client *ethclient.Client) uint64 {
	eigenPod, err := onchain.NewEigenPod(common.HexToAddress(eigenpodAddress), client)
	PanicOnError(err)

	timestamp, err := eigenPod.LastCheckpointTimestamp(nil)
	PanicOnError(err)
	return timestamp
}

func computeSlotImmediatelyPriorToTimestamp(timestampSeconds uint64, genesis time.Time) uint64 {
	return uint64(math.Floor(float64(timestampSeconds)-float64(genesis.Unix())) / 12)
}

// search through beacon state for validators whose withdrawal address is set to eigenpod.
func findAllValidatorsForEigenpod(eigenpodPubKey string, beaconState *spec.VersionedBeaconState) []ValidatorWithIndex {
	allValidators, err := beaconState.Validators()
	PanicOnError(err)

	expectedCredentials := sha256.Sum256([]byte(eigenpodPubKey))
	expectedCredentials[0] = 0 // the first byte of withdrawal credentials is set to 0.

	var outputValidators []ValidatorWithIndex = []ValidatorWithIndex{}
	var i uint64 = 0
	maxValidators := uint64(len(allValidators))
	for i = 0; i < maxValidators; i++ {
		validator := allValidators[i]
		if validator == nil {
			continue
		}
		if bytes.Equal(expectedCredentials[:], validator.WithdrawalCredentials) {
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
	PanicOnError(err)

	var validatorInfo []onchain.IEigenPodValidatorInfo = []onchain.IEigenPodValidatorInfo{}

	// TODO: batch/multicall
	for i := 0; i < len(allValidators); i++ {
		pubKeyHash := sha256.Sum256((allValidators[i]).Validator.PublicKey[:])
		info, err := eigenPod.ValidatorPubkeyHashToInfo(nil, pubKeyHash)
		PanicOnError(err)
		validatorInfo = append(validatorInfo, info)
	}

	return validatorInfo
}

// Stub for the execute function
func execute(ctx context.Context, eigenpodAddress, beacon_node_uri, node string, out *string) {
	eth, err := ethclient.Dial(node)
	if err != nil {
		fmt.Printf("ERROR: Invalid node - Failed to connect to `%s`.\n\n", node)
		PanicOnError(err)
	}
	beaconClient, err := getBeaconClient(beacon_node_uri)
	PanicOnError(err)

	lastCheckpoint := lastCheckpointedForEigenpod(eigenpodAddress, eth)

	// fetch the beacon state which corresponds to the slot immediately prior to this timestamp.
	genesis, err := beaconClient.GetChainGenesisTime(ctx)
	PanicOnError(err)

	slot := computeSlotImmediatelyPriorToTimestamp(lastCheckpoint, genesis)
	header, err := beaconClient.GetBeaconHeader(ctx, strconv.FormatUint(slot, 10))
	PanicOnError(err)

	beaconState, err := beaconClient.GetBeaconState(ctx, strconv.FormatUint(uint64(header.Header.Message.Slot), 10))
	PanicOnError(err)

	// filter through the beaconState's validators, and select only ones that have withdrawal address set to `eigenpod`.
	allValidatorsForEigenpod := findAllValidatorsForEigenpod(eigenpodAddress, beaconState)
	allValidatorInfo := getOnchainValidatorInfo(eth, eigenpodAddress, allValidatorsForEigenpod)

	// for each validator, request RPC information from the eigenpod (using the pubKeyHash), and;
	//			- we want all un-checkpointed, non-withdrawn validators that belong to this eigenpoint.
	//			- determine the validator's index.
	var checkpointValidatorIndices = []uint64{}
	for i := 0; i < len(allValidatorsForEigenpod); i++ {
		validator := allValidatorsForEigenpod[i]
		validatorInfo := allValidatorInfo[i]

		notCheckpointed := validatorInfo.LastCheckpointedAt != lastCheckpoint
		notWithdrawn := validatorInfo.Status != 2 // (TODO: does `abigen` generate a constant for this enum?)

		if notCheckpointed && notWithdrawn {
			checkpointValidatorIndices = append(checkpointValidatorIndices, validator.Index)
		}
	}

	proofs, err := eigenpodproofs.NewEigenPodProofs(1 /* ETH */, 300 /* oracleStateCacheExpirySeconds - 5min */)
	if err != nil {
		panic(err)
	}

	res, err := proofs.ProveCheckpointProofs(header.Header.Message, beaconState, checkpointValidatorIndices)
	PanicOnError(err)

	jsonString, err := json.Marshal(res)
	PanicOnError(err)

	if out != nil {
		os.WriteFile(*out, jsonString, os.ModePerm)
		PanicOnError(err)
	} else {
		fmt.Print(jsonString)
	}
}

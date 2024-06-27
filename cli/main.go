package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"

	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/onchain"
	client "github.com/attestantio/go-eth2-client"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/layr-labs/eigenlayer-backend/common/beacon"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	eigenpod := flag.String("eigenpod", "", "The onchain address of your eigenpod contract")
	beacon := flag.String("beacon", "", "URI to a functioning beacon node RPC")
	node := flag.String("node", "", "URI to a functioning execution-layer RPC")
	out := flag.String("output", "", "Output path for the proof. (defaults to stdout)")
	help := flag.Bool("help", false, "Prints the help message and exits.")

	flag.StringVar(eigenpod, "e", "", "The onchain address of your eigenpod contract (shorthand)")
	flag.StringVar(beacon, "b", "", "URI to a functioning beacon node RPC (shorthand)")
	flag.StringVar(node, "n", "", "URI to a functioning execution-layer RPC (shorthand)")
	flag.StringVar(node, "o", "", "Output path for the proof. (defaults to stdout)")

	flag.Parse()

	if help != nil && *help {
		// TODO: help.
		flag.Usage()
		log.Fatal("Showing help.")
	}

	if *eigenpod == "" || *beacon == "" || *node == "" {
		flag.Usage()
		log.Fatal("Must specify: -eigenpod, -beacon, and -node.")
	}

	execute(*eigenpod, *beacon, *node, out)
}

func getBeaconClient(beacon_uri string) (*client.Service, error) {
	return beacon.NewBeaconClient(beacon_uri)
}

func lastCheckpointedForEigenpod(eigenpod string) uint64 {
	panic("unimplemented")
}

func computeSlotImmediatelyPriorToTimestamp(timestampSeconds uint64) uint64 {
	var genesisTimestampSeconds uint64 = 0 // TODO: get time for genesis block.
	return uint64(math.Floor(float64(timestampSeconds)-float64(genesisTimestampSeconds)) / 12)
}

func findAllValidatorsForEigenpod(eigenpod string, beaconState any) {
	// TODO: search through beacon state for validators whose withdrawal address is set to eigenpod.
	panic("unimplemented")
}

func batchGetValidatorInfo(client *ethclient.Client, eigenpodAddress string, allValidators any) []onchain.IEigenPodValidatorInfo {
	eigenPod, err := onchain.NewEigenPod(common.HexToAddress(eigenpodAddress), nil)
	panicOnError(err)

	var validatorInfo []onchain.IEigenPodValidatorInfo = []onchain.IEigenPodValidatorInfo{}

	for i := 0; i < len(allValidators); i++ {
		pubKey := allValidators[i].pubKey
		info, err := eigenPod.ValidatorPubkeyHashToInfo(pubKey)
		panicOnError(err)
		validatorInfo = append(validatorInfo, info)
	}

	return validatorInfo
}

func getBeaconState(beacon client.Service, slot uint64) {

}

// Stub for the execute function
func execute(eigenpod, beacon_node_uri, node string, out *string) {
	eth, err := ethclient.Dial(node)
	panicOnError(err)

	beaconClient, err := getBeaconClient(beacon_node_uri)

	// TODO: get last checkpoint timestamp from RPC.
	lastCheckpoint := lastCheckpointedForEigenpod(eigenpod)

	// TODO: fetch the beacon state which corresponds to the slot immediately prior to this timestamp.
	slot := computeSlotImmediatelyPriorToTimestamp(lastCheckpoint)
	beaconState, beaconHeader := getBeaconState(beaconClient, slot)

	// TODO: filter through the beaconState's validators, and select only ones that have `eigenpod` set to the validator address.
	allValidatorsForEigenpod := findAllValidatorsForEigenpod(eigenpod, beaconState)
	allValidatorInfo := batchGetValidatorInfo(eth, allValidatorsForEigenpod)

	// for each validator, request RPC information from the eigenpod (using the pubKeyHash), and;
	//			- we want all un-checkpointed, non-withdrawn validators that belong to this eigenpoint.
	//			- determine the validator's index.
	var checkpointValidatorIndices = []uint64{}
	for i := 0; i < len(allValidatorsForEigenpod); i++ {
		validator := allValidatorsForEigenpod[i]
		validatorInfo := allValidatorInfo[i]

		notCheckpointed := validatorInfo.MostRecentBalanceUpdateTimestamp != lastCheckpoint // TODO: determine from validatorInfo
		notWithdrawn := validatorInfo.Status != 2                                           // (TODO: does `abigen` generate a constant for this enum?)

		if notCheckpointed && notWithdrawn {
			checkpointValidatorIndices = append(checkpointValidatorIndices, validator.index)
		}
	}

	proofs, err := eigenpodproofs.NewEigenPodProofs(1 /* ETH */, 300 /* oracleStateCacheExpirySeconds - 5min */)
	if err != nil {
		panic(err)
	}

	res, err := proofs.ProveCheckpointProofs(beaconHeader, beaconState, checkpointValidatorIndices)

	jsonString, err := json.Marshal(res)
	panicOnError(err)

	if out != nil {
		ioutil.WriteFile(*out, jsonString, os.ModePerm)
		panicOnError(err)
	} else {
		fmt.Print(jsonString)
	}
}

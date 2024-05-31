package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"

	beacon "github.com/Layr-Labs/eigenpod-proofs-generation/beacon"
	"github.com/Layr-Labs/eigenpod-proofs-generation/common"
)

var VALIDATOR_INDEX uint64 = 61068 //this is the index of a validator that has a full withdrawal

// this needs to be hand crafted. If you want the root of the header at the slot x,
// then look for entry in (x)%SLOTS_PER_HISTORICAL_ROOT in the block_roots.

var BEACON_BLOCK_HEADER_TO_VERIFY_INDEX uint64 = 2262

const FIRST_CAPELLA_SLOT_GOERLI = uint64(5193728)
const FIRST_CAPELLA_SLOT_MAINNET = uint64(6209536)

var GOERLI_CHAIN_ID uint64 = 5

func main() {
	// Example optional flag
	var newBalance int
	flag.IntVar(&newBalance, "newBalance", -1, "new balance")

	// Parse the flags
	flag.Parse()

	// Using flag.Args() to get positional arguments
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("No command provided")
		return
	}

	switch args[0] {
	//use this for withdrawal credentials and balance update proofs
	case "ValidatorFieldsProof":

		index, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		changeBalance, err := strconv.ParseBool(args[2])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		oracleStateFile := args[3]
		stateFile := args[4]
		outputFile := args[5]

		GenerateValidatorFieldsProof(oracleStateFile, stateFile, index, changeBalance, uint64(newBalance), outputFile)
	default:
		fmt.Println("Unknown command:", args[0])
		os.Exit(0)
	}
}

func GenerateValidatorFieldsProof(oracleBlockHeaderFile string, stateFile string, index uint64, changeBalance bool, newBalance uint64, output string) error {
	var state deneb.BeaconState
	var oracleBeaconBlockHeader phase0.BeaconBlockHeader
	SetupValidatorProof(oracleBlockHeaderFile, stateFile, index, changeBalance, newBalance, 0, &state, &oracleBeaconBlockHeader)

	validatorIndex := phase0.ValidatorIndex(index)

	beaconStateRoot, _ := state.HashTreeRoot()

	latestBlockHeaderRoot, err := oracleBeaconBlockHeader.HashTreeRoot()
	if err != nil {
		fmt.Println("Error with HashTreeRoot of latestBlockHeader", err)
		return err
	}

	epp, err := eigenpodproofs.NewEigenPodProofs(GOERLI_CHAIN_ID, 1000)
	if err != nil {
		fmt.Println("Error creating EPP object", err)
		return err
	}

	var versionedState spec.VersionedBeaconState
	versionedState.Deneb = &state

	stateRootProof, validatorFieldsProof, _ := ProveValidatorFields(epp, &oracleBeaconBlockHeader, &versionedState, uint64(validatorIndex))

	proofs := WithdrawalCredentialProofs{
		ValidatorIndex:                         uint64(validatorIndex),
		BeaconStateRoot:                        "0x" + hex.EncodeToString(beaconStateRoot[:]),
		LatestBlockHeaderRoot:                  "0x" + hex.EncodeToString(latestBlockHeaderRoot[:]),
		WithdrawalCredentialProof:              ConvertBytesToStrings(validatorFieldsProof),
		ValidatorFields:                        GetValidatorFields(state.Validators[validatorIndex]),
		StateRootAgainstLatestBlockHeaderProof: ConvertBytesToStrings(stateRootProof.Proof),
	}

	proofData, err := json.Marshal(proofs)
	if err != nil {
		fmt.Println("error")
		return err
	}

	_ = ioutil.WriteFile(output, proofData, 0644)

	return nil
}

func ProveValidatorFields(epp *eigenpodproofs.EigenPodProofs, oracleBlockHeader *phase0.BeaconBlockHeader, oracleBeaconState *spec.VersionedBeaconState, validatorIndex uint64) (*eigenpodproofs.StateRootProof, common.Proof, error) {
	oracleBeaconStateSlot, err := oracleBeaconState.Slot()
	if err != nil {
		return nil, nil, err
	}
	oracleBeaconStateValidators, err := oracleBeaconState.Validators()
	if err != nil {
		return nil, nil, err
	}

	stateRootProof := &eigenpodproofs.StateRootProof{}
	// Get beacon state top level roots
	beaconStateTopLevelRoots, err := epp.ComputeBeaconStateTopLevelRoots(oracleBeaconState)
	if err != nil {
		return nil, nil, err
	}

	// Get beacon state root. TODO: Combine this cheaply with compute beacon state top level roots
	stateRootProof.BeaconStateRoot = oracleBlockHeader.StateRoot
	if err != nil {
		return nil, nil, err
	}

	stateRootProof.Proof, err = beacon.ProveStateRootAgainstBlockHeader(oracleBlockHeader)

	if err != nil {
		return nil, nil, err
	}

	validatorFieldsProof, err := epp.ProveValidatorAgainstBeaconState(beaconStateTopLevelRoots, oracleBeaconStateSlot, oracleBeaconStateValidators, validatorIndex)

	if err != nil {
		return nil, nil, err
	}

	return stateRootProof, validatorFieldsProof, nil
}

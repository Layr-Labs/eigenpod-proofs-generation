package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	ssz "github.com/ferranbt/fastssz"

	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

type WithdrawalProofs struct {
	Slot                                   uint64   `json:"slot"`
	ValidatorIndex                         uint64   `json:"validatorIndex"`
	HistoricalSummaryIndex                 uint64   `json:"historicalSummaryIndex"`
	WithdrawalIndex                        uint64   `json:"withdrawalIndex"`
	BlockHeaderRootIndex                   uint64   `json:"blockHeaderRootIndex"`
	BeaconStateRoot                        string   `json:"beaconStateRoot"`
	SlotRoot                               string   `json:"slotRoot"`
	TimestampRoot                          string   `json:"timestampRoot"`
	BlockHeaderRoot                        string   `json:"blockHeaderRoot"`
	BlockBodyRoot                          string   `json:"blockBodyRoot"`
	ExecutionPayloadRoot                   string   `json:"executionPayloadRoot"`
	LatestBlockHeaderRoot                  string   `json:"latestBlockHeaderRoot"`
	SlotProof                              []string `json:"SlotProof"`
	WithdrawalProof                        []string `json:"WithdrawalProof"`
	ValidatorProof                         []string `json:"ValidatorProof"`
	TimestampProof                         []string `json:"TimestampProof"`
	ExecutionPayloadProof                  []string `json:"ExecutionPayloadProof"`
	ValidatorFields                        []string `json:"ValidatorFields"`
	WithdrawalFields                       []string `json:"WithdrawalFields"`
	StateRootAgainstLatestBlockHeaderProof []string `json:"StateRootAgainstLatestBlockHeaderProof"`
	HistoricalSummaryProof                 []string `json:"HistoricalSummaryProof"`
}

type WithdrawalCredentialProofs struct {
	ValidatorIndex                         uint64   `json:"validatorIndex"`
	BeaconStateRoot                        string   `json:"beaconStateRoot"`
	LatestBlockHeaderRoot                  string   `json:"latestBlockHeaderRoot"`
	WithdrawalCredentialProof              []string `json:"WithdrawalCredentialProof"`
	ValidatorFields                        []string `json:"ValidatorFields"`
	StateRootAgainstLatestBlockHeaderProof []string `json:"StateRootAgainstLatestBlockHeaderProof"`
}

type BalanceUpdateProofs struct {
	ValidatorIndex                         uint64   `json:"validatorIndex"`
	BeaconStateRoot                        string   `json:"beaconStateRoot"`
	SlotRoot                               string   `json:"slotRoot"`
	BalanceRoot                            string   `json:"balanceRoot"`
	LatestBlockHeaderRoot                  string   `json:"latestBlockHeaderRoot"`
	ValidatorBalanceProof                  []string `json:"ValidatorBalanceProof"`
	ValidatorFields                        []string `json:"ValidatorFields"`
	StateRootAgainstLatestBlockHeaderProof []string `json:"StateRootAgainstLatestBlockHeaderProof"`
	WithdrawalCredentialProof              []string `json:"WithdrawalCredentialProof"`
}

type InputDataBlockHeader struct {
	Data struct {
		Header struct {
			Message phase0.BeaconBlockHeader `json:"message"`
		} `json:"header"`
	} `json:"data"`
}

type InputDataBlock struct {
	Version string `json:"version"`
	Data    struct {
		Message   deneb.BeaconBlock `json:"message"`
		Signature string            `json:"signature"`
	} `json:"data"`
	Execution_optimistic bool `json:"execution_optimistic"`
	Finalized            bool `json:"finalized"`
}

type InputDataBlockCapella struct {
	Version string `json:"version"`
	Data    struct {
		Message   capella.BeaconBlock `json:"message"`
		Signature string              `json:"signature"`
	} `json:"data"`
	Execution_optimistic bool `json:"execution_optimistic"`
	Finalized            bool `json:"finalized"`
}

func SetupValidatorProof(oracleBlockHeaderFile string, stateFile string, validatorIndex uint64, changeBalance bool, newBalance uint64, incrementSlot uint64, state *deneb.BeaconState, oracleBlockHeader *phase0.BeaconBlockHeader) {
	//filename1 := "data/slot_58000/oracle_capella_beacon_state_58100.ssz" //this is the file for the repointed validator (either 61336 or 61068)
	//filename1 := "data/slot_209635/oracle_capella_beacon_state_209635.ssz" //this is the file for the slashed validator 61511

	stateJSON, err := eigenpodproofs.ParseDenebStateJSONFile(stateFile)
	if err != nil {
		fmt.Println("error with JSON parsing")
	}
	eigenpodproofs.ParseDenebBeaconStateFromJSON(*stateJSON, state)

	*oracleBlockHeader, err = ExtractBlockHeader(oracleBlockHeaderFile)
	if err != nil {
		fmt.Println("read error with header file")
	}

	//setting the withdrawal credentials of Validator
	byteArray, _ := hex.DecodeString("01000000000000000000000049c486e3f4303bc11c02f952fe5b08d0ab22d443")
	state.Validators[validatorIndex].WithdrawalCredentials = byteArray

	if incrementSlot > 0 {
		state.Slot = state.Slot + phase0.Slot(incrementSlot)
	}

	// 61336 is the validator we prove withdrawals for.  So we need a with/cred proof
	// that actually has some balance in it.  So we artificially set the balance
	if changeBalance {
		// state.Balances[validatorIndex] = 32000115173
		// state.Validators[validatorIndex].EffectiveBalance = 32000115173
		fmt.Println("new balance", newBalance)
		state.Balances[validatorIndex] = phase0.Gwei(newBalance)
		state.Validators[validatorIndex].EffectiveBalance = phase0.Gwei(newBalance)
	}

	newStateRoot, err := state.HashTreeRoot()
	if err != nil {
		fmt.Println("error with HashTreeRoot", err)
	}
	//Now that we've made these changes to "state", we need to update the oracleState.LatestBlockHeader.StateRoot
	oracleBlockHeader.StateRoot = newStateRoot
}

func ConvertBytesToStrings(b [][32]byte) []string {
	var s []string
	for _, v := range b {
		s = append(s, "0x"+hex.EncodeToString(v[:]))
	}
	return s
}

func GetValidatorFields(v *phase0.Validator) []string {
	var validatorFields []string
	hh := ssz.NewHasher()

	hh.PutBytes(v.PublicKey[:])
	validatorFields = append(validatorFields, "0x"+hex.EncodeToString(hh.Hash()))
	hh.Reset()

	hh.PutBytes(v.WithdrawalCredentials)
	validatorFields = append(validatorFields, "0x"+hex.EncodeToString(hh.Hash()))
	hh.Reset()

	hh.PutUint64(uint64(v.EffectiveBalance))
	validatorFields = append(validatorFields, "0x"+hex.EncodeToString(hh.Hash()))
	hh.Reset()

	hh.PutBool(v.Slashed)
	validatorFields = append(validatorFields, "0x"+hex.EncodeToString(hh.Hash()))
	hh.Reset()

	hh.PutUint64(uint64(v.ActivationEligibilityEpoch))
	validatorFields = append(validatorFields, "0x"+hex.EncodeToString(hh.Hash()))
	hh.Reset()

	hh.PutUint64(uint64(v.ActivationEpoch))
	validatorFields = append(validatorFields, "0x"+hex.EncodeToString(hh.Hash()))
	hh.Reset()

	hh.PutUint64(uint64(v.ExitEpoch))
	validatorFields = append(validatorFields, "0x"+hex.EncodeToString(hh.Hash()))
	hh.Reset()

	hh.PutUint64(uint64(v.WithdrawableEpoch))
	validatorFields = append(validatorFields, "0x"+hex.EncodeToString(hh.Hash()))
	hh.Reset()

	return validatorFields
}

func ExtractBlockHeader(blockHeaderFile string) (phase0.BeaconBlockHeader, error) {

	fileBytes, err := os.ReadFile(blockHeaderFile)
	if err != nil {
		return phase0.BeaconBlockHeader{}, err
	}
	// Decode JSON
	var inputData InputDataBlockHeader
	if err := json.Unmarshal(fileBytes, &inputData); err != nil {
		return phase0.BeaconBlockHeader{}, err
	}

	return inputData.Data.Header.Message, nil
}

func ExtractBlock(blockHeaderFile string) (deneb.BeaconBlock, error) {
	fileBytes, err := os.ReadFile(blockHeaderFile)
	if err != nil {
		return deneb.BeaconBlock{}, err
	}

	// Decode JSON
	var data InputDataBlock
	if err := json.Unmarshal(fileBytes, &data); err != nil {
		return deneb.BeaconBlock{}, err
	}

	// Extract block body
	return data.Data.Message, nil
}

func ExtractBlockCapella(blockHeaderFile string) (capella.BeaconBlock, error) {
	fileBytes, err := os.ReadFile(blockHeaderFile)
	if err != nil {
		return capella.BeaconBlock{}, err
	}

	// Decode JSON
	var data InputDataBlockCapella
	if err := json.Unmarshal(fileBytes, &data); err != nil {
		return capella.BeaconBlock{}, err
	}

	// Extract block body
	return data.Data.Message, nil
}

func GetWithdrawalFields(w *capella.Withdrawal) []string {
	var withdrawalFields []string
	hh := ssz.NewHasher()

	hh.PutUint64(uint64(w.Index))
	withdrawalFields = append(withdrawalFields, "0x"+hex.EncodeToString(hh.Hash()))
	hh.Reset()

	hh.PutUint64(uint64(w.ValidatorIndex))
	withdrawalFields = append(withdrawalFields, "0x"+hex.EncodeToString(hh.Hash()))
	hh.Reset()

	hh.PutBytes(w.Address[:])
	withdrawalFields = append(withdrawalFields, "0x"+hex.EncodeToString(hh.Hash()))
	hh.Reset()

	hh.PutUint64(uint64(w.Amount))
	withdrawalFields = append(withdrawalFields, "0x"+hex.EncodeToString(hh.Hash()))
	hh.Reset()

	return withdrawalFields
}

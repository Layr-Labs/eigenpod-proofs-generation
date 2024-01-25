package beacon

import (
	"math/big"

	"github.com/Layr-Labs/eigenpod-proofs-generation/common"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/rs/zerolog/log"
)

func ComputeValidatorTreeLeaves(validators []*phase0.Validator) ([]phase0.Root, error) {
	validatorNodeList := make([]phase0.Root, len(validators))
	for i := 0; i < len(validators); i++ {
		validatorRoot, err := validators[i].HashTreeRoot()
		if err != nil {
			return nil, err
		}
		validatorNodeList[i] = phase0.Root(validatorRoot)
	}

	return validatorNodeList, nil
}

func ProveValidatorBalanceAgainstValidatorBalanceList(balances []phase0.Gwei, validatorIndex uint64) (common.Proof, error) {

	buf := []byte{}

	for i := 0; i < len(balances); i++ {
		buf = ssz.MarshalUint64(buf, uint64(balances[i]))
	}
	//pad the buffer with 0s to make it a multiple of 32 bytes
	if rest := len(buf) % 32; rest != 0 {
		buf = append(buf, zeroBytes[:32-rest]...)
	}

	// //now we divide the buffer into 32 byte chunks
	numLeaves := len(buf) / 32
	balanceRootList := make([]phase0.Root, numLeaves)
	for i := 0; i < numLeaves; i++ {
		copy(balanceRootList[i][:], buf[i*32:(i+1)*32])
	}

	//refer to beaconstate_ssz.go in go-eth2-client
	numLayers := uint64(common.GetDepth(ssz.CalculateLimit(1099511627776, uint64(len(balances)), 8)))
	validatorBalanceIndex := uint64(validatorIndex / 4)
	proof, err := common.GetProof(balanceRootList, validatorBalanceIndex, numLayers)

	if err != nil {
		log.Debug().AnErr("error", err).Msg("error getting proof")
		return nil, err
	}

	//append the length of the balance array to the proof
	//convert big endian to little endian
	balanceListLenLE := common.BigToLittleEndian(big.NewInt(int64(len(balances))))

	proof = append(proof, balanceListLenLE)
	return proof, nil
}

func GetBalanceRoots(balances []phase0.Gwei) ([]phase0.Root, error) {
	buf := []byte{}

	for i := 0; i < len(balances); i++ {
		buf = ssz.MarshalUint64(buf, uint64(balances[i]))
	}

	// //now we divide the buffer into 32 byte chunks
	numLeaves := len(buf) / 32
	balanceRootList := make([]phase0.Root, numLeaves)
	for i := 0; i < numLeaves; i++ {
		copy(balanceRootList[i][:], buf[i*32:(i+1)*32])
	}
	return balanceRootList, nil
}

func ProveBlockRootAgainstBeaconStateViaHistoricalSummaries(beaconStateTopLevelRoots *BeaconStateTopLevelRoots, historicalSummaries []*capella.HistoricalSummary, historicalBlockRoots []phase0.Root, historicalSummaryIndex uint64, blockRootIndex uint64) ([][32]byte, error) {
	// prove the historical summaries against the beacon state
	historicalSummariesListAgainstBeaconState, err := ProveBeaconTopLevelRootAgainstBeaconState(beaconStateTopLevelRoots, HistoricalSummaryListIndex)
	if err != nil {
		return nil, err
	}

	// prove the historical summary against the historical summaries list
	historicalSummaryAgainstHistoricalSummariesListProof, err := ProveHistoricalSummaryAgainstHistoricalSummariesList(historicalSummaries, historicalSummaryIndex)
	if err != nil {
		return nil, err
	}

	// prove the block roots against the historical summary
	historicalSummary := historicalSummaries[historicalSummaryIndex]
	blockRootsListAgainstHistoricalSummaryProof, err := ProveBlockRootListAgainstHistoricalSummary(historicalSummary)
	if err != nil {
		return nil, err
	}

	// historical block roots are incklude the really old block root you wanna prove
	beaconBlockHeaderRootsProof, err := common.GetProof(historicalBlockRoots, blockRootIndex, BlockRootsMerkleSubtreeNumLayers)
	if err != nil {
		return nil, err
	}

	blockHeaderProof := append(beaconBlockHeaderRootsProof, blockRootsListAgainstHistoricalSummaryProof...)
	blockHeaderProof = append(blockHeaderProof, historicalSummaryAgainstHistoricalSummariesListProof...)
	blockHeaderProof = append(blockHeaderProof, historicalSummariesListAgainstBeaconState...)

	return blockHeaderProof, nil
}

func ProveHistoricalSummaryAgainstHistoricalSummariesList(historicalSummaries []*capella.HistoricalSummary, historicalSummaryIndex uint64) (common.Proof, error) {
	historicalSummaryNodeList := make([]phase0.Root, len(historicalSummaries))
	for i := 0; i < len(historicalSummaries); i++ {
		historicalSummaryRoot, err := historicalSummaries[i].HashTreeRoot()
		if err != nil {
			return nil, err
		}
		historicalSummaryNodeList[i] = phase0.Root(historicalSummaryRoot)
	}

	proof, err := common.GetProof(historicalSummaryNodeList, historicalSummaryIndex, HistoricalSummaryListMerkleSubtreeNumLayers)

	if err != nil {
		return nil, err
	}
	//append the length of the validator array to the proof
	//convert big endian to little endian
	historicalSummaryListLenLE := common.BigToLittleEndian(big.NewInt(int64(len(historicalSummaries))))

	proof = append(proof, historicalSummaryListLenLE)
	return proof, nil
}

func ProveBlockRootListAgainstHistoricalSummary(historicalSummary *capella.HistoricalSummary) (common.Proof, error) {
	//historical summary container is a struct with two fields - block_summary_roots and state_summary_roots.  We want to prove
	// the block roots field so the proof is only 1 layer deep and is just the state summary root
	proof := make(common.Proof, 1)
	proof[0] = historicalSummary.StateSummaryRoot

	return proof, nil
}

func ProveBlockRootAgainstBlockRootsList(blockRoots []phase0.Root, blockRootIndex uint64) (common.Proof, error) {
	proof, err := common.GetProof(blockRoots, blockRootIndex, BlockRootsMerkleSubtreeNumLayers)
	if err != nil {
		return nil, err
	}
	return proof, nil
}

func ProveWithdrawalAgainstExecutionPayload(
	executionPayloadFieldRoots []phase0.Root,
	withdrawals []*capella.Withdrawal,
	withdrawalIndex uint8,
) ([][32]byte, error) {
	// prove withdrawal list against the execution payload
	withdrawalListAgainstExecutionPayloadProof, err := ProveWithdrawalListAgainstExecutionPayload(executionPayloadFieldRoots)
	if err != nil {
		return nil, err
	}

	// prove the withdrawal against the withdrawal list
	withdrawalAgainstWithdrawalListProof, err := ProveWithdrawalAgainstWithdrawalList(withdrawals, withdrawalIndex)
	if err != nil {
		return nil, err
	}

	//NOTE: Ensure that these proofs are being appended in the right order
	fullWithdrawalProof := append(withdrawalAgainstWithdrawalListProof, withdrawalListAgainstExecutionPayloadProof...)
	return fullWithdrawalProof, nil
}

func ProveWithdrawalListAgainstExecutionPayload(executionPayloadFieldRoots []phase0.Root) (common.Proof, error) {
	return common.GetProof(executionPayloadFieldRoots, WithdrawalsIndex, common.CeilLog2(len(executionPayloadFieldRoots)))
}

func ProveWithdrawalAgainstWithdrawalList(withdrawals []*capella.Withdrawal, withdrawalIndex uint8) (common.Proof, error) {
	withdrawalNodeList := make([]phase0.Root, len(withdrawals))
	for i := 0; i < len(withdrawals); i++ {
		withdrawalRoot, err := withdrawals[i].HashTreeRoot()
		if err != nil {
			return nil, err
		}
		withdrawalNodeList[i] = phase0.Root(withdrawalRoot)
	}

	proof, err := common.GetProof(withdrawalNodeList, uint64(withdrawalIndex), WithdrawalListMerkleSubtreeNumLayers)
	if err != nil {
		return nil, err
	}
	//append the length of the withdrawals array to the proof
	//convert big endian to little endian
	withdrawalListLenLE := common.BigToLittleEndian(big.NewInt(int64(len(withdrawals))))

	proof = append(proof, withdrawalListLenLE)
	return proof, nil
}

func ProveTimestampAgainstExecutionPayload(executionPayloadFieldRoots []phase0.Root) (common.Proof, error) {
	return common.GetProof(executionPayloadFieldRoots, TimestampIndex, common.CeilLog2(len(executionPayloadFieldRoots)))
}

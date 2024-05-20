package txsubmitter

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Layr-Labs/eigenpod-proofs-generation/beacon"
	commonutils "github.com/Layr-Labs/eigenpod-proofs-generation/common_utils"
	contractEigenPod "github.com/Layr-Labs/eigensdk-go/contracts/bindings/EigenPod"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

type WithdrawalProofConfig struct {
	EigenPodAddress common.Address `json:"EIGENPOD_ADDRESS,required"`

	BeaconStateFiles struct {
		OracleStateFile       string `json:"ORACLE_STATE_FILE,required"`
		OracleBlockHeaderFile string `json:"ORACLE_BLOCK_HEADER_FILE,required"`
	}

	WithdrawalDetails struct {
		ValidatorIndices            []uint64 `json:"VALIDATOR_INDICES,required"`
		WithdrawalBlockHeaderFiles  []string `json:"WITHDRAWAL_BLOCK_HEADER_FILES,required"`
		WithdrawalBlockBodyFiles    []string `json:"WITHDRAWAL_BLOCK_BODY_FILES,required"`
		HistoricalSummaryStateFiles []string `json:"HISTORICAL_SUMMARY_STATE_FILES,required"`
	}
}

func (u *EigenPodProofTxSubmitter) GenerateVerifyAndProcessWithdrawalsTx(
	eigenpod common.Address,
	versionedOracleState *spec.VersionedBeaconState,
	oracleBeaconBlockHeader *phase0.BeaconBlockHeader,
	historicalSummaryStateBlockRoots [][]phase0.Root,
	withdrawalBlocks []*spec.VersionedSignedBeaconBlock,
	validatorIndices []uint64,
) (*types.Transaction, error) {
	withdrawalsProof, err := u.eigenPodProofs.ProveWithdrawals(
		oracleBeaconBlockHeader,
		versionedOracleState,
		historicalSummaryStateBlockRoots,
		withdrawalBlocks,
		validatorIndices)
	if err != nil {
		return nil, err
	}
	withdrawalProofs := make([]contractEigenPod.BeaconChainProofsWithdrawalProof, len(withdrawalsProof.WithdrawalProofs))
	for i, v := range withdrawalsProof.WithdrawalProofs {
		withdrawalProofs[i] = contractEigenPod.BeaconChainProofsWithdrawalProof{
			WithdrawalProof:                 v.WithdrawalProof.ToByteSlice(),
			SlotProof:                       v.SlotProof.ToByteSlice(),
			ExecutionPayloadProof:           v.ExecutionPayloadProof.ToByteSlice(),
			TimestampProof:                  v.TimestampProof.ToByteSlice(),
			HistoricalSummaryBlockRootProof: v.HistoricalSummaryBlockRootProof.ToByteSlice(),
			BlockRootIndex:                  v.BlockRootIndex,
			HistoricalSummaryIndex:          v.HistoricalSummaryIndex,
			WithdrawalIndex:                 v.WithdrawalIndex,
			BlockRoot:                       v.BlockRoot,
			SlotRoot:                        v.SlotRoot,
			TimestampRoot:                   v.TimestampRoot,
			ExecutionPayloadRoot:            v.ExecutionPayloadRoot,
		}
	}

	var validatorFields [][][32]byte
	for _, v := range withdrawalsProof.ValidatorFields {
		validatorFields = append(validatorFields, ConvertProofsToBytes32Array(v))
	}
	var withdrawalFields [][][32]byte
	for _, w := range withdrawalsProof.WithdrawalFields {
		withdrawalFields = append(withdrawalFields, ConvertProofsToBytes32Array(w))
	}
	var validatorFieldsProofs [][]byte
	for _, v := range withdrawalsProof.ValidatorFieldsProofs {
		validatorFieldsProofs = append(validatorFieldsProofs, v.ToByteSlice())
	}

	eigenPod, err := contractEigenPod.NewContractEigenPod(eigenpod, u.chainClient.Client)
	if err != nil {
		return nil, err
	}

	// update validator balance
	return eigenPod.VerifyAndProcessWithdrawals(
		u.chainClient.NoSendTransactOpts,
		withdrawalsProof.OracleTimestamp,
		contractEigenPod.BeaconChainProofsStateRootProof{
			BeaconStateRoot: withdrawalsProof.StateRootProof.BeaconStateRoot,
			Proof:           withdrawalsProof.StateRootProof.StateRootProof.ToByteSlice(),
		},
		withdrawalProofs,
		validatorFieldsProofs,
		validatorFields,
		withdrawalFields,
	)

}

func (u *EigenPodProofTxSubmitter) SubmitVerifyAndProcessWithdrawalsTx(withdrawalProofConfig string, submitTransaction bool) ([]byte, error) {
	ctx := context.Background()
	cfg, err := parseWithdrawalProofConfig(withdrawalProofConfig)
	if err != nil {
		log.Debug().AnErr("Error with parsing withdrawal proof config file", err)
		return nil, err
	}

	oracleBeaconBlockHeader, err := commonutils.ExtractBlockHeader(cfg.BeaconStateFiles.OracleBlockHeaderFile)
	if err != nil {
		log.Debug().AnErr("Error with parsing header file", err)
		return nil, err
	}

	oracleStateJSON, err := commonutils.ParseDenebStateJSONFile(cfg.BeaconStateFiles.OracleStateFile)
	var oracleState deneb.BeaconState
	if err != nil {
		log.Debug().AnErr("GenerateWithdrawalFieldsProof: error with JSON parsing state file", err)
		return nil, err
	}
	commonutils.ParseDenebBeaconStateFromJSON(oracleStateJSON, &oracleState)

	versionedOracleState, err := beacon.CreateVersionedState(&oracleState)

	historicalSummaryStateBlockRoots := make([][]phase0.Root, 0)
	for _, file := range cfg.WithdrawalDetails.HistoricalSummaryStateFiles {
		historicalSummaryStateJSON, err := commonutils.ParseDenebStateJSONFile(file)
		var historicalSummaryState deneb.BeaconState
		if err != nil {
			log.Debug().AnErr("GenerateWithdrawalFieldsProof: error with JSON parsing historical summary state file", err)
			return nil, err
		}
		commonutils.ParseDenebBeaconStateFromJSON(historicalSummaryStateJSON, &historicalSummaryState)

		historicalSummaryStateBlockRoots = append(historicalSummaryStateBlockRoots, historicalSummaryState.BlockRoots)
	}

	withdrawalBlocks := make([]*spec.VersionedSignedBeaconBlock, 0)
	for _, file := range cfg.WithdrawalDetails.WithdrawalBlockBodyFiles {
		block, err := commonutils.ExtractBlock(file)
		if err != nil {
			log.Debug().AnErr("Error with parsing block file", err)
			return nil, err
		}
		versionedSignedBlock, err := beacon.CreateVersionedSignedBlock(block)
		withdrawalBlocks = append(withdrawalBlocks, &versionedSignedBlock)
	}

	withdrawalTx, err := u.GenerateVerifyAndProcessWithdrawalsTx(cfg.EigenPodAddress, &versionedOracleState, &oracleBeaconBlockHeader, historicalSummaryStateBlockRoots, withdrawalBlocks, cfg.WithdrawalDetails.ValidatorIndices)
	if err != nil {
		log.Debug().AnErr("Error with generating withdrawal transaction", err)
		return nil, err
	}

	if !submitTransaction {
		_, err = u.chainClient.EstimateGasPriceAndLimitAndSendTx(ctx, withdrawalTx, "withdraw")
		if err != nil {
			return nil, fmt.Errorf("failed to execute withdrawal transaction: %w", err)
		}
	}
	return withdrawalTx.Data(), nil
}

func parseWithdrawalProofConfig(filePath string) (*WithdrawalProofConfig, error) {

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Debug().AnErr("Error with reading file", err)
		return nil, err
	}

	var cfg WithdrawalProofConfig
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Debug().Msg("error with JSON unmarshalling")
		return nil, err
	}
	return &cfg, nil
}

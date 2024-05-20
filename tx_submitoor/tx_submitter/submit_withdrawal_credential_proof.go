package txsubmitter

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
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

type WithdrawalCredentialProofConfig struct {
	EigenPodAddress  common.Address `json:"EIGENPOD_ADDRESS,required"`
	ValidatorIndices []uint64       `json:"VALIDATOR_INDICES,required"`

	BeaconStateFiles struct {
		OracleStateFile       string `json:"ORACLE_STATE_FILE,required"`
		OracleBlockHeaderFile string `json:"ORACLE_BLOCK_HEADER_FILE,required"`
	}
}

func (u *EigenPodProofTxSubmitter) GenerateVerifyWithdrawalCredentialsTx(
	eigenpod common.Address,
	versionedOracleState *spec.VersionedBeaconState,
	oracleBeaconBlockHeader *phase0.BeaconBlockHeader,
	validatorIndices []uint64,
) (*types.Transaction, error) {
	verifyWithdrawalCredentialsParams, err := u.eigenPodProofs.ProveValidatorContainers(oracleBeaconBlockHeader, versionedOracleState, validatorIndices)
	if err != nil {
		return nil, err
	}

	validatorFieldsBytes := make([][][32]byte, len(verifyWithdrawalCredentialsParams.ValidatorFields))
	for i, v := range verifyWithdrawalCredentialsParams.ValidatorFields {
		validatorFieldsBytes[i] = ConvertProofsToBytes32Array(v)
	}

	validatorIndicesBigInt := make([]*big.Int, len(validatorIndices))
	for i, v := range validatorIndices {
		validatorIndicesBigInt[i] = big.NewInt(int64(v))
	}
	validatorFieldsProofs := make([][]byte, len(verifyWithdrawalCredentialsParams.ValidatorFieldsProofs))
	for i, v := range verifyWithdrawalCredentialsParams.ValidatorFieldsProofs {
		validatorFieldsProofs[i] = v.ToByteSlice()
	}

	eigenPod, err := contractEigenPod.NewContractEigenPod(eigenpod, u.chainClient.Client)
	if err != nil {
		return nil, err
	}

	// update validator balance
	return eigenPod.VerifyWithdrawalCredentials(
		u.chainClient.NoSendTransactOpts,
		verifyWithdrawalCredentialsParams.OracleTimestamp,
		contractEigenPod.BeaconChainProofsStateRootProof{
			BeaconStateRoot: verifyWithdrawalCredentialsParams.StateRootProof.BeaconStateRoot,
			Proof:           verifyWithdrawalCredentialsParams.StateRootProof.StateRootProof.ToByteSlice(),
		},
		validatorIndicesBigInt,
		validatorFieldsProofs,
		validatorFieldsBytes,
	)

}

func (u *EigenPodProofTxSubmitter) SubmitVerifyWithdrawalCredentialsTx(withdrawalCredentialsProofConfig string, submitTransaction bool) ([]byte, error) {
	ctx := context.Background()
	cfg, err := parseWithdrawalCredentialsProofConfig(withdrawalCredentialsProofConfig)
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
	if err != nil {
		log.Debug().AnErr("Error with creating versioned state", err)
		return nil, err
	}

	withdrawalCredentialsTx, err := u.GenerateVerifyWithdrawalCredentialsTx(cfg.EigenPodAddress, &versionedOracleState, &oracleBeaconBlockHeader, cfg.ValidatorIndices)
	if err != nil {
		log.Debug().AnErr("Error with generating withdrawal transaction", err)
		return nil, err
	}

	if submitTransaction {
		_, err = u.chainClient.EstimateGasPriceAndLimitAndSendTx(ctx, withdrawalCredentialsTx, "withdraw")
		if err != nil {
			return nil, fmt.Errorf("failed to execute withdrawal transaction: %w", err)
		}
	}
	return withdrawalCredentialsTx.Data(), nil
}

func parseWithdrawalCredentialsProofConfig(filePath string) (*WithdrawalCredentialProofConfig, error) {

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Debug().AnErr("Error with reading file", err)
		return nil, err
	}

	var cfg WithdrawalCredentialProofConfig
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Debug().Msg("error with JSON unmarshalling")
		return nil, err
	}
	return &cfg, nil
}

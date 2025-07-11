package utils

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"sort"
	"strconv"
	"strings"

	lo "github.com/samber/lo"

	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/EigenPod"
	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/fatih/color"
	"github.com/jbrower95/multicall-go"
)

const (
	ValidatorStatusInactive  = 0
	ValidatorStatusActive    = 1
	ValidatorStatusWithdrawn = 2
)

const (
	BLSWithdrawalPrefix         = 0
	ETH1WithdrawalPrefix        = 1
	CompoundingWithdrawalPrefix = 2
)

// Constants defined in predeploy EIPs:
// - EIP-7521: https://eips.ethereum.org/EIPS/eip-7251#constants
// - EIP-7002: https://eips.ethereum.org/EIPS/eip-7002#configuration
var (
	CONSOLIDATION_PREDEPLOY = common.HexToAddress("0x0000BBdDc7CE488642fb579F8B00f3a590007251")
	WITHDRAWAL_PREDEPLOY    = common.HexToAddress("0x00000961Ef480Eb55e80D19ad83579A64c007002")

	// 2**256 - 1
	EXCESS_INHIBITOR = new(big.Int).Sub(
		new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil),
		big.NewInt(1),
	)

	// Consolidation constants
	MAX_CONSOLIDATION_REQUESTS_PER_BLOCK      = big.NewInt(2)
	MIN_CONSOLIDATION_REQUEST_FEE             = big.NewInt(1)
	CONSOLIDATION_REQUEST_FEE_UPDATE_FRACTION = big.NewInt(17)

	// Withdrawal constants
	MAX_WITHDRAWAL_REQUESTS_PER_BLOCK      = big.NewInt(16)
	MIN_WITHDRAWAL_REQUEST_FEE             = big.NewInt(1)
	WITHDRAWAL_REQUEST_FEE_UPDATE_FRACTION = big.NewInt(17)
)

type Checkpoint struct {
	ProofsRemaining uint64
	StartedAt       uint64
}

type Validator struct {
	Slashed                             bool
	Index                               uint64
	Status                              int
	PublicKey                           string
	IsAwaitingActivationQueue           bool
	IsAwaitingWithdrawalCredentialProof bool
	EffectiveBalance                    uint64
	CurrentBalance                      uint64
	WithdrawalPrefix                    uint8
}

type Predeploy struct {
	Address common.Address
	Caller  bind.ContractCaller
}

func newPredeploy(address common.Address, backend bind.ContractBackend) *Predeploy {
	return &Predeploy{
		address,
		backend,
	}
}

func Panic(message string) {
	color.Red(fmt.Sprintf("error: %s\n\n", message))

	os.Exit(1)
}

func WeiToGwei(val *big.Int) *big.Float {
	return new(big.Float).Quo(
		new(big.Float).SetInt(val),
		big.NewFloat(params.GWei),
	)
}

func GweiToEther(val *big.Float) *big.Float {
	return new(big.Float).Quo(val, big.NewFloat(params.GWei))
}

func GweiToWei(val *big.Float) *big.Float {
	return new(big.Float).Mul(val, big.NewFloat(params.GWei))
}

func IGweiToWei(val *big.Int) *big.Int {
	return new(big.Int).Mul(val, big.NewInt(params.GWei))
}

func IweiToEther(val *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(val), big.NewFloat(params.Ether))
}

func PanicOnError(message string, err error) {
	if err != nil {
		color.Red(fmt.Sprintf("error: %s\n\n", message))

		info := color.New(color.FgRed, color.Italic)
		info.Printf(fmt.Sprintf("caused by: %s\n", err))

		os.Exit(1)
	}
}

func Chunk[T any](arr []T, chunkSize uint64) [][]T {
	// Validate the chunkSize to ensure it's positive
	if chunkSize <= 0 {
		panic("chunkSize must be greater than 0")
	}

	// Create a slice to hold the chunks
	var chunks [][]T

	// Loop through the input slice and create chunks
	arrLen := uint64(len(arr))
	for i := uint64(0); i < arrLen; i += chunkSize {
		end := uint64(i + chunkSize)
		if end > arrLen {
			end = arrLen
		}
		chunks = append(chunks, arr[i:end])
	}

	return chunks
}

type ConsolidatableValidatorSet struct {
	Sources []ValidatorWithIndex
	Targets []ValidatorWithIndex
}

type ValidatorWithIndex = struct {
	Validator *phase0.Validator
	Index     uint64
}

type ValidatorWithOnchainInfo = struct {
	Info      EigenPod.IEigenPodTypesValidatorInfo
	Validator *phase0.Validator
	Index     uint64
}

type ValidatorWithMaybeOnchainInfo = struct {
	Info      *EigenPod.IEigenPodTypesValidatorInfo
	Validator *phase0.Validator
	Index     uint64
}

type Owner = struct {
	FromAddress        common.Address
	PublicKey          *ecdsa.PublicKey
	TransactionOptions *bind.TransactOpts
	IsDryRun           bool
}

func StartCheckpoint(ctx context.Context, eigenpodAddress string, ownerPrivateKey string, chainId *big.Int, eth *ethclient.Client, forceCheckpoint bool, noSend bool) (*types.Transaction, error) {
	ownerAccount, err := PrepareAccount(&ownerPrivateKey, chainId, noSend)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	eigenPod, err := EigenPod.NewEigenPod(common.HexToAddress(eigenpodAddress), eth)
	if err != nil {
		return nil, fmt.Errorf("failed to reach eigenpod: %w", err)
	}

	revertIfNoBalance := !forceCheckpoint

	txn, err := eigenPod.StartCheckpoint(ownerAccount.TransactionOptions, revertIfNoBalance)
	if err != nil {
		if !forceCheckpoint {
			return nil, fmt.Errorf("failed to start checkpoint (try running again with `--force`): %w", err)
		}

		return nil, fmt.Errorf("failed to start checkpoint: %w", err)
	}

	return txn, nil
}

func GetBeaconClient(beaconUri string, verbose bool) (BeaconClient, error) {
	beaconClient, _, err := NewBeaconClient(beaconUri, verbose)
	return beaconClient, err
}

func GetCurrentCheckpoint(eigenpodAddress string, client *ethclient.Client) (uint64, error) {
	eigenPod, err := EigenPod.NewEigenPod(common.HexToAddress(eigenpodAddress), client)
	if err != nil {
		return 0, fmt.Errorf("failed to locate eigenpod. is your address correct?: %w", err)
	}

	timestamp, err := eigenPod.CurrentCheckpointTimestamp(nil)
	if err != nil {
		return 0, fmt.Errorf("failed to locate eigenpod. Is your address correct?: %w", err)
	}

	return timestamp, nil
}

// Fetch and return the current checkpoint timestamp for the pod
// If the checkpoint exists (timestamp != 0), also return the beacon state for the checkpoint
// If the checkpoint does not exist (timestamp == 0), return the head beacon state (i.e. the state we would use "if we start a checkpoint now")
func GetCheckpointTimestampAndBeaconState(
	ctx context.Context,
	eigenpodAddress string,
	eth *ethclient.Client,
	beaconClient BeaconClient,
) (uint64, *spec.VersionedBeaconState, error) {
	tracing := GetContextTracingCallbacks(ctx)

	tracing.OnStartSection("GetCurrentCheckpoint", map[string]string{})
	checkpointTimestamp, err := GetCurrentCheckpoint(eigenpodAddress, eth)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to fetch current checkpoint: %w", err)
	}
	tracing.OnEndSection()

	// stateId to look up beacon state. "head" by default (if we do not have a checkpoint)
	beaconStateId := "head"

	// If we have a checkpoint, get the state id for the checkpoint's block root
	if checkpointTimestamp != 0 {
		// Fetch the checkpoint's block root
		tracing.OnStartSection("GetCurrentCheckpointBlockRoot", map[string]string{})
		blockRoot, err := GetCurrentCheckpointBlockRoot(eigenpodAddress, eth)
		if err != nil {
			return 0, nil, fmt.Errorf("failed to fetch last checkpoint: %w", err)
		}
		if blockRoot == nil {
			return 0, nil, fmt.Errorf("failed to fetch last checkpoint - nil blockRoot")
		}
		// Block root should be nonzero because we have an active checkpoint
		rootBytes := *blockRoot
		if AllZero(rootBytes[:]) {
			return 0, nil, fmt.Errorf("failed to fetch last checkpoint - empty blockRoot")
		}
		tracing.OnEndSection()

		headerBlock := "0x" + hex.EncodeToString((*blockRoot)[:])
		tracing.OnStartSection("GetBeaconHeader", map[string]string{})
		header, err := beaconClient.GetBeaconHeader(ctx, headerBlock)
		if err != nil {
			return 0, nil, fmt.Errorf("failed to fetch beacon header (%s): %w", headerBlock, err)
		}
		tracing.OnEndSection()

		beaconStateId = strconv.FormatUint(uint64(header.Header.Message.Slot), 10)
	}

	tracing.OnStartSection("GetBeaconState", map[string]string{})
	beaconState, err := beaconClient.GetBeaconState(ctx, beaconStateId)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to fetch beacon state: %w", err)
	}
	tracing.OnEndSection()

	return checkpointTimestamp, beaconState, nil
}

func GetBeaconHeadState(ctx context.Context, beaconClient BeaconClient) (*spec.VersionedBeaconState, error) {
	tracing := GetContextTracingCallbacks(ctx)

	tracing.OnStartSection("GetBeaconHeadState", map[string]string{})
	headState, err := beaconClient.GetBeaconState(ctx, "head")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch beacon state: %w", err)
	}
	tracing.OnEndSection()

	return headState, nil
}

// GetEigenPodValidatorsByIndex looks up all validators for an eigenpod and returns a map
// from validator index to validator info
func GetEigenPodValidatorsByIndex(
	eigenpodAddress string,
	beaconState *spec.VersionedBeaconState,
) (map[uint64]*phase0.Validator, error) {
	allValidatorsForEigenPod, err := FindAllValidatorsForEigenpod(eigenpodAddress, beaconState)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch validators for eigenpod: %w", err)
	}

	allValidatorsMap := make(map[uint64]*phase0.Validator)
	for _, v := range allValidatorsForEigenPod {
		allValidatorsMap[v.Index] = v.Validator
	}

	return allValidatorsMap, nil
}

func SortByStatus(validators map[string]Validator) ([]Validator, []Validator, []Validator, []Validator) {
	var awaitingActivationQueueValidators, inactiveValidators, activeValidators, withdrawnValidators []Validator

	// Iterate over all `validators` and sort them into inactive, active, or withdrawn.
	for _, validator := range validators {
		switch validator.Status {
		case ValidatorStatusInactive:
			if validator.IsAwaitingActivationQueue {
				awaitingActivationQueueValidators = append(awaitingActivationQueueValidators, validator)
			} else {
				inactiveValidators = append(inactiveValidators, validator)
			}
		case ValidatorStatusActive:
			activeValidators = append(activeValidators, validator)
		case ValidatorStatusWithdrawn:
			withdrawnValidators = append(withdrawnValidators, validator)
		}
	}

	// Sort each of these mappings in order of ascending Validator.Index
	sort.Slice(inactiveValidators, func(i, j int) bool {
		return inactiveValidators[i].Index < inactiveValidators[j].Index
	})
	sort.Slice(activeValidators, func(i, j int) bool {
		return activeValidators[i].Index < activeValidators[j].Index
	})
	sort.Slice(withdrawnValidators, func(i, j int) bool {
		return withdrawnValidators[i].Index < withdrawnValidators[j].Index
	})

	return awaitingActivationQueueValidators, inactiveValidators, activeValidators, withdrawnValidators
}

// search through beacon state for validators whose withdrawal address is set to eigenpod.
func FindAllValidatorsForEigenpod(eigenpodAddress string, beaconState *spec.VersionedBeaconState) ([]ValidatorWithIndex, error) {
	allValidators, err := beaconState.Validators()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch beacon state: %w", err)
	}

	eigenpod := common.HexToAddress(eigenpodAddress)

	var outputValidators []ValidatorWithIndex = []ValidatorWithIndex{}
	var i uint64 = 0
	maxValidators := uint64(len(allValidators))
	for i = 0; i < maxValidators; i++ {
		validator := allValidators[i]
		if validator == nil || (validator.WithdrawalCredentials[0] != 1 && validator.WithdrawalCredentials[0] != 2) { // withdrawalCredentials _need_ their first byte set to 1 or 2 to withdraw to an eigenpod on the execution layer.
			continue
		}
		// we check that the last 20 bytes of expectedCredentials matches validatorCredentials.
		// // first 12 bytes are not the pubKeyHash, see (https://github.com/Layr-Labs/eigenlayer-contracts/blob/d148952a2942a97a218a2ab70f9b9f1792796081/src/contracts/pods/EigenPod.sol#L663)
		validatorWithdrawalAddress := common.BytesToAddress(validator.WithdrawalCredentials[12:])

		if eigenpod.Cmp(validatorWithdrawalAddress) == 0 {
			outputValidators = append(outputValidators, ValidatorWithIndex{
				Validator: validator,
				Index:     i,
			})
		}
	}
	return outputValidators, nil
}

var zeroes = [16]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

func FetchMultipleOnchainValidatorInfoMulticalls(eigenpodAddress string, allValidators []*phase0.Validator) ([]*multicall.MultiCallMetaData[EigenPod.IEigenPodTypesValidatorInfo], error) {
	eigenpodAbi, err := abi.JSON(strings.NewReader(EigenPod.EigenPodABI))
	if err != nil {
		return nil, fmt.Errorf("failed to load eigenpod abi: %s", err)
	}

	type MulticallAndError struct {
		Multicall *multicall.MultiCallMetaData[EigenPod.IEigenPodTypesValidatorInfo]
		Error     error
	}

	requests := lo.Map(allValidators, func(validator *phase0.Validator, index int) MulticallAndError {
		pubKeyHash := sha256.Sum256(
			append(
				validator.PublicKey[:],
				zeroes[:]...,
			),
		)

		mc, err := multicall.Describe[EigenPod.IEigenPodTypesValidatorInfo](
			common.HexToAddress(eigenpodAddress),
			eigenpodAbi,
			"validatorPubkeyHashToInfo",
			pubKeyHash)

		return MulticallAndError{
			Multicall: mc,
			Error:     err,
		}
	})

	errs := []error{}
	for _, mc := range requests {
		if mc.Error != nil {
			errs = append(errs, mc.Error)
		}
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to form request for validator info: %s", errors.Join(errs...))
	}

	allMulticalls := lo.Map(requests, func(mc MulticallAndError, _ int) *multicall.MultiCallMetaData[EigenPod.IEigenPodTypesValidatorInfo] {
		return mc.Multicall
	})
	return allMulticalls, nil
}

func FetchMultipleOnchainValidatorInfo(ctx context.Context, client *ethclient.Client, eigenpodAddress string, allValidators []ValidatorWithIndex) ([]ValidatorWithOnchainInfo, error) {
	allMulticalls, err := FetchMultipleOnchainValidatorInfoMulticalls(eigenpodAddress, lo.Map(allValidators, func(validator ValidatorWithIndex, i int) *phase0.Validator { return validator.Validator }))
	if err != nil {
		return nil, fmt.Errorf("failed to form multicalls: %s", err.Error())
	}

	mc, err := multicall.NewMulticallClient(ctx, client, &multicall.TMulticallClientOptions{
		MaxBatchSizeBytes: 4096,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to contact multicall: %s", err.Error())
	}

	results, err := multicall.DoMany(mc, allMulticalls...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch validator info: %s", err.Error())
	}

	if results == nil {
		return nil, errors.New("no results returned fetching validator info")
	}

	return lo.Map(*results, func(info *EigenPod.IEigenPodTypesValidatorInfo, i int) ValidatorWithOnchainInfo {
		return ValidatorWithOnchainInfo{
			Info:      *info,
			Validator: allValidators[i].Validator,
			Index:     allValidators[i].Index,
		}
	}), nil
}

func FetchMultipleOnchainValidatorInfoWithFailures(ctx context.Context, client *ethclient.Client, eigenpodAddress string, allValidators []ValidatorWithIndex) ([]ValidatorWithMaybeOnchainInfo, error) {
	allMulticalls, err := FetchMultipleOnchainValidatorInfoMulticalls(eigenpodAddress, lo.Map(allValidators, func(validator ValidatorWithIndex, i int) *phase0.Validator { return validator.Validator }))
	if err != nil {
		return nil, fmt.Errorf("failed to form multicalls: %s", err.Error())
	}

	mc, err := multicall.NewMulticallClient(ctx, client, &multicall.TMulticallClientOptions{
		MaxBatchSizeBytes: 4096,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to contact multicall: %s", err.Error())
	}

	results, err := multicall.DoManyAllowFailures(mc, allMulticalls...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch validator info: %s", err.Error())
	}

	if results == nil {
		return nil, errors.New("no results returned fetching validator info")
	}

	return lo.Map(*results, func(info multicall.TypedMulticall3Result[*EigenPod.IEigenPodTypesValidatorInfo], i int) ValidatorWithMaybeOnchainInfo {
		return ValidatorWithMaybeOnchainInfo{
			Info:      info.Value,
			Validator: allValidators[i].Validator,
			Index:     allValidators[i].Index,
		}
	}), nil
}

func GetCurrentCheckpointBlockRoot(eigenpodAddress string, eth *ethclient.Client) (*[32]byte, error) {
	eigenPod, err := EigenPod.NewEigenPod(common.HexToAddress(eigenpodAddress), eth)
	if err != nil {
		return nil, fmt.Errorf("failed to locate Eigenpod. Is your address correct?: %w", err)
	}

	checkpoint, err := eigenPod.CurrentCheckpoint(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to reach eigenpod: %w", err)
	}

	return &checkpoint.BeaconBlockRoot, nil
}

func IsAwaitingWithdrawalCredentialProof(validatorInfo EigenPod.IEigenPodTypesValidatorInfo, validator *phase0.Validator) bool {
	return (validatorInfo.Status == ValidatorStatusInactive) && validator.ExitEpoch == FAR_FUTURE_EPOCH && validator.ActivationEpoch != FAR_FUTURE_EPOCH
}

// this is a mapping from <chainId, genesis_fork_version>.
func ForkVersions() map[uint64]string {
	return map[uint64]string{
		11155111: "90000069", //sepolia (https://github.com/eth-clients/sepolia/blob/main/README.md?plain=1#L66C26-L66C36)
		560048:   "10000910", //hoodi (https://github.com/eth-clients/hoodi/blob/main/README.md)
		17000:    "01017000", //holesky (https://github.com/eth-clients/holesky/blob/main/README.md)
		1:        "00000000", // mainnet (https://github.com/eth-clients/mainnet)
	}
}

func GetEthClient(ctx context.Context, node string) (*ethclient.Client, *big.Int, error) {
	eth, err := ethclient.Dial(node)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to reach eth --node: %w", err)
	}

	chainId, err := eth.ChainID(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch chain id: %w", err)
	}

	if chainId == nil || (chainId.Int64() != 17000 && chainId.Int64() != 1 && chainId.Int64() != 560048) {
		return nil, nil, errors.New("this tool only supports the Holesky, Hoodi, and Mainnet Ethereum Networks")
	}
	return eth, chainId, nil
}

func GetClients(ctx context.Context, node, beaconNodeUri string, enableLogs bool) (*ethclient.Client, BeaconClient, *big.Int, error) {
	eth, chainId, err := GetEthClient(ctx, node)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to reach eth --node: %w", err)
	}

	beaconClient, err := GetBeaconClient(beaconNodeUri, enableLogs)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to reach beacon client: %w", err)
	}

	genesisForkVersion, err := beaconClient.GetGenesisForkVersion(ctx)
	expectedForkVersion := ForkVersions()[chainId.Uint64()]
	gotForkVersion := hex.EncodeToString((*genesisForkVersion)[:])
	if err != nil || expectedForkVersion != gotForkVersion {
		return nil, nil, nil, fmt.Errorf("check that both nodes correspond to the same network and try again (expected genesis_fork_version: %s, got %s)", expectedForkVersion, gotForkVersion)
	}

	return eth, beaconClient, chainId, nil
}

func CastBalanceProofs(proofs []*eigenpodproofs.BalanceProof) []EigenPod.BeaconChainProofsBalanceProof {
	out := []EigenPod.BeaconChainProofsBalanceProof{}

	for i := 0; i < len(proofs); i++ {
		proof := proofs[i]
		out = append(out, EigenPod.BeaconChainProofsBalanceProof{
			PubkeyHash:  proof.PubkeyHash,
			BalanceRoot: proof.BalanceRoot,
			Proof:       proof.Proof.ToByteSlice(),
		})
	}

	return out
}

func PanicIfNoConsent(prompt string) {
	color.New(color.Bold).Printf("%s - Do you want to proceed? (y/n): ", prompt)
	var reply string

	fmt.Scanln(&reply)
	if reply == "y" {
		return
	} else {
		Panic("abort.")
	}
}

func PrepareAccount(owner *string, chainID *big.Int, noSend bool) (*Owner, error) {
	isSimulatingGas := owner != nil && *owner != ""
	senderPk, err := func() (string, error) {
		// if we're trying to send a transaction, make sure we were supplied a private key
		if owner == nil || *owner == "" {
			if !noSend {
				return "", errors.New("no private key supplied")
			}

			return "372d94b8645091147a5dfc10a454d0d539773d2431293bf0a195b44fa5ddbb33", nil // this is a RANDOM private key used as a default value. do not use.
		}

		return *owner, nil
	}()

	if err != nil {
		return nil, err
	}

	// trim leading "0x" if needed
	if senderPk[0:2] == "0x" {
		senderPk = senderPk[2:]
	}

	privateKey, err := crypto.HexToECDSA(senderPk)
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, err
	}

	auth.NoSend = noSend
	if noSend && !isSimulatingGas {
		auth.GasPrice = nil             // big.NewInt(10)  // Gas price to use for the transaction execution (nil = gas price oracle)
		auth.GasFeeCap = big.NewInt(10) // big.NewInt(10) // Gas fee cap to use for the 1559 transaction execution (nil = gas price oracle)
		auth.GasTipCap = big.NewInt(2)  // big.NewInt(2) // Gas priority fee cap to use for the 1559 transaction execution (nil = gas price oracle)
		auth.GasLimit = 21000
	}

	return &Owner{
		FromAddress:        fromAddress,
		PublicKey:          publicKeyECDSA,
		TransactionOptions: auth,
		IsDryRun:           noSend,
	}, nil
}

// golang was a mistake. these types are literally identical :'(
func CastValidatorFields(proof [][]eigenpodproofs.Bytes32) [][][32]byte {
	result := make([][][32]byte, len(proof))

	for i, slice := range proof {
		result[i] = make([][32]byte, len(slice))
		for j, bytes := range slice {
			result[i][j] = bytes
		}
	}

	return result
}

func Uint64ArrayToBigIntArray(nums []uint64) []*big.Int {
	out := []*big.Int{}
	for i := 0; i < len(nums); i++ {
		bigInt := new(big.Int).SetUint64(nums[i])
		out = append(out, bigInt)
	}
	return out
}

func AllZero(s []byte) bool {
	for _, v := range s {
		if v != 0 {
			return false
		}
	}
	return true
}

func WriteOutputToFileOrStdout(output []byte, out *string) error {
	if out != nil && *out != "" {
		err := os.WriteFile(*out, output, os.ModePerm)
		if err != nil {
			return err
		}
		PanicOnError("failed to write to disk", err)
		color.Green("Wrote output to %s\n", *out)
	} else {
		fmt.Println(string(output))
	}

	return nil
}

func SelectCheckpointableValidators(
	client *ethclient.Client,
	eigenpodAddress string,
	validators []ValidatorWithOnchainInfo,
	lastCheckpoint uint64,
) ([]ValidatorWithOnchainInfo, error) {
	var checkpointValidators = []ValidatorWithOnchainInfo{}
	for i := 0; i < len(validators); i++ {
		validator := validators[i]

		notCheckpointed := (validator.Info.LastCheckpointedAt != lastCheckpoint) || (validator.Info.LastCheckpointedAt == 0)
		isActive := validator.Info.Status == ValidatorStatusActive

		if notCheckpointed && isActive {
			checkpointValidators = append(checkpointValidators, validator)
		}
	}
	return checkpointValidators, nil
}

// (https://github.com/Layr-Labs/eigenlayer-contracts/blob/d148952a2942a97a218a2ab70f9b9f1792796081/src/contracts/libraries/BeaconChainProofs.sol#L64)
const FAR_FUTURE_EPOCH = math.MaxUint64

// Validators whose deposits have been processed but are awaiting activation on the beacon chain
// If the validator has 32 ETH effective balance, they should
func SelectAwaitingActivationValidators(
	client *ethclient.Client,
	eigenpodAddress string,
	validators []ValidatorWithOnchainInfo,
) ([]ValidatorWithOnchainInfo, error) {
	var awaitingActivationValidators = []ValidatorWithOnchainInfo{}
	for i := 0; i < len(validators); i++ {
		validator := validators[i]

		if validator.Validator.ActivationEpoch == FAR_FUTURE_EPOCH {
			awaitingActivationValidators = append(awaitingActivationValidators, validator)
		}
	}
	return awaitingActivationValidators, nil
}

func SelectAwaitingCredentialValidators(
	client *ethclient.Client,
	eigenpodAddress string,
	validators []ValidatorWithOnchainInfo,
) ([]ValidatorWithOnchainInfo, error) {
	var awaitingCredentialValidators = []ValidatorWithOnchainInfo{}
	for i := 0; i < len(validators); i++ {
		validator := validators[i]

		if (validator.Info.Status == ValidatorStatusInactive) &&
			(validator.Validator.ExitEpoch == FAR_FUTURE_EPOCH) &&
			(validator.Validator.ActivationEpoch != FAR_FUTURE_EPOCH) {
			awaitingCredentialValidators = append(awaitingCredentialValidators, validator)
		}
	}
	return awaitingCredentialValidators, nil
}

func SelectActiveValidators(
	client *ethclient.Client,
	eigenpodAddress string,
	validators []ValidatorWithOnchainInfo,
) ([]ValidatorWithOnchainInfo, error) {
	var activeValidators = []ValidatorWithOnchainInfo{}
	for i := 0; i < len(validators); i++ {
		validator := validators[i]
		if validator.Info.Status == ValidatorStatusActive {
			activeValidators = append(activeValidators, validator)
		}
	}
	return activeValidators, nil
}

/// PREDEPLOY UTILS

// CurrentConsolidationFee queries the EIP-7521 predeploy to get the current per-request fee in wei
func CurrentConsolidationFee(client *ethclient.Client) (*big.Int, error) {
	ctx := context.Background()

	predeploy := newPredeploy(CONSOLIDATION_PREDEPLOY, client)

	blockNum, err := client.BlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting latest block number: %w", err)
	}

	msg := ethereum.CallMsg{
		From: CONSOLIDATION_PREDEPLOY,
		To:   &CONSOLIDATION_PREDEPLOY,
		Data: []byte{},
	}

	result, err := predeploy.Caller.CallContract(ctx, msg, new(big.Int).SetUint64(blockNum))
	if err != nil {
		return nil, fmt.Errorf("error calling consolidation predeploy: %w", err)
	}

	if len(result) != 32 {
		return nil, fmt.Errorf("predeploy error: expected 32 byte result, got %d bytes", len(result))
	}

	return new(big.Int).SetBytes(result), nil
}

// CurrentWithdrawalFee queries the EIP-7002 predeploy to get the current per-request fee in wei
func CurrentWithdrawalFee(client *ethclient.Client) (*big.Int, error) {
	ctx := context.Background()

	predeploy := newPredeploy(WITHDRAWAL_PREDEPLOY, client)

	blockNum, err := client.BlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting latest block number: %w", err)
	}

	msg := ethereum.CallMsg{
		From: WITHDRAWAL_PREDEPLOY,
		To:   &WITHDRAWAL_PREDEPLOY,
		Data: []byte{},
	}

	result, err := predeploy.Caller.CallContract(ctx, msg, new(big.Int).SetUint64(blockNum))
	if err != nil {
		return nil, fmt.Errorf("error calling withdrawal predeploy: %w", err)
	}

	if len(result) != 32 {
		return nil, fmt.Errorf("predeploy error: expected 32 byte result, got %d bytes", len(result))
	}

	return new(big.Int).SetBytes(result), nil
}

// GetExcessConsolidationRequests reads EIP-7521 predeploy storage slot 0 to get the number of excess
// requests currently in the queue.
func GetExcessConsolidationRequests(client *ethclient.Client) (*big.Int, error) {
	ctx := context.Background()

	// Excess requests are at storage slot 0
	result, err := client.StorageAt(ctx, CONSOLIDATION_PREDEPLOY, common.Hash{}, nil)
	if err != nil {
		return nil, fmt.Errorf("error reading storage: %w", err)
	}

	return new(big.Int).SetBytes(result), nil
}

// GetExcessConsolidationRequests reads EIP-7002 predeploy storage slot 0 to get the number of excess
// requests currently in the queue.
func GetExcessWithdrawalRequests(client *ethclient.Client) (*big.Int, error) {
	ctx := context.Background()

	// Excess requests are at storage slot 0
	result, err := client.StorageAt(ctx, WITHDRAWAL_PREDEPLOY, common.Hash{}, nil)
	if err != nil {
		return nil, fmt.Errorf("error reading storage: %w", err)
	}

	return new(big.Int).SetBytes(result), nil
}

type PredeployFeeInfo struct {
	CurrentQueueSize *big.Int
	FeePerRequest    *big.Int
	TotalFee         *big.Int
	OverestimateFee  *big.Int
}

// EstimateConsolidationFeePerChunk estimates the consolidation fee required to submit multiple "chunks"
// of consolidation requests. This method assumes the requests are submitted in sequential blocks,
// with the system address updating the fee between each block.
//
// The returned array is the "fee per request" for each chunk
func EstimateConsolidationFeePerChunk(
	client *ethclient.Client,
	requestChunks [][]EigenPod.IEigenPodTypesConsolidationRequest,
) ([]*PredeployFeeInfo, error) {
	curExcess, err := GetExcessConsolidationRequests(client)
	if err != nil {
		return nil, fmt.Errorf("error getting excess requests from predeploy: %w", err)
	}

	// feesPerChunk := make([]*big.Int, 0, len(requestChunks))
	chunkInfo := make([]*PredeployFeeInfo, 0, len(requestChunks))

	for _, chunk := range requestChunks {
		fee := fakeExponential(
			MIN_CONSOLIDATION_REQUEST_FEE,
			curExcess,
			CONSOLIDATION_REQUEST_FEE_UPDATE_FRACTION,
		)

		totalFee := new(big.Int).Mul(
			big.NewInt(int64(len(chunk))),
			fee,
		)

		chunkInfo = append(chunkInfo, &PredeployFeeInfo{
			CurrentQueueSize: curExcess,
			FeePerRequest:    fee,
			TotalFee:         totalFee,
			OverestimateFee:  totalFee,
		})

		// Simulate excess update across block boundary
		// new_excess = max(0, curExcess + len(chunk) - MAX_DEQUEUE)
		reqCount := big.NewInt(int64(len(chunk)))
		curExcess.Add(curExcess, reqCount)
		curExcess.Sub(curExcess, MAX_CONSOLIDATION_REQUESTS_PER_BLOCK)

		if curExcess.Sign() == -1 {
			curExcess.SetUint64(0)
		}
	}

	return chunkInfo, nil
}

// EstimateWithdrawalFeePerChunk estimates the withdrawal fee required to submit multiple "chunks"
// of withdrawal requests. This method assumes the requests are submitted in sequential blocks,
// with the system address updating the fee between each block.
//
// The returned array is the "fee per request" for each chunk
func EstimateWithdrawalFeePerChunk(
	client *ethclient.Client,
	requestChunks [][]EigenPod.IEigenPodTypesWithdrawalRequest,
) ([]*PredeployFeeInfo, error) {
	curExcess, err := GetExcessWithdrawalRequests(client)
	if err != nil {
		return nil, fmt.Errorf("error getting excess requests from predeploy: %w", err)
	}

	// feesPerChunk := make([]*big.Int, 0, len(requestChunks))
	chunkInfo := make([]*PredeployFeeInfo, 0, len(requestChunks))

	for _, chunk := range requestChunks {
		fee := fakeExponential(
			MIN_WITHDRAWAL_REQUEST_FEE,
			curExcess,
			WITHDRAWAL_REQUEST_FEE_UPDATE_FRACTION,
		)

		totalFee := new(big.Int).Mul(
			big.NewInt(int64(len(chunk))),
			fee,
		)

		chunkInfo = append(chunkInfo, &PredeployFeeInfo{
			CurrentQueueSize: curExcess,
			FeePerRequest:    fee,
			TotalFee:         totalFee,
			OverestimateFee:  totalFee,
		})

		// Simulate excess update across block boundary
		// new_excess = max(0, curExcess + len(chunk) - MAX_DEQUEUE)
		reqCount := big.NewInt(int64(len(chunk)))
		curExcess.Add(curExcess, reqCount)
		curExcess.Sub(curExcess, MAX_WITHDRAWAL_REQUESTS_PER_BLOCK)

		if curExcess.Sign() == -1 {
			curExcess.SetUint64(0)
		}
	}

	return chunkInfo, nil
}

// GetConsolidationFeeInfoForRequest retrieves info on the current fees required
// to submit consolidation requests.
func GetConsolidationFeeInfoForRequest(
	client *ethclient.Client,
	request []EigenPod.IEigenPodTypesConsolidationRequest,
	overestimateFactor float64,
) (*PredeployFeeInfo, error) {
	if overestimateFactor == 0 {
		overestimateFactor = 1
	}

	curExcess, err := GetExcessConsolidationRequests(client)
	if err != nil {
		return nil, fmt.Errorf("error getting excess requests from predeploy: %w", err)
	}

	curFee, err := CurrentConsolidationFee(client)
	if err != nil {
		return nil, fmt.Errorf("error getting current consolidation fee: %w", err)
	}

	totalFee := new(big.Int).Mul(
		big.NewInt(int64(len(request))),
		curFee,
	)

	overestimatedFeeFloat := new(big.Float).Mul(
		new(big.Float).SetInt(totalFee),
		big.NewFloat(overestimateFactor), // e.g., 1.5
	)

	overestimatedFee := new(big.Int)
	overestimatedFeeFloat.Int(overestimatedFee) // rounds down

	return &PredeployFeeInfo{
		CurrentQueueSize: curExcess,
		FeePerRequest:    curFee,
		TotalFee:         totalFee,
		OverestimateFee:  overestimatedFee,
	}, nil
}

// GetWithdrawalFeeInfoForRequest retrieves info on the current fees required
// to submit withdrawal requests.
func GetWithdrawalFeeInfoForRequest(
	client *ethclient.Client,
	request []EigenPod.IEigenPodTypesWithdrawalRequest,
	overestimateFactor float64,
) (*PredeployFeeInfo, error) {
	if overestimateFactor == 0 {
		overestimateFactor = 1
	}

	curExcess, err := GetExcessWithdrawalRequests(client)
	if err != nil {
		return nil, fmt.Errorf("error getting excess requests from predeploy: %w", err)
	}

	curFee, err := CurrentWithdrawalFee(client)
	if err != nil {
		return nil, fmt.Errorf("error getting current withdrawal fee: %w", err)
	}

	totalFee := new(big.Int).Mul(
		big.NewInt(int64(len(request))),
		curFee,
	)

	overestimatedFeeFloat := new(big.Float).Mul(
		new(big.Float).SetInt(totalFee),
		big.NewFloat(overestimateFactor), // e.g., 1.5
	)

	overestimatedFee := new(big.Int)
	overestimatedFeeFloat.Int(overestimatedFee) // rounds down

	return &PredeployFeeInfo{
		CurrentQueueSize: curExcess,
		FeePerRequest:    curFee,
		TotalFee:         totalFee,
		OverestimateFee:  overestimatedFee,
	}, nil
}

// fakeExponential implements the EIP-7521 and EIP-7002 exponential fee increase formula
// (see https://eips.ethereum.org/EIPS/eip-7251#fee-calculation)
func fakeExponential(factor, numerator, denominator *big.Int) *big.Int {
	i := big.NewInt(1)
	output := big.NewInt(0)
	numAccum := new(big.Int).Mul(factor, denominator) // numerator_accum = factor * denominator

	for numAccum.Cmp(big.NewInt(0)) > 0 {
		output.Add(output, numAccum)

		// numAccum = (numAccum * numerator) // (denominator * i)
		numAccum.Mul(numAccum, numerator)
		denomProd := new(big.Int).Mul(denominator, i)
		numAccum.Div(numAccum, denomProd)

		i.Add(i, big.NewInt(1))
	}

	return output.Div(output, denominator)
}

// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package onchain

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// BeaconChainProofsBalanceContainerProof is an auto generated low-level Go binding around an user-defined struct.
type BeaconChainProofsBalanceContainerProof struct {
	BalanceContainerRoot [32]byte
	Proof                []byte
}

// BeaconChainProofsBalanceProof is an auto generated low-level Go binding around an user-defined struct.
type BeaconChainProofsBalanceProof struct {
	PubkeyHash  [32]byte
	BalanceRoot [32]byte
	Proof       []byte
}

// BeaconChainProofsStateRootProof is an auto generated low-level Go binding around an user-defined struct.
type BeaconChainProofsStateRootProof struct {
	BeaconStateRoot [32]byte
	Proof           []byte
}

// BeaconChainProofsValidatorProof is an auto generated low-level Go binding around an user-defined struct.
type BeaconChainProofsValidatorProof struct {
	ValidatorFields [][32]byte
	Proof           []byte
}

// IEigenPodCheckpoint is an auto generated low-level Go binding around an user-defined struct.
type IEigenPodCheckpoint struct {
	BeaconBlockRoot   [32]byte
	ProofsRemaining   *big.Int
	PodBalanceGwei    uint64
	BalanceDeltasGwei *big.Int
}

// IEigenPodValidatorInfo is an auto generated low-level Go binding around an user-defined struct.
type IEigenPodValidatorInfo struct {
	ValidatorIndex      uint64
	RestakedBalanceGwei uint64
	LastCheckpointedAt  uint64
	Status              uint8
}

// EigenPodMetaData contains all meta data concerning the EigenPod contract.
var EigenPodMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_ethPOS\",\"type\":\"address\",\"internalType\":\"contractIETHPOSDeposit\"},{\"name\":\"_eigenPodManager\",\"type\":\"address\",\"internalType\":\"contractIEigenPodManager\"},{\"name\":\"_GENESIS_TIME\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"GENESIS_TIME\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"activeValidatorCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"checkpointBalanceExitedGwei\",\"inputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"currentCheckpoint\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIEigenPod.Checkpoint\",\"components\":[{\"name\":\"beaconBlockRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proofsRemaining\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"podBalanceGwei\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"balanceDeltasGwei\",\"type\":\"int128\",\"internalType\":\"int128\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"currentCheckpointTimestamp\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenPodManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenPodManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ethPOS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIETHPOSDeposit\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getParentBlockRoot\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_podOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"lastCheckpointTimestamp\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"podOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proofSubmitter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recoverTokens\",\"inputs\":[{\"name\":\"tokenList\",\"type\":\"address[]\",\"internalType\":\"contractIERC20[]\"},{\"name\":\"amountsToWithdraw\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setProofSubmitter\",\"inputs\":[{\"name\":\"newProofSubmitter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stake\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"depositDataRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"startCheckpoint\",\"inputs\":[{\"name\":\"revertIfNoBalance\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"validatorPubkeyHashToInfo\",\"inputs\":[{\"name\":\"validatorPubkeyHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIEigenPod.ValidatorInfo\",\"components\":[{\"name\":\"validatorIndex\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"restakedBalanceGwei\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"lastCheckpointedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumIEigenPod.VALIDATOR_STATUS\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"validatorPubkeyToInfo\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIEigenPod.ValidatorInfo\",\"components\":[{\"name\":\"validatorIndex\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"restakedBalanceGwei\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"lastCheckpointedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumIEigenPod.VALIDATOR_STATUS\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"validatorStatus\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumIEigenPod.VALIDATOR_STATUS\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"validatorStatus\",\"inputs\":[{\"name\":\"pubkeyHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumIEigenPod.VALIDATOR_STATUS\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyCheckpointProofs\",\"inputs\":[{\"name\":\"balanceContainerProof\",\"type\":\"tuple\",\"internalType\":\"structBeaconChainProofs.BalanceContainerProof\",\"components\":[{\"name\":\"balanceContainerRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"proofs\",\"type\":\"tuple[]\",\"internalType\":\"structBeaconChainProofs.BalanceProof[]\",\"components\":[{\"name\":\"pubkeyHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"balanceRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyStaleBalance\",\"inputs\":[{\"name\":\"beaconTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"stateRootProof\",\"type\":\"tuple\",\"internalType\":\"structBeaconChainProofs.StateRootProof\",\"components\":[{\"name\":\"beaconStateRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structBeaconChainProofs.ValidatorProof\",\"components\":[{\"name\":\"validatorFields\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"proof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyWithdrawalCredentials\",\"inputs\":[{\"name\":\"beaconTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"stateRootProof\",\"type\":\"tuple\",\"internalType\":\"structBeaconChainProofs.StateRootProof\",\"components\":[{\"name\":\"beaconStateRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"validatorIndices\",\"type\":\"uint40[]\",\"internalType\":\"uint40[]\"},{\"name\":\"validatorFieldsProofs\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"validatorFields\",\"type\":\"bytes32[][]\",\"internalType\":\"bytes32[][]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawRestakedBeaconChainETH\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amountWei\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawableRestakedExecutionLayerGwei\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"CheckpointCreated\",\"inputs\":[{\"name\":\"checkpointTimestamp\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"beaconBlockRoot\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"validatorCount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CheckpointFinalized\",\"inputs\":[{\"name\":\"checkpointTimestamp\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"totalShareDeltaWei\",\"type\":\"int256\",\"indexed\":false,\"internalType\":\"int256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EigenPodStaked\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NonBeaconChainETHReceived\",\"inputs\":[{\"name\":\"amountReceived\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProofSubmitterUpdated\",\"inputs\":[{\"name\":\"prevProofSubmitter\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newProofSubmitter\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RestakedBeaconChainETHWithdrawn\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorBalanceUpdated\",\"inputs\":[{\"name\":\"validatorIndex\",\"type\":\"uint40\",\"indexed\":false,\"internalType\":\"uint40\"},{\"name\":\"balanceTimestamp\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"newValidatorBalanceGwei\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorCheckpointed\",\"inputs\":[{\"name\":\"checkpointTimestamp\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"validatorIndex\",\"type\":\"uint40\",\"indexed\":true,\"internalType\":\"uint40\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorRestaked\",\"inputs\":[{\"name\":\"validatorIndex\",\"type\":\"uint40\",\"indexed\":false,\"internalType\":\"uint40\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorWithdrawn\",\"inputs\":[{\"name\":\"checkpointTimestamp\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"validatorIndex\",\"type\":\"uint40\",\"indexed\":true,\"internalType\":\"uint40\"}],\"anonymous\":false}]",
}

// EigenPodABI is the input ABI used to generate the binding from.
// Deprecated: Use EigenPodMetaData.ABI instead.
var EigenPodABI = EigenPodMetaData.ABI

// EigenPod is an auto generated Go binding around an Ethereum contract.
type EigenPod struct {
	EigenPodCaller     // Read-only binding to the contract
	EigenPodTransactor // Write-only binding to the contract
	EigenPodFilterer   // Log filterer for contract events
}

// EigenPodCaller is an auto generated read-only Go binding around an Ethereum contract.
type EigenPodCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EigenPodTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EigenPodTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EigenPodFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EigenPodFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EigenPodSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EigenPodSession struct {
	Contract     *EigenPod         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EigenPodCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EigenPodCallerSession struct {
	Contract *EigenPodCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// EigenPodTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EigenPodTransactorSession struct {
	Contract     *EigenPodTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// EigenPodRaw is an auto generated low-level Go binding around an Ethereum contract.
type EigenPodRaw struct {
	Contract *EigenPod // Generic contract binding to access the raw methods on
}

// EigenPodCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EigenPodCallerRaw struct {
	Contract *EigenPodCaller // Generic read-only contract binding to access the raw methods on
}

// EigenPodTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EigenPodTransactorRaw struct {
	Contract *EigenPodTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEigenPod creates a new instance of EigenPod, bound to a specific deployed contract.
func NewEigenPod(address common.Address, backend bind.ContractBackend) (*EigenPod, error) {
	contract, err := bindEigenPod(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EigenPod{EigenPodCaller: EigenPodCaller{contract: contract}, EigenPodTransactor: EigenPodTransactor{contract: contract}, EigenPodFilterer: EigenPodFilterer{contract: contract}}, nil
}

// NewEigenPodCaller creates a new read-only instance of EigenPod, bound to a specific deployed contract.
func NewEigenPodCaller(address common.Address, caller bind.ContractCaller) (*EigenPodCaller, error) {
	contract, err := bindEigenPod(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EigenPodCaller{contract: contract}, nil
}

// NewEigenPodTransactor creates a new write-only instance of EigenPod, bound to a specific deployed contract.
func NewEigenPodTransactor(address common.Address, transactor bind.ContractTransactor) (*EigenPodTransactor, error) {
	contract, err := bindEigenPod(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EigenPodTransactor{contract: contract}, nil
}

// NewEigenPodFilterer creates a new log filterer instance of EigenPod, bound to a specific deployed contract.
func NewEigenPodFilterer(address common.Address, filterer bind.ContractFilterer) (*EigenPodFilterer, error) {
	contract, err := bindEigenPod(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EigenPodFilterer{contract: contract}, nil
}

// bindEigenPod binds a generic wrapper to an already deployed contract.
func bindEigenPod(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EigenPodMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EigenPod *EigenPodRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EigenPod.Contract.EigenPodCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EigenPod *EigenPodRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EigenPod.Contract.EigenPodTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EigenPod *EigenPodRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EigenPod.Contract.EigenPodTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EigenPod *EigenPodCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EigenPod.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EigenPod *EigenPodTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EigenPod.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EigenPod *EigenPodTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EigenPod.Contract.contract.Transact(opts, method, params...)
}

// GENESISTIME is a free data retrieval call binding the contract method 0xf2882461.
//
// Solidity: function GENESIS_TIME() view returns(uint64)
func (_EigenPod *EigenPodCaller) GENESISTIME(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _EigenPod.contract.Call(opts, &out, "GENESIS_TIME")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GENESISTIME is a free data retrieval call binding the contract method 0xf2882461.
//
// Solidity: function GENESIS_TIME() view returns(uint64)
func (_EigenPod *EigenPodSession) GENESISTIME() (uint64, error) {
	return _EigenPod.Contract.GENESISTIME(&_EigenPod.CallOpts)
}

// GENESISTIME is a free data retrieval call binding the contract method 0xf2882461.
//
// Solidity: function GENESIS_TIME() view returns(uint64)
func (_EigenPod *EigenPodCallerSession) GENESISTIME() (uint64, error) {
	return _EigenPod.Contract.GENESISTIME(&_EigenPod.CallOpts)
}

// ActiveValidatorCount is a free data retrieval call binding the contract method 0x2340e8d3.
//
// Solidity: function activeValidatorCount() view returns(uint256)
func (_EigenPod *EigenPodCaller) ActiveValidatorCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EigenPod.contract.Call(opts, &out, "activeValidatorCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ActiveValidatorCount is a free data retrieval call binding the contract method 0x2340e8d3.
//
// Solidity: function activeValidatorCount() view returns(uint256)
func (_EigenPod *EigenPodSession) ActiveValidatorCount() (*big.Int, error) {
	return _EigenPod.Contract.ActiveValidatorCount(&_EigenPod.CallOpts)
}

// ActiveValidatorCount is a free data retrieval call binding the contract method 0x2340e8d3.
//
// Solidity: function activeValidatorCount() view returns(uint256)
func (_EigenPod *EigenPodCallerSession) ActiveValidatorCount() (*big.Int, error) {
	return _EigenPod.Contract.ActiveValidatorCount(&_EigenPod.CallOpts)
}

// CheckpointBalanceExitedGwei is a free data retrieval call binding the contract method 0x52396a59.
//
// Solidity: function checkpointBalanceExitedGwei(uint64 ) view returns(uint64)
func (_EigenPod *EigenPodCaller) CheckpointBalanceExitedGwei(opts *bind.CallOpts, arg0 uint64) (uint64, error) {
	var out []interface{}
	err := _EigenPod.contract.Call(opts, &out, "checkpointBalanceExitedGwei", arg0)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// CheckpointBalanceExitedGwei is a free data retrieval call binding the contract method 0x52396a59.
//
// Solidity: function checkpointBalanceExitedGwei(uint64 ) view returns(uint64)
func (_EigenPod *EigenPodSession) CheckpointBalanceExitedGwei(arg0 uint64) (uint64, error) {
	return _EigenPod.Contract.CheckpointBalanceExitedGwei(&_EigenPod.CallOpts, arg0)
}

// CheckpointBalanceExitedGwei is a free data retrieval call binding the contract method 0x52396a59.
//
// Solidity: function checkpointBalanceExitedGwei(uint64 ) view returns(uint64)
func (_EigenPod *EigenPodCallerSession) CheckpointBalanceExitedGwei(arg0 uint64) (uint64, error) {
	return _EigenPod.Contract.CheckpointBalanceExitedGwei(&_EigenPod.CallOpts, arg0)
}

// CurrentCheckpoint is a free data retrieval call binding the contract method 0x47d28372.
//
// Solidity: function currentCheckpoint() view returns((bytes32,uint24,uint64,int128))
func (_EigenPod *EigenPodCaller) CurrentCheckpoint(opts *bind.CallOpts) (IEigenPodCheckpoint, error) {
	var out []interface{}
	err := _EigenPod.contract.Call(opts, &out, "currentCheckpoint")

	if err != nil {
		return *new(IEigenPodCheckpoint), err
	}

	out0 := *abi.ConvertType(out[0], new(IEigenPodCheckpoint)).(*IEigenPodCheckpoint)

	return out0, err

}

// CurrentCheckpoint is a free data retrieval call binding the contract method 0x47d28372.
//
// Solidity: function currentCheckpoint() view returns((bytes32,uint24,uint64,int128))
func (_EigenPod *EigenPodSession) CurrentCheckpoint() (IEigenPodCheckpoint, error) {
	return _EigenPod.Contract.CurrentCheckpoint(&_EigenPod.CallOpts)
}

// CurrentCheckpoint is a free data retrieval call binding the contract method 0x47d28372.
//
// Solidity: function currentCheckpoint() view returns((bytes32,uint24,uint64,int128))
func (_EigenPod *EigenPodCallerSession) CurrentCheckpoint() (IEigenPodCheckpoint, error) {
	return _EigenPod.Contract.CurrentCheckpoint(&_EigenPod.CallOpts)
}

// CurrentCheckpointTimestamp is a free data retrieval call binding the contract method 0x42ecff2a.
//
// Solidity: function currentCheckpointTimestamp() view returns(uint64)
func (_EigenPod *EigenPodCaller) CurrentCheckpointTimestamp(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _EigenPod.contract.Call(opts, &out, "currentCheckpointTimestamp")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// CurrentCheckpointTimestamp is a free data retrieval call binding the contract method 0x42ecff2a.
//
// Solidity: function currentCheckpointTimestamp() view returns(uint64)
func (_EigenPod *EigenPodSession) CurrentCheckpointTimestamp() (uint64, error) {
	return _EigenPod.Contract.CurrentCheckpointTimestamp(&_EigenPod.CallOpts)
}

// CurrentCheckpointTimestamp is a free data retrieval call binding the contract method 0x42ecff2a.
//
// Solidity: function currentCheckpointTimestamp() view returns(uint64)
func (_EigenPod *EigenPodCallerSession) CurrentCheckpointTimestamp() (uint64, error) {
	return _EigenPod.Contract.CurrentCheckpointTimestamp(&_EigenPod.CallOpts)
}

// EigenPodManager is a free data retrieval call binding the contract method 0x4665bcda.
//
// Solidity: function eigenPodManager() view returns(address)
func (_EigenPod *EigenPodCaller) EigenPodManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EigenPod.contract.Call(opts, &out, "eigenPodManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenPodManager is a free data retrieval call binding the contract method 0x4665bcda.
//
// Solidity: function eigenPodManager() view returns(address)
func (_EigenPod *EigenPodSession) EigenPodManager() (common.Address, error) {
	return _EigenPod.Contract.EigenPodManager(&_EigenPod.CallOpts)
}

// EigenPodManager is a free data retrieval call binding the contract method 0x4665bcda.
//
// Solidity: function eigenPodManager() view returns(address)
func (_EigenPod *EigenPodCallerSession) EigenPodManager() (common.Address, error) {
	return _EigenPod.Contract.EigenPodManager(&_EigenPod.CallOpts)
}

// EthPOS is a free data retrieval call binding the contract method 0x74cdd798.
//
// Solidity: function ethPOS() view returns(address)
func (_EigenPod *EigenPodCaller) EthPOS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EigenPod.contract.Call(opts, &out, "ethPOS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EthPOS is a free data retrieval call binding the contract method 0x74cdd798.
//
// Solidity: function ethPOS() view returns(address)
func (_EigenPod *EigenPodSession) EthPOS() (common.Address, error) {
	return _EigenPod.Contract.EthPOS(&_EigenPod.CallOpts)
}

// EthPOS is a free data retrieval call binding the contract method 0x74cdd798.
//
// Solidity: function ethPOS() view returns(address)
func (_EigenPod *EigenPodCallerSession) EthPOS() (common.Address, error) {
	return _EigenPod.Contract.EthPOS(&_EigenPod.CallOpts)
}

// GetParentBlockRoot is a free data retrieval call binding the contract method 0x6c0d2d5a.
//
// Solidity: function getParentBlockRoot(uint64 timestamp) view returns(bytes32)
func (_EigenPod *EigenPodCaller) GetParentBlockRoot(opts *bind.CallOpts, timestamp uint64) ([32]byte, error) {
	var out []interface{}
	err := _EigenPod.contract.Call(opts, &out, "getParentBlockRoot", timestamp)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetParentBlockRoot is a free data retrieval call binding the contract method 0x6c0d2d5a.
//
// Solidity: function getParentBlockRoot(uint64 timestamp) view returns(bytes32)
func (_EigenPod *EigenPodSession) GetParentBlockRoot(timestamp uint64) ([32]byte, error) {
	return _EigenPod.Contract.GetParentBlockRoot(&_EigenPod.CallOpts, timestamp)
}

// GetParentBlockRoot is a free data retrieval call binding the contract method 0x6c0d2d5a.
//
// Solidity: function getParentBlockRoot(uint64 timestamp) view returns(bytes32)
func (_EigenPod *EigenPodCallerSession) GetParentBlockRoot(timestamp uint64) ([32]byte, error) {
	return _EigenPod.Contract.GetParentBlockRoot(&_EigenPod.CallOpts, timestamp)
}

// LastCheckpointTimestamp is a free data retrieval call binding the contract method 0xee94d67c.
//
// Solidity: function lastCheckpointTimestamp() view returns(uint64)
func (_EigenPod *EigenPodCaller) LastCheckpointTimestamp(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _EigenPod.contract.Call(opts, &out, "lastCheckpointTimestamp")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// LastCheckpointTimestamp is a free data retrieval call binding the contract method 0xee94d67c.
//
// Solidity: function lastCheckpointTimestamp() view returns(uint64)
func (_EigenPod *EigenPodSession) LastCheckpointTimestamp() (uint64, error) {
	return _EigenPod.Contract.LastCheckpointTimestamp(&_EigenPod.CallOpts)
}

// LastCheckpointTimestamp is a free data retrieval call binding the contract method 0xee94d67c.
//
// Solidity: function lastCheckpointTimestamp() view returns(uint64)
func (_EigenPod *EigenPodCallerSession) LastCheckpointTimestamp() (uint64, error) {
	return _EigenPod.Contract.LastCheckpointTimestamp(&_EigenPod.CallOpts)
}

// PodOwner is a free data retrieval call binding the contract method 0x0b18ff66.
//
// Solidity: function podOwner() view returns(address)
func (_EigenPod *EigenPodCaller) PodOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EigenPod.contract.Call(opts, &out, "podOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PodOwner is a free data retrieval call binding the contract method 0x0b18ff66.
//
// Solidity: function podOwner() view returns(address)
func (_EigenPod *EigenPodSession) PodOwner() (common.Address, error) {
	return _EigenPod.Contract.PodOwner(&_EigenPod.CallOpts)
}

// PodOwner is a free data retrieval call binding the contract method 0x0b18ff66.
//
// Solidity: function podOwner() view returns(address)
func (_EigenPod *EigenPodCallerSession) PodOwner() (common.Address, error) {
	return _EigenPod.Contract.PodOwner(&_EigenPod.CallOpts)
}

// ProofSubmitter is a free data retrieval call binding the contract method 0x58753357.
//
// Solidity: function proofSubmitter() view returns(address)
func (_EigenPod *EigenPodCaller) ProofSubmitter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EigenPod.contract.Call(opts, &out, "proofSubmitter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ProofSubmitter is a free data retrieval call binding the contract method 0x58753357.
//
// Solidity: function proofSubmitter() view returns(address)
func (_EigenPod *EigenPodSession) ProofSubmitter() (common.Address, error) {
	return _EigenPod.Contract.ProofSubmitter(&_EigenPod.CallOpts)
}

// ProofSubmitter is a free data retrieval call binding the contract method 0x58753357.
//
// Solidity: function proofSubmitter() view returns(address)
func (_EigenPod *EigenPodCallerSession) ProofSubmitter() (common.Address, error) {
	return _EigenPod.Contract.ProofSubmitter(&_EigenPod.CallOpts)
}

// ValidatorPubkeyHashToInfo is a free data retrieval call binding the contract method 0x6fcd0e53.
//
// Solidity: function validatorPubkeyHashToInfo(bytes32 validatorPubkeyHash) view returns((uint64,uint64,uint64,uint8))
func (_EigenPod *EigenPodCaller) ValidatorPubkeyHashToInfo(opts *bind.CallOpts, validatorPubkeyHash [32]byte) (IEigenPodValidatorInfo, error) {
	var out []interface{}
	err := _EigenPod.contract.Call(opts, &out, "validatorPubkeyHashToInfo", validatorPubkeyHash)

	if err != nil {
		return *new(IEigenPodValidatorInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IEigenPodValidatorInfo)).(*IEigenPodValidatorInfo)

	return out0, err

}

// ValidatorPubkeyHashToInfo is a free data retrieval call binding the contract method 0x6fcd0e53.
//
// Solidity: function validatorPubkeyHashToInfo(bytes32 validatorPubkeyHash) view returns((uint64,uint64,uint64,uint8))
func (_EigenPod *EigenPodSession) ValidatorPubkeyHashToInfo(validatorPubkeyHash [32]byte) (IEigenPodValidatorInfo, error) {
	return _EigenPod.Contract.ValidatorPubkeyHashToInfo(&_EigenPod.CallOpts, validatorPubkeyHash)
}

// ValidatorPubkeyHashToInfo is a free data retrieval call binding the contract method 0x6fcd0e53.
//
// Solidity: function validatorPubkeyHashToInfo(bytes32 validatorPubkeyHash) view returns((uint64,uint64,uint64,uint8))
func (_EigenPod *EigenPodCallerSession) ValidatorPubkeyHashToInfo(validatorPubkeyHash [32]byte) (IEigenPodValidatorInfo, error) {
	return _EigenPod.Contract.ValidatorPubkeyHashToInfo(&_EigenPod.CallOpts, validatorPubkeyHash)
}

// ValidatorPubkeyToInfo is a free data retrieval call binding the contract method 0xb522538a.
//
// Solidity: function validatorPubkeyToInfo(bytes validatorPubkey) view returns((uint64,uint64,uint64,uint8))
func (_EigenPod *EigenPodCaller) ValidatorPubkeyToInfo(opts *bind.CallOpts, validatorPubkey []byte) (IEigenPodValidatorInfo, error) {
	var out []interface{}
	err := _EigenPod.contract.Call(opts, &out, "validatorPubkeyToInfo", validatorPubkey)

	if err != nil {
		return *new(IEigenPodValidatorInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IEigenPodValidatorInfo)).(*IEigenPodValidatorInfo)

	return out0, err

}

// ValidatorPubkeyToInfo is a free data retrieval call binding the contract method 0xb522538a.
//
// Solidity: function validatorPubkeyToInfo(bytes validatorPubkey) view returns((uint64,uint64,uint64,uint8))
func (_EigenPod *EigenPodSession) ValidatorPubkeyToInfo(validatorPubkey []byte) (IEigenPodValidatorInfo, error) {
	return _EigenPod.Contract.ValidatorPubkeyToInfo(&_EigenPod.CallOpts, validatorPubkey)
}

// ValidatorPubkeyToInfo is a free data retrieval call binding the contract method 0xb522538a.
//
// Solidity: function validatorPubkeyToInfo(bytes validatorPubkey) view returns((uint64,uint64,uint64,uint8))
func (_EigenPod *EigenPodCallerSession) ValidatorPubkeyToInfo(validatorPubkey []byte) (IEigenPodValidatorInfo, error) {
	return _EigenPod.Contract.ValidatorPubkeyToInfo(&_EigenPod.CallOpts, validatorPubkey)
}

// ValidatorStatus is a free data retrieval call binding the contract method 0x58eaee79.
//
// Solidity: function validatorStatus(bytes validatorPubkey) view returns(uint8)
func (_EigenPod *EigenPodCaller) ValidatorStatus(opts *bind.CallOpts, validatorPubkey []byte) (uint8, error) {
	var out []interface{}
	err := _EigenPod.contract.Call(opts, &out, "validatorStatus", validatorPubkey)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// ValidatorStatus is a free data retrieval call binding the contract method 0x58eaee79.
//
// Solidity: function validatorStatus(bytes validatorPubkey) view returns(uint8)
func (_EigenPod *EigenPodSession) ValidatorStatus(validatorPubkey []byte) (uint8, error) {
	return _EigenPod.Contract.ValidatorStatus(&_EigenPod.CallOpts, validatorPubkey)
}

// ValidatorStatus is a free data retrieval call binding the contract method 0x58eaee79.
//
// Solidity: function validatorStatus(bytes validatorPubkey) view returns(uint8)
func (_EigenPod *EigenPodCallerSession) ValidatorStatus(validatorPubkey []byte) (uint8, error) {
	return _EigenPod.Contract.ValidatorStatus(&_EigenPod.CallOpts, validatorPubkey)
}

// ValidatorStatus0 is a free data retrieval call binding the contract method 0x7439841f.
//
// Solidity: function validatorStatus(bytes32 pubkeyHash) view returns(uint8)
func (_EigenPod *EigenPodCaller) ValidatorStatus0(opts *bind.CallOpts, pubkeyHash [32]byte) (uint8, error) {
	var out []interface{}
	err := _EigenPod.contract.Call(opts, &out, "validatorStatus0", pubkeyHash)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// ValidatorStatus0 is a free data retrieval call binding the contract method 0x7439841f.
//
// Solidity: function validatorStatus(bytes32 pubkeyHash) view returns(uint8)
func (_EigenPod *EigenPodSession) ValidatorStatus0(pubkeyHash [32]byte) (uint8, error) {
	return _EigenPod.Contract.ValidatorStatus0(&_EigenPod.CallOpts, pubkeyHash)
}

// ValidatorStatus0 is a free data retrieval call binding the contract method 0x7439841f.
//
// Solidity: function validatorStatus(bytes32 pubkeyHash) view returns(uint8)
func (_EigenPod *EigenPodCallerSession) ValidatorStatus0(pubkeyHash [32]byte) (uint8, error) {
	return _EigenPod.Contract.ValidatorStatus0(&_EigenPod.CallOpts, pubkeyHash)
}

// WithdrawableRestakedExecutionLayerGwei is a free data retrieval call binding the contract method 0x3474aa16.
//
// Solidity: function withdrawableRestakedExecutionLayerGwei() view returns(uint64)
func (_EigenPod *EigenPodCaller) WithdrawableRestakedExecutionLayerGwei(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _EigenPod.contract.Call(opts, &out, "withdrawableRestakedExecutionLayerGwei")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// WithdrawableRestakedExecutionLayerGwei is a free data retrieval call binding the contract method 0x3474aa16.
//
// Solidity: function withdrawableRestakedExecutionLayerGwei() view returns(uint64)
func (_EigenPod *EigenPodSession) WithdrawableRestakedExecutionLayerGwei() (uint64, error) {
	return _EigenPod.Contract.WithdrawableRestakedExecutionLayerGwei(&_EigenPod.CallOpts)
}

// WithdrawableRestakedExecutionLayerGwei is a free data retrieval call binding the contract method 0x3474aa16.
//
// Solidity: function withdrawableRestakedExecutionLayerGwei() view returns(uint64)
func (_EigenPod *EigenPodCallerSession) WithdrawableRestakedExecutionLayerGwei() (uint64, error) {
	return _EigenPod.Contract.WithdrawableRestakedExecutionLayerGwei(&_EigenPod.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _podOwner) returns()
func (_EigenPod *EigenPodTransactor) Initialize(opts *bind.TransactOpts, _podOwner common.Address) (*types.Transaction, error) {
	return _EigenPod.contract.Transact(opts, "initialize", _podOwner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _podOwner) returns()
func (_EigenPod *EigenPodSession) Initialize(_podOwner common.Address) (*types.Transaction, error) {
	return _EigenPod.Contract.Initialize(&_EigenPod.TransactOpts, _podOwner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _podOwner) returns()
func (_EigenPod *EigenPodTransactorSession) Initialize(_podOwner common.Address) (*types.Transaction, error) {
	return _EigenPod.Contract.Initialize(&_EigenPod.TransactOpts, _podOwner)
}

// RecoverTokens is a paid mutator transaction binding the contract method 0xdda3346c.
//
// Solidity: function recoverTokens(address[] tokenList, uint256[] amountsToWithdraw, address recipient) returns()
func (_EigenPod *EigenPodTransactor) RecoverTokens(opts *bind.TransactOpts, tokenList []common.Address, amountsToWithdraw []*big.Int, recipient common.Address) (*types.Transaction, error) {
	return _EigenPod.contract.Transact(opts, "recoverTokens", tokenList, amountsToWithdraw, recipient)
}

// RecoverTokens is a paid mutator transaction binding the contract method 0xdda3346c.
//
// Solidity: function recoverTokens(address[] tokenList, uint256[] amountsToWithdraw, address recipient) returns()
func (_EigenPod *EigenPodSession) RecoverTokens(tokenList []common.Address, amountsToWithdraw []*big.Int, recipient common.Address) (*types.Transaction, error) {
	return _EigenPod.Contract.RecoverTokens(&_EigenPod.TransactOpts, tokenList, amountsToWithdraw, recipient)
}

// RecoverTokens is a paid mutator transaction binding the contract method 0xdda3346c.
//
// Solidity: function recoverTokens(address[] tokenList, uint256[] amountsToWithdraw, address recipient) returns()
func (_EigenPod *EigenPodTransactorSession) RecoverTokens(tokenList []common.Address, amountsToWithdraw []*big.Int, recipient common.Address) (*types.Transaction, error) {
	return _EigenPod.Contract.RecoverTokens(&_EigenPod.TransactOpts, tokenList, amountsToWithdraw, recipient)
}

// SetProofSubmitter is a paid mutator transaction binding the contract method 0xd06d5587.
//
// Solidity: function setProofSubmitter(address newProofSubmitter) returns()
func (_EigenPod *EigenPodTransactor) SetProofSubmitter(opts *bind.TransactOpts, newProofSubmitter common.Address) (*types.Transaction, error) {
	return _EigenPod.contract.Transact(opts, "setProofSubmitter", newProofSubmitter)
}

// SetProofSubmitter is a paid mutator transaction binding the contract method 0xd06d5587.
//
// Solidity: function setProofSubmitter(address newProofSubmitter) returns()
func (_EigenPod *EigenPodSession) SetProofSubmitter(newProofSubmitter common.Address) (*types.Transaction, error) {
	return _EigenPod.Contract.SetProofSubmitter(&_EigenPod.TransactOpts, newProofSubmitter)
}

// SetProofSubmitter is a paid mutator transaction binding the contract method 0xd06d5587.
//
// Solidity: function setProofSubmitter(address newProofSubmitter) returns()
func (_EigenPod *EigenPodTransactorSession) SetProofSubmitter(newProofSubmitter common.Address) (*types.Transaction, error) {
	return _EigenPod.Contract.SetProofSubmitter(&_EigenPod.TransactOpts, newProofSubmitter)
}

// Stake is a paid mutator transaction binding the contract method 0x9b4e4634.
//
// Solidity: function stake(bytes pubkey, bytes signature, bytes32 depositDataRoot) payable returns()
func (_EigenPod *EigenPodTransactor) Stake(opts *bind.TransactOpts, pubkey []byte, signature []byte, depositDataRoot [32]byte) (*types.Transaction, error) {
	return _EigenPod.contract.Transact(opts, "stake", pubkey, signature, depositDataRoot)
}

// Stake is a paid mutator transaction binding the contract method 0x9b4e4634.
//
// Solidity: function stake(bytes pubkey, bytes signature, bytes32 depositDataRoot) payable returns()
func (_EigenPod *EigenPodSession) Stake(pubkey []byte, signature []byte, depositDataRoot [32]byte) (*types.Transaction, error) {
	return _EigenPod.Contract.Stake(&_EigenPod.TransactOpts, pubkey, signature, depositDataRoot)
}

// Stake is a paid mutator transaction binding the contract method 0x9b4e4634.
//
// Solidity: function stake(bytes pubkey, bytes signature, bytes32 depositDataRoot) payable returns()
func (_EigenPod *EigenPodTransactorSession) Stake(pubkey []byte, signature []byte, depositDataRoot [32]byte) (*types.Transaction, error) {
	return _EigenPod.Contract.Stake(&_EigenPod.TransactOpts, pubkey, signature, depositDataRoot)
}

// StartCheckpoint is a paid mutator transaction binding the contract method 0x88676cad.
//
// Solidity: function startCheckpoint(bool revertIfNoBalance) returns()
func (_EigenPod *EigenPodTransactor) StartCheckpoint(opts *bind.TransactOpts, revertIfNoBalance bool) (*types.Transaction, error) {
	return _EigenPod.contract.Transact(opts, "startCheckpoint", revertIfNoBalance)
}

// StartCheckpoint is a paid mutator transaction binding the contract method 0x88676cad.
//
// Solidity: function startCheckpoint(bool revertIfNoBalance) returns()
func (_EigenPod *EigenPodSession) StartCheckpoint(revertIfNoBalance bool) (*types.Transaction, error) {
	return _EigenPod.Contract.StartCheckpoint(&_EigenPod.TransactOpts, revertIfNoBalance)
}

// StartCheckpoint is a paid mutator transaction binding the contract method 0x88676cad.
//
// Solidity: function startCheckpoint(bool revertIfNoBalance) returns()
func (_EigenPod *EigenPodTransactorSession) StartCheckpoint(revertIfNoBalance bool) (*types.Transaction, error) {
	return _EigenPod.Contract.StartCheckpoint(&_EigenPod.TransactOpts, revertIfNoBalance)
}

// VerifyCheckpointProofs is a paid mutator transaction binding the contract method 0xf074ba62.
//
// Solidity: function verifyCheckpointProofs((bytes32,bytes) balanceContainerProof, (bytes32,bytes32,bytes)[] proofs) returns()
func (_EigenPod *EigenPodTransactor) VerifyCheckpointProofs(opts *bind.TransactOpts, balanceContainerProof BeaconChainProofsBalanceContainerProof, proofs []BeaconChainProofsBalanceProof) (*types.Transaction, error) {
	return _EigenPod.contract.Transact(opts, "verifyCheckpointProofs", balanceContainerProof, proofs)
}

// VerifyCheckpointProofs is a paid mutator transaction binding the contract method 0xf074ba62.
//
// Solidity: function verifyCheckpointProofs((bytes32,bytes) balanceContainerProof, (bytes32,bytes32,bytes)[] proofs) returns()
func (_EigenPod *EigenPodSession) VerifyCheckpointProofs(balanceContainerProof BeaconChainProofsBalanceContainerProof, proofs []BeaconChainProofsBalanceProof) (*types.Transaction, error) {
	return _EigenPod.Contract.VerifyCheckpointProofs(&_EigenPod.TransactOpts, balanceContainerProof, proofs)
}

// VerifyCheckpointProofs is a paid mutator transaction binding the contract method 0xf074ba62.
//
// Solidity: function verifyCheckpointProofs((bytes32,bytes) balanceContainerProof, (bytes32,bytes32,bytes)[] proofs) returns()
func (_EigenPod *EigenPodTransactorSession) VerifyCheckpointProofs(balanceContainerProof BeaconChainProofsBalanceContainerProof, proofs []BeaconChainProofsBalanceProof) (*types.Transaction, error) {
	return _EigenPod.Contract.VerifyCheckpointProofs(&_EigenPod.TransactOpts, balanceContainerProof, proofs)
}

// VerifyStaleBalance is a paid mutator transaction binding the contract method 0x039157d2.
//
// Solidity: function verifyStaleBalance(uint64 beaconTimestamp, (bytes32,bytes) stateRootProof, (bytes32[],bytes) proof) returns()
func (_EigenPod *EigenPodTransactor) VerifyStaleBalance(opts *bind.TransactOpts, beaconTimestamp uint64, stateRootProof BeaconChainProofsStateRootProof, proof BeaconChainProofsValidatorProof) (*types.Transaction, error) {
	return _EigenPod.contract.Transact(opts, "verifyStaleBalance", beaconTimestamp, stateRootProof, proof)
}

// VerifyStaleBalance is a paid mutator transaction binding the contract method 0x039157d2.
//
// Solidity: function verifyStaleBalance(uint64 beaconTimestamp, (bytes32,bytes) stateRootProof, (bytes32[],bytes) proof) returns()
func (_EigenPod *EigenPodSession) VerifyStaleBalance(beaconTimestamp uint64, stateRootProof BeaconChainProofsStateRootProof, proof BeaconChainProofsValidatorProof) (*types.Transaction, error) {
	return _EigenPod.Contract.VerifyStaleBalance(&_EigenPod.TransactOpts, beaconTimestamp, stateRootProof, proof)
}

// VerifyStaleBalance is a paid mutator transaction binding the contract method 0x039157d2.
//
// Solidity: function verifyStaleBalance(uint64 beaconTimestamp, (bytes32,bytes) stateRootProof, (bytes32[],bytes) proof) returns()
func (_EigenPod *EigenPodTransactorSession) VerifyStaleBalance(beaconTimestamp uint64, stateRootProof BeaconChainProofsStateRootProof, proof BeaconChainProofsValidatorProof) (*types.Transaction, error) {
	return _EigenPod.Contract.VerifyStaleBalance(&_EigenPod.TransactOpts, beaconTimestamp, stateRootProof, proof)
}

// VerifyWithdrawalCredentials is a paid mutator transaction binding the contract method 0x3f65cf19.
//
// Solidity: function verifyWithdrawalCredentials(uint64 beaconTimestamp, (bytes32,bytes) stateRootProof, uint40[] validatorIndices, bytes[] validatorFieldsProofs, bytes32[][] validatorFields) returns()
func (_EigenPod *EigenPodTransactor) VerifyWithdrawalCredentials(opts *bind.TransactOpts, beaconTimestamp uint64, stateRootProof BeaconChainProofsStateRootProof, validatorIndices []*big.Int, validatorFieldsProofs [][]byte, validatorFields [][][32]byte) (*types.Transaction, error) {
	return _EigenPod.contract.Transact(opts, "verifyWithdrawalCredentials", beaconTimestamp, stateRootProof, validatorIndices, validatorFieldsProofs, validatorFields)
}

// VerifyWithdrawalCredentials is a paid mutator transaction binding the contract method 0x3f65cf19.
//
// Solidity: function verifyWithdrawalCredentials(uint64 beaconTimestamp, (bytes32,bytes) stateRootProof, uint40[] validatorIndices, bytes[] validatorFieldsProofs, bytes32[][] validatorFields) returns()
func (_EigenPod *EigenPodSession) VerifyWithdrawalCredentials(beaconTimestamp uint64, stateRootProof BeaconChainProofsStateRootProof, validatorIndices []*big.Int, validatorFieldsProofs [][]byte, validatorFields [][][32]byte) (*types.Transaction, error) {
	return _EigenPod.Contract.VerifyWithdrawalCredentials(&_EigenPod.TransactOpts, beaconTimestamp, stateRootProof, validatorIndices, validatorFieldsProofs, validatorFields)
}

// VerifyWithdrawalCredentials is a paid mutator transaction binding the contract method 0x3f65cf19.
//
// Solidity: function verifyWithdrawalCredentials(uint64 beaconTimestamp, (bytes32,bytes) stateRootProof, uint40[] validatorIndices, bytes[] validatorFieldsProofs, bytes32[][] validatorFields) returns()
func (_EigenPod *EigenPodTransactorSession) VerifyWithdrawalCredentials(beaconTimestamp uint64, stateRootProof BeaconChainProofsStateRootProof, validatorIndices []*big.Int, validatorFieldsProofs [][]byte, validatorFields [][][32]byte) (*types.Transaction, error) {
	return _EigenPod.Contract.VerifyWithdrawalCredentials(&_EigenPod.TransactOpts, beaconTimestamp, stateRootProof, validatorIndices, validatorFieldsProofs, validatorFields)
}

// WithdrawRestakedBeaconChainETH is a paid mutator transaction binding the contract method 0xc4907442.
//
// Solidity: function withdrawRestakedBeaconChainETH(address recipient, uint256 amountWei) returns()
func (_EigenPod *EigenPodTransactor) WithdrawRestakedBeaconChainETH(opts *bind.TransactOpts, recipient common.Address, amountWei *big.Int) (*types.Transaction, error) {
	return _EigenPod.contract.Transact(opts, "withdrawRestakedBeaconChainETH", recipient, amountWei)
}

// WithdrawRestakedBeaconChainETH is a paid mutator transaction binding the contract method 0xc4907442.
//
// Solidity: function withdrawRestakedBeaconChainETH(address recipient, uint256 amountWei) returns()
func (_EigenPod *EigenPodSession) WithdrawRestakedBeaconChainETH(recipient common.Address, amountWei *big.Int) (*types.Transaction, error) {
	return _EigenPod.Contract.WithdrawRestakedBeaconChainETH(&_EigenPod.TransactOpts, recipient, amountWei)
}

// WithdrawRestakedBeaconChainETH is a paid mutator transaction binding the contract method 0xc4907442.
//
// Solidity: function withdrawRestakedBeaconChainETH(address recipient, uint256 amountWei) returns()
func (_EigenPod *EigenPodTransactorSession) WithdrawRestakedBeaconChainETH(recipient common.Address, amountWei *big.Int) (*types.Transaction, error) {
	return _EigenPod.Contract.WithdrawRestakedBeaconChainETH(&_EigenPod.TransactOpts, recipient, amountWei)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_EigenPod *EigenPodTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EigenPod.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_EigenPod *EigenPodSession) Receive() (*types.Transaction, error) {
	return _EigenPod.Contract.Receive(&_EigenPod.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_EigenPod *EigenPodTransactorSession) Receive() (*types.Transaction, error) {
	return _EigenPod.Contract.Receive(&_EigenPod.TransactOpts)
}

// EigenPodCheckpointCreatedIterator is returned from FilterCheckpointCreated and is used to iterate over the raw logs and unpacked data for CheckpointCreated events raised by the EigenPod contract.
type EigenPodCheckpointCreatedIterator struct {
	Event *EigenPodCheckpointCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EigenPodCheckpointCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EigenPodCheckpointCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EigenPodCheckpointCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EigenPodCheckpointCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EigenPodCheckpointCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EigenPodCheckpointCreated represents a CheckpointCreated event raised by the EigenPod contract.
type EigenPodCheckpointCreated struct {
	CheckpointTimestamp uint64
	BeaconBlockRoot     [32]byte
	ValidatorCount      *big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterCheckpointCreated is a free log retrieval operation binding the contract event 0x575796133bbed337e5b39aa49a30dc2556a91e0c6c2af4b7b886ae77ebef1076.
//
// Solidity: event CheckpointCreated(uint64 indexed checkpointTimestamp, bytes32 indexed beaconBlockRoot, uint256 validatorCount)
func (_EigenPod *EigenPodFilterer) FilterCheckpointCreated(opts *bind.FilterOpts, checkpointTimestamp []uint64, beaconBlockRoot [][32]byte) (*EigenPodCheckpointCreatedIterator, error) {

	var checkpointTimestampRule []interface{}
	for _, checkpointTimestampItem := range checkpointTimestamp {
		checkpointTimestampRule = append(checkpointTimestampRule, checkpointTimestampItem)
	}
	var beaconBlockRootRule []interface{}
	for _, beaconBlockRootItem := range beaconBlockRoot {
		beaconBlockRootRule = append(beaconBlockRootRule, beaconBlockRootItem)
	}

	logs, sub, err := _EigenPod.contract.FilterLogs(opts, "CheckpointCreated", checkpointTimestampRule, beaconBlockRootRule)
	if err != nil {
		return nil, err
	}
	return &EigenPodCheckpointCreatedIterator{contract: _EigenPod.contract, event: "CheckpointCreated", logs: logs, sub: sub}, nil
}

// WatchCheckpointCreated is a free log subscription operation binding the contract event 0x575796133bbed337e5b39aa49a30dc2556a91e0c6c2af4b7b886ae77ebef1076.
//
// Solidity: event CheckpointCreated(uint64 indexed checkpointTimestamp, bytes32 indexed beaconBlockRoot, uint256 validatorCount)
func (_EigenPod *EigenPodFilterer) WatchCheckpointCreated(opts *bind.WatchOpts, sink chan<- *EigenPodCheckpointCreated, checkpointTimestamp []uint64, beaconBlockRoot [][32]byte) (event.Subscription, error) {

	var checkpointTimestampRule []interface{}
	for _, checkpointTimestampItem := range checkpointTimestamp {
		checkpointTimestampRule = append(checkpointTimestampRule, checkpointTimestampItem)
	}
	var beaconBlockRootRule []interface{}
	for _, beaconBlockRootItem := range beaconBlockRoot {
		beaconBlockRootRule = append(beaconBlockRootRule, beaconBlockRootItem)
	}

	logs, sub, err := _EigenPod.contract.WatchLogs(opts, "CheckpointCreated", checkpointTimestampRule, beaconBlockRootRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EigenPodCheckpointCreated)
				if err := _EigenPod.contract.UnpackLog(event, "CheckpointCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCheckpointCreated is a log parse operation binding the contract event 0x575796133bbed337e5b39aa49a30dc2556a91e0c6c2af4b7b886ae77ebef1076.
//
// Solidity: event CheckpointCreated(uint64 indexed checkpointTimestamp, bytes32 indexed beaconBlockRoot, uint256 validatorCount)
func (_EigenPod *EigenPodFilterer) ParseCheckpointCreated(log types.Log) (*EigenPodCheckpointCreated, error) {
	event := new(EigenPodCheckpointCreated)
	if err := _EigenPod.contract.UnpackLog(event, "CheckpointCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EigenPodCheckpointFinalizedIterator is returned from FilterCheckpointFinalized and is used to iterate over the raw logs and unpacked data for CheckpointFinalized events raised by the EigenPod contract.
type EigenPodCheckpointFinalizedIterator struct {
	Event *EigenPodCheckpointFinalized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EigenPodCheckpointFinalizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EigenPodCheckpointFinalized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EigenPodCheckpointFinalized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EigenPodCheckpointFinalizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EigenPodCheckpointFinalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EigenPodCheckpointFinalized represents a CheckpointFinalized event raised by the EigenPod contract.
type EigenPodCheckpointFinalized struct {
	CheckpointTimestamp uint64
	TotalShareDeltaWei  *big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterCheckpointFinalized is a free log retrieval operation binding the contract event 0x525408c201bc1576eb44116f6478f1c2a54775b19a043bcfdc708364f74f8e44.
//
// Solidity: event CheckpointFinalized(uint64 indexed checkpointTimestamp, int256 totalShareDeltaWei)
func (_EigenPod *EigenPodFilterer) FilterCheckpointFinalized(opts *bind.FilterOpts, checkpointTimestamp []uint64) (*EigenPodCheckpointFinalizedIterator, error) {

	var checkpointTimestampRule []interface{}
	for _, checkpointTimestampItem := range checkpointTimestamp {
		checkpointTimestampRule = append(checkpointTimestampRule, checkpointTimestampItem)
	}

	logs, sub, err := _EigenPod.contract.FilterLogs(opts, "CheckpointFinalized", checkpointTimestampRule)
	if err != nil {
		return nil, err
	}
	return &EigenPodCheckpointFinalizedIterator{contract: _EigenPod.contract, event: "CheckpointFinalized", logs: logs, sub: sub}, nil
}

// WatchCheckpointFinalized is a free log subscription operation binding the contract event 0x525408c201bc1576eb44116f6478f1c2a54775b19a043bcfdc708364f74f8e44.
//
// Solidity: event CheckpointFinalized(uint64 indexed checkpointTimestamp, int256 totalShareDeltaWei)
func (_EigenPod *EigenPodFilterer) WatchCheckpointFinalized(opts *bind.WatchOpts, sink chan<- *EigenPodCheckpointFinalized, checkpointTimestamp []uint64) (event.Subscription, error) {

	var checkpointTimestampRule []interface{}
	for _, checkpointTimestampItem := range checkpointTimestamp {
		checkpointTimestampRule = append(checkpointTimestampRule, checkpointTimestampItem)
	}

	logs, sub, err := _EigenPod.contract.WatchLogs(opts, "CheckpointFinalized", checkpointTimestampRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EigenPodCheckpointFinalized)
				if err := _EigenPod.contract.UnpackLog(event, "CheckpointFinalized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCheckpointFinalized is a log parse operation binding the contract event 0x525408c201bc1576eb44116f6478f1c2a54775b19a043bcfdc708364f74f8e44.
//
// Solidity: event CheckpointFinalized(uint64 indexed checkpointTimestamp, int256 totalShareDeltaWei)
func (_EigenPod *EigenPodFilterer) ParseCheckpointFinalized(log types.Log) (*EigenPodCheckpointFinalized, error) {
	event := new(EigenPodCheckpointFinalized)
	if err := _EigenPod.contract.UnpackLog(event, "CheckpointFinalized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EigenPodEigenPodStakedIterator is returned from FilterEigenPodStaked and is used to iterate over the raw logs and unpacked data for EigenPodStaked events raised by the EigenPod contract.
type EigenPodEigenPodStakedIterator struct {
	Event *EigenPodEigenPodStaked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EigenPodEigenPodStakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EigenPodEigenPodStaked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EigenPodEigenPodStaked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EigenPodEigenPodStakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EigenPodEigenPodStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EigenPodEigenPodStaked represents a EigenPodStaked event raised by the EigenPod contract.
type EigenPodEigenPodStaked struct {
	Pubkey []byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEigenPodStaked is a free log retrieval operation binding the contract event 0x606865b7934a25d4aed43f6cdb426403353fa4b3009c4d228407474581b01e23.
//
// Solidity: event EigenPodStaked(bytes pubkey)
func (_EigenPod *EigenPodFilterer) FilterEigenPodStaked(opts *bind.FilterOpts) (*EigenPodEigenPodStakedIterator, error) {

	logs, sub, err := _EigenPod.contract.FilterLogs(opts, "EigenPodStaked")
	if err != nil {
		return nil, err
	}
	return &EigenPodEigenPodStakedIterator{contract: _EigenPod.contract, event: "EigenPodStaked", logs: logs, sub: sub}, nil
}

// WatchEigenPodStaked is a free log subscription operation binding the contract event 0x606865b7934a25d4aed43f6cdb426403353fa4b3009c4d228407474581b01e23.
//
// Solidity: event EigenPodStaked(bytes pubkey)
func (_EigenPod *EigenPodFilterer) WatchEigenPodStaked(opts *bind.WatchOpts, sink chan<- *EigenPodEigenPodStaked) (event.Subscription, error) {

	logs, sub, err := _EigenPod.contract.WatchLogs(opts, "EigenPodStaked")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EigenPodEigenPodStaked)
				if err := _EigenPod.contract.UnpackLog(event, "EigenPodStaked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseEigenPodStaked is a log parse operation binding the contract event 0x606865b7934a25d4aed43f6cdb426403353fa4b3009c4d228407474581b01e23.
//
// Solidity: event EigenPodStaked(bytes pubkey)
func (_EigenPod *EigenPodFilterer) ParseEigenPodStaked(log types.Log) (*EigenPodEigenPodStaked, error) {
	event := new(EigenPodEigenPodStaked)
	if err := _EigenPod.contract.UnpackLog(event, "EigenPodStaked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EigenPodInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the EigenPod contract.
type EigenPodInitializedIterator struct {
	Event *EigenPodInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EigenPodInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EigenPodInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EigenPodInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EigenPodInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EigenPodInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EigenPodInitialized represents a Initialized event raised by the EigenPod contract.
type EigenPodInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_EigenPod *EigenPodFilterer) FilterInitialized(opts *bind.FilterOpts) (*EigenPodInitializedIterator, error) {

	logs, sub, err := _EigenPod.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &EigenPodInitializedIterator{contract: _EigenPod.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_EigenPod *EigenPodFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *EigenPodInitialized) (event.Subscription, error) {

	logs, sub, err := _EigenPod.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EigenPodInitialized)
				if err := _EigenPod.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_EigenPod *EigenPodFilterer) ParseInitialized(log types.Log) (*EigenPodInitialized, error) {
	event := new(EigenPodInitialized)
	if err := _EigenPod.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EigenPodNonBeaconChainETHReceivedIterator is returned from FilterNonBeaconChainETHReceived and is used to iterate over the raw logs and unpacked data for NonBeaconChainETHReceived events raised by the EigenPod contract.
type EigenPodNonBeaconChainETHReceivedIterator struct {
	Event *EigenPodNonBeaconChainETHReceived // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EigenPodNonBeaconChainETHReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EigenPodNonBeaconChainETHReceived)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EigenPodNonBeaconChainETHReceived)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EigenPodNonBeaconChainETHReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EigenPodNonBeaconChainETHReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EigenPodNonBeaconChainETHReceived represents a NonBeaconChainETHReceived event raised by the EigenPod contract.
type EigenPodNonBeaconChainETHReceived struct {
	AmountReceived *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterNonBeaconChainETHReceived is a free log retrieval operation binding the contract event 0x6fdd3dbdb173299608c0aa9f368735857c8842b581f8389238bf05bd04b3bf49.
//
// Solidity: event NonBeaconChainETHReceived(uint256 amountReceived)
func (_EigenPod *EigenPodFilterer) FilterNonBeaconChainETHReceived(opts *bind.FilterOpts) (*EigenPodNonBeaconChainETHReceivedIterator, error) {

	logs, sub, err := _EigenPod.contract.FilterLogs(opts, "NonBeaconChainETHReceived")
	if err != nil {
		return nil, err
	}
	return &EigenPodNonBeaconChainETHReceivedIterator{contract: _EigenPod.contract, event: "NonBeaconChainETHReceived", logs: logs, sub: sub}, nil
}

// WatchNonBeaconChainETHReceived is a free log subscription operation binding the contract event 0x6fdd3dbdb173299608c0aa9f368735857c8842b581f8389238bf05bd04b3bf49.
//
// Solidity: event NonBeaconChainETHReceived(uint256 amountReceived)
func (_EigenPod *EigenPodFilterer) WatchNonBeaconChainETHReceived(opts *bind.WatchOpts, sink chan<- *EigenPodNonBeaconChainETHReceived) (event.Subscription, error) {

	logs, sub, err := _EigenPod.contract.WatchLogs(opts, "NonBeaconChainETHReceived")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EigenPodNonBeaconChainETHReceived)
				if err := _EigenPod.contract.UnpackLog(event, "NonBeaconChainETHReceived", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNonBeaconChainETHReceived is a log parse operation binding the contract event 0x6fdd3dbdb173299608c0aa9f368735857c8842b581f8389238bf05bd04b3bf49.
//
// Solidity: event NonBeaconChainETHReceived(uint256 amountReceived)
func (_EigenPod *EigenPodFilterer) ParseNonBeaconChainETHReceived(log types.Log) (*EigenPodNonBeaconChainETHReceived, error) {
	event := new(EigenPodNonBeaconChainETHReceived)
	if err := _EigenPod.contract.UnpackLog(event, "NonBeaconChainETHReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EigenPodProofSubmitterUpdatedIterator is returned from FilterProofSubmitterUpdated and is used to iterate over the raw logs and unpacked data for ProofSubmitterUpdated events raised by the EigenPod contract.
type EigenPodProofSubmitterUpdatedIterator struct {
	Event *EigenPodProofSubmitterUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EigenPodProofSubmitterUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EigenPodProofSubmitterUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EigenPodProofSubmitterUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EigenPodProofSubmitterUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EigenPodProofSubmitterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EigenPodProofSubmitterUpdated represents a ProofSubmitterUpdated event raised by the EigenPod contract.
type EigenPodProofSubmitterUpdated struct {
	PrevProofSubmitter common.Address
	NewProofSubmitter  common.Address
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterProofSubmitterUpdated is a free log retrieval operation binding the contract event 0xfb8129080a19d34dceac04ba253fc50304dc86c729bd63cdca4a969ad19a5eac.
//
// Solidity: event ProofSubmitterUpdated(address prevProofSubmitter, address newProofSubmitter)
func (_EigenPod *EigenPodFilterer) FilterProofSubmitterUpdated(opts *bind.FilterOpts) (*EigenPodProofSubmitterUpdatedIterator, error) {

	logs, sub, err := _EigenPod.contract.FilterLogs(opts, "ProofSubmitterUpdated")
	if err != nil {
		return nil, err
	}
	return &EigenPodProofSubmitterUpdatedIterator{contract: _EigenPod.contract, event: "ProofSubmitterUpdated", logs: logs, sub: sub}, nil
}

// WatchProofSubmitterUpdated is a free log subscription operation binding the contract event 0xfb8129080a19d34dceac04ba253fc50304dc86c729bd63cdca4a969ad19a5eac.
//
// Solidity: event ProofSubmitterUpdated(address prevProofSubmitter, address newProofSubmitter)
func (_EigenPod *EigenPodFilterer) WatchProofSubmitterUpdated(opts *bind.WatchOpts, sink chan<- *EigenPodProofSubmitterUpdated) (event.Subscription, error) {

	logs, sub, err := _EigenPod.contract.WatchLogs(opts, "ProofSubmitterUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EigenPodProofSubmitterUpdated)
				if err := _EigenPod.contract.UnpackLog(event, "ProofSubmitterUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseProofSubmitterUpdated is a log parse operation binding the contract event 0xfb8129080a19d34dceac04ba253fc50304dc86c729bd63cdca4a969ad19a5eac.
//
// Solidity: event ProofSubmitterUpdated(address prevProofSubmitter, address newProofSubmitter)
func (_EigenPod *EigenPodFilterer) ParseProofSubmitterUpdated(log types.Log) (*EigenPodProofSubmitterUpdated, error) {
	event := new(EigenPodProofSubmitterUpdated)
	if err := _EigenPod.contract.UnpackLog(event, "ProofSubmitterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EigenPodRestakedBeaconChainETHWithdrawnIterator is returned from FilterRestakedBeaconChainETHWithdrawn and is used to iterate over the raw logs and unpacked data for RestakedBeaconChainETHWithdrawn events raised by the EigenPod contract.
type EigenPodRestakedBeaconChainETHWithdrawnIterator struct {
	Event *EigenPodRestakedBeaconChainETHWithdrawn // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EigenPodRestakedBeaconChainETHWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EigenPodRestakedBeaconChainETHWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EigenPodRestakedBeaconChainETHWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EigenPodRestakedBeaconChainETHWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EigenPodRestakedBeaconChainETHWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EigenPodRestakedBeaconChainETHWithdrawn represents a RestakedBeaconChainETHWithdrawn event raised by the EigenPod contract.
type EigenPodRestakedBeaconChainETHWithdrawn struct {
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRestakedBeaconChainETHWithdrawn is a free log retrieval operation binding the contract event 0x8947fd2ce07ef9cc302c4e8f0461015615d91ce851564839e91cc804c2f49d8e.
//
// Solidity: event RestakedBeaconChainETHWithdrawn(address indexed recipient, uint256 amount)
func (_EigenPod *EigenPodFilterer) FilterRestakedBeaconChainETHWithdrawn(opts *bind.FilterOpts, recipient []common.Address) (*EigenPodRestakedBeaconChainETHWithdrawnIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _EigenPod.contract.FilterLogs(opts, "RestakedBeaconChainETHWithdrawn", recipientRule)
	if err != nil {
		return nil, err
	}
	return &EigenPodRestakedBeaconChainETHWithdrawnIterator{contract: _EigenPod.contract, event: "RestakedBeaconChainETHWithdrawn", logs: logs, sub: sub}, nil
}

// WatchRestakedBeaconChainETHWithdrawn is a free log subscription operation binding the contract event 0x8947fd2ce07ef9cc302c4e8f0461015615d91ce851564839e91cc804c2f49d8e.
//
// Solidity: event RestakedBeaconChainETHWithdrawn(address indexed recipient, uint256 amount)
func (_EigenPod *EigenPodFilterer) WatchRestakedBeaconChainETHWithdrawn(opts *bind.WatchOpts, sink chan<- *EigenPodRestakedBeaconChainETHWithdrawn, recipient []common.Address) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _EigenPod.contract.WatchLogs(opts, "RestakedBeaconChainETHWithdrawn", recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EigenPodRestakedBeaconChainETHWithdrawn)
				if err := _EigenPod.contract.UnpackLog(event, "RestakedBeaconChainETHWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRestakedBeaconChainETHWithdrawn is a log parse operation binding the contract event 0x8947fd2ce07ef9cc302c4e8f0461015615d91ce851564839e91cc804c2f49d8e.
//
// Solidity: event RestakedBeaconChainETHWithdrawn(address indexed recipient, uint256 amount)
func (_EigenPod *EigenPodFilterer) ParseRestakedBeaconChainETHWithdrawn(log types.Log) (*EigenPodRestakedBeaconChainETHWithdrawn, error) {
	event := new(EigenPodRestakedBeaconChainETHWithdrawn)
	if err := _EigenPod.contract.UnpackLog(event, "RestakedBeaconChainETHWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EigenPodValidatorBalanceUpdatedIterator is returned from FilterValidatorBalanceUpdated and is used to iterate over the raw logs and unpacked data for ValidatorBalanceUpdated events raised by the EigenPod contract.
type EigenPodValidatorBalanceUpdatedIterator struct {
	Event *EigenPodValidatorBalanceUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EigenPodValidatorBalanceUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EigenPodValidatorBalanceUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EigenPodValidatorBalanceUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EigenPodValidatorBalanceUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EigenPodValidatorBalanceUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EigenPodValidatorBalanceUpdated represents a ValidatorBalanceUpdated event raised by the EigenPod contract.
type EigenPodValidatorBalanceUpdated struct {
	ValidatorIndex          *big.Int
	BalanceTimestamp        uint64
	NewValidatorBalanceGwei uint64
	Raw                     types.Log // Blockchain specific contextual infos
}

// FilterValidatorBalanceUpdated is a free log retrieval operation binding the contract event 0x0e5fac175b83177cc047381e030d8fb3b42b37bd1c025e22c280facad62c32df.
//
// Solidity: event ValidatorBalanceUpdated(uint40 validatorIndex, uint64 balanceTimestamp, uint64 newValidatorBalanceGwei)
func (_EigenPod *EigenPodFilterer) FilterValidatorBalanceUpdated(opts *bind.FilterOpts) (*EigenPodValidatorBalanceUpdatedIterator, error) {

	logs, sub, err := _EigenPod.contract.FilterLogs(opts, "ValidatorBalanceUpdated")
	if err != nil {
		return nil, err
	}
	return &EigenPodValidatorBalanceUpdatedIterator{contract: _EigenPod.contract, event: "ValidatorBalanceUpdated", logs: logs, sub: sub}, nil
}

// WatchValidatorBalanceUpdated is a free log subscription operation binding the contract event 0x0e5fac175b83177cc047381e030d8fb3b42b37bd1c025e22c280facad62c32df.
//
// Solidity: event ValidatorBalanceUpdated(uint40 validatorIndex, uint64 balanceTimestamp, uint64 newValidatorBalanceGwei)
func (_EigenPod *EigenPodFilterer) WatchValidatorBalanceUpdated(opts *bind.WatchOpts, sink chan<- *EigenPodValidatorBalanceUpdated) (event.Subscription, error) {

	logs, sub, err := _EigenPod.contract.WatchLogs(opts, "ValidatorBalanceUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EigenPodValidatorBalanceUpdated)
				if err := _EigenPod.contract.UnpackLog(event, "ValidatorBalanceUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseValidatorBalanceUpdated is a log parse operation binding the contract event 0x0e5fac175b83177cc047381e030d8fb3b42b37bd1c025e22c280facad62c32df.
//
// Solidity: event ValidatorBalanceUpdated(uint40 validatorIndex, uint64 balanceTimestamp, uint64 newValidatorBalanceGwei)
func (_EigenPod *EigenPodFilterer) ParseValidatorBalanceUpdated(log types.Log) (*EigenPodValidatorBalanceUpdated, error) {
	event := new(EigenPodValidatorBalanceUpdated)
	if err := _EigenPod.contract.UnpackLog(event, "ValidatorBalanceUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EigenPodValidatorCheckpointedIterator is returned from FilterValidatorCheckpointed and is used to iterate over the raw logs and unpacked data for ValidatorCheckpointed events raised by the EigenPod contract.
type EigenPodValidatorCheckpointedIterator struct {
	Event *EigenPodValidatorCheckpointed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EigenPodValidatorCheckpointedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EigenPodValidatorCheckpointed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EigenPodValidatorCheckpointed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EigenPodValidatorCheckpointedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EigenPodValidatorCheckpointedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EigenPodValidatorCheckpointed represents a ValidatorCheckpointed event raised by the EigenPod contract.
type EigenPodValidatorCheckpointed struct {
	CheckpointTimestamp uint64
	ValidatorIndex      *big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterValidatorCheckpointed is a free log retrieval operation binding the contract event 0xa91c59033c3423e18b54d0acecebb4972f9ea95aedf5f4cae3b677b02eaf3a3f.
//
// Solidity: event ValidatorCheckpointed(uint64 indexed checkpointTimestamp, uint40 indexed validatorIndex)
func (_EigenPod *EigenPodFilterer) FilterValidatorCheckpointed(opts *bind.FilterOpts, checkpointTimestamp []uint64, validatorIndex []*big.Int) (*EigenPodValidatorCheckpointedIterator, error) {

	var checkpointTimestampRule []interface{}
	for _, checkpointTimestampItem := range checkpointTimestamp {
		checkpointTimestampRule = append(checkpointTimestampRule, checkpointTimestampItem)
	}
	var validatorIndexRule []interface{}
	for _, validatorIndexItem := range validatorIndex {
		validatorIndexRule = append(validatorIndexRule, validatorIndexItem)
	}

	logs, sub, err := _EigenPod.contract.FilterLogs(opts, "ValidatorCheckpointed", checkpointTimestampRule, validatorIndexRule)
	if err != nil {
		return nil, err
	}
	return &EigenPodValidatorCheckpointedIterator{contract: _EigenPod.contract, event: "ValidatorCheckpointed", logs: logs, sub: sub}, nil
}

// WatchValidatorCheckpointed is a free log subscription operation binding the contract event 0xa91c59033c3423e18b54d0acecebb4972f9ea95aedf5f4cae3b677b02eaf3a3f.
//
// Solidity: event ValidatorCheckpointed(uint64 indexed checkpointTimestamp, uint40 indexed validatorIndex)
func (_EigenPod *EigenPodFilterer) WatchValidatorCheckpointed(opts *bind.WatchOpts, sink chan<- *EigenPodValidatorCheckpointed, checkpointTimestamp []uint64, validatorIndex []*big.Int) (event.Subscription, error) {

	var checkpointTimestampRule []interface{}
	for _, checkpointTimestampItem := range checkpointTimestamp {
		checkpointTimestampRule = append(checkpointTimestampRule, checkpointTimestampItem)
	}
	var validatorIndexRule []interface{}
	for _, validatorIndexItem := range validatorIndex {
		validatorIndexRule = append(validatorIndexRule, validatorIndexItem)
	}

	logs, sub, err := _EigenPod.contract.WatchLogs(opts, "ValidatorCheckpointed", checkpointTimestampRule, validatorIndexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EigenPodValidatorCheckpointed)
				if err := _EigenPod.contract.UnpackLog(event, "ValidatorCheckpointed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseValidatorCheckpointed is a log parse operation binding the contract event 0xa91c59033c3423e18b54d0acecebb4972f9ea95aedf5f4cae3b677b02eaf3a3f.
//
// Solidity: event ValidatorCheckpointed(uint64 indexed checkpointTimestamp, uint40 indexed validatorIndex)
func (_EigenPod *EigenPodFilterer) ParseValidatorCheckpointed(log types.Log) (*EigenPodValidatorCheckpointed, error) {
	event := new(EigenPodValidatorCheckpointed)
	if err := _EigenPod.contract.UnpackLog(event, "ValidatorCheckpointed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EigenPodValidatorRestakedIterator is returned from FilterValidatorRestaked and is used to iterate over the raw logs and unpacked data for ValidatorRestaked events raised by the EigenPod contract.
type EigenPodValidatorRestakedIterator struct {
	Event *EigenPodValidatorRestaked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EigenPodValidatorRestakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EigenPodValidatorRestaked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EigenPodValidatorRestaked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EigenPodValidatorRestakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EigenPodValidatorRestakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EigenPodValidatorRestaked represents a ValidatorRestaked event raised by the EigenPod contract.
type EigenPodValidatorRestaked struct {
	ValidatorIndex *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterValidatorRestaked is a free log retrieval operation binding the contract event 0x2d0800bbc377ea54a08c5db6a87aafff5e3e9c8fead0eda110e40e0c10441449.
//
// Solidity: event ValidatorRestaked(uint40 validatorIndex)
func (_EigenPod *EigenPodFilterer) FilterValidatorRestaked(opts *bind.FilterOpts) (*EigenPodValidatorRestakedIterator, error) {

	logs, sub, err := _EigenPod.contract.FilterLogs(opts, "ValidatorRestaked")
	if err != nil {
		return nil, err
	}
	return &EigenPodValidatorRestakedIterator{contract: _EigenPod.contract, event: "ValidatorRestaked", logs: logs, sub: sub}, nil
}

// WatchValidatorRestaked is a free log subscription operation binding the contract event 0x2d0800bbc377ea54a08c5db6a87aafff5e3e9c8fead0eda110e40e0c10441449.
//
// Solidity: event ValidatorRestaked(uint40 validatorIndex)
func (_EigenPod *EigenPodFilterer) WatchValidatorRestaked(opts *bind.WatchOpts, sink chan<- *EigenPodValidatorRestaked) (event.Subscription, error) {

	logs, sub, err := _EigenPod.contract.WatchLogs(opts, "ValidatorRestaked")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EigenPodValidatorRestaked)
				if err := _EigenPod.contract.UnpackLog(event, "ValidatorRestaked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseValidatorRestaked is a log parse operation binding the contract event 0x2d0800bbc377ea54a08c5db6a87aafff5e3e9c8fead0eda110e40e0c10441449.
//
// Solidity: event ValidatorRestaked(uint40 validatorIndex)
func (_EigenPod *EigenPodFilterer) ParseValidatorRestaked(log types.Log) (*EigenPodValidatorRestaked, error) {
	event := new(EigenPodValidatorRestaked)
	if err := _EigenPod.contract.UnpackLog(event, "ValidatorRestaked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EigenPodValidatorWithdrawnIterator is returned from FilterValidatorWithdrawn and is used to iterate over the raw logs and unpacked data for ValidatorWithdrawn events raised by the EigenPod contract.
type EigenPodValidatorWithdrawnIterator struct {
	Event *EigenPodValidatorWithdrawn // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EigenPodValidatorWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EigenPodValidatorWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EigenPodValidatorWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EigenPodValidatorWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EigenPodValidatorWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EigenPodValidatorWithdrawn represents a ValidatorWithdrawn event raised by the EigenPod contract.
type EigenPodValidatorWithdrawn struct {
	CheckpointTimestamp uint64
	ValidatorIndex      *big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterValidatorWithdrawn is a free log retrieval operation binding the contract event 0x2a02361ffa66cf2c2da4682c2355a6adcaa9f6c227b6e6563e68480f9587626a.
//
// Solidity: event ValidatorWithdrawn(uint64 indexed checkpointTimestamp, uint40 indexed validatorIndex)
func (_EigenPod *EigenPodFilterer) FilterValidatorWithdrawn(opts *bind.FilterOpts, checkpointTimestamp []uint64, validatorIndex []*big.Int) (*EigenPodValidatorWithdrawnIterator, error) {

	var checkpointTimestampRule []interface{}
	for _, checkpointTimestampItem := range checkpointTimestamp {
		checkpointTimestampRule = append(checkpointTimestampRule, checkpointTimestampItem)
	}
	var validatorIndexRule []interface{}
	for _, validatorIndexItem := range validatorIndex {
		validatorIndexRule = append(validatorIndexRule, validatorIndexItem)
	}

	logs, sub, err := _EigenPod.contract.FilterLogs(opts, "ValidatorWithdrawn", checkpointTimestampRule, validatorIndexRule)
	if err != nil {
		return nil, err
	}
	return &EigenPodValidatorWithdrawnIterator{contract: _EigenPod.contract, event: "ValidatorWithdrawn", logs: logs, sub: sub}, nil
}

// WatchValidatorWithdrawn is a free log subscription operation binding the contract event 0x2a02361ffa66cf2c2da4682c2355a6adcaa9f6c227b6e6563e68480f9587626a.
//
// Solidity: event ValidatorWithdrawn(uint64 indexed checkpointTimestamp, uint40 indexed validatorIndex)
func (_EigenPod *EigenPodFilterer) WatchValidatorWithdrawn(opts *bind.WatchOpts, sink chan<- *EigenPodValidatorWithdrawn, checkpointTimestamp []uint64, validatorIndex []*big.Int) (event.Subscription, error) {

	var checkpointTimestampRule []interface{}
	for _, checkpointTimestampItem := range checkpointTimestamp {
		checkpointTimestampRule = append(checkpointTimestampRule, checkpointTimestampItem)
	}
	var validatorIndexRule []interface{}
	for _, validatorIndexItem := range validatorIndex {
		validatorIndexRule = append(validatorIndexRule, validatorIndexItem)
	}

	logs, sub, err := _EigenPod.contract.WatchLogs(opts, "ValidatorWithdrawn", checkpointTimestampRule, validatorIndexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EigenPodValidatorWithdrawn)
				if err := _EigenPod.contract.UnpackLog(event, "ValidatorWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseValidatorWithdrawn is a log parse operation binding the contract event 0x2a02361ffa66cf2c2da4682c2355a6adcaa9f6c227b6e6563e68480f9587626a.
//
// Solidity: event ValidatorWithdrawn(uint64 indexed checkpointTimestamp, uint40 indexed validatorIndex)
func (_EigenPod *EigenPodFilterer) ParseValidatorWithdrawn(log types.Log) (*EigenPodValidatorWithdrawn, error) {
	event := new(EigenPodValidatorWithdrawn)
	if err := _EigenPod.contract.UnpackLog(event, "ValidatorWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

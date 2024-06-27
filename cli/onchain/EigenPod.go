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

// BeaconChainProofsStateRootProof is an auto generated low-level Go binding around an user-defined struct.
type BeaconChainProofsStateRootProof struct {
	BeaconStateRoot [32]byte
	Proof           []byte
}

// BeaconChainProofsWithdrawalProof is an auto generated low-level Go binding around an user-defined struct.
type BeaconChainProofsWithdrawalProof struct {
	WithdrawalProof                 []byte
	SlotProof                       []byte
	ExecutionPayloadProof           []byte
	TimestampProof                  []byte
	HistoricalSummaryBlockRootProof []byte
	BlockRootIndex                  uint64
	HistoricalSummaryIndex          uint64
	WithdrawalIndex                 uint64
	BlockRoot                       [32]byte
	SlotRoot                        [32]byte
	TimestampRoot                   [32]byte
	ExecutionPayloadRoot            [32]byte
}

// IEigenPodValidatorInfo is an auto generated low-level Go binding around an user-defined struct.
type IEigenPodValidatorInfo struct {
	ValidatorIndex                   uint64
	RestakedBalanceGwei              uint64
	MostRecentBalanceUpdateTimestamp uint64
	Status                           uint8
}

// EigenPodMetaData contains all meta data concerning the EigenPod contract.
var EigenPodMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_ethPOS\",\"type\":\"address\",\"internalType\":\"contractIETHPOSDeposit\"},{\"name\":\"_delayedWithdrawalRouter\",\"type\":\"address\",\"internalType\":\"contractIDelayedWithdrawalRouter\"},{\"name\":\"_eigenPodManager\",\"type\":\"address\",\"internalType\":\"contractIEigenPodManager\"},{\"name\":\"_MAX_RESTAKED_BALANCE_GWEI_PER_VALIDATOR\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_GENESIS_TIME\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"GENESIS_TIME\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_RESTAKED_BALANCE_GWEI_PER_VALIDATOR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"activateRestaking\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"delayedWithdrawalRouter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIDelayedWithdrawalRouter\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenPodManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenPodManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ethPOS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIETHPOSDeposit\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasRestaked\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_podOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"mostRecentWithdrawalTimestamp\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nonBeaconChainETHBalanceWei\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"podOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"provenWithdrawal\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recoverTokens\",\"inputs\":[{\"name\":\"tokenList\",\"type\":\"address[]\",\"internalType\":\"contractIERC20[]\"},{\"name\":\"amountsToWithdraw\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stake\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"depositDataRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"sumOfPartialWithdrawalsClaimedGwei\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"validatorPubkeyHashToInfo\",\"inputs\":[{\"name\":\"validatorPubkeyHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIEigenPod.ValidatorInfo\",\"components\":[{\"name\":\"validatorIndex\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"restakedBalanceGwei\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"mostRecentBalanceUpdateTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumIEigenPod.VALIDATOR_STATUS\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"validatorPubkeyToInfo\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIEigenPod.ValidatorInfo\",\"components\":[{\"name\":\"validatorIndex\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"restakedBalanceGwei\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"mostRecentBalanceUpdateTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumIEigenPod.VALIDATOR_STATUS\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"validatorStatus\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumIEigenPod.VALIDATOR_STATUS\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"validatorStatus\",\"inputs\":[{\"name\":\"pubkeyHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumIEigenPod.VALIDATOR_STATUS\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyAndProcessWithdrawals\",\"inputs\":[{\"name\":\"oracleTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"stateRootProof\",\"type\":\"tuple\",\"internalType\":\"structBeaconChainProofs.StateRootProof\",\"components\":[{\"name\":\"beaconStateRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"withdrawalProofs\",\"type\":\"tuple[]\",\"internalType\":\"structBeaconChainProofs.WithdrawalProof[]\",\"components\":[{\"name\":\"withdrawalProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"slotProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"executionPayloadProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"timestampProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"historicalSummaryBlockRootProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"blockRootIndex\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"historicalSummaryIndex\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"withdrawalIndex\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blockRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"slotRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"timestampRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"executionPayloadRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"validatorFieldsProofs\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"validatorFields\",\"type\":\"bytes32[][]\",\"internalType\":\"bytes32[][]\"},{\"name\":\"withdrawalFields\",\"type\":\"bytes32[][]\",\"internalType\":\"bytes32[][]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyBalanceUpdates\",\"inputs\":[{\"name\":\"oracleTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"validatorIndices\",\"type\":\"uint40[]\",\"internalType\":\"uint40[]\"},{\"name\":\"stateRootProof\",\"type\":\"tuple\",\"internalType\":\"structBeaconChainProofs.StateRootProof\",\"components\":[{\"name\":\"beaconStateRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"validatorFieldsProofs\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"validatorFields\",\"type\":\"bytes32[][]\",\"internalType\":\"bytes32[][]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyWithdrawalCredentials\",\"inputs\":[{\"name\":\"oracleTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"stateRootProof\",\"type\":\"tuple\",\"internalType\":\"structBeaconChainProofs.StateRootProof\",\"components\":[{\"name\":\"beaconStateRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"validatorIndices\",\"type\":\"uint40[]\",\"internalType\":\"uint40[]\"},{\"name\":\"validatorFieldsProofs\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"validatorFields\",\"type\":\"bytes32[][]\",\"internalType\":\"bytes32[][]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawBeforeRestaking\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawNonBeaconChainETHBalanceWei\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amountToWithdraw\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawRestakedBeaconChainETH\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amountWei\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawableRestakedExecutionLayerGwei\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"EigenPodStaked\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FullWithdrawalRedeemed\",\"inputs\":[{\"name\":\"validatorIndex\",\"type\":\"uint40\",\"indexed\":false,\"internalType\":\"uint40\"},{\"name\":\"withdrawalTimestamp\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"withdrawalAmountGwei\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NonBeaconChainETHReceived\",\"inputs\":[{\"name\":\"amountReceived\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NonBeaconChainETHWithdrawn\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amountWithdrawn\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PartialWithdrawalRedeemed\",\"inputs\":[{\"name\":\"validatorIndex\",\"type\":\"uint40\",\"indexed\":false,\"internalType\":\"uint40\"},{\"name\":\"withdrawalTimestamp\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"partialWithdrawalAmountGwei\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RestakedBeaconChainETHWithdrawn\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RestakingActivated\",\"inputs\":[{\"name\":\"podOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorBalanceUpdated\",\"inputs\":[{\"name\":\"validatorIndex\",\"type\":\"uint40\",\"indexed\":false,\"internalType\":\"uint40\"},{\"name\":\"balanceTimestamp\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"newValidatorBalanceGwei\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorRestaked\",\"inputs\":[{\"name\":\"validatorIndex\",\"type\":\"uint40\",\"indexed\":false,\"internalType\":\"uint40\"}],\"anonymous\":false}]",
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

// MAXRESTAKEDBALANCEGWEIPERVALIDATOR is a free data retrieval call binding the contract method 0x1d905d5c.
//
// Solidity: function MAX_RESTAKED_BALANCE_GWEI_PER_VALIDATOR() view returns(uint64)
func (_EigenPod *EigenPodCaller) MAXRESTAKEDBALANCEGWEIPERVALIDATOR(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _EigenPod.contract.Call(opts, &out, "MAX_RESTAKED_BALANCE_GWEI_PER_VALIDATOR")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// MAXRESTAKEDBALANCEGWEIPERVALIDATOR is a free data retrieval call binding the contract method 0x1d905d5c.
//
// Solidity: function MAX_RESTAKED_BALANCE_GWEI_PER_VALIDATOR() view returns(uint64)
func (_EigenPod *EigenPodSession) MAXRESTAKEDBALANCEGWEIPERVALIDATOR() (uint64, error) {
	return _EigenPod.Contract.MAXRESTAKEDBALANCEGWEIPERVALIDATOR(&_EigenPod.CallOpts)
}

// MAXRESTAKEDBALANCEGWEIPERVALIDATOR is a free data retrieval call binding the contract method 0x1d905d5c.
//
// Solidity: function MAX_RESTAKED_BALANCE_GWEI_PER_VALIDATOR() view returns(uint64)
func (_EigenPod *EigenPodCallerSession) MAXRESTAKEDBALANCEGWEIPERVALIDATOR() (uint64, error) {
	return _EigenPod.Contract.MAXRESTAKEDBALANCEGWEIPERVALIDATOR(&_EigenPod.CallOpts)
}

// DelayedWithdrawalRouter is a free data retrieval call binding the contract method 0x1a5057be.
//
// Solidity: function delayedWithdrawalRouter() view returns(address)
func (_EigenPod *EigenPodCaller) DelayedWithdrawalRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EigenPod.contract.Call(opts, &out, "delayedWithdrawalRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DelayedWithdrawalRouter is a free data retrieval call binding the contract method 0x1a5057be.
//
// Solidity: function delayedWithdrawalRouter() view returns(address)
func (_EigenPod *EigenPodSession) DelayedWithdrawalRouter() (common.Address, error) {
	return _EigenPod.Contract.DelayedWithdrawalRouter(&_EigenPod.CallOpts)
}

// DelayedWithdrawalRouter is a free data retrieval call binding the contract method 0x1a5057be.
//
// Solidity: function delayedWithdrawalRouter() view returns(address)
func (_EigenPod *EigenPodCallerSession) DelayedWithdrawalRouter() (common.Address, error) {
	return _EigenPod.Contract.DelayedWithdrawalRouter(&_EigenPod.CallOpts)
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

// HasRestaked is a free data retrieval call binding the contract method 0x3106ab53.
//
// Solidity: function hasRestaked() view returns(bool)
func (_EigenPod *EigenPodCaller) HasRestaked(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _EigenPod.contract.Call(opts, &out, "hasRestaked")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRestaked is a free data retrieval call binding the contract method 0x3106ab53.
//
// Solidity: function hasRestaked() view returns(bool)
func (_EigenPod *EigenPodSession) HasRestaked() (bool, error) {
	return _EigenPod.Contract.HasRestaked(&_EigenPod.CallOpts)
}

// HasRestaked is a free data retrieval call binding the contract method 0x3106ab53.
//
// Solidity: function hasRestaked() view returns(bool)
func (_EigenPod *EigenPodCallerSession) HasRestaked() (bool, error) {
	return _EigenPod.Contract.HasRestaked(&_EigenPod.CallOpts)
}

// MostRecentWithdrawalTimestamp is a free data retrieval call binding the contract method 0x87e0d289.
//
// Solidity: function mostRecentWithdrawalTimestamp() view returns(uint64)
func (_EigenPod *EigenPodCaller) MostRecentWithdrawalTimestamp(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _EigenPod.contract.Call(opts, &out, "mostRecentWithdrawalTimestamp")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// MostRecentWithdrawalTimestamp is a free data retrieval call binding the contract method 0x87e0d289.
//
// Solidity: function mostRecentWithdrawalTimestamp() view returns(uint64)
func (_EigenPod *EigenPodSession) MostRecentWithdrawalTimestamp() (uint64, error) {
	return _EigenPod.Contract.MostRecentWithdrawalTimestamp(&_EigenPod.CallOpts)
}

// MostRecentWithdrawalTimestamp is a free data retrieval call binding the contract method 0x87e0d289.
//
// Solidity: function mostRecentWithdrawalTimestamp() view returns(uint64)
func (_EigenPod *EigenPodCallerSession) MostRecentWithdrawalTimestamp() (uint64, error) {
	return _EigenPod.Contract.MostRecentWithdrawalTimestamp(&_EigenPod.CallOpts)
}

// NonBeaconChainETHBalanceWei is a free data retrieval call binding the contract method 0xfe80b087.
//
// Solidity: function nonBeaconChainETHBalanceWei() view returns(uint256)
func (_EigenPod *EigenPodCaller) NonBeaconChainETHBalanceWei(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EigenPod.contract.Call(opts, &out, "nonBeaconChainETHBalanceWei")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NonBeaconChainETHBalanceWei is a free data retrieval call binding the contract method 0xfe80b087.
//
// Solidity: function nonBeaconChainETHBalanceWei() view returns(uint256)
func (_EigenPod *EigenPodSession) NonBeaconChainETHBalanceWei() (*big.Int, error) {
	return _EigenPod.Contract.NonBeaconChainETHBalanceWei(&_EigenPod.CallOpts)
}

// NonBeaconChainETHBalanceWei is a free data retrieval call binding the contract method 0xfe80b087.
//
// Solidity: function nonBeaconChainETHBalanceWei() view returns(uint256)
func (_EigenPod *EigenPodCallerSession) NonBeaconChainETHBalanceWei() (*big.Int, error) {
	return _EigenPod.Contract.NonBeaconChainETHBalanceWei(&_EigenPod.CallOpts)
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

// ProvenWithdrawal is a free data retrieval call binding the contract method 0x34bea20a.
//
// Solidity: function provenWithdrawal(bytes32 , uint64 ) view returns(bool)
func (_EigenPod *EigenPodCaller) ProvenWithdrawal(opts *bind.CallOpts, arg0 [32]byte, arg1 uint64) (bool, error) {
	var out []interface{}
	err := _EigenPod.contract.Call(opts, &out, "provenWithdrawal", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ProvenWithdrawal is a free data retrieval call binding the contract method 0x34bea20a.
//
// Solidity: function provenWithdrawal(bytes32 , uint64 ) view returns(bool)
func (_EigenPod *EigenPodSession) ProvenWithdrawal(arg0 [32]byte, arg1 uint64) (bool, error) {
	return _EigenPod.Contract.ProvenWithdrawal(&_EigenPod.CallOpts, arg0, arg1)
}

// ProvenWithdrawal is a free data retrieval call binding the contract method 0x34bea20a.
//
// Solidity: function provenWithdrawal(bytes32 , uint64 ) view returns(bool)
func (_EigenPod *EigenPodCallerSession) ProvenWithdrawal(arg0 [32]byte, arg1 uint64) (bool, error) {
	return _EigenPod.Contract.ProvenWithdrawal(&_EigenPod.CallOpts, arg0, arg1)
}

// SumOfPartialWithdrawalsClaimedGwei is a free data retrieval call binding the contract method 0x5d3f65b6.
//
// Solidity: function sumOfPartialWithdrawalsClaimedGwei() view returns(uint64)
func (_EigenPod *EigenPodCaller) SumOfPartialWithdrawalsClaimedGwei(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _EigenPod.contract.Call(opts, &out, "sumOfPartialWithdrawalsClaimedGwei")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// SumOfPartialWithdrawalsClaimedGwei is a free data retrieval call binding the contract method 0x5d3f65b6.
//
// Solidity: function sumOfPartialWithdrawalsClaimedGwei() view returns(uint64)
func (_EigenPod *EigenPodSession) SumOfPartialWithdrawalsClaimedGwei() (uint64, error) {
	return _EigenPod.Contract.SumOfPartialWithdrawalsClaimedGwei(&_EigenPod.CallOpts)
}

// SumOfPartialWithdrawalsClaimedGwei is a free data retrieval call binding the contract method 0x5d3f65b6.
//
// Solidity: function sumOfPartialWithdrawalsClaimedGwei() view returns(uint64)
func (_EigenPod *EigenPodCallerSession) SumOfPartialWithdrawalsClaimedGwei() (uint64, error) {
	return _EigenPod.Contract.SumOfPartialWithdrawalsClaimedGwei(&_EigenPod.CallOpts)
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

// ActivateRestaking is a paid mutator transaction binding the contract method 0x0cd4649e.
//
// Solidity: function activateRestaking() returns()
func (_EigenPod *EigenPodTransactor) ActivateRestaking(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EigenPod.contract.Transact(opts, "activateRestaking")
}

// ActivateRestaking is a paid mutator transaction binding the contract method 0x0cd4649e.
//
// Solidity: function activateRestaking() returns()
func (_EigenPod *EigenPodSession) ActivateRestaking() (*types.Transaction, error) {
	return _EigenPod.Contract.ActivateRestaking(&_EigenPod.TransactOpts)
}

// ActivateRestaking is a paid mutator transaction binding the contract method 0x0cd4649e.
//
// Solidity: function activateRestaking() returns()
func (_EigenPod *EigenPodTransactorSession) ActivateRestaking() (*types.Transaction, error) {
	return _EigenPod.Contract.ActivateRestaking(&_EigenPod.TransactOpts)
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

// VerifyAndProcessWithdrawals is a paid mutator transaction binding the contract method 0xe251ef52.
//
// Solidity: function verifyAndProcessWithdrawals(uint64 oracleTimestamp, (bytes32,bytes) stateRootProof, (bytes,bytes,bytes,bytes,bytes,uint64,uint64,uint64,bytes32,bytes32,bytes32,bytes32)[] withdrawalProofs, bytes[] validatorFieldsProofs, bytes32[][] validatorFields, bytes32[][] withdrawalFields) returns()
func (_EigenPod *EigenPodTransactor) VerifyAndProcessWithdrawals(opts *bind.TransactOpts, oracleTimestamp uint64, stateRootProof BeaconChainProofsStateRootProof, withdrawalProofs []BeaconChainProofsWithdrawalProof, validatorFieldsProofs [][]byte, validatorFields [][][32]byte, withdrawalFields [][][32]byte) (*types.Transaction, error) {
	return _EigenPod.contract.Transact(opts, "verifyAndProcessWithdrawals", oracleTimestamp, stateRootProof, withdrawalProofs, validatorFieldsProofs, validatorFields, withdrawalFields)
}

// VerifyAndProcessWithdrawals is a paid mutator transaction binding the contract method 0xe251ef52.
//
// Solidity: function verifyAndProcessWithdrawals(uint64 oracleTimestamp, (bytes32,bytes) stateRootProof, (bytes,bytes,bytes,bytes,bytes,uint64,uint64,uint64,bytes32,bytes32,bytes32,bytes32)[] withdrawalProofs, bytes[] validatorFieldsProofs, bytes32[][] validatorFields, bytes32[][] withdrawalFields) returns()
func (_EigenPod *EigenPodSession) VerifyAndProcessWithdrawals(oracleTimestamp uint64, stateRootProof BeaconChainProofsStateRootProof, withdrawalProofs []BeaconChainProofsWithdrawalProof, validatorFieldsProofs [][]byte, validatorFields [][][32]byte, withdrawalFields [][][32]byte) (*types.Transaction, error) {
	return _EigenPod.Contract.VerifyAndProcessWithdrawals(&_EigenPod.TransactOpts, oracleTimestamp, stateRootProof, withdrawalProofs, validatorFieldsProofs, validatorFields, withdrawalFields)
}

// VerifyAndProcessWithdrawals is a paid mutator transaction binding the contract method 0xe251ef52.
//
// Solidity: function verifyAndProcessWithdrawals(uint64 oracleTimestamp, (bytes32,bytes) stateRootProof, (bytes,bytes,bytes,bytes,bytes,uint64,uint64,uint64,bytes32,bytes32,bytes32,bytes32)[] withdrawalProofs, bytes[] validatorFieldsProofs, bytes32[][] validatorFields, bytes32[][] withdrawalFields) returns()
func (_EigenPod *EigenPodTransactorSession) VerifyAndProcessWithdrawals(oracleTimestamp uint64, stateRootProof BeaconChainProofsStateRootProof, withdrawalProofs []BeaconChainProofsWithdrawalProof, validatorFieldsProofs [][]byte, validatorFields [][][32]byte, withdrawalFields [][][32]byte) (*types.Transaction, error) {
	return _EigenPod.Contract.VerifyAndProcessWithdrawals(&_EigenPod.TransactOpts, oracleTimestamp, stateRootProof, withdrawalProofs, validatorFieldsProofs, validatorFields, withdrawalFields)
}

// VerifyBalanceUpdates is a paid mutator transaction binding the contract method 0xa50600f4.
//
// Solidity: function verifyBalanceUpdates(uint64 oracleTimestamp, uint40[] validatorIndices, (bytes32,bytes) stateRootProof, bytes[] validatorFieldsProofs, bytes32[][] validatorFields) returns()
func (_EigenPod *EigenPodTransactor) VerifyBalanceUpdates(opts *bind.TransactOpts, oracleTimestamp uint64, validatorIndices []*big.Int, stateRootProof BeaconChainProofsStateRootProof, validatorFieldsProofs [][]byte, validatorFields [][][32]byte) (*types.Transaction, error) {
	return _EigenPod.contract.Transact(opts, "verifyBalanceUpdates", oracleTimestamp, validatorIndices, stateRootProof, validatorFieldsProofs, validatorFields)
}

// VerifyBalanceUpdates is a paid mutator transaction binding the contract method 0xa50600f4.
//
// Solidity: function verifyBalanceUpdates(uint64 oracleTimestamp, uint40[] validatorIndices, (bytes32,bytes) stateRootProof, bytes[] validatorFieldsProofs, bytes32[][] validatorFields) returns()
func (_EigenPod *EigenPodSession) VerifyBalanceUpdates(oracleTimestamp uint64, validatorIndices []*big.Int, stateRootProof BeaconChainProofsStateRootProof, validatorFieldsProofs [][]byte, validatorFields [][][32]byte) (*types.Transaction, error) {
	return _EigenPod.Contract.VerifyBalanceUpdates(&_EigenPod.TransactOpts, oracleTimestamp, validatorIndices, stateRootProof, validatorFieldsProofs, validatorFields)
}

// VerifyBalanceUpdates is a paid mutator transaction binding the contract method 0xa50600f4.
//
// Solidity: function verifyBalanceUpdates(uint64 oracleTimestamp, uint40[] validatorIndices, (bytes32,bytes) stateRootProof, bytes[] validatorFieldsProofs, bytes32[][] validatorFields) returns()
func (_EigenPod *EigenPodTransactorSession) VerifyBalanceUpdates(oracleTimestamp uint64, validatorIndices []*big.Int, stateRootProof BeaconChainProofsStateRootProof, validatorFieldsProofs [][]byte, validatorFields [][][32]byte) (*types.Transaction, error) {
	return _EigenPod.Contract.VerifyBalanceUpdates(&_EigenPod.TransactOpts, oracleTimestamp, validatorIndices, stateRootProof, validatorFieldsProofs, validatorFields)
}

// VerifyWithdrawalCredentials is a paid mutator transaction binding the contract method 0x3f65cf19.
//
// Solidity: function verifyWithdrawalCredentials(uint64 oracleTimestamp, (bytes32,bytes) stateRootProof, uint40[] validatorIndices, bytes[] validatorFieldsProofs, bytes32[][] validatorFields) returns()
func (_EigenPod *EigenPodTransactor) VerifyWithdrawalCredentials(opts *bind.TransactOpts, oracleTimestamp uint64, stateRootProof BeaconChainProofsStateRootProof, validatorIndices []*big.Int, validatorFieldsProofs [][]byte, validatorFields [][][32]byte) (*types.Transaction, error) {
	return _EigenPod.contract.Transact(opts, "verifyWithdrawalCredentials", oracleTimestamp, stateRootProof, validatorIndices, validatorFieldsProofs, validatorFields)
}

// VerifyWithdrawalCredentials is a paid mutator transaction binding the contract method 0x3f65cf19.
//
// Solidity: function verifyWithdrawalCredentials(uint64 oracleTimestamp, (bytes32,bytes) stateRootProof, uint40[] validatorIndices, bytes[] validatorFieldsProofs, bytes32[][] validatorFields) returns()
func (_EigenPod *EigenPodSession) VerifyWithdrawalCredentials(oracleTimestamp uint64, stateRootProof BeaconChainProofsStateRootProof, validatorIndices []*big.Int, validatorFieldsProofs [][]byte, validatorFields [][][32]byte) (*types.Transaction, error) {
	return _EigenPod.Contract.VerifyWithdrawalCredentials(&_EigenPod.TransactOpts, oracleTimestamp, stateRootProof, validatorIndices, validatorFieldsProofs, validatorFields)
}

// VerifyWithdrawalCredentials is a paid mutator transaction binding the contract method 0x3f65cf19.
//
// Solidity: function verifyWithdrawalCredentials(uint64 oracleTimestamp, (bytes32,bytes) stateRootProof, uint40[] validatorIndices, bytes[] validatorFieldsProofs, bytes32[][] validatorFields) returns()
func (_EigenPod *EigenPodTransactorSession) VerifyWithdrawalCredentials(oracleTimestamp uint64, stateRootProof BeaconChainProofsStateRootProof, validatorIndices []*big.Int, validatorFieldsProofs [][]byte, validatorFields [][][32]byte) (*types.Transaction, error) {
	return _EigenPod.Contract.VerifyWithdrawalCredentials(&_EigenPod.TransactOpts, oracleTimestamp, stateRootProof, validatorIndices, validatorFieldsProofs, validatorFields)
}

// WithdrawBeforeRestaking is a paid mutator transaction binding the contract method 0xbaa7145a.
//
// Solidity: function withdrawBeforeRestaking() returns()
func (_EigenPod *EigenPodTransactor) WithdrawBeforeRestaking(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EigenPod.contract.Transact(opts, "withdrawBeforeRestaking")
}

// WithdrawBeforeRestaking is a paid mutator transaction binding the contract method 0xbaa7145a.
//
// Solidity: function withdrawBeforeRestaking() returns()
func (_EigenPod *EigenPodSession) WithdrawBeforeRestaking() (*types.Transaction, error) {
	return _EigenPod.Contract.WithdrawBeforeRestaking(&_EigenPod.TransactOpts)
}

// WithdrawBeforeRestaking is a paid mutator transaction binding the contract method 0xbaa7145a.
//
// Solidity: function withdrawBeforeRestaking() returns()
func (_EigenPod *EigenPodTransactorSession) WithdrawBeforeRestaking() (*types.Transaction, error) {
	return _EigenPod.Contract.WithdrawBeforeRestaking(&_EigenPod.TransactOpts)
}

// WithdrawNonBeaconChainETHBalanceWei is a paid mutator transaction binding the contract method 0xe2c83445.
//
// Solidity: function withdrawNonBeaconChainETHBalanceWei(address recipient, uint256 amountToWithdraw) returns()
func (_EigenPod *EigenPodTransactor) WithdrawNonBeaconChainETHBalanceWei(opts *bind.TransactOpts, recipient common.Address, amountToWithdraw *big.Int) (*types.Transaction, error) {
	return _EigenPod.contract.Transact(opts, "withdrawNonBeaconChainETHBalanceWei", recipient, amountToWithdraw)
}

// WithdrawNonBeaconChainETHBalanceWei is a paid mutator transaction binding the contract method 0xe2c83445.
//
// Solidity: function withdrawNonBeaconChainETHBalanceWei(address recipient, uint256 amountToWithdraw) returns()
func (_EigenPod *EigenPodSession) WithdrawNonBeaconChainETHBalanceWei(recipient common.Address, amountToWithdraw *big.Int) (*types.Transaction, error) {
	return _EigenPod.Contract.WithdrawNonBeaconChainETHBalanceWei(&_EigenPod.TransactOpts, recipient, amountToWithdraw)
}

// WithdrawNonBeaconChainETHBalanceWei is a paid mutator transaction binding the contract method 0xe2c83445.
//
// Solidity: function withdrawNonBeaconChainETHBalanceWei(address recipient, uint256 amountToWithdraw) returns()
func (_EigenPod *EigenPodTransactorSession) WithdrawNonBeaconChainETHBalanceWei(recipient common.Address, amountToWithdraw *big.Int) (*types.Transaction, error) {
	return _EigenPod.Contract.WithdrawNonBeaconChainETHBalanceWei(&_EigenPod.TransactOpts, recipient, amountToWithdraw)
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

// EigenPodFullWithdrawalRedeemedIterator is returned from FilterFullWithdrawalRedeemed and is used to iterate over the raw logs and unpacked data for FullWithdrawalRedeemed events raised by the EigenPod contract.
type EigenPodFullWithdrawalRedeemedIterator struct {
	Event *EigenPodFullWithdrawalRedeemed // Event containing the contract specifics and raw log

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
func (it *EigenPodFullWithdrawalRedeemedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EigenPodFullWithdrawalRedeemed)
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
		it.Event = new(EigenPodFullWithdrawalRedeemed)
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
func (it *EigenPodFullWithdrawalRedeemedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EigenPodFullWithdrawalRedeemedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EigenPodFullWithdrawalRedeemed represents a FullWithdrawalRedeemed event raised by the EigenPod contract.
type EigenPodFullWithdrawalRedeemed struct {
	ValidatorIndex       *big.Int
	WithdrawalTimestamp  uint64
	Recipient            common.Address
	WithdrawalAmountGwei uint64
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterFullWithdrawalRedeemed is a free log retrieval operation binding the contract event 0xb76a93bb649ece524688f1a01d184e0bbebcda58eae80c28a898bec3fb5a0963.
//
// Solidity: event FullWithdrawalRedeemed(uint40 validatorIndex, uint64 withdrawalTimestamp, address indexed recipient, uint64 withdrawalAmountGwei)
func (_EigenPod *EigenPodFilterer) FilterFullWithdrawalRedeemed(opts *bind.FilterOpts, recipient []common.Address) (*EigenPodFullWithdrawalRedeemedIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _EigenPod.contract.FilterLogs(opts, "FullWithdrawalRedeemed", recipientRule)
	if err != nil {
		return nil, err
	}
	return &EigenPodFullWithdrawalRedeemedIterator{contract: _EigenPod.contract, event: "FullWithdrawalRedeemed", logs: logs, sub: sub}, nil
}

// WatchFullWithdrawalRedeemed is a free log subscription operation binding the contract event 0xb76a93bb649ece524688f1a01d184e0bbebcda58eae80c28a898bec3fb5a0963.
//
// Solidity: event FullWithdrawalRedeemed(uint40 validatorIndex, uint64 withdrawalTimestamp, address indexed recipient, uint64 withdrawalAmountGwei)
func (_EigenPod *EigenPodFilterer) WatchFullWithdrawalRedeemed(opts *bind.WatchOpts, sink chan<- *EigenPodFullWithdrawalRedeemed, recipient []common.Address) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _EigenPod.contract.WatchLogs(opts, "FullWithdrawalRedeemed", recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EigenPodFullWithdrawalRedeemed)
				if err := _EigenPod.contract.UnpackLog(event, "FullWithdrawalRedeemed", log); err != nil {
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

// ParseFullWithdrawalRedeemed is a log parse operation binding the contract event 0xb76a93bb649ece524688f1a01d184e0bbebcda58eae80c28a898bec3fb5a0963.
//
// Solidity: event FullWithdrawalRedeemed(uint40 validatorIndex, uint64 withdrawalTimestamp, address indexed recipient, uint64 withdrawalAmountGwei)
func (_EigenPod *EigenPodFilterer) ParseFullWithdrawalRedeemed(log types.Log) (*EigenPodFullWithdrawalRedeemed, error) {
	event := new(EigenPodFullWithdrawalRedeemed)
	if err := _EigenPod.contract.UnpackLog(event, "FullWithdrawalRedeemed", log); err != nil {
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

// EigenPodNonBeaconChainETHWithdrawnIterator is returned from FilterNonBeaconChainETHWithdrawn and is used to iterate over the raw logs and unpacked data for NonBeaconChainETHWithdrawn events raised by the EigenPod contract.
type EigenPodNonBeaconChainETHWithdrawnIterator struct {
	Event *EigenPodNonBeaconChainETHWithdrawn // Event containing the contract specifics and raw log

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
func (it *EigenPodNonBeaconChainETHWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EigenPodNonBeaconChainETHWithdrawn)
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
		it.Event = new(EigenPodNonBeaconChainETHWithdrawn)
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
func (it *EigenPodNonBeaconChainETHWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EigenPodNonBeaconChainETHWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EigenPodNonBeaconChainETHWithdrawn represents a NonBeaconChainETHWithdrawn event raised by the EigenPod contract.
type EigenPodNonBeaconChainETHWithdrawn struct {
	Recipient       common.Address
	AmountWithdrawn *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterNonBeaconChainETHWithdrawn is a free log retrieval operation binding the contract event 0x30420aacd028abb3c1fd03aba253ae725d6ddd52d16c9ac4cb5742cd43f53096.
//
// Solidity: event NonBeaconChainETHWithdrawn(address indexed recipient, uint256 amountWithdrawn)
func (_EigenPod *EigenPodFilterer) FilterNonBeaconChainETHWithdrawn(opts *bind.FilterOpts, recipient []common.Address) (*EigenPodNonBeaconChainETHWithdrawnIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _EigenPod.contract.FilterLogs(opts, "NonBeaconChainETHWithdrawn", recipientRule)
	if err != nil {
		return nil, err
	}
	return &EigenPodNonBeaconChainETHWithdrawnIterator{contract: _EigenPod.contract, event: "NonBeaconChainETHWithdrawn", logs: logs, sub: sub}, nil
}

// WatchNonBeaconChainETHWithdrawn is a free log subscription operation binding the contract event 0x30420aacd028abb3c1fd03aba253ae725d6ddd52d16c9ac4cb5742cd43f53096.
//
// Solidity: event NonBeaconChainETHWithdrawn(address indexed recipient, uint256 amountWithdrawn)
func (_EigenPod *EigenPodFilterer) WatchNonBeaconChainETHWithdrawn(opts *bind.WatchOpts, sink chan<- *EigenPodNonBeaconChainETHWithdrawn, recipient []common.Address) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _EigenPod.contract.WatchLogs(opts, "NonBeaconChainETHWithdrawn", recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EigenPodNonBeaconChainETHWithdrawn)
				if err := _EigenPod.contract.UnpackLog(event, "NonBeaconChainETHWithdrawn", log); err != nil {
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

// ParseNonBeaconChainETHWithdrawn is a log parse operation binding the contract event 0x30420aacd028abb3c1fd03aba253ae725d6ddd52d16c9ac4cb5742cd43f53096.
//
// Solidity: event NonBeaconChainETHWithdrawn(address indexed recipient, uint256 amountWithdrawn)
func (_EigenPod *EigenPodFilterer) ParseNonBeaconChainETHWithdrawn(log types.Log) (*EigenPodNonBeaconChainETHWithdrawn, error) {
	event := new(EigenPodNonBeaconChainETHWithdrawn)
	if err := _EigenPod.contract.UnpackLog(event, "NonBeaconChainETHWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EigenPodPartialWithdrawalRedeemedIterator is returned from FilterPartialWithdrawalRedeemed and is used to iterate over the raw logs and unpacked data for PartialWithdrawalRedeemed events raised by the EigenPod contract.
type EigenPodPartialWithdrawalRedeemedIterator struct {
	Event *EigenPodPartialWithdrawalRedeemed // Event containing the contract specifics and raw log

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
func (it *EigenPodPartialWithdrawalRedeemedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EigenPodPartialWithdrawalRedeemed)
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
		it.Event = new(EigenPodPartialWithdrawalRedeemed)
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
func (it *EigenPodPartialWithdrawalRedeemedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EigenPodPartialWithdrawalRedeemedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EigenPodPartialWithdrawalRedeemed represents a PartialWithdrawalRedeemed event raised by the EigenPod contract.
type EigenPodPartialWithdrawalRedeemed struct {
	ValidatorIndex              *big.Int
	WithdrawalTimestamp         uint64
	Recipient                   common.Address
	PartialWithdrawalAmountGwei uint64
	Raw                         types.Log // Blockchain specific contextual infos
}

// FilterPartialWithdrawalRedeemed is a free log retrieval operation binding the contract event 0x8a7335714231dbd551aaba6314f4a97a14c201e53a3e25e1140325cdf67d7a4e.
//
// Solidity: event PartialWithdrawalRedeemed(uint40 validatorIndex, uint64 withdrawalTimestamp, address indexed recipient, uint64 partialWithdrawalAmountGwei)
func (_EigenPod *EigenPodFilterer) FilterPartialWithdrawalRedeemed(opts *bind.FilterOpts, recipient []common.Address) (*EigenPodPartialWithdrawalRedeemedIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _EigenPod.contract.FilterLogs(opts, "PartialWithdrawalRedeemed", recipientRule)
	if err != nil {
		return nil, err
	}
	return &EigenPodPartialWithdrawalRedeemedIterator{contract: _EigenPod.contract, event: "PartialWithdrawalRedeemed", logs: logs, sub: sub}, nil
}

// WatchPartialWithdrawalRedeemed is a free log subscription operation binding the contract event 0x8a7335714231dbd551aaba6314f4a97a14c201e53a3e25e1140325cdf67d7a4e.
//
// Solidity: event PartialWithdrawalRedeemed(uint40 validatorIndex, uint64 withdrawalTimestamp, address indexed recipient, uint64 partialWithdrawalAmountGwei)
func (_EigenPod *EigenPodFilterer) WatchPartialWithdrawalRedeemed(opts *bind.WatchOpts, sink chan<- *EigenPodPartialWithdrawalRedeemed, recipient []common.Address) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _EigenPod.contract.WatchLogs(opts, "PartialWithdrawalRedeemed", recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EigenPodPartialWithdrawalRedeemed)
				if err := _EigenPod.contract.UnpackLog(event, "PartialWithdrawalRedeemed", log); err != nil {
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

// ParsePartialWithdrawalRedeemed is a log parse operation binding the contract event 0x8a7335714231dbd551aaba6314f4a97a14c201e53a3e25e1140325cdf67d7a4e.
//
// Solidity: event PartialWithdrawalRedeemed(uint40 validatorIndex, uint64 withdrawalTimestamp, address indexed recipient, uint64 partialWithdrawalAmountGwei)
func (_EigenPod *EigenPodFilterer) ParsePartialWithdrawalRedeemed(log types.Log) (*EigenPodPartialWithdrawalRedeemed, error) {
	event := new(EigenPodPartialWithdrawalRedeemed)
	if err := _EigenPod.contract.UnpackLog(event, "PartialWithdrawalRedeemed", log); err != nil {
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

// EigenPodRestakingActivatedIterator is returned from FilterRestakingActivated and is used to iterate over the raw logs and unpacked data for RestakingActivated events raised by the EigenPod contract.
type EigenPodRestakingActivatedIterator struct {
	Event *EigenPodRestakingActivated // Event containing the contract specifics and raw log

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
func (it *EigenPodRestakingActivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EigenPodRestakingActivated)
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
		it.Event = new(EigenPodRestakingActivated)
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
func (it *EigenPodRestakingActivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EigenPodRestakingActivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EigenPodRestakingActivated represents a RestakingActivated event raised by the EigenPod contract.
type EigenPodRestakingActivated struct {
	PodOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterRestakingActivated is a free log retrieval operation binding the contract event 0xca8dfc8c5e0a67a74501c072a3325f685259bebbae7cfd230ab85198a78b70cd.
//
// Solidity: event RestakingActivated(address indexed podOwner)
func (_EigenPod *EigenPodFilterer) FilterRestakingActivated(opts *bind.FilterOpts, podOwner []common.Address) (*EigenPodRestakingActivatedIterator, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _EigenPod.contract.FilterLogs(opts, "RestakingActivated", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return &EigenPodRestakingActivatedIterator{contract: _EigenPod.contract, event: "RestakingActivated", logs: logs, sub: sub}, nil
}

// WatchRestakingActivated is a free log subscription operation binding the contract event 0xca8dfc8c5e0a67a74501c072a3325f685259bebbae7cfd230ab85198a78b70cd.
//
// Solidity: event RestakingActivated(address indexed podOwner)
func (_EigenPod *EigenPodFilterer) WatchRestakingActivated(opts *bind.WatchOpts, sink chan<- *EigenPodRestakingActivated, podOwner []common.Address) (event.Subscription, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _EigenPod.contract.WatchLogs(opts, "RestakingActivated", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EigenPodRestakingActivated)
				if err := _EigenPod.contract.UnpackLog(event, "RestakingActivated", log); err != nil {
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

// ParseRestakingActivated is a log parse operation binding the contract event 0xca8dfc8c5e0a67a74501c072a3325f685259bebbae7cfd230ab85198a78b70cd.
//
// Solidity: event RestakingActivated(address indexed podOwner)
func (_EigenPod *EigenPodFilterer) ParseRestakingActivated(log types.Log) (*EigenPodRestakingActivated, error) {
	event := new(EigenPodRestakingActivated)
	if err := _EigenPod.contract.UnpackLog(event, "RestakingActivated", log); err != nil {
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

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

// EigenPodManagerMetaData contains all meta data concerning the EigenPodManager contract.
var EigenPodManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_ethPOS\",\"type\":\"address\",\"internalType\":\"contractIETHPOSDeposit\"},{\"name\":\"_eigenPodBeacon\",\"type\":\"address\",\"internalType\":\"contractIBeacon\"},{\"name\":\"_strategyManager\",\"type\":\"address\",\"internalType\":\"contractIStrategyManager\"},{\"name\":\"_slasher\",\"type\":\"address\",\"internalType\":\"contractISlasher\"},{\"name\":\"_delegationManager\",\"type\":\"address\",\"internalType\":\"contractIDelegationManager\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addShares\",\"inputs\":[{\"name\":\"podOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"beaconChainETHStrategy\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIStrategy\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"createPod\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"delegationManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIDelegationManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenPodBeacon\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIBeacon\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ethPOS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIETHPOSDeposit\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPod\",\"inputs\":[{\"name\":\"podOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenPod\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasPod\",\"inputs\":[{\"name\":\"podOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_pauserRegistry\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"},{\"name\":\"_initPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"numPods\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ownerToPod\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenPod\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pauseAll\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pauserRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"podOwnerShares\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recordBeaconChainETHBalanceUpdate\",\"inputs\":[{\"name\":\"podOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"sharesDelta\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeShares\",\"inputs\":[{\"name\":\"podOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPauserRegistry\",\"inputs\":[{\"name\":\"newPauserRegistry\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"slasher\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractISlasher\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"stake\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"depositDataRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"strategyManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIStrategyManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawSharesAsTokens\",\"inputs\":[{\"name\":\"podOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"destination\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"BeaconChainETHDeposited\",\"inputs\":[{\"name\":\"podOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BeaconChainETHWithdrawalCompleted\",\"inputs\":[{\"name\":\"podOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"nonce\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"},{\"name\":\"delegatedAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"withdrawer\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"withdrawalRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NewTotalShares\",\"inputs\":[{\"name\":\"podOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newTotalShares\",\"type\":\"int256\",\"indexed\":false,\"internalType\":\"int256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PauserRegistrySet\",\"inputs\":[{\"name\":\"pauserRegistry\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIPauserRegistry\"},{\"name\":\"newPauserRegistry\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIPauserRegistry\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PodDeployed\",\"inputs\":[{\"name\":\"eigenPod\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"podOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PodSharesUpdated\",\"inputs\":[{\"name\":\"podOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sharesDelta\",\"type\":\"int256\",\"indexed\":false,\"internalType\":\"int256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
}

// EigenPodManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use EigenPodManagerMetaData.ABI instead.
var EigenPodManagerABI = EigenPodManagerMetaData.ABI

// EigenPodManager is an auto generated Go binding around an Ethereum contract.
type EigenPodManager struct {
	EigenPodManagerCaller     // Read-only binding to the contract
	EigenPodManagerTransactor // Write-only binding to the contract
	EigenPodManagerFilterer   // Log filterer for contract events
}

// EigenPodManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type EigenPodManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EigenPodManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EigenPodManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EigenPodManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EigenPodManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EigenPodManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EigenPodManagerSession struct {
	Contract     *EigenPodManager  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EigenPodManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EigenPodManagerCallerSession struct {
	Contract *EigenPodManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// EigenPodManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EigenPodManagerTransactorSession struct {
	Contract     *EigenPodManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// EigenPodManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type EigenPodManagerRaw struct {
	Contract *EigenPodManager // Generic contract binding to access the raw methods on
}

// EigenPodManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EigenPodManagerCallerRaw struct {
	Contract *EigenPodManagerCaller // Generic read-only contract binding to access the raw methods on
}

// EigenPodManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EigenPodManagerTransactorRaw struct {
	Contract *EigenPodManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEigenPodManager creates a new instance of EigenPodManager, bound to a specific deployed contract.
func NewEigenPodManager(address common.Address, backend bind.ContractBackend) (*EigenPodManager, error) {
	contract, err := bindEigenPodManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EigenPodManager{EigenPodManagerCaller: EigenPodManagerCaller{contract: contract}, EigenPodManagerTransactor: EigenPodManagerTransactor{contract: contract}, EigenPodManagerFilterer: EigenPodManagerFilterer{contract: contract}}, nil
}

// NewEigenPodManagerCaller creates a new read-only instance of EigenPodManager, bound to a specific deployed contract.
func NewEigenPodManagerCaller(address common.Address, caller bind.ContractCaller) (*EigenPodManagerCaller, error) {
	contract, err := bindEigenPodManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EigenPodManagerCaller{contract: contract}, nil
}

// NewEigenPodManagerTransactor creates a new write-only instance of EigenPodManager, bound to a specific deployed contract.
func NewEigenPodManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*EigenPodManagerTransactor, error) {
	contract, err := bindEigenPodManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EigenPodManagerTransactor{contract: contract}, nil
}

// NewEigenPodManagerFilterer creates a new log filterer instance of EigenPodManager, bound to a specific deployed contract.
func NewEigenPodManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*EigenPodManagerFilterer, error) {
	contract, err := bindEigenPodManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EigenPodManagerFilterer{contract: contract}, nil
}

// bindEigenPodManager binds a generic wrapper to an already deployed contract.
func bindEigenPodManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EigenPodManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EigenPodManager *EigenPodManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EigenPodManager.Contract.EigenPodManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EigenPodManager *EigenPodManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EigenPodManager.Contract.EigenPodManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EigenPodManager *EigenPodManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EigenPodManager.Contract.EigenPodManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EigenPodManager *EigenPodManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EigenPodManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EigenPodManager *EigenPodManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EigenPodManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EigenPodManager *EigenPodManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EigenPodManager.Contract.contract.Transact(opts, method, params...)
}

// BeaconChainETHStrategy is a free data retrieval call binding the contract method 0x9104c319.
//
// Solidity: function beaconChainETHStrategy() view returns(address)
func (_EigenPodManager *EigenPodManagerCaller) BeaconChainETHStrategy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EigenPodManager.contract.Call(opts, &out, "beaconChainETHStrategy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BeaconChainETHStrategy is a free data retrieval call binding the contract method 0x9104c319.
//
// Solidity: function beaconChainETHStrategy() view returns(address)
func (_EigenPodManager *EigenPodManagerSession) BeaconChainETHStrategy() (common.Address, error) {
	return _EigenPodManager.Contract.BeaconChainETHStrategy(&_EigenPodManager.CallOpts)
}

// BeaconChainETHStrategy is a free data retrieval call binding the contract method 0x9104c319.
//
// Solidity: function beaconChainETHStrategy() view returns(address)
func (_EigenPodManager *EigenPodManagerCallerSession) BeaconChainETHStrategy() (common.Address, error) {
	return _EigenPodManager.Contract.BeaconChainETHStrategy(&_EigenPodManager.CallOpts)
}

// DelegationManager is a free data retrieval call binding the contract method 0xea4d3c9b.
//
// Solidity: function delegationManager() view returns(address)
func (_EigenPodManager *EigenPodManagerCaller) DelegationManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EigenPodManager.contract.Call(opts, &out, "delegationManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DelegationManager is a free data retrieval call binding the contract method 0xea4d3c9b.
//
// Solidity: function delegationManager() view returns(address)
func (_EigenPodManager *EigenPodManagerSession) DelegationManager() (common.Address, error) {
	return _EigenPodManager.Contract.DelegationManager(&_EigenPodManager.CallOpts)
}

// DelegationManager is a free data retrieval call binding the contract method 0xea4d3c9b.
//
// Solidity: function delegationManager() view returns(address)
func (_EigenPodManager *EigenPodManagerCallerSession) DelegationManager() (common.Address, error) {
	return _EigenPodManager.Contract.DelegationManager(&_EigenPodManager.CallOpts)
}

// EigenPodBeacon is a free data retrieval call binding the contract method 0x292b7b2b.
//
// Solidity: function eigenPodBeacon() view returns(address)
func (_EigenPodManager *EigenPodManagerCaller) EigenPodBeacon(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EigenPodManager.contract.Call(opts, &out, "eigenPodBeacon")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenPodBeacon is a free data retrieval call binding the contract method 0x292b7b2b.
//
// Solidity: function eigenPodBeacon() view returns(address)
func (_EigenPodManager *EigenPodManagerSession) EigenPodBeacon() (common.Address, error) {
	return _EigenPodManager.Contract.EigenPodBeacon(&_EigenPodManager.CallOpts)
}

// EigenPodBeacon is a free data retrieval call binding the contract method 0x292b7b2b.
//
// Solidity: function eigenPodBeacon() view returns(address)
func (_EigenPodManager *EigenPodManagerCallerSession) EigenPodBeacon() (common.Address, error) {
	return _EigenPodManager.Contract.EigenPodBeacon(&_EigenPodManager.CallOpts)
}

// EthPOS is a free data retrieval call binding the contract method 0x74cdd798.
//
// Solidity: function ethPOS() view returns(address)
func (_EigenPodManager *EigenPodManagerCaller) EthPOS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EigenPodManager.contract.Call(opts, &out, "ethPOS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EthPOS is a free data retrieval call binding the contract method 0x74cdd798.
//
// Solidity: function ethPOS() view returns(address)
func (_EigenPodManager *EigenPodManagerSession) EthPOS() (common.Address, error) {
	return _EigenPodManager.Contract.EthPOS(&_EigenPodManager.CallOpts)
}

// EthPOS is a free data retrieval call binding the contract method 0x74cdd798.
//
// Solidity: function ethPOS() view returns(address)
func (_EigenPodManager *EigenPodManagerCallerSession) EthPOS() (common.Address, error) {
	return _EigenPodManager.Contract.EthPOS(&_EigenPodManager.CallOpts)
}

// GetPod is a free data retrieval call binding the contract method 0xa38406a3.
//
// Solidity: function getPod(address podOwner) view returns(address)
func (_EigenPodManager *EigenPodManagerCaller) GetPod(opts *bind.CallOpts, podOwner common.Address) (common.Address, error) {
	var out []interface{}
	err := _EigenPodManager.contract.Call(opts, &out, "getPod", podOwner)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetPod is a free data retrieval call binding the contract method 0xa38406a3.
//
// Solidity: function getPod(address podOwner) view returns(address)
func (_EigenPodManager *EigenPodManagerSession) GetPod(podOwner common.Address) (common.Address, error) {
	return _EigenPodManager.Contract.GetPod(&_EigenPodManager.CallOpts, podOwner)
}

// GetPod is a free data retrieval call binding the contract method 0xa38406a3.
//
// Solidity: function getPod(address podOwner) view returns(address)
func (_EigenPodManager *EigenPodManagerCallerSession) GetPod(podOwner common.Address) (common.Address, error) {
	return _EigenPodManager.Contract.GetPod(&_EigenPodManager.CallOpts, podOwner)
}

// HasPod is a free data retrieval call binding the contract method 0xf6848d24.
//
// Solidity: function hasPod(address podOwner) view returns(bool)
func (_EigenPodManager *EigenPodManagerCaller) HasPod(opts *bind.CallOpts, podOwner common.Address) (bool, error) {
	var out []interface{}
	err := _EigenPodManager.contract.Call(opts, &out, "hasPod", podOwner)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasPod is a free data retrieval call binding the contract method 0xf6848d24.
//
// Solidity: function hasPod(address podOwner) view returns(bool)
func (_EigenPodManager *EigenPodManagerSession) HasPod(podOwner common.Address) (bool, error) {
	return _EigenPodManager.Contract.HasPod(&_EigenPodManager.CallOpts, podOwner)
}

// HasPod is a free data retrieval call binding the contract method 0xf6848d24.
//
// Solidity: function hasPod(address podOwner) view returns(bool)
func (_EigenPodManager *EigenPodManagerCallerSession) HasPod(podOwner common.Address) (bool, error) {
	return _EigenPodManager.Contract.HasPod(&_EigenPodManager.CallOpts, podOwner)
}

// NumPods is a free data retrieval call binding the contract method 0xa6a509be.
//
// Solidity: function numPods() view returns(uint256)
func (_EigenPodManager *EigenPodManagerCaller) NumPods(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EigenPodManager.contract.Call(opts, &out, "numPods")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NumPods is a free data retrieval call binding the contract method 0xa6a509be.
//
// Solidity: function numPods() view returns(uint256)
func (_EigenPodManager *EigenPodManagerSession) NumPods() (*big.Int, error) {
	return _EigenPodManager.Contract.NumPods(&_EigenPodManager.CallOpts)
}

// NumPods is a free data retrieval call binding the contract method 0xa6a509be.
//
// Solidity: function numPods() view returns(uint256)
func (_EigenPodManager *EigenPodManagerCallerSession) NumPods() (*big.Int, error) {
	return _EigenPodManager.Contract.NumPods(&_EigenPodManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_EigenPodManager *EigenPodManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EigenPodManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_EigenPodManager *EigenPodManagerSession) Owner() (common.Address, error) {
	return _EigenPodManager.Contract.Owner(&_EigenPodManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_EigenPodManager *EigenPodManagerCallerSession) Owner() (common.Address, error) {
	return _EigenPodManager.Contract.Owner(&_EigenPodManager.CallOpts)
}

// OwnerToPod is a free data retrieval call binding the contract method 0x9ba06275.
//
// Solidity: function ownerToPod(address ) view returns(address)
func (_EigenPodManager *EigenPodManagerCaller) OwnerToPod(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _EigenPodManager.contract.Call(opts, &out, "ownerToPod", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerToPod is a free data retrieval call binding the contract method 0x9ba06275.
//
// Solidity: function ownerToPod(address ) view returns(address)
func (_EigenPodManager *EigenPodManagerSession) OwnerToPod(arg0 common.Address) (common.Address, error) {
	return _EigenPodManager.Contract.OwnerToPod(&_EigenPodManager.CallOpts, arg0)
}

// OwnerToPod is a free data retrieval call binding the contract method 0x9ba06275.
//
// Solidity: function ownerToPod(address ) view returns(address)
func (_EigenPodManager *EigenPodManagerCallerSession) OwnerToPod(arg0 common.Address) (common.Address, error) {
	return _EigenPodManager.Contract.OwnerToPod(&_EigenPodManager.CallOpts, arg0)
}

// Paused is a free data retrieval call binding the contract method 0x5ac86ab7.
//
// Solidity: function paused(uint8 index) view returns(bool)
func (_EigenPodManager *EigenPodManagerCaller) Paused(opts *bind.CallOpts, index uint8) (bool, error) {
	var out []interface{}
	err := _EigenPodManager.contract.Call(opts, &out, "paused", index)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5ac86ab7.
//
// Solidity: function paused(uint8 index) view returns(bool)
func (_EigenPodManager *EigenPodManagerSession) Paused(index uint8) (bool, error) {
	return _EigenPodManager.Contract.Paused(&_EigenPodManager.CallOpts, index)
}

// Paused is a free data retrieval call binding the contract method 0x5ac86ab7.
//
// Solidity: function paused(uint8 index) view returns(bool)
func (_EigenPodManager *EigenPodManagerCallerSession) Paused(index uint8) (bool, error) {
	return _EigenPodManager.Contract.Paused(&_EigenPodManager.CallOpts, index)
}

// Paused0 is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint256)
func (_EigenPodManager *EigenPodManagerCaller) Paused0(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EigenPodManager.contract.Call(opts, &out, "paused0")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Paused0 is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint256)
func (_EigenPodManager *EigenPodManagerSession) Paused0() (*big.Int, error) {
	return _EigenPodManager.Contract.Paused0(&_EigenPodManager.CallOpts)
}

// Paused0 is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint256)
func (_EigenPodManager *EigenPodManagerCallerSession) Paused0() (*big.Int, error) {
	return _EigenPodManager.Contract.Paused0(&_EigenPodManager.CallOpts)
}

// PauserRegistry is a free data retrieval call binding the contract method 0x886f1195.
//
// Solidity: function pauserRegistry() view returns(address)
func (_EigenPodManager *EigenPodManagerCaller) PauserRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EigenPodManager.contract.Call(opts, &out, "pauserRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PauserRegistry is a free data retrieval call binding the contract method 0x886f1195.
//
// Solidity: function pauserRegistry() view returns(address)
func (_EigenPodManager *EigenPodManagerSession) PauserRegistry() (common.Address, error) {
	return _EigenPodManager.Contract.PauserRegistry(&_EigenPodManager.CallOpts)
}

// PauserRegistry is a free data retrieval call binding the contract method 0x886f1195.
//
// Solidity: function pauserRegistry() view returns(address)
func (_EigenPodManager *EigenPodManagerCallerSession) PauserRegistry() (common.Address, error) {
	return _EigenPodManager.Contract.PauserRegistry(&_EigenPodManager.CallOpts)
}

// PodOwnerShares is a free data retrieval call binding the contract method 0x60f4062b.
//
// Solidity: function podOwnerShares(address ) view returns(int256)
func (_EigenPodManager *EigenPodManagerCaller) PodOwnerShares(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _EigenPodManager.contract.Call(opts, &out, "podOwnerShares", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PodOwnerShares is a free data retrieval call binding the contract method 0x60f4062b.
//
// Solidity: function podOwnerShares(address ) view returns(int256)
func (_EigenPodManager *EigenPodManagerSession) PodOwnerShares(arg0 common.Address) (*big.Int, error) {
	return _EigenPodManager.Contract.PodOwnerShares(&_EigenPodManager.CallOpts, arg0)
}

// PodOwnerShares is a free data retrieval call binding the contract method 0x60f4062b.
//
// Solidity: function podOwnerShares(address ) view returns(int256)
func (_EigenPodManager *EigenPodManagerCallerSession) PodOwnerShares(arg0 common.Address) (*big.Int, error) {
	return _EigenPodManager.Contract.PodOwnerShares(&_EigenPodManager.CallOpts, arg0)
}

// Slasher is a free data retrieval call binding the contract method 0xb1344271.
//
// Solidity: function slasher() view returns(address)
func (_EigenPodManager *EigenPodManagerCaller) Slasher(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EigenPodManager.contract.Call(opts, &out, "slasher")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Slasher is a free data retrieval call binding the contract method 0xb1344271.
//
// Solidity: function slasher() view returns(address)
func (_EigenPodManager *EigenPodManagerSession) Slasher() (common.Address, error) {
	return _EigenPodManager.Contract.Slasher(&_EigenPodManager.CallOpts)
}

// Slasher is a free data retrieval call binding the contract method 0xb1344271.
//
// Solidity: function slasher() view returns(address)
func (_EigenPodManager *EigenPodManagerCallerSession) Slasher() (common.Address, error) {
	return _EigenPodManager.Contract.Slasher(&_EigenPodManager.CallOpts)
}

// StrategyManager is a free data retrieval call binding the contract method 0x39b70e38.
//
// Solidity: function strategyManager() view returns(address)
func (_EigenPodManager *EigenPodManagerCaller) StrategyManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EigenPodManager.contract.Call(opts, &out, "strategyManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StrategyManager is a free data retrieval call binding the contract method 0x39b70e38.
//
// Solidity: function strategyManager() view returns(address)
func (_EigenPodManager *EigenPodManagerSession) StrategyManager() (common.Address, error) {
	return _EigenPodManager.Contract.StrategyManager(&_EigenPodManager.CallOpts)
}

// StrategyManager is a free data retrieval call binding the contract method 0x39b70e38.
//
// Solidity: function strategyManager() view returns(address)
func (_EigenPodManager *EigenPodManagerCallerSession) StrategyManager() (common.Address, error) {
	return _EigenPodManager.Contract.StrategyManager(&_EigenPodManager.CallOpts)
}

// AddShares is a paid mutator transaction binding the contract method 0x0e81073c.
//
// Solidity: function addShares(address podOwner, uint256 shares) returns(uint256)
func (_EigenPodManager *EigenPodManagerTransactor) AddShares(opts *bind.TransactOpts, podOwner common.Address, shares *big.Int) (*types.Transaction, error) {
	return _EigenPodManager.contract.Transact(opts, "addShares", podOwner, shares)
}

// AddShares is a paid mutator transaction binding the contract method 0x0e81073c.
//
// Solidity: function addShares(address podOwner, uint256 shares) returns(uint256)
func (_EigenPodManager *EigenPodManagerSession) AddShares(podOwner common.Address, shares *big.Int) (*types.Transaction, error) {
	return _EigenPodManager.Contract.AddShares(&_EigenPodManager.TransactOpts, podOwner, shares)
}

// AddShares is a paid mutator transaction binding the contract method 0x0e81073c.
//
// Solidity: function addShares(address podOwner, uint256 shares) returns(uint256)
func (_EigenPodManager *EigenPodManagerTransactorSession) AddShares(podOwner common.Address, shares *big.Int) (*types.Transaction, error) {
	return _EigenPodManager.Contract.AddShares(&_EigenPodManager.TransactOpts, podOwner, shares)
}

// CreatePod is a paid mutator transaction binding the contract method 0x84d81062.
//
// Solidity: function createPod() returns(address)
func (_EigenPodManager *EigenPodManagerTransactor) CreatePod(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EigenPodManager.contract.Transact(opts, "createPod")
}

// CreatePod is a paid mutator transaction binding the contract method 0x84d81062.
//
// Solidity: function createPod() returns(address)
func (_EigenPodManager *EigenPodManagerSession) CreatePod() (*types.Transaction, error) {
	return _EigenPodManager.Contract.CreatePod(&_EigenPodManager.TransactOpts)
}

// CreatePod is a paid mutator transaction binding the contract method 0x84d81062.
//
// Solidity: function createPod() returns(address)
func (_EigenPodManager *EigenPodManagerTransactorSession) CreatePod() (*types.Transaction, error) {
	return _EigenPodManager.Contract.CreatePod(&_EigenPodManager.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x1794bb3c.
//
// Solidity: function initialize(address initialOwner, address _pauserRegistry, uint256 _initPausedStatus) returns()
func (_EigenPodManager *EigenPodManagerTransactor) Initialize(opts *bind.TransactOpts, initialOwner common.Address, _pauserRegistry common.Address, _initPausedStatus *big.Int) (*types.Transaction, error) {
	return _EigenPodManager.contract.Transact(opts, "initialize", initialOwner, _pauserRegistry, _initPausedStatus)
}

// Initialize is a paid mutator transaction binding the contract method 0x1794bb3c.
//
// Solidity: function initialize(address initialOwner, address _pauserRegistry, uint256 _initPausedStatus) returns()
func (_EigenPodManager *EigenPodManagerSession) Initialize(initialOwner common.Address, _pauserRegistry common.Address, _initPausedStatus *big.Int) (*types.Transaction, error) {
	return _EigenPodManager.Contract.Initialize(&_EigenPodManager.TransactOpts, initialOwner, _pauserRegistry, _initPausedStatus)
}

// Initialize is a paid mutator transaction binding the contract method 0x1794bb3c.
//
// Solidity: function initialize(address initialOwner, address _pauserRegistry, uint256 _initPausedStatus) returns()
func (_EigenPodManager *EigenPodManagerTransactorSession) Initialize(initialOwner common.Address, _pauserRegistry common.Address, _initPausedStatus *big.Int) (*types.Transaction, error) {
	return _EigenPodManager.Contract.Initialize(&_EigenPodManager.TransactOpts, initialOwner, _pauserRegistry, _initPausedStatus)
}

// Pause is a paid mutator transaction binding the contract method 0x136439dd.
//
// Solidity: function pause(uint256 newPausedStatus) returns()
func (_EigenPodManager *EigenPodManagerTransactor) Pause(opts *bind.TransactOpts, newPausedStatus *big.Int) (*types.Transaction, error) {
	return _EigenPodManager.contract.Transact(opts, "pause", newPausedStatus)
}

// Pause is a paid mutator transaction binding the contract method 0x136439dd.
//
// Solidity: function pause(uint256 newPausedStatus) returns()
func (_EigenPodManager *EigenPodManagerSession) Pause(newPausedStatus *big.Int) (*types.Transaction, error) {
	return _EigenPodManager.Contract.Pause(&_EigenPodManager.TransactOpts, newPausedStatus)
}

// Pause is a paid mutator transaction binding the contract method 0x136439dd.
//
// Solidity: function pause(uint256 newPausedStatus) returns()
func (_EigenPodManager *EigenPodManagerTransactorSession) Pause(newPausedStatus *big.Int) (*types.Transaction, error) {
	return _EigenPodManager.Contract.Pause(&_EigenPodManager.TransactOpts, newPausedStatus)
}

// PauseAll is a paid mutator transaction binding the contract method 0x595c6a67.
//
// Solidity: function pauseAll() returns()
func (_EigenPodManager *EigenPodManagerTransactor) PauseAll(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EigenPodManager.contract.Transact(opts, "pauseAll")
}

// PauseAll is a paid mutator transaction binding the contract method 0x595c6a67.
//
// Solidity: function pauseAll() returns()
func (_EigenPodManager *EigenPodManagerSession) PauseAll() (*types.Transaction, error) {
	return _EigenPodManager.Contract.PauseAll(&_EigenPodManager.TransactOpts)
}

// PauseAll is a paid mutator transaction binding the contract method 0x595c6a67.
//
// Solidity: function pauseAll() returns()
func (_EigenPodManager *EigenPodManagerTransactorSession) PauseAll() (*types.Transaction, error) {
	return _EigenPodManager.Contract.PauseAll(&_EigenPodManager.TransactOpts)
}

// RecordBeaconChainETHBalanceUpdate is a paid mutator transaction binding the contract method 0xc2c51c40.
//
// Solidity: function recordBeaconChainETHBalanceUpdate(address podOwner, int256 sharesDelta) returns()
func (_EigenPodManager *EigenPodManagerTransactor) RecordBeaconChainETHBalanceUpdate(opts *bind.TransactOpts, podOwner common.Address, sharesDelta *big.Int) (*types.Transaction, error) {
	return _EigenPodManager.contract.Transact(opts, "recordBeaconChainETHBalanceUpdate", podOwner, sharesDelta)
}

// RecordBeaconChainETHBalanceUpdate is a paid mutator transaction binding the contract method 0xc2c51c40.
//
// Solidity: function recordBeaconChainETHBalanceUpdate(address podOwner, int256 sharesDelta) returns()
func (_EigenPodManager *EigenPodManagerSession) RecordBeaconChainETHBalanceUpdate(podOwner common.Address, sharesDelta *big.Int) (*types.Transaction, error) {
	return _EigenPodManager.Contract.RecordBeaconChainETHBalanceUpdate(&_EigenPodManager.TransactOpts, podOwner, sharesDelta)
}

// RecordBeaconChainETHBalanceUpdate is a paid mutator transaction binding the contract method 0xc2c51c40.
//
// Solidity: function recordBeaconChainETHBalanceUpdate(address podOwner, int256 sharesDelta) returns()
func (_EigenPodManager *EigenPodManagerTransactorSession) RecordBeaconChainETHBalanceUpdate(podOwner common.Address, sharesDelta *big.Int) (*types.Transaction, error) {
	return _EigenPodManager.Contract.RecordBeaconChainETHBalanceUpdate(&_EigenPodManager.TransactOpts, podOwner, sharesDelta)
}

// RemoveShares is a paid mutator transaction binding the contract method 0xbeffbb89.
//
// Solidity: function removeShares(address podOwner, uint256 shares) returns()
func (_EigenPodManager *EigenPodManagerTransactor) RemoveShares(opts *bind.TransactOpts, podOwner common.Address, shares *big.Int) (*types.Transaction, error) {
	return _EigenPodManager.contract.Transact(opts, "removeShares", podOwner, shares)
}

// RemoveShares is a paid mutator transaction binding the contract method 0xbeffbb89.
//
// Solidity: function removeShares(address podOwner, uint256 shares) returns()
func (_EigenPodManager *EigenPodManagerSession) RemoveShares(podOwner common.Address, shares *big.Int) (*types.Transaction, error) {
	return _EigenPodManager.Contract.RemoveShares(&_EigenPodManager.TransactOpts, podOwner, shares)
}

// RemoveShares is a paid mutator transaction binding the contract method 0xbeffbb89.
//
// Solidity: function removeShares(address podOwner, uint256 shares) returns()
func (_EigenPodManager *EigenPodManagerTransactorSession) RemoveShares(podOwner common.Address, shares *big.Int) (*types.Transaction, error) {
	return _EigenPodManager.Contract.RemoveShares(&_EigenPodManager.TransactOpts, podOwner, shares)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_EigenPodManager *EigenPodManagerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EigenPodManager.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_EigenPodManager *EigenPodManagerSession) RenounceOwnership() (*types.Transaction, error) {
	return _EigenPodManager.Contract.RenounceOwnership(&_EigenPodManager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_EigenPodManager *EigenPodManagerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _EigenPodManager.Contract.RenounceOwnership(&_EigenPodManager.TransactOpts)
}

// SetPauserRegistry is a paid mutator transaction binding the contract method 0x10d67a2f.
//
// Solidity: function setPauserRegistry(address newPauserRegistry) returns()
func (_EigenPodManager *EigenPodManagerTransactor) SetPauserRegistry(opts *bind.TransactOpts, newPauserRegistry common.Address) (*types.Transaction, error) {
	return _EigenPodManager.contract.Transact(opts, "setPauserRegistry", newPauserRegistry)
}

// SetPauserRegistry is a paid mutator transaction binding the contract method 0x10d67a2f.
//
// Solidity: function setPauserRegistry(address newPauserRegistry) returns()
func (_EigenPodManager *EigenPodManagerSession) SetPauserRegistry(newPauserRegistry common.Address) (*types.Transaction, error) {
	return _EigenPodManager.Contract.SetPauserRegistry(&_EigenPodManager.TransactOpts, newPauserRegistry)
}

// SetPauserRegistry is a paid mutator transaction binding the contract method 0x10d67a2f.
//
// Solidity: function setPauserRegistry(address newPauserRegistry) returns()
func (_EigenPodManager *EigenPodManagerTransactorSession) SetPauserRegistry(newPauserRegistry common.Address) (*types.Transaction, error) {
	return _EigenPodManager.Contract.SetPauserRegistry(&_EigenPodManager.TransactOpts, newPauserRegistry)
}

// Stake is a paid mutator transaction binding the contract method 0x9b4e4634.
//
// Solidity: function stake(bytes pubkey, bytes signature, bytes32 depositDataRoot) payable returns()
func (_EigenPodManager *EigenPodManagerTransactor) Stake(opts *bind.TransactOpts, pubkey []byte, signature []byte, depositDataRoot [32]byte) (*types.Transaction, error) {
	return _EigenPodManager.contract.Transact(opts, "stake", pubkey, signature, depositDataRoot)
}

// Stake is a paid mutator transaction binding the contract method 0x9b4e4634.
//
// Solidity: function stake(bytes pubkey, bytes signature, bytes32 depositDataRoot) payable returns()
func (_EigenPodManager *EigenPodManagerSession) Stake(pubkey []byte, signature []byte, depositDataRoot [32]byte) (*types.Transaction, error) {
	return _EigenPodManager.Contract.Stake(&_EigenPodManager.TransactOpts, pubkey, signature, depositDataRoot)
}

// Stake is a paid mutator transaction binding the contract method 0x9b4e4634.
//
// Solidity: function stake(bytes pubkey, bytes signature, bytes32 depositDataRoot) payable returns()
func (_EigenPodManager *EigenPodManagerTransactorSession) Stake(pubkey []byte, signature []byte, depositDataRoot [32]byte) (*types.Transaction, error) {
	return _EigenPodManager.Contract.Stake(&_EigenPodManager.TransactOpts, pubkey, signature, depositDataRoot)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_EigenPodManager *EigenPodManagerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _EigenPodManager.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_EigenPodManager *EigenPodManagerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _EigenPodManager.Contract.TransferOwnership(&_EigenPodManager.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_EigenPodManager *EigenPodManagerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _EigenPodManager.Contract.TransferOwnership(&_EigenPodManager.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0xfabc1cbc.
//
// Solidity: function unpause(uint256 newPausedStatus) returns()
func (_EigenPodManager *EigenPodManagerTransactor) Unpause(opts *bind.TransactOpts, newPausedStatus *big.Int) (*types.Transaction, error) {
	return _EigenPodManager.contract.Transact(opts, "unpause", newPausedStatus)
}

// Unpause is a paid mutator transaction binding the contract method 0xfabc1cbc.
//
// Solidity: function unpause(uint256 newPausedStatus) returns()
func (_EigenPodManager *EigenPodManagerSession) Unpause(newPausedStatus *big.Int) (*types.Transaction, error) {
	return _EigenPodManager.Contract.Unpause(&_EigenPodManager.TransactOpts, newPausedStatus)
}

// Unpause is a paid mutator transaction binding the contract method 0xfabc1cbc.
//
// Solidity: function unpause(uint256 newPausedStatus) returns()
func (_EigenPodManager *EigenPodManagerTransactorSession) Unpause(newPausedStatus *big.Int) (*types.Transaction, error) {
	return _EigenPodManager.Contract.Unpause(&_EigenPodManager.TransactOpts, newPausedStatus)
}

// WithdrawSharesAsTokens is a paid mutator transaction binding the contract method 0x387b1300.
//
// Solidity: function withdrawSharesAsTokens(address podOwner, address destination, uint256 shares) returns()
func (_EigenPodManager *EigenPodManagerTransactor) WithdrawSharesAsTokens(opts *bind.TransactOpts, podOwner common.Address, destination common.Address, shares *big.Int) (*types.Transaction, error) {
	return _EigenPodManager.contract.Transact(opts, "withdrawSharesAsTokens", podOwner, destination, shares)
}

// WithdrawSharesAsTokens is a paid mutator transaction binding the contract method 0x387b1300.
//
// Solidity: function withdrawSharesAsTokens(address podOwner, address destination, uint256 shares) returns()
func (_EigenPodManager *EigenPodManagerSession) WithdrawSharesAsTokens(podOwner common.Address, destination common.Address, shares *big.Int) (*types.Transaction, error) {
	return _EigenPodManager.Contract.WithdrawSharesAsTokens(&_EigenPodManager.TransactOpts, podOwner, destination, shares)
}

// WithdrawSharesAsTokens is a paid mutator transaction binding the contract method 0x387b1300.
//
// Solidity: function withdrawSharesAsTokens(address podOwner, address destination, uint256 shares) returns()
func (_EigenPodManager *EigenPodManagerTransactorSession) WithdrawSharesAsTokens(podOwner common.Address, destination common.Address, shares *big.Int) (*types.Transaction, error) {
	return _EigenPodManager.Contract.WithdrawSharesAsTokens(&_EigenPodManager.TransactOpts, podOwner, destination, shares)
}

// EigenPodManagerBeaconChainETHDepositedIterator is returned from FilterBeaconChainETHDeposited and is used to iterate over the raw logs and unpacked data for BeaconChainETHDeposited events raised by the EigenPodManager contract.
type EigenPodManagerBeaconChainETHDepositedIterator struct {
	Event *EigenPodManagerBeaconChainETHDeposited // Event containing the contract specifics and raw log

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
func (it *EigenPodManagerBeaconChainETHDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EigenPodManagerBeaconChainETHDeposited)
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
		it.Event = new(EigenPodManagerBeaconChainETHDeposited)
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
func (it *EigenPodManagerBeaconChainETHDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EigenPodManagerBeaconChainETHDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EigenPodManagerBeaconChainETHDeposited represents a BeaconChainETHDeposited event raised by the EigenPodManager contract.
type EigenPodManagerBeaconChainETHDeposited struct {
	PodOwner common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterBeaconChainETHDeposited is a free log retrieval operation binding the contract event 0x35a85cabc603f48abb2b71d9fbd8adea7c449d7f0be900ae7a2986ea369c3d0d.
//
// Solidity: event BeaconChainETHDeposited(address indexed podOwner, uint256 amount)
func (_EigenPodManager *EigenPodManagerFilterer) FilterBeaconChainETHDeposited(opts *bind.FilterOpts, podOwner []common.Address) (*EigenPodManagerBeaconChainETHDepositedIterator, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _EigenPodManager.contract.FilterLogs(opts, "BeaconChainETHDeposited", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return &EigenPodManagerBeaconChainETHDepositedIterator{contract: _EigenPodManager.contract, event: "BeaconChainETHDeposited", logs: logs, sub: sub}, nil
}

// WatchBeaconChainETHDeposited is a free log subscription operation binding the contract event 0x35a85cabc603f48abb2b71d9fbd8adea7c449d7f0be900ae7a2986ea369c3d0d.
//
// Solidity: event BeaconChainETHDeposited(address indexed podOwner, uint256 amount)
func (_EigenPodManager *EigenPodManagerFilterer) WatchBeaconChainETHDeposited(opts *bind.WatchOpts, sink chan<- *EigenPodManagerBeaconChainETHDeposited, podOwner []common.Address) (event.Subscription, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _EigenPodManager.contract.WatchLogs(opts, "BeaconChainETHDeposited", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EigenPodManagerBeaconChainETHDeposited)
				if err := _EigenPodManager.contract.UnpackLog(event, "BeaconChainETHDeposited", log); err != nil {
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

// ParseBeaconChainETHDeposited is a log parse operation binding the contract event 0x35a85cabc603f48abb2b71d9fbd8adea7c449d7f0be900ae7a2986ea369c3d0d.
//
// Solidity: event BeaconChainETHDeposited(address indexed podOwner, uint256 amount)
func (_EigenPodManager *EigenPodManagerFilterer) ParseBeaconChainETHDeposited(log types.Log) (*EigenPodManagerBeaconChainETHDeposited, error) {
	event := new(EigenPodManagerBeaconChainETHDeposited)
	if err := _EigenPodManager.contract.UnpackLog(event, "BeaconChainETHDeposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EigenPodManagerBeaconChainETHWithdrawalCompletedIterator is returned from FilterBeaconChainETHWithdrawalCompleted and is used to iterate over the raw logs and unpacked data for BeaconChainETHWithdrawalCompleted events raised by the EigenPodManager contract.
type EigenPodManagerBeaconChainETHWithdrawalCompletedIterator struct {
	Event *EigenPodManagerBeaconChainETHWithdrawalCompleted // Event containing the contract specifics and raw log

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
func (it *EigenPodManagerBeaconChainETHWithdrawalCompletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EigenPodManagerBeaconChainETHWithdrawalCompleted)
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
		it.Event = new(EigenPodManagerBeaconChainETHWithdrawalCompleted)
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
func (it *EigenPodManagerBeaconChainETHWithdrawalCompletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EigenPodManagerBeaconChainETHWithdrawalCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EigenPodManagerBeaconChainETHWithdrawalCompleted represents a BeaconChainETHWithdrawalCompleted event raised by the EigenPodManager contract.
type EigenPodManagerBeaconChainETHWithdrawalCompleted struct {
	PodOwner         common.Address
	Shares           *big.Int
	Nonce            *big.Int
	DelegatedAddress common.Address
	Withdrawer       common.Address
	WithdrawalRoot   [32]byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterBeaconChainETHWithdrawalCompleted is a free log retrieval operation binding the contract event 0xa6bab1d55a361fcea2eee2bc9491e4f01e6cf333df03c9c4f2c144466429f7d6.
//
// Solidity: event BeaconChainETHWithdrawalCompleted(address indexed podOwner, uint256 shares, uint96 nonce, address delegatedAddress, address withdrawer, bytes32 withdrawalRoot)
func (_EigenPodManager *EigenPodManagerFilterer) FilterBeaconChainETHWithdrawalCompleted(opts *bind.FilterOpts, podOwner []common.Address) (*EigenPodManagerBeaconChainETHWithdrawalCompletedIterator, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _EigenPodManager.contract.FilterLogs(opts, "BeaconChainETHWithdrawalCompleted", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return &EigenPodManagerBeaconChainETHWithdrawalCompletedIterator{contract: _EigenPodManager.contract, event: "BeaconChainETHWithdrawalCompleted", logs: logs, sub: sub}, nil
}

// WatchBeaconChainETHWithdrawalCompleted is a free log subscription operation binding the contract event 0xa6bab1d55a361fcea2eee2bc9491e4f01e6cf333df03c9c4f2c144466429f7d6.
//
// Solidity: event BeaconChainETHWithdrawalCompleted(address indexed podOwner, uint256 shares, uint96 nonce, address delegatedAddress, address withdrawer, bytes32 withdrawalRoot)
func (_EigenPodManager *EigenPodManagerFilterer) WatchBeaconChainETHWithdrawalCompleted(opts *bind.WatchOpts, sink chan<- *EigenPodManagerBeaconChainETHWithdrawalCompleted, podOwner []common.Address) (event.Subscription, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _EigenPodManager.contract.WatchLogs(opts, "BeaconChainETHWithdrawalCompleted", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EigenPodManagerBeaconChainETHWithdrawalCompleted)
				if err := _EigenPodManager.contract.UnpackLog(event, "BeaconChainETHWithdrawalCompleted", log); err != nil {
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

// ParseBeaconChainETHWithdrawalCompleted is a log parse operation binding the contract event 0xa6bab1d55a361fcea2eee2bc9491e4f01e6cf333df03c9c4f2c144466429f7d6.
//
// Solidity: event BeaconChainETHWithdrawalCompleted(address indexed podOwner, uint256 shares, uint96 nonce, address delegatedAddress, address withdrawer, bytes32 withdrawalRoot)
func (_EigenPodManager *EigenPodManagerFilterer) ParseBeaconChainETHWithdrawalCompleted(log types.Log) (*EigenPodManagerBeaconChainETHWithdrawalCompleted, error) {
	event := new(EigenPodManagerBeaconChainETHWithdrawalCompleted)
	if err := _EigenPodManager.contract.UnpackLog(event, "BeaconChainETHWithdrawalCompleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EigenPodManagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the EigenPodManager contract.
type EigenPodManagerInitializedIterator struct {
	Event *EigenPodManagerInitialized // Event containing the contract specifics and raw log

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
func (it *EigenPodManagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EigenPodManagerInitialized)
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
		it.Event = new(EigenPodManagerInitialized)
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
func (it *EigenPodManagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EigenPodManagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EigenPodManagerInitialized represents a Initialized event raised by the EigenPodManager contract.
type EigenPodManagerInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_EigenPodManager *EigenPodManagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*EigenPodManagerInitializedIterator, error) {

	logs, sub, err := _EigenPodManager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &EigenPodManagerInitializedIterator{contract: _EigenPodManager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_EigenPodManager *EigenPodManagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *EigenPodManagerInitialized) (event.Subscription, error) {

	logs, sub, err := _EigenPodManager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EigenPodManagerInitialized)
				if err := _EigenPodManager.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_EigenPodManager *EigenPodManagerFilterer) ParseInitialized(log types.Log) (*EigenPodManagerInitialized, error) {
	event := new(EigenPodManagerInitialized)
	if err := _EigenPodManager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EigenPodManagerNewTotalSharesIterator is returned from FilterNewTotalShares and is used to iterate over the raw logs and unpacked data for NewTotalShares events raised by the EigenPodManager contract.
type EigenPodManagerNewTotalSharesIterator struct {
	Event *EigenPodManagerNewTotalShares // Event containing the contract specifics and raw log

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
func (it *EigenPodManagerNewTotalSharesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EigenPodManagerNewTotalShares)
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
		it.Event = new(EigenPodManagerNewTotalShares)
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
func (it *EigenPodManagerNewTotalSharesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EigenPodManagerNewTotalSharesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EigenPodManagerNewTotalShares represents a NewTotalShares event raised by the EigenPodManager contract.
type EigenPodManagerNewTotalShares struct {
	PodOwner       common.Address
	NewTotalShares *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterNewTotalShares is a free log retrieval operation binding the contract event 0xd4def76d6d2bed6f14d5cd9af73cc2913d618d00edde42432e81c09bfe077098.
//
// Solidity: event NewTotalShares(address indexed podOwner, int256 newTotalShares)
func (_EigenPodManager *EigenPodManagerFilterer) FilterNewTotalShares(opts *bind.FilterOpts, podOwner []common.Address) (*EigenPodManagerNewTotalSharesIterator, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _EigenPodManager.contract.FilterLogs(opts, "NewTotalShares", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return &EigenPodManagerNewTotalSharesIterator{contract: _EigenPodManager.contract, event: "NewTotalShares", logs: logs, sub: sub}, nil
}

// WatchNewTotalShares is a free log subscription operation binding the contract event 0xd4def76d6d2bed6f14d5cd9af73cc2913d618d00edde42432e81c09bfe077098.
//
// Solidity: event NewTotalShares(address indexed podOwner, int256 newTotalShares)
func (_EigenPodManager *EigenPodManagerFilterer) WatchNewTotalShares(opts *bind.WatchOpts, sink chan<- *EigenPodManagerNewTotalShares, podOwner []common.Address) (event.Subscription, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _EigenPodManager.contract.WatchLogs(opts, "NewTotalShares", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EigenPodManagerNewTotalShares)
				if err := _EigenPodManager.contract.UnpackLog(event, "NewTotalShares", log); err != nil {
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

// ParseNewTotalShares is a log parse operation binding the contract event 0xd4def76d6d2bed6f14d5cd9af73cc2913d618d00edde42432e81c09bfe077098.
//
// Solidity: event NewTotalShares(address indexed podOwner, int256 newTotalShares)
func (_EigenPodManager *EigenPodManagerFilterer) ParseNewTotalShares(log types.Log) (*EigenPodManagerNewTotalShares, error) {
	event := new(EigenPodManagerNewTotalShares)
	if err := _EigenPodManager.contract.UnpackLog(event, "NewTotalShares", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EigenPodManagerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the EigenPodManager contract.
type EigenPodManagerOwnershipTransferredIterator struct {
	Event *EigenPodManagerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *EigenPodManagerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EigenPodManagerOwnershipTransferred)
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
		it.Event = new(EigenPodManagerOwnershipTransferred)
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
func (it *EigenPodManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EigenPodManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EigenPodManagerOwnershipTransferred represents a OwnershipTransferred event raised by the EigenPodManager contract.
type EigenPodManagerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_EigenPodManager *EigenPodManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*EigenPodManagerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _EigenPodManager.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &EigenPodManagerOwnershipTransferredIterator{contract: _EigenPodManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_EigenPodManager *EigenPodManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *EigenPodManagerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _EigenPodManager.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EigenPodManagerOwnershipTransferred)
				if err := _EigenPodManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_EigenPodManager *EigenPodManagerFilterer) ParseOwnershipTransferred(log types.Log) (*EigenPodManagerOwnershipTransferred, error) {
	event := new(EigenPodManagerOwnershipTransferred)
	if err := _EigenPodManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EigenPodManagerPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the EigenPodManager contract.
type EigenPodManagerPausedIterator struct {
	Event *EigenPodManagerPaused // Event containing the contract specifics and raw log

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
func (it *EigenPodManagerPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EigenPodManagerPaused)
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
		it.Event = new(EigenPodManagerPaused)
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
func (it *EigenPodManagerPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EigenPodManagerPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EigenPodManagerPaused represents a Paused event raised by the EigenPodManager contract.
type EigenPodManagerPaused struct {
	Account         common.Address
	NewPausedStatus *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0xab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d.
//
// Solidity: event Paused(address indexed account, uint256 newPausedStatus)
func (_EigenPodManager *EigenPodManagerFilterer) FilterPaused(opts *bind.FilterOpts, account []common.Address) (*EigenPodManagerPausedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _EigenPodManager.contract.FilterLogs(opts, "Paused", accountRule)
	if err != nil {
		return nil, err
	}
	return &EigenPodManagerPausedIterator{contract: _EigenPodManager.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0xab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d.
//
// Solidity: event Paused(address indexed account, uint256 newPausedStatus)
func (_EigenPodManager *EigenPodManagerFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *EigenPodManagerPaused, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _EigenPodManager.contract.WatchLogs(opts, "Paused", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EigenPodManagerPaused)
				if err := _EigenPodManager.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0xab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d.
//
// Solidity: event Paused(address indexed account, uint256 newPausedStatus)
func (_EigenPodManager *EigenPodManagerFilterer) ParsePaused(log types.Log) (*EigenPodManagerPaused, error) {
	event := new(EigenPodManagerPaused)
	if err := _EigenPodManager.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EigenPodManagerPauserRegistrySetIterator is returned from FilterPauserRegistrySet and is used to iterate over the raw logs and unpacked data for PauserRegistrySet events raised by the EigenPodManager contract.
type EigenPodManagerPauserRegistrySetIterator struct {
	Event *EigenPodManagerPauserRegistrySet // Event containing the contract specifics and raw log

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
func (it *EigenPodManagerPauserRegistrySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EigenPodManagerPauserRegistrySet)
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
		it.Event = new(EigenPodManagerPauserRegistrySet)
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
func (it *EigenPodManagerPauserRegistrySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EigenPodManagerPauserRegistrySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EigenPodManagerPauserRegistrySet represents a PauserRegistrySet event raised by the EigenPodManager contract.
type EigenPodManagerPauserRegistrySet struct {
	PauserRegistry    common.Address
	NewPauserRegistry common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterPauserRegistrySet is a free log retrieval operation binding the contract event 0x6e9fcd539896fca60e8b0f01dd580233e48a6b0f7df013b89ba7f565869acdb6.
//
// Solidity: event PauserRegistrySet(address pauserRegistry, address newPauserRegistry)
func (_EigenPodManager *EigenPodManagerFilterer) FilterPauserRegistrySet(opts *bind.FilterOpts) (*EigenPodManagerPauserRegistrySetIterator, error) {

	logs, sub, err := _EigenPodManager.contract.FilterLogs(opts, "PauserRegistrySet")
	if err != nil {
		return nil, err
	}
	return &EigenPodManagerPauserRegistrySetIterator{contract: _EigenPodManager.contract, event: "PauserRegistrySet", logs: logs, sub: sub}, nil
}

// WatchPauserRegistrySet is a free log subscription operation binding the contract event 0x6e9fcd539896fca60e8b0f01dd580233e48a6b0f7df013b89ba7f565869acdb6.
//
// Solidity: event PauserRegistrySet(address pauserRegistry, address newPauserRegistry)
func (_EigenPodManager *EigenPodManagerFilterer) WatchPauserRegistrySet(opts *bind.WatchOpts, sink chan<- *EigenPodManagerPauserRegistrySet) (event.Subscription, error) {

	logs, sub, err := _EigenPodManager.contract.WatchLogs(opts, "PauserRegistrySet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EigenPodManagerPauserRegistrySet)
				if err := _EigenPodManager.contract.UnpackLog(event, "PauserRegistrySet", log); err != nil {
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

// ParsePauserRegistrySet is a log parse operation binding the contract event 0x6e9fcd539896fca60e8b0f01dd580233e48a6b0f7df013b89ba7f565869acdb6.
//
// Solidity: event PauserRegistrySet(address pauserRegistry, address newPauserRegistry)
func (_EigenPodManager *EigenPodManagerFilterer) ParsePauserRegistrySet(log types.Log) (*EigenPodManagerPauserRegistrySet, error) {
	event := new(EigenPodManagerPauserRegistrySet)
	if err := _EigenPodManager.contract.UnpackLog(event, "PauserRegistrySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EigenPodManagerPodDeployedIterator is returned from FilterPodDeployed and is used to iterate over the raw logs and unpacked data for PodDeployed events raised by the EigenPodManager contract.
type EigenPodManagerPodDeployedIterator struct {
	Event *EigenPodManagerPodDeployed // Event containing the contract specifics and raw log

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
func (it *EigenPodManagerPodDeployedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EigenPodManagerPodDeployed)
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
		it.Event = new(EigenPodManagerPodDeployed)
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
func (it *EigenPodManagerPodDeployedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EigenPodManagerPodDeployedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EigenPodManagerPodDeployed represents a PodDeployed event raised by the EigenPodManager contract.
type EigenPodManagerPodDeployed struct {
	EigenPod common.Address
	PodOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterPodDeployed is a free log retrieval operation binding the contract event 0x21c99d0db02213c32fff5b05cf0a718ab5f858802b91498f80d82270289d856a.
//
// Solidity: event PodDeployed(address indexed eigenPod, address indexed podOwner)
func (_EigenPodManager *EigenPodManagerFilterer) FilterPodDeployed(opts *bind.FilterOpts, eigenPod []common.Address, podOwner []common.Address) (*EigenPodManagerPodDeployedIterator, error) {

	var eigenPodRule []interface{}
	for _, eigenPodItem := range eigenPod {
		eigenPodRule = append(eigenPodRule, eigenPodItem)
	}
	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _EigenPodManager.contract.FilterLogs(opts, "PodDeployed", eigenPodRule, podOwnerRule)
	if err != nil {
		return nil, err
	}
	return &EigenPodManagerPodDeployedIterator{contract: _EigenPodManager.contract, event: "PodDeployed", logs: logs, sub: sub}, nil
}

// WatchPodDeployed is a free log subscription operation binding the contract event 0x21c99d0db02213c32fff5b05cf0a718ab5f858802b91498f80d82270289d856a.
//
// Solidity: event PodDeployed(address indexed eigenPod, address indexed podOwner)
func (_EigenPodManager *EigenPodManagerFilterer) WatchPodDeployed(opts *bind.WatchOpts, sink chan<- *EigenPodManagerPodDeployed, eigenPod []common.Address, podOwner []common.Address) (event.Subscription, error) {

	var eigenPodRule []interface{}
	for _, eigenPodItem := range eigenPod {
		eigenPodRule = append(eigenPodRule, eigenPodItem)
	}
	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _EigenPodManager.contract.WatchLogs(opts, "PodDeployed", eigenPodRule, podOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EigenPodManagerPodDeployed)
				if err := _EigenPodManager.contract.UnpackLog(event, "PodDeployed", log); err != nil {
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

// ParsePodDeployed is a log parse operation binding the contract event 0x21c99d0db02213c32fff5b05cf0a718ab5f858802b91498f80d82270289d856a.
//
// Solidity: event PodDeployed(address indexed eigenPod, address indexed podOwner)
func (_EigenPodManager *EigenPodManagerFilterer) ParsePodDeployed(log types.Log) (*EigenPodManagerPodDeployed, error) {
	event := new(EigenPodManagerPodDeployed)
	if err := _EigenPodManager.contract.UnpackLog(event, "PodDeployed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EigenPodManagerPodSharesUpdatedIterator is returned from FilterPodSharesUpdated and is used to iterate over the raw logs and unpacked data for PodSharesUpdated events raised by the EigenPodManager contract.
type EigenPodManagerPodSharesUpdatedIterator struct {
	Event *EigenPodManagerPodSharesUpdated // Event containing the contract specifics and raw log

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
func (it *EigenPodManagerPodSharesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EigenPodManagerPodSharesUpdated)
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
		it.Event = new(EigenPodManagerPodSharesUpdated)
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
func (it *EigenPodManagerPodSharesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EigenPodManagerPodSharesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EigenPodManagerPodSharesUpdated represents a PodSharesUpdated event raised by the EigenPodManager contract.
type EigenPodManagerPodSharesUpdated struct {
	PodOwner    common.Address
	SharesDelta *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPodSharesUpdated is a free log retrieval operation binding the contract event 0x4e2b791dedccd9fb30141b088cabf5c14a8912b52f59375c95c010700b8c6193.
//
// Solidity: event PodSharesUpdated(address indexed podOwner, int256 sharesDelta)
func (_EigenPodManager *EigenPodManagerFilterer) FilterPodSharesUpdated(opts *bind.FilterOpts, podOwner []common.Address) (*EigenPodManagerPodSharesUpdatedIterator, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _EigenPodManager.contract.FilterLogs(opts, "PodSharesUpdated", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return &EigenPodManagerPodSharesUpdatedIterator{contract: _EigenPodManager.contract, event: "PodSharesUpdated", logs: logs, sub: sub}, nil
}

// WatchPodSharesUpdated is a free log subscription operation binding the contract event 0x4e2b791dedccd9fb30141b088cabf5c14a8912b52f59375c95c010700b8c6193.
//
// Solidity: event PodSharesUpdated(address indexed podOwner, int256 sharesDelta)
func (_EigenPodManager *EigenPodManagerFilterer) WatchPodSharesUpdated(opts *bind.WatchOpts, sink chan<- *EigenPodManagerPodSharesUpdated, podOwner []common.Address) (event.Subscription, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _EigenPodManager.contract.WatchLogs(opts, "PodSharesUpdated", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EigenPodManagerPodSharesUpdated)
				if err := _EigenPodManager.contract.UnpackLog(event, "PodSharesUpdated", log); err != nil {
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

// ParsePodSharesUpdated is a log parse operation binding the contract event 0x4e2b791dedccd9fb30141b088cabf5c14a8912b52f59375c95c010700b8c6193.
//
// Solidity: event PodSharesUpdated(address indexed podOwner, int256 sharesDelta)
func (_EigenPodManager *EigenPodManagerFilterer) ParsePodSharesUpdated(log types.Log) (*EigenPodManagerPodSharesUpdated, error) {
	event := new(EigenPodManagerPodSharesUpdated)
	if err := _EigenPodManager.contract.UnpackLog(event, "PodSharesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EigenPodManagerUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the EigenPodManager contract.
type EigenPodManagerUnpausedIterator struct {
	Event *EigenPodManagerUnpaused // Event containing the contract specifics and raw log

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
func (it *EigenPodManagerUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EigenPodManagerUnpaused)
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
		it.Event = new(EigenPodManagerUnpaused)
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
func (it *EigenPodManagerUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EigenPodManagerUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EigenPodManagerUnpaused represents a Unpaused event raised by the EigenPodManager contract.
type EigenPodManagerUnpaused struct {
	Account         common.Address
	NewPausedStatus *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x3582d1828e26bf56bd801502bc021ac0bc8afb57c826e4986b45593c8fad389c.
//
// Solidity: event Unpaused(address indexed account, uint256 newPausedStatus)
func (_EigenPodManager *EigenPodManagerFilterer) FilterUnpaused(opts *bind.FilterOpts, account []common.Address) (*EigenPodManagerUnpausedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _EigenPodManager.contract.FilterLogs(opts, "Unpaused", accountRule)
	if err != nil {
		return nil, err
	}
	return &EigenPodManagerUnpausedIterator{contract: _EigenPodManager.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x3582d1828e26bf56bd801502bc021ac0bc8afb57c826e4986b45593c8fad389c.
//
// Solidity: event Unpaused(address indexed account, uint256 newPausedStatus)
func (_EigenPodManager *EigenPodManagerFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *EigenPodManagerUnpaused, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _EigenPodManager.contract.WatchLogs(opts, "Unpaused", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EigenPodManagerUnpaused)
				if err := _EigenPodManager.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

// ParseUnpaused is a log parse operation binding the contract event 0x3582d1828e26bf56bd801502bc021ac0bc8afb57c826e4986b45593c8fad389c.
//
// Solidity: event Unpaused(address indexed account, uint256 newPausedStatus)
func (_EigenPodManager *EigenPodManagerFilterer) ParseUnpaused(log types.Log) (*EigenPodManagerUnpaused, error) {
	event := new(EigenPodManagerUnpaused)
	if err := _EigenPodManager.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

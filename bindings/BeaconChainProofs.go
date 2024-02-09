// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractBeaconChainProofs

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

// BeaconChainProofsMetaData contains all meta data concerning the BeaconChainProofs contract.
var BeaconChainProofsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"latestBlockRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"beaconStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"stateRootProof\",\"type\":\"bytes\"}],\"name\":\"verifyStateRootAgainstLatestBlockRoot\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"beaconStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"validatorFields\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"validatorFieldsProof\",\"type\":\"bytes\"},{\"internalType\":\"uint40\",\"name\":\"validatorIndex\",\"type\":\"uint40\"}],\"name\":\"verifyValidatorFields\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"beaconStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"withdrawalFields\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"withdrawalProof\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"slotProof\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"executionPayloadProof\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"timestampProof\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"historicalSummaryBlockRootProof\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"blockRootIndex\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"historicalSummaryIndex\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"withdrawalIndex\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"blockRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"slotRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"timestampRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"executionPayloadRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structBeaconChainProofs.WithdrawalProof\",\"name\":\"withdrawalProof\",\"type\":\"tuple\"}],\"name\":\"verifyWithdrawal\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// BeaconChainProofsABI is the input ABI used to generate the binding from.
// Deprecated: Use BeaconChainProofsMetaData.ABI instead.
var BeaconChainProofsABI = BeaconChainProofsMetaData.ABI

// BeaconChainProofs is an auto generated Go binding around an Ethereum contract.
type BeaconChainProofs struct {
	BeaconChainProofsCaller     // Read-only binding to the contract
	BeaconChainProofsTransactor // Write-only binding to the contract
	BeaconChainProofsFilterer   // Log filterer for contract events
}

// BeaconChainProofsCaller is an auto generated read-only Go binding around an Ethereum contract.
type BeaconChainProofsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BeaconChainProofsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BeaconChainProofsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BeaconChainProofsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BeaconChainProofsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BeaconChainProofsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BeaconChainProofsSession struct {
	Contract     *BeaconChainProofs // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// BeaconChainProofsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BeaconChainProofsCallerSession struct {
	Contract *BeaconChainProofsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// BeaconChainProofsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BeaconChainProofsTransactorSession struct {
	Contract     *BeaconChainProofsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// BeaconChainProofsRaw is an auto generated low-level Go binding around an Ethereum contract.
type BeaconChainProofsRaw struct {
	Contract *BeaconChainProofs // Generic contract binding to access the raw methods on
}

// BeaconChainProofsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BeaconChainProofsCallerRaw struct {
	Contract *BeaconChainProofsCaller // Generic read-only contract binding to access the raw methods on
}

// BeaconChainProofsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BeaconChainProofsTransactorRaw struct {
	Contract *BeaconChainProofsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBeaconChainProofs creates a new instance of BeaconChainProofs, bound to a specific deployed contract.
func NewBeaconChainProofs(address common.Address, backend bind.ContractBackend) (*BeaconChainProofs, error) {
	contract, err := bindBeaconChainProofs(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BeaconChainProofs{BeaconChainProofsCaller: BeaconChainProofsCaller{contract: contract}, BeaconChainProofsTransactor: BeaconChainProofsTransactor{contract: contract}, BeaconChainProofsFilterer: BeaconChainProofsFilterer{contract: contract}}, nil
}

// NewBeaconChainProofsCaller creates a new read-only instance of BeaconChainProofs, bound to a specific deployed contract.
func NewBeaconChainProofsCaller(address common.Address, caller bind.ContractCaller) (*BeaconChainProofsCaller, error) {
	contract, err := bindBeaconChainProofs(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BeaconChainProofsCaller{contract: contract}, nil
}

// NewBeaconChainProofsTransactor creates a new write-only instance of BeaconChainProofs, bound to a specific deployed contract.
func NewBeaconChainProofsTransactor(address common.Address, transactor bind.ContractTransactor) (*BeaconChainProofsTransactor, error) {
	contract, err := bindBeaconChainProofs(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BeaconChainProofsTransactor{contract: contract}, nil
}

// NewBeaconChainProofsFilterer creates a new log filterer instance of BeaconChainProofs, bound to a specific deployed contract.
func NewBeaconChainProofsFilterer(address common.Address, filterer bind.ContractFilterer) (*BeaconChainProofsFilterer, error) {
	contract, err := bindBeaconChainProofs(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BeaconChainProofsFilterer{contract: contract}, nil
}

// bindBeaconChainProofs binds a generic wrapper to an already deployed contract.
func bindBeaconChainProofs(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BeaconChainProofsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BeaconChainProofs *BeaconChainProofsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BeaconChainProofs.Contract.BeaconChainProofsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BeaconChainProofs *BeaconChainProofsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconChainProofs.Contract.BeaconChainProofsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BeaconChainProofs *BeaconChainProofsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BeaconChainProofs.Contract.BeaconChainProofsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BeaconChainProofs *BeaconChainProofsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BeaconChainProofs.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BeaconChainProofs *BeaconChainProofsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconChainProofs.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BeaconChainProofs *BeaconChainProofsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BeaconChainProofs.Contract.contract.Transact(opts, method, params...)
}

// VerifyStateRootAgainstLatestBlockRoot is a free data retrieval call binding the contract method 0x9cdee1f8.
//
// Solidity: function verifyStateRootAgainstLatestBlockRoot(bytes32 latestBlockRoot, bytes32 beaconStateRoot, bytes stateRootProof) view returns()
func (_BeaconChainProofs *BeaconChainProofsCaller) VerifyStateRootAgainstLatestBlockRoot(opts *bind.CallOpts, latestBlockRoot [32]byte, beaconStateRoot [32]byte, stateRootProof []byte) error {
	var out []interface{}
	err := _BeaconChainProofs.contract.Call(opts, &out, "verifyStateRootAgainstLatestBlockRoot", latestBlockRoot, beaconStateRoot, stateRootProof)

	if err != nil {
		return err
	}

	return err

}

// VerifyStateRootAgainstLatestBlockRoot is a free data retrieval call binding the contract method 0x9cdee1f8.
//
// Solidity: function verifyStateRootAgainstLatestBlockRoot(bytes32 latestBlockRoot, bytes32 beaconStateRoot, bytes stateRootProof) view returns()
func (_BeaconChainProofs *BeaconChainProofsSession) VerifyStateRootAgainstLatestBlockRoot(latestBlockRoot [32]byte, beaconStateRoot [32]byte, stateRootProof []byte) error {
	return _BeaconChainProofs.Contract.VerifyStateRootAgainstLatestBlockRoot(&_BeaconChainProofs.CallOpts, latestBlockRoot, beaconStateRoot, stateRootProof)
}

// VerifyStateRootAgainstLatestBlockRoot is a free data retrieval call binding the contract method 0x9cdee1f8.
//
// Solidity: function verifyStateRootAgainstLatestBlockRoot(bytes32 latestBlockRoot, bytes32 beaconStateRoot, bytes stateRootProof) view returns()
func (_BeaconChainProofs *BeaconChainProofsCallerSession) VerifyStateRootAgainstLatestBlockRoot(latestBlockRoot [32]byte, beaconStateRoot [32]byte, stateRootProof []byte) error {
	return _BeaconChainProofs.Contract.VerifyStateRootAgainstLatestBlockRoot(&_BeaconChainProofs.CallOpts, latestBlockRoot, beaconStateRoot, stateRootProof)
}

// VerifyValidatorFields is a free data retrieval call binding the contract method 0x256f222b.
//
// Solidity: function verifyValidatorFields(bytes32 beaconStateRoot, bytes32[] validatorFields, bytes validatorFieldsProof, uint40 validatorIndex) view returns()
func (_BeaconChainProofs *BeaconChainProofsCaller) VerifyValidatorFields(opts *bind.CallOpts, beaconStateRoot [32]byte, validatorFields [][32]byte, validatorFieldsProof []byte, validatorIndex *big.Int) error {
	var out []interface{}
	err := _BeaconChainProofs.contract.Call(opts, &out, "verifyValidatorFields", beaconStateRoot, validatorFields, validatorFieldsProof, validatorIndex)

	if err != nil {
		return err
	}

	return err

}

// VerifyValidatorFields is a free data retrieval call binding the contract method 0x256f222b.
//
// Solidity: function verifyValidatorFields(bytes32 beaconStateRoot, bytes32[] validatorFields, bytes validatorFieldsProof, uint40 validatorIndex) view returns()
func (_BeaconChainProofs *BeaconChainProofsSession) VerifyValidatorFields(beaconStateRoot [32]byte, validatorFields [][32]byte, validatorFieldsProof []byte, validatorIndex *big.Int) error {
	return _BeaconChainProofs.Contract.VerifyValidatorFields(&_BeaconChainProofs.CallOpts, beaconStateRoot, validatorFields, validatorFieldsProof, validatorIndex)
}

// VerifyValidatorFields is a free data retrieval call binding the contract method 0x256f222b.
//
// Solidity: function verifyValidatorFields(bytes32 beaconStateRoot, bytes32[] validatorFields, bytes validatorFieldsProof, uint40 validatorIndex) view returns()
func (_BeaconChainProofs *BeaconChainProofsCallerSession) VerifyValidatorFields(beaconStateRoot [32]byte, validatorFields [][32]byte, validatorFieldsProof []byte, validatorIndex *big.Int) error {
	return _BeaconChainProofs.Contract.VerifyValidatorFields(&_BeaconChainProofs.CallOpts, beaconStateRoot, validatorFields, validatorFieldsProof, validatorIndex)
}

// VerifyWithdrawal is a free data retrieval call binding the contract method 0xbf03617a.
//
// Solidity: function verifyWithdrawal(bytes32 beaconStateRoot, bytes32[] withdrawalFields, (bytes,bytes,bytes,bytes,bytes,uint64,uint64,uint64,bytes32,bytes32,bytes32,bytes32) withdrawalProof) view returns()
func (_BeaconChainProofs *BeaconChainProofsCaller) VerifyWithdrawal(opts *bind.CallOpts, beaconStateRoot [32]byte, withdrawalFields [][32]byte, withdrawalProof BeaconChainProofsWithdrawalProof) error {
	var out []interface{}
	err := _BeaconChainProofs.contract.Call(opts, &out, "verifyWithdrawal", beaconStateRoot, withdrawalFields, withdrawalProof)

	if err != nil {
		return err
	}

	return err

}

// VerifyWithdrawal is a free data retrieval call binding the contract method 0xbf03617a.
//
// Solidity: function verifyWithdrawal(bytes32 beaconStateRoot, bytes32[] withdrawalFields, (bytes,bytes,bytes,bytes,bytes,uint64,uint64,uint64,bytes32,bytes32,bytes32,bytes32) withdrawalProof) view returns()
func (_BeaconChainProofs *BeaconChainProofsSession) VerifyWithdrawal(beaconStateRoot [32]byte, withdrawalFields [][32]byte, withdrawalProof BeaconChainProofsWithdrawalProof) error {
	return _BeaconChainProofs.Contract.VerifyWithdrawal(&_BeaconChainProofs.CallOpts, beaconStateRoot, withdrawalFields, withdrawalProof)
}

// VerifyWithdrawal is a free data retrieval call binding the contract method 0xbf03617a.
//
// Solidity: function verifyWithdrawal(bytes32 beaconStateRoot, bytes32[] withdrawalFields, (bytes,bytes,bytes,bytes,bytes,uint64,uint64,uint64,bytes32,bytes32,bytes32,bytes32) withdrawalProof) view returns()
func (_BeaconChainProofs *BeaconChainProofsCallerSession) VerifyWithdrawal(beaconStateRoot [32]byte, withdrawalFields [][32]byte, withdrawalProof BeaconChainProofsWithdrawalProof) error {
	return _BeaconChainProofs.Contract.VerifyWithdrawal(&_BeaconChainProofs.CallOpts, beaconStateRoot, withdrawalFields, withdrawalProof)
}

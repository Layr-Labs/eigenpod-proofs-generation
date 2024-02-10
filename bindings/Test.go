// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package beaconchainproofs

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

// TestMetaData contains all meta data concerning the Test contract.
var TestMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"latestBlockRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"beaconStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"stateRootProof\",\"type\":\"bytes\"}],\"name\":\"verifyStateRootAgainstLatestBlockRoot\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"beaconStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"validatorFields\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"validatorFieldsProof\",\"type\":\"bytes\"},{\"internalType\":\"uint40\",\"name\":\"validatorIndex\",\"type\":\"uint40\"}],\"name\":\"verifyValidatorFields\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"beaconStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"withdrawalFields\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint64\",\"name\":\"denebForkTimestamp\",\"type\":\"uint64\"}],\"name\":\"verifyWithdrawal\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// TestABI is the input ABI used to generate the binding from.
// Deprecated: Use TestMetaData.ABI instead.
var TestABI = TestMetaData.ABI

// Test is an auto generated Go binding around an Ethereum contract.
type Test struct {
	TestCaller     // Read-only binding to the contract
	TestTransactor // Write-only binding to the contract
	TestFilterer   // Log filterer for contract events
}

// TestCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestSession struct {
	Contract     *Test             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestCallerSession struct {
	Contract *TestCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// TestTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestTransactorSession struct {
	Contract     *TestTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestRaw struct {
	Contract *Test // Generic contract binding to access the raw methods on
}

// TestCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestCallerRaw struct {
	Contract *TestCaller // Generic read-only contract binding to access the raw methods on
}

// TestTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestTransactorRaw struct {
	Contract *TestTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTest creates a new instance of Test, bound to a specific deployed contract.
func NewTest(address common.Address, backend bind.ContractBackend) (*Test, error) {
	contract, err := bindTest(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Test{TestCaller: TestCaller{contract: contract}, TestTransactor: TestTransactor{contract: contract}, TestFilterer: TestFilterer{contract: contract}}, nil
}

// NewTestCaller creates a new read-only instance of Test, bound to a specific deployed contract.
func NewTestCaller(address common.Address, caller bind.ContractCaller) (*TestCaller, error) {
	contract, err := bindTest(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestCaller{contract: contract}, nil
}

// NewTestTransactor creates a new write-only instance of Test, bound to a specific deployed contract.
func NewTestTransactor(address common.Address, transactor bind.ContractTransactor) (*TestTransactor, error) {
	contract, err := bindTest(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestTransactor{contract: contract}, nil
}

// NewTestFilterer creates a new log filterer instance of Test, bound to a specific deployed contract.
func NewTestFilterer(address common.Address, filterer bind.ContractFilterer) (*TestFilterer, error) {
	contract, err := bindTest(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestFilterer{contract: contract}, nil
}

// bindTest binds a generic wrapper to an already deployed contract.
func bindTest(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TestMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Test *TestRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Test.Contract.TestCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Test *TestRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Test.Contract.TestTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Test *TestRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Test.Contract.TestTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Test *TestCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Test.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Test *TestTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Test.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Test *TestTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Test.Contract.contract.Transact(opts, method, params...)
}

// VerifyStateRootAgainstLatestBlockRoot is a free data retrieval call binding the contract method 0x9cdee1f8.
//
// Solidity: function verifyStateRootAgainstLatestBlockRoot(bytes32 latestBlockRoot, bytes32 beaconStateRoot, bytes stateRootProof) view returns()
func (_Test *TestCaller) VerifyStateRootAgainstLatestBlockRoot(opts *bind.CallOpts, latestBlockRoot [32]byte, beaconStateRoot [32]byte, stateRootProof []byte) error {
	var out []interface{}
	err := _Test.contract.Call(opts, &out, "verifyStateRootAgainstLatestBlockRoot", latestBlockRoot, beaconStateRoot, stateRootProof)

	if err != nil {
		return err
	}

	return err

}

// VerifyStateRootAgainstLatestBlockRoot is a free data retrieval call binding the contract method 0x9cdee1f8.
//
// Solidity: function verifyStateRootAgainstLatestBlockRoot(bytes32 latestBlockRoot, bytes32 beaconStateRoot, bytes stateRootProof) view returns()
func (_Test *TestSession) VerifyStateRootAgainstLatestBlockRoot(latestBlockRoot [32]byte, beaconStateRoot [32]byte, stateRootProof []byte) error {
	return _Test.Contract.VerifyStateRootAgainstLatestBlockRoot(&_Test.CallOpts, latestBlockRoot, beaconStateRoot, stateRootProof)
}

// VerifyStateRootAgainstLatestBlockRoot is a free data retrieval call binding the contract method 0x9cdee1f8.
//
// Solidity: function verifyStateRootAgainstLatestBlockRoot(bytes32 latestBlockRoot, bytes32 beaconStateRoot, bytes stateRootProof) view returns()
func (_Test *TestCallerSession) VerifyStateRootAgainstLatestBlockRoot(latestBlockRoot [32]byte, beaconStateRoot [32]byte, stateRootProof []byte) error {
	return _Test.Contract.VerifyStateRootAgainstLatestBlockRoot(&_Test.CallOpts, latestBlockRoot, beaconStateRoot, stateRootProof)
}

// VerifyValidatorFields is a free data retrieval call binding the contract method 0x256f222b.
//
// Solidity: function verifyValidatorFields(bytes32 beaconStateRoot, bytes32[] validatorFields, bytes validatorFieldsProof, uint40 validatorIndex) view returns()
func (_Test *TestCaller) VerifyValidatorFields(opts *bind.CallOpts, beaconStateRoot [32]byte, validatorFields [][32]byte, validatorFieldsProof []byte, validatorIndex *big.Int) error {
	var out []interface{}
	err := _Test.contract.Call(opts, &out, "verifyValidatorFields", beaconStateRoot, validatorFields, validatorFieldsProof, validatorIndex)

	if err != nil {
		return err
	}

	return err

}

// VerifyValidatorFields is a free data retrieval call binding the contract method 0x256f222b.
//
// Solidity: function verifyValidatorFields(bytes32 beaconStateRoot, bytes32[] validatorFields, bytes validatorFieldsProof, uint40 validatorIndex) view returns()
func (_Test *TestSession) VerifyValidatorFields(beaconStateRoot [32]byte, validatorFields [][32]byte, validatorFieldsProof []byte, validatorIndex *big.Int) error {
	return _Test.Contract.VerifyValidatorFields(&_Test.CallOpts, beaconStateRoot, validatorFields, validatorFieldsProof, validatorIndex)
}

// VerifyValidatorFields is a free data retrieval call binding the contract method 0x256f222b.
//
// Solidity: function verifyValidatorFields(bytes32 beaconStateRoot, bytes32[] validatorFields, bytes validatorFieldsProof, uint40 validatorIndex) view returns()
func (_Test *TestCallerSession) VerifyValidatorFields(beaconStateRoot [32]byte, validatorFields [][32]byte, validatorFieldsProof []byte, validatorIndex *big.Int) error {
	return _Test.Contract.VerifyValidatorFields(&_Test.CallOpts, beaconStateRoot, validatorFields, validatorFieldsProof, validatorIndex)
}

// VerifyWithdrawal is a free data retrieval call binding the contract method 0xf4e0e350.
//
// Solidity: function verifyWithdrawal(bytes32 beaconStateRoot, bytes32[] withdrawalFields, uint64 denebForkTimestamp) view returns()
func (_Test *TestCaller) VerifyWithdrawal(opts *bind.CallOpts, beaconStateRoot [32]byte, withdrawalFields [][32]byte, denebForkTimestamp uint64) error {
	var out []interface{}
	err := _Test.contract.Call(opts, &out, "verifyWithdrawal", beaconStateRoot, withdrawalFields, denebForkTimestamp)

	if err != nil {
		return err
	}

	return err

}

// VerifyWithdrawal is a free data retrieval call binding the contract method 0xf4e0e350.
//
// Solidity: function verifyWithdrawal(bytes32 beaconStateRoot, bytes32[] withdrawalFields, uint64 denebForkTimestamp) view returns()
func (_Test *TestSession) VerifyWithdrawal(beaconStateRoot [32]byte, withdrawalFields [][32]byte, denebForkTimestamp uint64) error {
	return _Test.Contract.VerifyWithdrawal(&_Test.CallOpts, beaconStateRoot, withdrawalFields, denebForkTimestamp)
}

// VerifyWithdrawal is a free data retrieval call binding the contract method 0xf4e0e350.
//
// Solidity: function verifyWithdrawal(bytes32 beaconStateRoot, bytes32[] withdrawalFields, uint64 denebForkTimestamp) view returns()
func (_Test *TestCallerSession) VerifyWithdrawal(beaconStateRoot [32]byte, withdrawalFields [][32]byte, denebForkTimestamp uint64) error {
	return _Test.Contract.VerifyWithdrawal(&_Test.CallOpts, beaconStateRoot, withdrawalFields, denebForkTimestamp)
}

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

// BeaconChainProofsContractWithdrawalProof is an auto generated low-level Go binding around an user-defined struct.
type BeaconChainProofsContractWithdrawalProof struct {
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

// BeaconChainProofsTestMetaData contains all meta data concerning the BeaconChainProofsTest contract.
var BeaconChainProofsTestMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"latestBlockRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"beaconStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"stateRootProof\",\"type\":\"bytes\"}],\"name\":\"verifyStateRootAgainstLatestBlockRoot\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"beaconStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"validatorFields\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"validatorFieldsProof\",\"type\":\"bytes\"},{\"internalType\":\"uint40\",\"name\":\"validatorIndex\",\"type\":\"uint40\"}],\"name\":\"verifyValidatorFields\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"beaconStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"withdrawalFields\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"withdrawalProof\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"slotProof\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"executionPayloadProof\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"timestampProof\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"historicalSummaryBlockRootProof\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"blockRootIndex\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"historicalSummaryIndex\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"withdrawalIndex\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"blockRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"slotRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"timestampRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"executionPayloadRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structBeaconChainProofsContract.WithdrawalProof\",\"name\":\"withdrawalProof\",\"type\":\"tuple\"},{\"internalType\":\"uint64\",\"name\":\"denebForkTimestamp\",\"type\":\"uint64\"}],\"name\":\"verifyWithdrawal\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// BeaconChainProofsTestABI is the input ABI used to generate the binding from.
// Deprecated: Use BeaconChainProofsTestMetaData.ABI instead.
var BeaconChainProofsTestABI = BeaconChainProofsTestMetaData.ABI

// BeaconChainProofsTest is an auto generated Go binding around an Ethereum contract.
type BeaconChainProofsTest struct {
	BeaconChainProofsTestCaller     // Read-only binding to the contract
	BeaconChainProofsTestTransactor // Write-only binding to the contract
	BeaconChainProofsTestFilterer   // Log filterer for contract events
}

// BeaconChainProofsTestCaller is an auto generated read-only Go binding around an Ethereum contract.
type BeaconChainProofsTestCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BeaconChainProofsTestTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BeaconChainProofsTestTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BeaconChainProofsTestFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BeaconChainProofsTestFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BeaconChainProofsTestSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BeaconChainProofsTestSession struct {
	Contract     *BeaconChainProofsTest // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// BeaconChainProofsTestCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BeaconChainProofsTestCallerSession struct {
	Contract *BeaconChainProofsTestCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// BeaconChainProofsTestTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BeaconChainProofsTestTransactorSession struct {
	Contract     *BeaconChainProofsTestTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// BeaconChainProofsTestRaw is an auto generated low-level Go binding around an Ethereum contract.
type BeaconChainProofsTestRaw struct {
	Contract *BeaconChainProofsTest // Generic contract binding to access the raw methods on
}

// BeaconChainProofsTestCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BeaconChainProofsTestCallerRaw struct {
	Contract *BeaconChainProofsTestCaller // Generic read-only contract binding to access the raw methods on
}

// BeaconChainProofsTestTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BeaconChainProofsTestTransactorRaw struct {
	Contract *BeaconChainProofsTestTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBeaconChainProofsTest creates a new instance of BeaconChainProofsTest, bound to a specific deployed contract.
func NewBeaconChainProofsTest(address common.Address, backend bind.ContractBackend) (*BeaconChainProofsTest, error) {
	contract, err := bindBeaconChainProofsTest(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BeaconChainProofsTest{BeaconChainProofsTestCaller: BeaconChainProofsTestCaller{contract: contract}, BeaconChainProofsTestTransactor: BeaconChainProofsTestTransactor{contract: contract}, BeaconChainProofsTestFilterer: BeaconChainProofsTestFilterer{contract: contract}}, nil
}

// NewBeaconChainProofsTestCaller creates a new read-only instance of BeaconChainProofsTest, bound to a specific deployed contract.
func NewBeaconChainProofsTestCaller(address common.Address, caller bind.ContractCaller) (*BeaconChainProofsTestCaller, error) {
	contract, err := bindBeaconChainProofsTest(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BeaconChainProofsTestCaller{contract: contract}, nil
}

// NewBeaconChainProofsTestTransactor creates a new write-only instance of BeaconChainProofsTest, bound to a specific deployed contract.
func NewBeaconChainProofsTestTransactor(address common.Address, transactor bind.ContractTransactor) (*BeaconChainProofsTestTransactor, error) {
	contract, err := bindBeaconChainProofsTest(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BeaconChainProofsTestTransactor{contract: contract}, nil
}

// NewBeaconChainProofsTestFilterer creates a new log filterer instance of BeaconChainProofsTest, bound to a specific deployed contract.
func NewBeaconChainProofsTestFilterer(address common.Address, filterer bind.ContractFilterer) (*BeaconChainProofsTestFilterer, error) {
	contract, err := bindBeaconChainProofsTest(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BeaconChainProofsTestFilterer{contract: contract}, nil
}

// bindBeaconChainProofsTest binds a generic wrapper to an already deployed contract.
func bindBeaconChainProofsTest(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BeaconChainProofsTestMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BeaconChainProofsTest *BeaconChainProofsTestRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BeaconChainProofsTest.Contract.BeaconChainProofsTestCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BeaconChainProofsTest *BeaconChainProofsTestRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconChainProofsTest.Contract.BeaconChainProofsTestTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BeaconChainProofsTest *BeaconChainProofsTestRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BeaconChainProofsTest.Contract.BeaconChainProofsTestTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BeaconChainProofsTest *BeaconChainProofsTestCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BeaconChainProofsTest.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BeaconChainProofsTest *BeaconChainProofsTestTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconChainProofsTest.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BeaconChainProofsTest *BeaconChainProofsTestTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BeaconChainProofsTest.Contract.contract.Transact(opts, method, params...)
}

// VerifyStateRootAgainstLatestBlockRoot is a free data retrieval call binding the contract method 0x9cdee1f8.
//
// Solidity: function verifyStateRootAgainstLatestBlockRoot(bytes32 latestBlockRoot, bytes32 beaconStateRoot, bytes stateRootProof) view returns()
func (_BeaconChainProofsTest *BeaconChainProofsTestCaller) VerifyStateRootAgainstLatestBlockRoot(opts *bind.CallOpts, latestBlockRoot [32]byte, beaconStateRoot [32]byte, stateRootProof []byte) error {
	var out []interface{}
	err := _BeaconChainProofsTest.contract.Call(opts, &out, "verifyStateRootAgainstLatestBlockRoot", latestBlockRoot, beaconStateRoot, stateRootProof)

	if err != nil {
		return err
	}

	return err

}

// VerifyStateRootAgainstLatestBlockRoot is a free data retrieval call binding the contract method 0x9cdee1f8.
//
// Solidity: function verifyStateRootAgainstLatestBlockRoot(bytes32 latestBlockRoot, bytes32 beaconStateRoot, bytes stateRootProof) view returns()
func (_BeaconChainProofsTest *BeaconChainProofsTestSession) VerifyStateRootAgainstLatestBlockRoot(latestBlockRoot [32]byte, beaconStateRoot [32]byte, stateRootProof []byte) error {
	return _BeaconChainProofsTest.Contract.VerifyStateRootAgainstLatestBlockRoot(&_BeaconChainProofsTest.CallOpts, latestBlockRoot, beaconStateRoot, stateRootProof)
}

// VerifyStateRootAgainstLatestBlockRoot is a free data retrieval call binding the contract method 0x9cdee1f8.
//
// Solidity: function verifyStateRootAgainstLatestBlockRoot(bytes32 latestBlockRoot, bytes32 beaconStateRoot, bytes stateRootProof) view returns()
func (_BeaconChainProofsTest *BeaconChainProofsTestCallerSession) VerifyStateRootAgainstLatestBlockRoot(latestBlockRoot [32]byte, beaconStateRoot [32]byte, stateRootProof []byte) error {
	return _BeaconChainProofsTest.Contract.VerifyStateRootAgainstLatestBlockRoot(&_BeaconChainProofsTest.CallOpts, latestBlockRoot, beaconStateRoot, stateRootProof)
}

// VerifyValidatorFields is a free data retrieval call binding the contract method 0x256f222b.
//
// Solidity: function verifyValidatorFields(bytes32 beaconStateRoot, bytes32[] validatorFields, bytes validatorFieldsProof, uint40 validatorIndex) view returns()
func (_BeaconChainProofsTest *BeaconChainProofsTestCaller) VerifyValidatorFields(opts *bind.CallOpts, beaconStateRoot [32]byte, validatorFields [][32]byte, validatorFieldsProof []byte, validatorIndex *big.Int) error {
	var out []interface{}
	err := _BeaconChainProofsTest.contract.Call(opts, &out, "verifyValidatorFields", beaconStateRoot, validatorFields, validatorFieldsProof, validatorIndex)

	if err != nil {
		return err
	}

	return err

}

// VerifyValidatorFields is a free data retrieval call binding the contract method 0x256f222b.
//
// Solidity: function verifyValidatorFields(bytes32 beaconStateRoot, bytes32[] validatorFields, bytes validatorFieldsProof, uint40 validatorIndex) view returns()
func (_BeaconChainProofsTest *BeaconChainProofsTestSession) VerifyValidatorFields(beaconStateRoot [32]byte, validatorFields [][32]byte, validatorFieldsProof []byte, validatorIndex *big.Int) error {
	return _BeaconChainProofsTest.Contract.VerifyValidatorFields(&_BeaconChainProofsTest.CallOpts, beaconStateRoot, validatorFields, validatorFieldsProof, validatorIndex)
}

// VerifyValidatorFields is a free data retrieval call binding the contract method 0x256f222b.
//
// Solidity: function verifyValidatorFields(bytes32 beaconStateRoot, bytes32[] validatorFields, bytes validatorFieldsProof, uint40 validatorIndex) view returns()
func (_BeaconChainProofsTest *BeaconChainProofsTestCallerSession) VerifyValidatorFields(beaconStateRoot [32]byte, validatorFields [][32]byte, validatorFieldsProof []byte, validatorIndex *big.Int) error {
	return _BeaconChainProofsTest.Contract.VerifyValidatorFields(&_BeaconChainProofsTest.CallOpts, beaconStateRoot, validatorFields, validatorFieldsProof, validatorIndex)
}

// VerifyWithdrawal is a free data retrieval call binding the contract method 0x17941319.
//
// Solidity: function verifyWithdrawal(bytes32 beaconStateRoot, bytes32[] withdrawalFields, (bytes,bytes,bytes,bytes,bytes,uint64,uint64,uint64,bytes32,bytes32,bytes32,bytes32) withdrawalProof, uint64 denebForkTimestamp) view returns()
func (_BeaconChainProofsTest *BeaconChainProofsTestCaller) VerifyWithdrawal(opts *bind.CallOpts, beaconStateRoot [32]byte, withdrawalFields [][32]byte, withdrawalProof BeaconChainProofsContractWithdrawalProof, denebForkTimestamp uint64) error {
	var out []interface{}
	err := _BeaconChainProofsTest.contract.Call(opts, &out, "verifyWithdrawal", beaconStateRoot, withdrawalFields, withdrawalProof, denebForkTimestamp)

	if err != nil {
		return err
	}

	return err

}

// VerifyWithdrawal is a free data retrieval call binding the contract method 0x17941319.
//
// Solidity: function verifyWithdrawal(bytes32 beaconStateRoot, bytes32[] withdrawalFields, (bytes,bytes,bytes,bytes,bytes,uint64,uint64,uint64,bytes32,bytes32,bytes32,bytes32) withdrawalProof, uint64 denebForkTimestamp) view returns()
func (_BeaconChainProofsTest *BeaconChainProofsTestSession) VerifyWithdrawal(beaconStateRoot [32]byte, withdrawalFields [][32]byte, withdrawalProof BeaconChainProofsContractWithdrawalProof, denebForkTimestamp uint64) error {
	return _BeaconChainProofsTest.Contract.VerifyWithdrawal(&_BeaconChainProofsTest.CallOpts, beaconStateRoot, withdrawalFields, withdrawalProof, denebForkTimestamp)
}

// VerifyWithdrawal is a free data retrieval call binding the contract method 0x17941319.
//
// Solidity: function verifyWithdrawal(bytes32 beaconStateRoot, bytes32[] withdrawalFields, (bytes,bytes,bytes,bytes,bytes,uint64,uint64,uint64,bytes32,bytes32,bytes32,bytes32) withdrawalProof, uint64 denebForkTimestamp) view returns()
func (_BeaconChainProofsTest *BeaconChainProofsTestCallerSession) VerifyWithdrawal(beaconStateRoot [32]byte, withdrawalFields [][32]byte, withdrawalProof BeaconChainProofsContractWithdrawalProof, denebForkTimestamp uint64) error {
	return _BeaconChainProofsTest.Contract.VerifyWithdrawal(&_BeaconChainProofsTest.CallOpts, beaconStateRoot, withdrawalFields, withdrawalProof, denebForkTimestamp)
}

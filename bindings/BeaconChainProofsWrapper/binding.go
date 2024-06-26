// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractBeaconChainProofsWrapper

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

// ContractBeaconChainProofsWrapperMetaData contains all meta data concerning the ContractBeaconChainProofsWrapper contract.
var ContractBeaconChainProofsWrapperMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"verifyBalanceContainer\",\"inputs\":[{\"name\":\"beaconBlockRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structBeaconChainProofs.BalanceContainerProof\",\"components\":[{\"name\":\"balanceContainerRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyStateRoot\",\"inputs\":[{\"name\":\"beaconBlockRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structBeaconChainProofs.StateRootProof\",\"components\":[{\"name\":\"beaconStateRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyValidatorBalance\",\"inputs\":[{\"name\":\"balanceContainerRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorIndex\",\"type\":\"uint40\",\"internalType\":\"uint40\"},{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structBeaconChainProofs.BalanceProof\",\"components\":[{\"name\":\"pubkeyHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"balanceRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyValidatorFields\",\"inputs\":[{\"name\":\"beaconStateRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorFields\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"validatorFieldsProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validatorIndex\",\"type\":\"uint40\",\"internalType\":\"uint40\"}],\"outputs\":[],\"stateMutability\":\"view\"}]",
	Bin: "0x608060405234801561001057600080fd5b50611006806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80630d361f3a14610051578063256f222b1461006657806331f60d4c146100795780639030a9bb1461008c575b600080fd5b61006461005f366004610c0c565b61009f565b005b610064610074366004610cb6565b6100ad565b610064610087366004610d70565b6100c3565b61006461009a366004610c0c565b6100d4565b6100a982826100de565b5050565b6100bb868686868686610264565b505050505050565b6100ce83838361047b565b50505050565b6100a982826105f0565b6100ea60056003610de4565b6100f5906020610dfc565b6101026020830183610e1b565b90501461018a5760405162461bcd60e51b8152602060048201526044602482018190527f426561636f6e436861696e50726f6f66732e76657269667942616c616e636543908201527f6f6e7461696e65723a2050726f6f662068617320696e636f7272656374206c656064820152630dccee8d60e31b608482015260a4015b60405180910390fd5b606c6101db61019c6020840184610e1b565b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525087925050853590508461074b565b61025f5760405162461bcd60e51b815260206004820152604960248201527f426561636f6e436861696e50726f6f66732e76657269667942616c616e63654360448201527f6f6e7461696e65723a20696e76616c69642062616c616e636520636f6e7461696064820152683732b910383937b7b360b91b608482015260a401610181565b505050565b600884146102df5760405162461bcd60e51b815260206004820152604e6024820152600080516020610fb183398151915260448201527f724669656c64733a2056616c696461746f72206669656c64732068617320696e60648201526d0c6dee4e4cac6e840d8cadccee8d60931b608482015260a401610181565b60056102ed60286001610de4565b6102f79190610de4565b610302906020610dfc565b82146103705760405162461bcd60e51b81526020600482015260436024820152600080516020610fb183398151915260448201527f724669656c64733a2050726f6f662068617320696e636f7272656374206c656e6064820152620cee8d60eb1b608482015260a401610181565b60006103ae86868080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525061076392505050565b9050600064ffffffffff83166103c660286001610de4565b600b901b17905061041185858080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508c925086915085905061074b565b6104715760405162461bcd60e51b815260206004820152603d6024820152600080516020610fb183398151915260448201527f724669656c64733a20496e76616c6964206d65726b6c652070726f6f660000006064820152608401610181565b5050505050505050565b600061048960266001610de4565b610494906020610dfc565b6104a16040840184610e1b565b9050146105125760405162461bcd60e51b815260206004820152604460248201819052600080516020610fb1833981519152908201527f7242616c616e63653a2050726f6f662068617320696e636f7272656374206c656064820152630dccee8d60e31b608482015260a401610181565b600061051f600485610e78565b64ffffffffff1690506105796105386040850185610e1b565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250899250505060208601358461074b565b6105d95760405162461bcd60e51b815260206004820152603e6024820152600080516020610fb183398151915260448201527f7242616c616e63653a20496e76616c6964206d65726b6c652070726f6f6600006064820152608401610181565b6105e7836020013585610a11565b95945050505050565b6105fc60036020610dfc565b6106096020830183610e1b565b90501461067e5760405162461bcd60e51b815260206004820152603d60248201527f426561636f6e436861696e50726f6f66732e7665726966795374617465526f6f60448201527f743a2050726f6f662068617320696e636f7272656374206c656e6774680000006064820152608401610181565b6106ce61068e6020830183610e1b565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508692505084359050600361074b565b6100a95760405162461bcd60e51b815260206004820152604260248201527f426561636f6e436861696e50726f6f66732e7665726966795374617465526f6f60448201527f743a20496e76616c696420737461746520726f6f74206d65726b6c652070726f60648201526137b360f11b608482015260a401610181565b600083610759868585610aa8565b1495945050505050565b600080600283516107749190610e9c565b905060008167ffffffffffffffff81111561079157610791610eb0565b6040519080825280602002602001820160405280156107ba578160200160208202803683370190505b50905060005b828110156108c1576002856107d58383610dfc565b815181106107e5576107e5610ec6565b6020026020010151868360026107fb9190610dfc565b610806906001610de4565b8151811061081657610816610ec6565b6020026020010151604051602001610838929190918252602082015260400190565b60408051601f198184030181529082905261085291610edc565b602060405180830381855afa15801561086f573d6000803e3d6000fd5b5050506040513d601f19601f820116820180604052508101906108929190610f17565b8282815181106108a4576108a4610ec6565b6020908102919091010152806108b981610f30565b9150506107c0565b506108cd600283610e9c565b91505b81156109ed5760005b828110156109da576002826108ee8383610dfc565b815181106108fe576108fe610ec6565b6020026020010151838360026109149190610dfc565b61091f906001610de4565b8151811061092f5761092f610ec6565b6020026020010151604051602001610951929190918252602082015260400190565b60408051601f198184030181529082905261096b91610edc565b602060405180830381855afa158015610988573d6000803e3d6000fd5b5050506040513d601f19601f820116820180604052508101906109ab9190610f17565b8282815181106109bd576109bd610ec6565b6020908102919091010152806109d281610f30565b9150506108d9565b506109e6600283610e9c565b91506108d0565b80600081518110610a0057610a00610ec6565b602002602001015192505050919050565b600080610a1f600484610f4b565b610a2a906040610f6f565b64ffffffffff169050610aa084821b60f881901c60e882901c61ff00161760d882901c62ff0000161760c882901c63ff000000161764ff0000000060b883901c161765ff000000000060a883901c161766ff000000000000609883901c161767ff0000000000000060889290921c919091161790565b949350505050565b60008351600014158015610ac7575060208451610ac59190610f9c565b155b610b565760405162461bcd60e51b815260206004820152605460248201527f4d65726b6c652e70726f63657373496e636c7573696f6e50726f6f665368613260448201527f35363a2070726f6f66206c656e6774682073686f756c642062652061206e6f6e60648201527316bd32b9379036bab63a34b836329037b310199960611b608482015260a401610181565b604080516020808201909252848152905b85518111610bea57610b7a600285610f9c565b610bad578151600052808601516020526020826040600060026107d05a03fa610ba257600080fd5b600284049350610bd8565b8086015160005281516020526020826040600060026107d05a03fa610bd157600080fd5b6002840493505b610be3602082610de4565b9050610b67565b5051949350505050565b600060408284031215610c0657600080fd5b50919050565b60008060408385031215610c1f57600080fd5b82359150602083013567ffffffffffffffff811115610c3d57600080fd5b610c4985828601610bf4565b9150509250929050565b60008083601f840112610c6557600080fd5b50813567ffffffffffffffff811115610c7d57600080fd5b602083019150836020828501011115610c9557600080fd5b9250929050565b803564ffffffffff81168114610cb157600080fd5b919050565b60008060008060008060808789031215610ccf57600080fd5b86359550602087013567ffffffffffffffff80821115610cee57600080fd5b818901915089601f830112610d0257600080fd5b813581811115610d1157600080fd5b8a60208260051b8501011115610d2657600080fd5b602083019750809650506040890135915080821115610d4457600080fd5b50610d5189828a01610c53565b9094509250610d64905060608801610c9c565b90509295509295509295565b600080600060608486031215610d8557600080fd5b83359250610d9560208501610c9c565b9150604084013567ffffffffffffffff811115610db157600080fd5b840160608187031215610dc357600080fd5b809150509250925092565b634e487b7160e01b600052601160045260246000fd5b60008219821115610df757610df7610dce565b500190565b6000816000190483118215151615610e1657610e16610dce565b500290565b6000808335601e19843603018112610e3257600080fd5b83018035915067ffffffffffffffff821115610e4d57600080fd5b602001915036819003821315610c9557600080fd5b634e487b7160e01b600052601260045260246000fd5b600064ffffffffff80841680610e9057610e90610e62565b92169190910492915050565b600082610eab57610eab610e62565b500490565b634e487b7160e01b600052604160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b6000825160005b81811015610efd5760208186018101518583015201610ee3565b81811115610f0c576000828501525b509190910192915050565b600060208284031215610f2957600080fd5b5051919050565b6000600019821415610f4457610f44610dce565b5060010190565b600064ffffffffff80841680610f6357610f63610e62565b92169190910692915050565b600064ffffffffff80831681851681830481118215151615610f9357610f93610dce565b02949350505050565b600082610fab57610fab610e62565b50069056fe426561636f6e436861696e50726f6f66732e76657269667956616c696461746fa26469706673582212205da3570ae5f836192d8a7d9f6b33b02a0c8db469459648fccff73b9b51138a3164736f6c634300080c0033",
}

// ContractBeaconChainProofsWrapperABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractBeaconChainProofsWrapperMetaData.ABI instead.
var ContractBeaconChainProofsWrapperABI = ContractBeaconChainProofsWrapperMetaData.ABI

// ContractBeaconChainProofsWrapperBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractBeaconChainProofsWrapperMetaData.Bin instead.
var ContractBeaconChainProofsWrapperBin = ContractBeaconChainProofsWrapperMetaData.Bin

// DeployContractBeaconChainProofsWrapper deploys a new Ethereum contract, binding an instance of ContractBeaconChainProofsWrapper to it.
func DeployContractBeaconChainProofsWrapper(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ContractBeaconChainProofsWrapper, error) {
	parsed, err := ContractBeaconChainProofsWrapperMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractBeaconChainProofsWrapperBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractBeaconChainProofsWrapper{ContractBeaconChainProofsWrapperCaller: ContractBeaconChainProofsWrapperCaller{contract: contract}, ContractBeaconChainProofsWrapperTransactor: ContractBeaconChainProofsWrapperTransactor{contract: contract}, ContractBeaconChainProofsWrapperFilterer: ContractBeaconChainProofsWrapperFilterer{contract: contract}}, nil
}

// ContractBeaconChainProofsWrapper is an auto generated Go binding around an Ethereum contract.
type ContractBeaconChainProofsWrapper struct {
	ContractBeaconChainProofsWrapperCaller     // Read-only binding to the contract
	ContractBeaconChainProofsWrapperTransactor // Write-only binding to the contract
	ContractBeaconChainProofsWrapperFilterer   // Log filterer for contract events
}

// ContractBeaconChainProofsWrapperCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractBeaconChainProofsWrapperCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractBeaconChainProofsWrapperTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractBeaconChainProofsWrapperTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractBeaconChainProofsWrapperFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractBeaconChainProofsWrapperFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractBeaconChainProofsWrapperSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractBeaconChainProofsWrapperSession struct {
	Contract     *ContractBeaconChainProofsWrapper // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                     // Call options to use throughout this session
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// ContractBeaconChainProofsWrapperCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractBeaconChainProofsWrapperCallerSession struct {
	Contract *ContractBeaconChainProofsWrapperCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                           // Call options to use throughout this session
}

// ContractBeaconChainProofsWrapperTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractBeaconChainProofsWrapperTransactorSession struct {
	Contract     *ContractBeaconChainProofsWrapperTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                           // Transaction auth options to use throughout this session
}

// ContractBeaconChainProofsWrapperRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractBeaconChainProofsWrapperRaw struct {
	Contract *ContractBeaconChainProofsWrapper // Generic contract binding to access the raw methods on
}

// ContractBeaconChainProofsWrapperCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractBeaconChainProofsWrapperCallerRaw struct {
	Contract *ContractBeaconChainProofsWrapperCaller // Generic read-only contract binding to access the raw methods on
}

// ContractBeaconChainProofsWrapperTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractBeaconChainProofsWrapperTransactorRaw struct {
	Contract *ContractBeaconChainProofsWrapperTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractBeaconChainProofsWrapper creates a new instance of ContractBeaconChainProofsWrapper, bound to a specific deployed contract.
func NewContractBeaconChainProofsWrapper(address common.Address, backend bind.ContractBackend) (*ContractBeaconChainProofsWrapper, error) {
	contract, err := bindContractBeaconChainProofsWrapper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractBeaconChainProofsWrapper{ContractBeaconChainProofsWrapperCaller: ContractBeaconChainProofsWrapperCaller{contract: contract}, ContractBeaconChainProofsWrapperTransactor: ContractBeaconChainProofsWrapperTransactor{contract: contract}, ContractBeaconChainProofsWrapperFilterer: ContractBeaconChainProofsWrapperFilterer{contract: contract}}, nil
}

// NewContractBeaconChainProofsWrapperCaller creates a new read-only instance of ContractBeaconChainProofsWrapper, bound to a specific deployed contract.
func NewContractBeaconChainProofsWrapperCaller(address common.Address, caller bind.ContractCaller) (*ContractBeaconChainProofsWrapperCaller, error) {
	contract, err := bindContractBeaconChainProofsWrapper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractBeaconChainProofsWrapperCaller{contract: contract}, nil
}

// NewContractBeaconChainProofsWrapperTransactor creates a new write-only instance of ContractBeaconChainProofsWrapper, bound to a specific deployed contract.
func NewContractBeaconChainProofsWrapperTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractBeaconChainProofsWrapperTransactor, error) {
	contract, err := bindContractBeaconChainProofsWrapper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractBeaconChainProofsWrapperTransactor{contract: contract}, nil
}

// NewContractBeaconChainProofsWrapperFilterer creates a new log filterer instance of ContractBeaconChainProofsWrapper, bound to a specific deployed contract.
func NewContractBeaconChainProofsWrapperFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractBeaconChainProofsWrapperFilterer, error) {
	contract, err := bindContractBeaconChainProofsWrapper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractBeaconChainProofsWrapperFilterer{contract: contract}, nil
}

// bindContractBeaconChainProofsWrapper binds a generic wrapper to an already deployed contract.
func bindContractBeaconChainProofsWrapper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractBeaconChainProofsWrapperMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractBeaconChainProofsWrapper *ContractBeaconChainProofsWrapperRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractBeaconChainProofsWrapper.Contract.ContractBeaconChainProofsWrapperCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractBeaconChainProofsWrapper *ContractBeaconChainProofsWrapperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractBeaconChainProofsWrapper.Contract.ContractBeaconChainProofsWrapperTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractBeaconChainProofsWrapper *ContractBeaconChainProofsWrapperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractBeaconChainProofsWrapper.Contract.ContractBeaconChainProofsWrapperTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractBeaconChainProofsWrapper *ContractBeaconChainProofsWrapperCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractBeaconChainProofsWrapper.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractBeaconChainProofsWrapper *ContractBeaconChainProofsWrapperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractBeaconChainProofsWrapper.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractBeaconChainProofsWrapper *ContractBeaconChainProofsWrapperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractBeaconChainProofsWrapper.Contract.contract.Transact(opts, method, params...)
}

// VerifyBalanceContainer is a free data retrieval call binding the contract method 0x0d361f3a.
//
// Solidity: function verifyBalanceContainer(bytes32 beaconBlockRoot, (bytes32,bytes) proof) view returns()
func (_ContractBeaconChainProofsWrapper *ContractBeaconChainProofsWrapperCaller) VerifyBalanceContainer(opts *bind.CallOpts, beaconBlockRoot [32]byte, proof BeaconChainProofsBalanceContainerProof) error {
	var out []interface{}
	err := _ContractBeaconChainProofsWrapper.contract.Call(opts, &out, "verifyBalanceContainer", beaconBlockRoot, proof)

	if err != nil {
		return err
	}

	return err

}

// VerifyBalanceContainer is a free data retrieval call binding the contract method 0x0d361f3a.
//
// Solidity: function verifyBalanceContainer(bytes32 beaconBlockRoot, (bytes32,bytes) proof) view returns()
func (_ContractBeaconChainProofsWrapper *ContractBeaconChainProofsWrapperSession) VerifyBalanceContainer(beaconBlockRoot [32]byte, proof BeaconChainProofsBalanceContainerProof) error {
	return _ContractBeaconChainProofsWrapper.Contract.VerifyBalanceContainer(&_ContractBeaconChainProofsWrapper.CallOpts, beaconBlockRoot, proof)
}

// VerifyBalanceContainer is a free data retrieval call binding the contract method 0x0d361f3a.
//
// Solidity: function verifyBalanceContainer(bytes32 beaconBlockRoot, (bytes32,bytes) proof) view returns()
func (_ContractBeaconChainProofsWrapper *ContractBeaconChainProofsWrapperCallerSession) VerifyBalanceContainer(beaconBlockRoot [32]byte, proof BeaconChainProofsBalanceContainerProof) error {
	return _ContractBeaconChainProofsWrapper.Contract.VerifyBalanceContainer(&_ContractBeaconChainProofsWrapper.CallOpts, beaconBlockRoot, proof)
}

// VerifyStateRoot is a free data retrieval call binding the contract method 0x9030a9bb.
//
// Solidity: function verifyStateRoot(bytes32 beaconBlockRoot, (bytes32,bytes) proof) view returns()
func (_ContractBeaconChainProofsWrapper *ContractBeaconChainProofsWrapperCaller) VerifyStateRoot(opts *bind.CallOpts, beaconBlockRoot [32]byte, proof BeaconChainProofsStateRootProof) error {
	var out []interface{}
	err := _ContractBeaconChainProofsWrapper.contract.Call(opts, &out, "verifyStateRoot", beaconBlockRoot, proof)

	if err != nil {
		return err
	}

	return err

}

// VerifyStateRoot is a free data retrieval call binding the contract method 0x9030a9bb.
//
// Solidity: function verifyStateRoot(bytes32 beaconBlockRoot, (bytes32,bytes) proof) view returns()
func (_ContractBeaconChainProofsWrapper *ContractBeaconChainProofsWrapperSession) VerifyStateRoot(beaconBlockRoot [32]byte, proof BeaconChainProofsStateRootProof) error {
	return _ContractBeaconChainProofsWrapper.Contract.VerifyStateRoot(&_ContractBeaconChainProofsWrapper.CallOpts, beaconBlockRoot, proof)
}

// VerifyStateRoot is a free data retrieval call binding the contract method 0x9030a9bb.
//
// Solidity: function verifyStateRoot(bytes32 beaconBlockRoot, (bytes32,bytes) proof) view returns()
func (_ContractBeaconChainProofsWrapper *ContractBeaconChainProofsWrapperCallerSession) VerifyStateRoot(beaconBlockRoot [32]byte, proof BeaconChainProofsStateRootProof) error {
	return _ContractBeaconChainProofsWrapper.Contract.VerifyStateRoot(&_ContractBeaconChainProofsWrapper.CallOpts, beaconBlockRoot, proof)
}

// VerifyValidatorBalance is a free data retrieval call binding the contract method 0x31f60d4c.
//
// Solidity: function verifyValidatorBalance(bytes32 balanceContainerRoot, uint40 validatorIndex, (bytes32,bytes32,bytes) proof) view returns()
func (_ContractBeaconChainProofsWrapper *ContractBeaconChainProofsWrapperCaller) VerifyValidatorBalance(opts *bind.CallOpts, balanceContainerRoot [32]byte, validatorIndex *big.Int, proof BeaconChainProofsBalanceProof) error {
	var out []interface{}
	err := _ContractBeaconChainProofsWrapper.contract.Call(opts, &out, "verifyValidatorBalance", balanceContainerRoot, validatorIndex, proof)

	if err != nil {
		return err
	}

	return err

}

// VerifyValidatorBalance is a free data retrieval call binding the contract method 0x31f60d4c.
//
// Solidity: function verifyValidatorBalance(bytes32 balanceContainerRoot, uint40 validatorIndex, (bytes32,bytes32,bytes) proof) view returns()
func (_ContractBeaconChainProofsWrapper *ContractBeaconChainProofsWrapperSession) VerifyValidatorBalance(balanceContainerRoot [32]byte, validatorIndex *big.Int, proof BeaconChainProofsBalanceProof) error {
	return _ContractBeaconChainProofsWrapper.Contract.VerifyValidatorBalance(&_ContractBeaconChainProofsWrapper.CallOpts, balanceContainerRoot, validatorIndex, proof)
}

// VerifyValidatorBalance is a free data retrieval call binding the contract method 0x31f60d4c.
//
// Solidity: function verifyValidatorBalance(bytes32 balanceContainerRoot, uint40 validatorIndex, (bytes32,bytes32,bytes) proof) view returns()
func (_ContractBeaconChainProofsWrapper *ContractBeaconChainProofsWrapperCallerSession) VerifyValidatorBalance(balanceContainerRoot [32]byte, validatorIndex *big.Int, proof BeaconChainProofsBalanceProof) error {
	return _ContractBeaconChainProofsWrapper.Contract.VerifyValidatorBalance(&_ContractBeaconChainProofsWrapper.CallOpts, balanceContainerRoot, validatorIndex, proof)
}

// VerifyValidatorFields is a free data retrieval call binding the contract method 0x256f222b.
//
// Solidity: function verifyValidatorFields(bytes32 beaconStateRoot, bytes32[] validatorFields, bytes validatorFieldsProof, uint40 validatorIndex) view returns()
func (_ContractBeaconChainProofsWrapper *ContractBeaconChainProofsWrapperCaller) VerifyValidatorFields(opts *bind.CallOpts, beaconStateRoot [32]byte, validatorFields [][32]byte, validatorFieldsProof []byte, validatorIndex *big.Int) error {
	var out []interface{}
	err := _ContractBeaconChainProofsWrapper.contract.Call(opts, &out, "verifyValidatorFields", beaconStateRoot, validatorFields, validatorFieldsProof, validatorIndex)

	if err != nil {
		return err
	}

	return err

}

// VerifyValidatorFields is a free data retrieval call binding the contract method 0x256f222b.
//
// Solidity: function verifyValidatorFields(bytes32 beaconStateRoot, bytes32[] validatorFields, bytes validatorFieldsProof, uint40 validatorIndex) view returns()
func (_ContractBeaconChainProofsWrapper *ContractBeaconChainProofsWrapperSession) VerifyValidatorFields(beaconStateRoot [32]byte, validatorFields [][32]byte, validatorFieldsProof []byte, validatorIndex *big.Int) error {
	return _ContractBeaconChainProofsWrapper.Contract.VerifyValidatorFields(&_ContractBeaconChainProofsWrapper.CallOpts, beaconStateRoot, validatorFields, validatorFieldsProof, validatorIndex)
}

// VerifyValidatorFields is a free data retrieval call binding the contract method 0x256f222b.
//
// Solidity: function verifyValidatorFields(bytes32 beaconStateRoot, bytes32[] validatorFields, bytes validatorFieldsProof, uint40 validatorIndex) view returns()
func (_ContractBeaconChainProofsWrapper *ContractBeaconChainProofsWrapperCallerSession) VerifyValidatorFields(beaconStateRoot [32]byte, validatorFields [][32]byte, validatorFieldsProof []byte, validatorIndex *big.Int) error {
	return _ContractBeaconChainProofsWrapper.Contract.VerifyValidatorFields(&_ContractBeaconChainProofsWrapper.CallOpts, beaconStateRoot, validatorFields, validatorFieldsProof, validatorIndex)
}

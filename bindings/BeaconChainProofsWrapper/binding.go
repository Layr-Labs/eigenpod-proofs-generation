// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package BeaconChainProofsWrapper

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

// BeaconChainProofsWrapperMetaData contains all meta data concerning the BeaconChainProofsWrapper contract.
var BeaconChainProofsWrapperMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"BALANCE_CONTAINER_INDEX\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"BALANCE_TREE_HEIGHT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"BEACON_BLOCK_HEADER_TREE_HEIGHT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"DENEB_BEACON_STATE_TREE_HEIGHT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"FAR_FUTURE_EPOCH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"PECTRA_BEACON_STATE_TREE_HEIGHT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"PECTRA_FORK_TIMESTAMP\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"SECONDS_PER_EPOCH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"SECONDS_PER_SLOT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"SLOTS_PER_EPOCH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"STATE_ROOT_INDEX\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"VALIDATOR_CONTAINER_INDEX\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"VALIDATOR_FIELDS_LENGTH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"VALIDATOR_TREE_HEIGHT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getActivationEpoch\",\"inputs\":[{\"name\":\"validatorFields\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getBalanceAtIndex\",\"inputs\":[{\"name\":\"balanceRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorIndex\",\"type\":\"uint40\",\"internalType\":\"uint40\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getBeaconStateTreeHeight\",\"inputs\":[{\"name\":\"proofTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getEffectiveBalanceGwei\",\"inputs\":[{\"name\":\"validatorFields\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getExitEpoch\",\"inputs\":[{\"name\":\"validatorFields\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getPubkeyHash\",\"inputs\":[{\"name\":\"validatorFields\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getWithdrawalCredentials\",\"inputs\":[{\"name\":\"validatorFields\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"isValidatorSlashed\",\"inputs\":[{\"name\":\"validatorFields\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"verifyBalanceContainer\",\"inputs\":[{\"name\":\"proofTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"beaconBlockRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structBeaconChainProofs.BalanceContainerProof\",\"components\":[{\"name\":\"balanceContainerRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyStateRoot\",\"inputs\":[{\"name\":\"beaconBlockRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structBeaconChainProofs.StateRootProof\",\"components\":[{\"name\":\"beaconStateRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyValidatorBalance\",\"inputs\":[{\"name\":\"balanceContainerRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorIndex\",\"type\":\"uint40\",\"internalType\":\"uint40\"},{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structBeaconChainProofs.BalanceProof\",\"components\":[{\"name\":\"pubkeyHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"balanceRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyValidatorFields\",\"inputs\":[{\"name\":\"proofTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"beaconStateRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorFields\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"validatorFieldsProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validatorIndex\",\"type\":\"uint40\",\"internalType\":\"uint40\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"error\",\"name\":\"InvalidProof\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidProofLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidProofLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidValidatorFieldsLength\",\"inputs\":[]}]",
	Bin: "0x6080604052348015600e575f5ffd5b5061114b8061001c5f395ff3fe608060405234801561000f575f5ffd5b5060043610610187575f3560e01c806366efbf4e116100d9578063a9ccd48711610093578063bf7836461161006e578063bf783646146102e7578063da30e279146102ee578063ec158777146102f5578063f17e9b4014610302575f5ffd5b8063a9ccd487146102c1578063aaa645a6146102d4578063b100b899146101e0575f5ffd5b806366efbf4e1461027e5780638ad5b4ff146102855780639030a9bb1461028c57806399ca22101461029f578063a38f2e7e146102b2578063a4cc5882146102b9575f5ffd5b8063304b9071116101445780634027da191161011f5780634027da191461023b578063423fe16f146102455780634534711b1461025857806360249fda1461026b575f5ffd5b8063304b9071146101fa57806331f60d4c146102155780633d6c9e1814610228575f5ffd5b8063043c35d31461018b5780630b9448ce146101a157806310c6e4c3146101c457806319312e29146101cb5780631d5c7b1c146101e05780632e808427146101e7575b5f5ffd5b60265b6040519081526020015b60405180910390f35b6101b46101af366004610c39565b610309565b6040519015158152602001610198565b602861018e565b6101de6101d9366004610d2f565b610319565b005b600361018e565b61018e6101f5366004610c39565b610329565b600c5b6040516001600160401b039091168152602001610198565b6101fd610223366004610d95565b610333565b61018e610236366004610ded565b610347565b63672a41006101fd565b6101fd610253366004610c39565b610351565b6101fd610266366004610c39565b61035b565b6101fd610279366004610e06565b610365565b600561018e565b600b61018e565b6101de61029a366004610e30565b610377565b61018e6102ad366004610c39565b610385565b600861018e565b6101fd61038f565b6101fd6102cf366004610c39565b6103a1565b6101de6102e2366004610eb7565b6103ab565b600661018e565b60206101fd565b6001600160401b036101fd565b600c61018e565b5f610313826103c3565b92915050565b6103248383836103eb565b505050565b5f610313826104b3565b5f61033f8484846104d5565b949350505050565b5f610313826105b3565b5f610313826105d8565b5f610313826105fc565b5f6103708383610613565b9392505050565b610381828261063f565b5050565b5f610313826106e4565b5f61039c600c6020610f93565b905090565b5f610313826106f8565b6103ba8787878787878761070f565b50505050505050565b5f816003815181106103d7576103d7610fbc565b60200260200101515f5f1b14159050919050565b5f6103f5846105b3565b9050610402816003610fd0565b61040d906020610fe3565b61041a6020840184610ffa565b90501461043a576040516313717da960e21b815260040160405180910390fd5b6003811b600c1761048f6104516020850185610ffa565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284375f92019190915250889250508635905084610843565b6104ac576040516309bde33960e01b815260040160405180910390fd5b5050505050565b5f815f815181106104c6576104c6610fbc565b60200260200101519050919050565b5f6104e260266001610fd0565b6104ed906020610fe3565b6104fa6040840184610ffa565b90501461051a576040516313717da960e21b815260040160405180910390fd5b5f610526600485611050565b64ffffffffff16905061057f61053f6040850185610ffa565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284375f920191909152508992505050602086013584610843565b61059c576040516309bde33960e01b815260040160405180910390fd5b6105aa836020013585610613565b95945050505050565b5f63672a41006001600160401b03831611156105d0576006610313565b600592915050565b5f610313826006815181106105ef576105ef610fbc565b602002602001015161085a565b5f610313826002815181106105ef576105ef610fbc565b5f80610620600484611079565b61062b9060406110a2565b64ffffffffff16905061033f84821b61085a565b61064b60036020610fe3565b6106586020830183610ffa565b905014610678576040516313717da960e21b815260040160405180910390fd5b6106c76106886020830183610ffa565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284375f9201919091525086925050843590506003610843565b610381576040516309bde33960e01b815260040160405180910390fd5b5f816001815181106104c6576104c6610fbc565b5f610313826005815181106105ef576105ef610fbc565b600884146107305760405163200591bd60e01b815260040160405180910390fd5b5f61073a886105b3565b90508061074960286001610fd0565b6107539190610fd0565b61075e906020610fe3565b831461077d576040516313717da960e21b815260040160405180910390fd5b5f6107b98787808060200260200160405190810160405280939291908181526020018383602002808284375f920191909152506108c192505050565b90505f64ffffffffff84166107d060286001610fd0565b600b901b17905061081a86868080601f0160208091040260200160405190810160405280939291908181526020018383808284375f920191909152508d9250869150859050610843565b610837576040516309bde33960e01b815260040160405180910390fd5b50505050505050505050565b5f83610850868585610b51565b1495945050505050565b60f881901c60e882901c61ff00161760d882901c62ff0000161760c882901c63ff000000161764ff0000000060b883901c161765ff000000000060a883901c161766ff000000000000609883901c161767ff0000000000000060889290921c919091161790565b5f5f600283516108d191906110c2565b90505f816001600160401b038111156108ec576108ec610c25565b604051908082528060200260200182016040528015610915578160200160208202803683370190505b5090505f5b82811015610a0f5760028561092f8383610fe3565b8151811061093f5761093f610fbc565b6020026020010151868360026109559190610fe3565b610960906001610fd0565b8151811061097057610970610fbc565b6020026020010151604051602001610992929190918252602082015260400190565b60408051601f19818403018152908290526109ac916110d5565b602060405180830381855afa1580156109c7573d5f5f3e3d5ffd5b5050506040513d601f19601f820116820180604052508101906109ea91906110eb565b8282815181106109fc576109fc610fbc565b602090810291909101015260010161091a565b50610a1b6002836110c2565b91505b8115610b2e575f5b82811015610b1b57600282610a3b8383610fe3565b81518110610a4b57610a4b610fbc565b602002602001015183836002610a619190610fe3565b610a6c906001610fd0565b81518110610a7c57610a7c610fbc565b6020026020010151604051602001610a9e929190918252602082015260400190565b60408051601f1981840301815290829052610ab8916110d5565b602060405180830381855afa158015610ad3573d5f5f3e3d5ffd5b5050506040513d601f19601f82011682018060405250810190610af691906110eb565b828281518110610b0857610b08610fbc565b6020908102919091010152600101610a26565b50610b276002836110c2565b9150610a1e565b805f81518110610b4057610b40610fbc565b602002602001015192505050919050565b5f83515f14158015610b6e575060208451610b6c9190611102565b155b610b8b576040516313717da960e21b815260040160405180910390fd5b604080516020808201909252848152905b85518111610c1b57610baf600285611102565b5f03610be15781515f528086015160205260208260405f60026107d05a03fa610bd6575f5ffd5b600284049350610c09565b808601515f52815160205260208260405f60026107d05a03fa610c02575f5ffd5b6002840493505b610c14602082610fd0565b9050610b9c565b5051949350505050565b634e487b7160e01b5f52604160045260245ffd5b5f60208284031215610c49575f5ffd5b81356001600160401b03811115610c5e575f5ffd5b8201601f81018413610c6e575f5ffd5b80356001600160401b03811115610c8757610c87610c25565b8060051b604051601f19603f83011681018181106001600160401b0382111715610cb357610cb3610c25565b604052918252602081840181019290810187841115610cd0575f5ffd5b6020850194505b83851015610cf357843580825260209586019590935001610cd7565b509695505050505050565b80356001600160401b0381168114610d14575f5ffd5b919050565b5f60408284031215610d29575f5ffd5b50919050565b5f5f5f60608486031215610d41575f5ffd5b610d4a84610cfe565b92506020840135915060408401356001600160401b03811115610d6b575f5ffd5b610d7786828701610d19565b9150509250925092565b803564ffffffffff81168114610d14575f5ffd5b5f5f5f60608486031215610da7575f5ffd5b83359250610db760208501610d81565b915060408401356001600160401b03811115610dd1575f5ffd5b840160608187031215610de2575f5ffd5b809150509250925092565b5f60208284031215610dfd575f5ffd5b61037082610cfe565b5f5f60408385031215610e17575f5ffd5b82359150610e2760208401610d81565b90509250929050565b5f5f60408385031215610e41575f5ffd5b8235915060208301356001600160401b03811115610e5d575f5ffd5b610e6985828601610d19565b9150509250929050565b5f5f83601f840112610e83575f5ffd5b5081356001600160401b03811115610e99575f5ffd5b602083019150836020828501011115610eb0575f5ffd5b9250929050565b5f5f5f5f5f5f5f60a0888a031215610ecd575f5ffd5b610ed688610cfe565b96506020880135955060408801356001600160401b03811115610ef7575f5ffd5b8801601f81018a13610f07575f5ffd5b80356001600160401b03811115610f1c575f5ffd5b8a60208260051b8401011115610f30575f5ffd5b6020919091019550935060608801356001600160401b03811115610f52575f5ffd5b610f5e8a828b01610e73565b9094509250610f71905060808901610d81565b905092959891949750929550565b634e487b7160e01b5f52601160045260245ffd5b6001600160401b038181168382160290811690818114610fb557610fb5610f7f565b5092915050565b634e487b7160e01b5f52603260045260245ffd5b8082018082111561031357610313610f7f565b808202811582820484141761031357610313610f7f565b5f5f8335601e1984360301811261100f575f5ffd5b8301803591506001600160401b03821115611028575f5ffd5b602001915036819003821315610eb0575f5ffd5b634e487b7160e01b5f52601260045260245ffd5b5f64ffffffffff8316806110665761106661103c565b8064ffffffffff84160491505092915050565b5f64ffffffffff83168061108f5761108f61103c565b8064ffffffffff84160691505092915050565b64ffffffffff8181168382160290811690818114610fb557610fb5610f7f565b5f826110d0576110d061103c565b500490565b5f82518060208501845e5f920191825250919050565b5f602082840312156110fb575f5ffd5b5051919050565b5f826111105761111061103c565b50069056fea264697066735822122066cd9d3d9ad06b6e337f65402707ca8e4f78ad100191893801f4ab77cca969a064736f6c634300081b0033",
}

// BeaconChainProofsWrapperABI is the input ABI used to generate the binding from.
// Deprecated: Use BeaconChainProofsWrapperMetaData.ABI instead.
var BeaconChainProofsWrapperABI = BeaconChainProofsWrapperMetaData.ABI

// BeaconChainProofsWrapperBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BeaconChainProofsWrapperMetaData.Bin instead.
var BeaconChainProofsWrapperBin = BeaconChainProofsWrapperMetaData.Bin

// DeployBeaconChainProofsWrapper deploys a new Ethereum contract, binding an instance of BeaconChainProofsWrapper to it.
func DeployBeaconChainProofsWrapper(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BeaconChainProofsWrapper, error) {
	parsed, err := BeaconChainProofsWrapperMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BeaconChainProofsWrapperBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BeaconChainProofsWrapper{BeaconChainProofsWrapperCaller: BeaconChainProofsWrapperCaller{contract: contract}, BeaconChainProofsWrapperTransactor: BeaconChainProofsWrapperTransactor{contract: contract}, BeaconChainProofsWrapperFilterer: BeaconChainProofsWrapperFilterer{contract: contract}}, nil
}

// BeaconChainProofsWrapper is an auto generated Go binding around an Ethereum contract.
type BeaconChainProofsWrapper struct {
	BeaconChainProofsWrapperCaller     // Read-only binding to the contract
	BeaconChainProofsWrapperTransactor // Write-only binding to the contract
	BeaconChainProofsWrapperFilterer   // Log filterer for contract events
}

// BeaconChainProofsWrapperCaller is an auto generated read-only Go binding around an Ethereum contract.
type BeaconChainProofsWrapperCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BeaconChainProofsWrapperTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BeaconChainProofsWrapperTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BeaconChainProofsWrapperFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BeaconChainProofsWrapperFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BeaconChainProofsWrapperSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BeaconChainProofsWrapperSession struct {
	Contract     *BeaconChainProofsWrapper // Generic contract binding to set the session for
	CallOpts     bind.CallOpts             // Call options to use throughout this session
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// BeaconChainProofsWrapperCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BeaconChainProofsWrapperCallerSession struct {
	Contract *BeaconChainProofsWrapperCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                   // Call options to use throughout this session
}

// BeaconChainProofsWrapperTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BeaconChainProofsWrapperTransactorSession struct {
	Contract     *BeaconChainProofsWrapperTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// BeaconChainProofsWrapperRaw is an auto generated low-level Go binding around an Ethereum contract.
type BeaconChainProofsWrapperRaw struct {
	Contract *BeaconChainProofsWrapper // Generic contract binding to access the raw methods on
}

// BeaconChainProofsWrapperCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BeaconChainProofsWrapperCallerRaw struct {
	Contract *BeaconChainProofsWrapperCaller // Generic read-only contract binding to access the raw methods on
}

// BeaconChainProofsWrapperTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BeaconChainProofsWrapperTransactorRaw struct {
	Contract *BeaconChainProofsWrapperTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBeaconChainProofsWrapper creates a new instance of BeaconChainProofsWrapper, bound to a specific deployed contract.
func NewBeaconChainProofsWrapper(address common.Address, backend bind.ContractBackend) (*BeaconChainProofsWrapper, error) {
	contract, err := bindBeaconChainProofsWrapper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BeaconChainProofsWrapper{BeaconChainProofsWrapperCaller: BeaconChainProofsWrapperCaller{contract: contract}, BeaconChainProofsWrapperTransactor: BeaconChainProofsWrapperTransactor{contract: contract}, BeaconChainProofsWrapperFilterer: BeaconChainProofsWrapperFilterer{contract: contract}}, nil
}

// NewBeaconChainProofsWrapperCaller creates a new read-only instance of BeaconChainProofsWrapper, bound to a specific deployed contract.
func NewBeaconChainProofsWrapperCaller(address common.Address, caller bind.ContractCaller) (*BeaconChainProofsWrapperCaller, error) {
	contract, err := bindBeaconChainProofsWrapper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BeaconChainProofsWrapperCaller{contract: contract}, nil
}

// NewBeaconChainProofsWrapperTransactor creates a new write-only instance of BeaconChainProofsWrapper, bound to a specific deployed contract.
func NewBeaconChainProofsWrapperTransactor(address common.Address, transactor bind.ContractTransactor) (*BeaconChainProofsWrapperTransactor, error) {
	contract, err := bindBeaconChainProofsWrapper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BeaconChainProofsWrapperTransactor{contract: contract}, nil
}

// NewBeaconChainProofsWrapperFilterer creates a new log filterer instance of BeaconChainProofsWrapper, bound to a specific deployed contract.
func NewBeaconChainProofsWrapperFilterer(address common.Address, filterer bind.ContractFilterer) (*BeaconChainProofsWrapperFilterer, error) {
	contract, err := bindBeaconChainProofsWrapper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BeaconChainProofsWrapperFilterer{contract: contract}, nil
}

// bindBeaconChainProofsWrapper binds a generic wrapper to an already deployed contract.
func bindBeaconChainProofsWrapper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BeaconChainProofsWrapperMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BeaconChainProofsWrapper.Contract.BeaconChainProofsWrapperCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconChainProofsWrapper.Contract.BeaconChainProofsWrapperTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BeaconChainProofsWrapper.Contract.BeaconChainProofsWrapperTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BeaconChainProofsWrapper.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconChainProofsWrapper.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BeaconChainProofsWrapper.Contract.contract.Transact(opts, method, params...)
}

// BALANCECONTAINERINDEX is a free data retrieval call binding the contract method 0xf17e9b40.
//
// Solidity: function BALANCE_CONTAINER_INDEX() pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCaller) BALANCECONTAINERINDEX(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BeaconChainProofsWrapper.contract.Call(opts, &out, "BALANCE_CONTAINER_INDEX")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BALANCECONTAINERINDEX is a free data retrieval call binding the contract method 0xf17e9b40.
//
// Solidity: function BALANCE_CONTAINER_INDEX() pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperSession) BALANCECONTAINERINDEX() (*big.Int, error) {
	return _BeaconChainProofsWrapper.Contract.BALANCECONTAINERINDEX(&_BeaconChainProofsWrapper.CallOpts)
}

// BALANCECONTAINERINDEX is a free data retrieval call binding the contract method 0xf17e9b40.
//
// Solidity: function BALANCE_CONTAINER_INDEX() pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCallerSession) BALANCECONTAINERINDEX() (*big.Int, error) {
	return _BeaconChainProofsWrapper.Contract.BALANCECONTAINERINDEX(&_BeaconChainProofsWrapper.CallOpts)
}

// BALANCETREEHEIGHT is a free data retrieval call binding the contract method 0x043c35d3.
//
// Solidity: function BALANCE_TREE_HEIGHT() pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCaller) BALANCETREEHEIGHT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BeaconChainProofsWrapper.contract.Call(opts, &out, "BALANCE_TREE_HEIGHT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BALANCETREEHEIGHT is a free data retrieval call binding the contract method 0x043c35d3.
//
// Solidity: function BALANCE_TREE_HEIGHT() pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperSession) BALANCETREEHEIGHT() (*big.Int, error) {
	return _BeaconChainProofsWrapper.Contract.BALANCETREEHEIGHT(&_BeaconChainProofsWrapper.CallOpts)
}

// BALANCETREEHEIGHT is a free data retrieval call binding the contract method 0x043c35d3.
//
// Solidity: function BALANCE_TREE_HEIGHT() pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCallerSession) BALANCETREEHEIGHT() (*big.Int, error) {
	return _BeaconChainProofsWrapper.Contract.BALANCETREEHEIGHT(&_BeaconChainProofsWrapper.CallOpts)
}

// BEACONBLOCKHEADERTREEHEIGHT is a free data retrieval call binding the contract method 0x1d5c7b1c.
//
// Solidity: function BEACON_BLOCK_HEADER_TREE_HEIGHT() pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCaller) BEACONBLOCKHEADERTREEHEIGHT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BeaconChainProofsWrapper.contract.Call(opts, &out, "BEACON_BLOCK_HEADER_TREE_HEIGHT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BEACONBLOCKHEADERTREEHEIGHT is a free data retrieval call binding the contract method 0x1d5c7b1c.
//
// Solidity: function BEACON_BLOCK_HEADER_TREE_HEIGHT() pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperSession) BEACONBLOCKHEADERTREEHEIGHT() (*big.Int, error) {
	return _BeaconChainProofsWrapper.Contract.BEACONBLOCKHEADERTREEHEIGHT(&_BeaconChainProofsWrapper.CallOpts)
}

// BEACONBLOCKHEADERTREEHEIGHT is a free data retrieval call binding the contract method 0x1d5c7b1c.
//
// Solidity: function BEACON_BLOCK_HEADER_TREE_HEIGHT() pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCallerSession) BEACONBLOCKHEADERTREEHEIGHT() (*big.Int, error) {
	return _BeaconChainProofsWrapper.Contract.BEACONBLOCKHEADERTREEHEIGHT(&_BeaconChainProofsWrapper.CallOpts)
}

// DENEBBEACONSTATETREEHEIGHT is a free data retrieval call binding the contract method 0x66efbf4e.
//
// Solidity: function DENEB_BEACON_STATE_TREE_HEIGHT() pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCaller) DENEBBEACONSTATETREEHEIGHT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BeaconChainProofsWrapper.contract.Call(opts, &out, "DENEB_BEACON_STATE_TREE_HEIGHT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DENEBBEACONSTATETREEHEIGHT is a free data retrieval call binding the contract method 0x66efbf4e.
//
// Solidity: function DENEB_BEACON_STATE_TREE_HEIGHT() pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperSession) DENEBBEACONSTATETREEHEIGHT() (*big.Int, error) {
	return _BeaconChainProofsWrapper.Contract.DENEBBEACONSTATETREEHEIGHT(&_BeaconChainProofsWrapper.CallOpts)
}

// DENEBBEACONSTATETREEHEIGHT is a free data retrieval call binding the contract method 0x66efbf4e.
//
// Solidity: function DENEB_BEACON_STATE_TREE_HEIGHT() pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCallerSession) DENEBBEACONSTATETREEHEIGHT() (*big.Int, error) {
	return _BeaconChainProofsWrapper.Contract.DENEBBEACONSTATETREEHEIGHT(&_BeaconChainProofsWrapper.CallOpts)
}

// FARFUTUREEPOCH is a free data retrieval call binding the contract method 0xec158777.
//
// Solidity: function FAR_FUTURE_EPOCH() pure returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCaller) FARFUTUREEPOCH(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _BeaconChainProofsWrapper.contract.Call(opts, &out, "FAR_FUTURE_EPOCH")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// FARFUTUREEPOCH is a free data retrieval call binding the contract method 0xec158777.
//
// Solidity: function FAR_FUTURE_EPOCH() pure returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperSession) FARFUTUREEPOCH() (uint64, error) {
	return _BeaconChainProofsWrapper.Contract.FARFUTUREEPOCH(&_BeaconChainProofsWrapper.CallOpts)
}

// FARFUTUREEPOCH is a free data retrieval call binding the contract method 0xec158777.
//
// Solidity: function FAR_FUTURE_EPOCH() pure returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCallerSession) FARFUTUREEPOCH() (uint64, error) {
	return _BeaconChainProofsWrapper.Contract.FARFUTUREEPOCH(&_BeaconChainProofsWrapper.CallOpts)
}

// PECTRABEACONSTATETREEHEIGHT is a free data retrieval call binding the contract method 0xbf783646.
//
// Solidity: function PECTRA_BEACON_STATE_TREE_HEIGHT() pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCaller) PECTRABEACONSTATETREEHEIGHT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BeaconChainProofsWrapper.contract.Call(opts, &out, "PECTRA_BEACON_STATE_TREE_HEIGHT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PECTRABEACONSTATETREEHEIGHT is a free data retrieval call binding the contract method 0xbf783646.
//
// Solidity: function PECTRA_BEACON_STATE_TREE_HEIGHT() pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperSession) PECTRABEACONSTATETREEHEIGHT() (*big.Int, error) {
	return _BeaconChainProofsWrapper.Contract.PECTRABEACONSTATETREEHEIGHT(&_BeaconChainProofsWrapper.CallOpts)
}

// PECTRABEACONSTATETREEHEIGHT is a free data retrieval call binding the contract method 0xbf783646.
//
// Solidity: function PECTRA_BEACON_STATE_TREE_HEIGHT() pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCallerSession) PECTRABEACONSTATETREEHEIGHT() (*big.Int, error) {
	return _BeaconChainProofsWrapper.Contract.PECTRABEACONSTATETREEHEIGHT(&_BeaconChainProofsWrapper.CallOpts)
}

// PECTRAFORKTIMESTAMP is a free data retrieval call binding the contract method 0x4027da19.
//
// Solidity: function PECTRA_FORK_TIMESTAMP() pure returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCaller) PECTRAFORKTIMESTAMP(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _BeaconChainProofsWrapper.contract.Call(opts, &out, "PECTRA_FORK_TIMESTAMP")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// PECTRAFORKTIMESTAMP is a free data retrieval call binding the contract method 0x4027da19.
//
// Solidity: function PECTRA_FORK_TIMESTAMP() pure returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperSession) PECTRAFORKTIMESTAMP() (uint64, error) {
	return _BeaconChainProofsWrapper.Contract.PECTRAFORKTIMESTAMP(&_BeaconChainProofsWrapper.CallOpts)
}

// PECTRAFORKTIMESTAMP is a free data retrieval call binding the contract method 0x4027da19.
//
// Solidity: function PECTRA_FORK_TIMESTAMP() pure returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCallerSession) PECTRAFORKTIMESTAMP() (uint64, error) {
	return _BeaconChainProofsWrapper.Contract.PECTRAFORKTIMESTAMP(&_BeaconChainProofsWrapper.CallOpts)
}

// SECONDSPEREPOCH is a free data retrieval call binding the contract method 0xa4cc5882.
//
// Solidity: function SECONDS_PER_EPOCH() pure returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCaller) SECONDSPEREPOCH(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _BeaconChainProofsWrapper.contract.Call(opts, &out, "SECONDS_PER_EPOCH")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// SECONDSPEREPOCH is a free data retrieval call binding the contract method 0xa4cc5882.
//
// Solidity: function SECONDS_PER_EPOCH() pure returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperSession) SECONDSPEREPOCH() (uint64, error) {
	return _BeaconChainProofsWrapper.Contract.SECONDSPEREPOCH(&_BeaconChainProofsWrapper.CallOpts)
}

// SECONDSPEREPOCH is a free data retrieval call binding the contract method 0xa4cc5882.
//
// Solidity: function SECONDS_PER_EPOCH() pure returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCallerSession) SECONDSPEREPOCH() (uint64, error) {
	return _BeaconChainProofsWrapper.Contract.SECONDSPEREPOCH(&_BeaconChainProofsWrapper.CallOpts)
}

// SECONDSPERSLOT is a free data retrieval call binding the contract method 0x304b9071.
//
// Solidity: function SECONDS_PER_SLOT() pure returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCaller) SECONDSPERSLOT(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _BeaconChainProofsWrapper.contract.Call(opts, &out, "SECONDS_PER_SLOT")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// SECONDSPERSLOT is a free data retrieval call binding the contract method 0x304b9071.
//
// Solidity: function SECONDS_PER_SLOT() pure returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperSession) SECONDSPERSLOT() (uint64, error) {
	return _BeaconChainProofsWrapper.Contract.SECONDSPERSLOT(&_BeaconChainProofsWrapper.CallOpts)
}

// SECONDSPERSLOT is a free data retrieval call binding the contract method 0x304b9071.
//
// Solidity: function SECONDS_PER_SLOT() pure returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCallerSession) SECONDSPERSLOT() (uint64, error) {
	return _BeaconChainProofsWrapper.Contract.SECONDSPERSLOT(&_BeaconChainProofsWrapper.CallOpts)
}

// SLOTSPEREPOCH is a free data retrieval call binding the contract method 0xda30e279.
//
// Solidity: function SLOTS_PER_EPOCH() pure returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCaller) SLOTSPEREPOCH(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _BeaconChainProofsWrapper.contract.Call(opts, &out, "SLOTS_PER_EPOCH")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// SLOTSPEREPOCH is a free data retrieval call binding the contract method 0xda30e279.
//
// Solidity: function SLOTS_PER_EPOCH() pure returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperSession) SLOTSPEREPOCH() (uint64, error) {
	return _BeaconChainProofsWrapper.Contract.SLOTSPEREPOCH(&_BeaconChainProofsWrapper.CallOpts)
}

// SLOTSPEREPOCH is a free data retrieval call binding the contract method 0xda30e279.
//
// Solidity: function SLOTS_PER_EPOCH() pure returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCallerSession) SLOTSPEREPOCH() (uint64, error) {
	return _BeaconChainProofsWrapper.Contract.SLOTSPEREPOCH(&_BeaconChainProofsWrapper.CallOpts)
}

// STATEROOTINDEX is a free data retrieval call binding the contract method 0xb100b899.
//
// Solidity: function STATE_ROOT_INDEX() pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCaller) STATEROOTINDEX(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BeaconChainProofsWrapper.contract.Call(opts, &out, "STATE_ROOT_INDEX")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// STATEROOTINDEX is a free data retrieval call binding the contract method 0xb100b899.
//
// Solidity: function STATE_ROOT_INDEX() pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperSession) STATEROOTINDEX() (*big.Int, error) {
	return _BeaconChainProofsWrapper.Contract.STATEROOTINDEX(&_BeaconChainProofsWrapper.CallOpts)
}

// STATEROOTINDEX is a free data retrieval call binding the contract method 0xb100b899.
//
// Solidity: function STATE_ROOT_INDEX() pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCallerSession) STATEROOTINDEX() (*big.Int, error) {
	return _BeaconChainProofsWrapper.Contract.STATEROOTINDEX(&_BeaconChainProofsWrapper.CallOpts)
}

// VALIDATORCONTAINERINDEX is a free data retrieval call binding the contract method 0x8ad5b4ff.
//
// Solidity: function VALIDATOR_CONTAINER_INDEX() pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCaller) VALIDATORCONTAINERINDEX(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BeaconChainProofsWrapper.contract.Call(opts, &out, "VALIDATOR_CONTAINER_INDEX")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VALIDATORCONTAINERINDEX is a free data retrieval call binding the contract method 0x8ad5b4ff.
//
// Solidity: function VALIDATOR_CONTAINER_INDEX() pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperSession) VALIDATORCONTAINERINDEX() (*big.Int, error) {
	return _BeaconChainProofsWrapper.Contract.VALIDATORCONTAINERINDEX(&_BeaconChainProofsWrapper.CallOpts)
}

// VALIDATORCONTAINERINDEX is a free data retrieval call binding the contract method 0x8ad5b4ff.
//
// Solidity: function VALIDATOR_CONTAINER_INDEX() pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCallerSession) VALIDATORCONTAINERINDEX() (*big.Int, error) {
	return _BeaconChainProofsWrapper.Contract.VALIDATORCONTAINERINDEX(&_BeaconChainProofsWrapper.CallOpts)
}

// VALIDATORFIELDSLENGTH is a free data retrieval call binding the contract method 0xa38f2e7e.
//
// Solidity: function VALIDATOR_FIELDS_LENGTH() pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCaller) VALIDATORFIELDSLENGTH(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BeaconChainProofsWrapper.contract.Call(opts, &out, "VALIDATOR_FIELDS_LENGTH")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VALIDATORFIELDSLENGTH is a free data retrieval call binding the contract method 0xa38f2e7e.
//
// Solidity: function VALIDATOR_FIELDS_LENGTH() pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperSession) VALIDATORFIELDSLENGTH() (*big.Int, error) {
	return _BeaconChainProofsWrapper.Contract.VALIDATORFIELDSLENGTH(&_BeaconChainProofsWrapper.CallOpts)
}

// VALIDATORFIELDSLENGTH is a free data retrieval call binding the contract method 0xa38f2e7e.
//
// Solidity: function VALIDATOR_FIELDS_LENGTH() pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCallerSession) VALIDATORFIELDSLENGTH() (*big.Int, error) {
	return _BeaconChainProofsWrapper.Contract.VALIDATORFIELDSLENGTH(&_BeaconChainProofsWrapper.CallOpts)
}

// VALIDATORTREEHEIGHT is a free data retrieval call binding the contract method 0x10c6e4c3.
//
// Solidity: function VALIDATOR_TREE_HEIGHT() pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCaller) VALIDATORTREEHEIGHT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BeaconChainProofsWrapper.contract.Call(opts, &out, "VALIDATOR_TREE_HEIGHT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VALIDATORTREEHEIGHT is a free data retrieval call binding the contract method 0x10c6e4c3.
//
// Solidity: function VALIDATOR_TREE_HEIGHT() pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperSession) VALIDATORTREEHEIGHT() (*big.Int, error) {
	return _BeaconChainProofsWrapper.Contract.VALIDATORTREEHEIGHT(&_BeaconChainProofsWrapper.CallOpts)
}

// VALIDATORTREEHEIGHT is a free data retrieval call binding the contract method 0x10c6e4c3.
//
// Solidity: function VALIDATOR_TREE_HEIGHT() pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCallerSession) VALIDATORTREEHEIGHT() (*big.Int, error) {
	return _BeaconChainProofsWrapper.Contract.VALIDATORTREEHEIGHT(&_BeaconChainProofsWrapper.CallOpts)
}

// GetActivationEpoch is a free data retrieval call binding the contract method 0xa9ccd487.
//
// Solidity: function getActivationEpoch(bytes32[] validatorFields) pure returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCaller) GetActivationEpoch(opts *bind.CallOpts, validatorFields [][32]byte) (uint64, error) {
	var out []interface{}
	err := _BeaconChainProofsWrapper.contract.Call(opts, &out, "getActivationEpoch", validatorFields)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetActivationEpoch is a free data retrieval call binding the contract method 0xa9ccd487.
//
// Solidity: function getActivationEpoch(bytes32[] validatorFields) pure returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperSession) GetActivationEpoch(validatorFields [][32]byte) (uint64, error) {
	return _BeaconChainProofsWrapper.Contract.GetActivationEpoch(&_BeaconChainProofsWrapper.CallOpts, validatorFields)
}

// GetActivationEpoch is a free data retrieval call binding the contract method 0xa9ccd487.
//
// Solidity: function getActivationEpoch(bytes32[] validatorFields) pure returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCallerSession) GetActivationEpoch(validatorFields [][32]byte) (uint64, error) {
	return _BeaconChainProofsWrapper.Contract.GetActivationEpoch(&_BeaconChainProofsWrapper.CallOpts, validatorFields)
}

// GetBalanceAtIndex is a free data retrieval call binding the contract method 0x60249fda.
//
// Solidity: function getBalanceAtIndex(bytes32 balanceRoot, uint40 validatorIndex) pure returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCaller) GetBalanceAtIndex(opts *bind.CallOpts, balanceRoot [32]byte, validatorIndex *big.Int) (uint64, error) {
	var out []interface{}
	err := _BeaconChainProofsWrapper.contract.Call(opts, &out, "getBalanceAtIndex", balanceRoot, validatorIndex)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetBalanceAtIndex is a free data retrieval call binding the contract method 0x60249fda.
//
// Solidity: function getBalanceAtIndex(bytes32 balanceRoot, uint40 validatorIndex) pure returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperSession) GetBalanceAtIndex(balanceRoot [32]byte, validatorIndex *big.Int) (uint64, error) {
	return _BeaconChainProofsWrapper.Contract.GetBalanceAtIndex(&_BeaconChainProofsWrapper.CallOpts, balanceRoot, validatorIndex)
}

// GetBalanceAtIndex is a free data retrieval call binding the contract method 0x60249fda.
//
// Solidity: function getBalanceAtIndex(bytes32 balanceRoot, uint40 validatorIndex) pure returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCallerSession) GetBalanceAtIndex(balanceRoot [32]byte, validatorIndex *big.Int) (uint64, error) {
	return _BeaconChainProofsWrapper.Contract.GetBalanceAtIndex(&_BeaconChainProofsWrapper.CallOpts, balanceRoot, validatorIndex)
}

// GetBeaconStateTreeHeight is a free data retrieval call binding the contract method 0x3d6c9e18.
//
// Solidity: function getBeaconStateTreeHeight(uint64 proofTimestamp) pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCaller) GetBeaconStateTreeHeight(opts *bind.CallOpts, proofTimestamp uint64) (*big.Int, error) {
	var out []interface{}
	err := _BeaconChainProofsWrapper.contract.Call(opts, &out, "getBeaconStateTreeHeight", proofTimestamp)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBeaconStateTreeHeight is a free data retrieval call binding the contract method 0x3d6c9e18.
//
// Solidity: function getBeaconStateTreeHeight(uint64 proofTimestamp) pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperSession) GetBeaconStateTreeHeight(proofTimestamp uint64) (*big.Int, error) {
	return _BeaconChainProofsWrapper.Contract.GetBeaconStateTreeHeight(&_BeaconChainProofsWrapper.CallOpts, proofTimestamp)
}

// GetBeaconStateTreeHeight is a free data retrieval call binding the contract method 0x3d6c9e18.
//
// Solidity: function getBeaconStateTreeHeight(uint64 proofTimestamp) pure returns(uint256)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCallerSession) GetBeaconStateTreeHeight(proofTimestamp uint64) (*big.Int, error) {
	return _BeaconChainProofsWrapper.Contract.GetBeaconStateTreeHeight(&_BeaconChainProofsWrapper.CallOpts, proofTimestamp)
}

// GetEffectiveBalanceGwei is a free data retrieval call binding the contract method 0x4534711b.
//
// Solidity: function getEffectiveBalanceGwei(bytes32[] validatorFields) pure returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCaller) GetEffectiveBalanceGwei(opts *bind.CallOpts, validatorFields [][32]byte) (uint64, error) {
	var out []interface{}
	err := _BeaconChainProofsWrapper.contract.Call(opts, &out, "getEffectiveBalanceGwei", validatorFields)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetEffectiveBalanceGwei is a free data retrieval call binding the contract method 0x4534711b.
//
// Solidity: function getEffectiveBalanceGwei(bytes32[] validatorFields) pure returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperSession) GetEffectiveBalanceGwei(validatorFields [][32]byte) (uint64, error) {
	return _BeaconChainProofsWrapper.Contract.GetEffectiveBalanceGwei(&_BeaconChainProofsWrapper.CallOpts, validatorFields)
}

// GetEffectiveBalanceGwei is a free data retrieval call binding the contract method 0x4534711b.
//
// Solidity: function getEffectiveBalanceGwei(bytes32[] validatorFields) pure returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCallerSession) GetEffectiveBalanceGwei(validatorFields [][32]byte) (uint64, error) {
	return _BeaconChainProofsWrapper.Contract.GetEffectiveBalanceGwei(&_BeaconChainProofsWrapper.CallOpts, validatorFields)
}

// GetExitEpoch is a free data retrieval call binding the contract method 0x423fe16f.
//
// Solidity: function getExitEpoch(bytes32[] validatorFields) pure returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCaller) GetExitEpoch(opts *bind.CallOpts, validatorFields [][32]byte) (uint64, error) {
	var out []interface{}
	err := _BeaconChainProofsWrapper.contract.Call(opts, &out, "getExitEpoch", validatorFields)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetExitEpoch is a free data retrieval call binding the contract method 0x423fe16f.
//
// Solidity: function getExitEpoch(bytes32[] validatorFields) pure returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperSession) GetExitEpoch(validatorFields [][32]byte) (uint64, error) {
	return _BeaconChainProofsWrapper.Contract.GetExitEpoch(&_BeaconChainProofsWrapper.CallOpts, validatorFields)
}

// GetExitEpoch is a free data retrieval call binding the contract method 0x423fe16f.
//
// Solidity: function getExitEpoch(bytes32[] validatorFields) pure returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCallerSession) GetExitEpoch(validatorFields [][32]byte) (uint64, error) {
	return _BeaconChainProofsWrapper.Contract.GetExitEpoch(&_BeaconChainProofsWrapper.CallOpts, validatorFields)
}

// GetPubkeyHash is a free data retrieval call binding the contract method 0x2e808427.
//
// Solidity: function getPubkeyHash(bytes32[] validatorFields) pure returns(bytes32)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCaller) GetPubkeyHash(opts *bind.CallOpts, validatorFields [][32]byte) ([32]byte, error) {
	var out []interface{}
	err := _BeaconChainProofsWrapper.contract.Call(opts, &out, "getPubkeyHash", validatorFields)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetPubkeyHash is a free data retrieval call binding the contract method 0x2e808427.
//
// Solidity: function getPubkeyHash(bytes32[] validatorFields) pure returns(bytes32)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperSession) GetPubkeyHash(validatorFields [][32]byte) ([32]byte, error) {
	return _BeaconChainProofsWrapper.Contract.GetPubkeyHash(&_BeaconChainProofsWrapper.CallOpts, validatorFields)
}

// GetPubkeyHash is a free data retrieval call binding the contract method 0x2e808427.
//
// Solidity: function getPubkeyHash(bytes32[] validatorFields) pure returns(bytes32)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCallerSession) GetPubkeyHash(validatorFields [][32]byte) ([32]byte, error) {
	return _BeaconChainProofsWrapper.Contract.GetPubkeyHash(&_BeaconChainProofsWrapper.CallOpts, validatorFields)
}

// GetWithdrawalCredentials is a free data retrieval call binding the contract method 0x99ca2210.
//
// Solidity: function getWithdrawalCredentials(bytes32[] validatorFields) pure returns(bytes32)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCaller) GetWithdrawalCredentials(opts *bind.CallOpts, validatorFields [][32]byte) ([32]byte, error) {
	var out []interface{}
	err := _BeaconChainProofsWrapper.contract.Call(opts, &out, "getWithdrawalCredentials", validatorFields)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetWithdrawalCredentials is a free data retrieval call binding the contract method 0x99ca2210.
//
// Solidity: function getWithdrawalCredentials(bytes32[] validatorFields) pure returns(bytes32)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperSession) GetWithdrawalCredentials(validatorFields [][32]byte) ([32]byte, error) {
	return _BeaconChainProofsWrapper.Contract.GetWithdrawalCredentials(&_BeaconChainProofsWrapper.CallOpts, validatorFields)
}

// GetWithdrawalCredentials is a free data retrieval call binding the contract method 0x99ca2210.
//
// Solidity: function getWithdrawalCredentials(bytes32[] validatorFields) pure returns(bytes32)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCallerSession) GetWithdrawalCredentials(validatorFields [][32]byte) ([32]byte, error) {
	return _BeaconChainProofsWrapper.Contract.GetWithdrawalCredentials(&_BeaconChainProofsWrapper.CallOpts, validatorFields)
}

// IsValidatorSlashed is a free data retrieval call binding the contract method 0x0b9448ce.
//
// Solidity: function isValidatorSlashed(bytes32[] validatorFields) pure returns(bool)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCaller) IsValidatorSlashed(opts *bind.CallOpts, validatorFields [][32]byte) (bool, error) {
	var out []interface{}
	err := _BeaconChainProofsWrapper.contract.Call(opts, &out, "isValidatorSlashed", validatorFields)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidatorSlashed is a free data retrieval call binding the contract method 0x0b9448ce.
//
// Solidity: function isValidatorSlashed(bytes32[] validatorFields) pure returns(bool)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperSession) IsValidatorSlashed(validatorFields [][32]byte) (bool, error) {
	return _BeaconChainProofsWrapper.Contract.IsValidatorSlashed(&_BeaconChainProofsWrapper.CallOpts, validatorFields)
}

// IsValidatorSlashed is a free data retrieval call binding the contract method 0x0b9448ce.
//
// Solidity: function isValidatorSlashed(bytes32[] validatorFields) pure returns(bool)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCallerSession) IsValidatorSlashed(validatorFields [][32]byte) (bool, error) {
	return _BeaconChainProofsWrapper.Contract.IsValidatorSlashed(&_BeaconChainProofsWrapper.CallOpts, validatorFields)
}

// VerifyBalanceContainer is a free data retrieval call binding the contract method 0x19312e29.
//
// Solidity: function verifyBalanceContainer(uint64 proofTimestamp, bytes32 beaconBlockRoot, (bytes32,bytes) proof) view returns()
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCaller) VerifyBalanceContainer(opts *bind.CallOpts, proofTimestamp uint64, beaconBlockRoot [32]byte, proof BeaconChainProofsBalanceContainerProof) error {
	var out []interface{}
	err := _BeaconChainProofsWrapper.contract.Call(opts, &out, "verifyBalanceContainer", proofTimestamp, beaconBlockRoot, proof)

	if err != nil {
		return err
	}

	return err

}

// VerifyBalanceContainer is a free data retrieval call binding the contract method 0x19312e29.
//
// Solidity: function verifyBalanceContainer(uint64 proofTimestamp, bytes32 beaconBlockRoot, (bytes32,bytes) proof) view returns()
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperSession) VerifyBalanceContainer(proofTimestamp uint64, beaconBlockRoot [32]byte, proof BeaconChainProofsBalanceContainerProof) error {
	return _BeaconChainProofsWrapper.Contract.VerifyBalanceContainer(&_BeaconChainProofsWrapper.CallOpts, proofTimestamp, beaconBlockRoot, proof)
}

// VerifyBalanceContainer is a free data retrieval call binding the contract method 0x19312e29.
//
// Solidity: function verifyBalanceContainer(uint64 proofTimestamp, bytes32 beaconBlockRoot, (bytes32,bytes) proof) view returns()
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCallerSession) VerifyBalanceContainer(proofTimestamp uint64, beaconBlockRoot [32]byte, proof BeaconChainProofsBalanceContainerProof) error {
	return _BeaconChainProofsWrapper.Contract.VerifyBalanceContainer(&_BeaconChainProofsWrapper.CallOpts, proofTimestamp, beaconBlockRoot, proof)
}

// VerifyStateRoot is a free data retrieval call binding the contract method 0x9030a9bb.
//
// Solidity: function verifyStateRoot(bytes32 beaconBlockRoot, (bytes32,bytes) proof) view returns()
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCaller) VerifyStateRoot(opts *bind.CallOpts, beaconBlockRoot [32]byte, proof BeaconChainProofsStateRootProof) error {
	var out []interface{}
	err := _BeaconChainProofsWrapper.contract.Call(opts, &out, "verifyStateRoot", beaconBlockRoot, proof)

	if err != nil {
		return err
	}

	return err

}

// VerifyStateRoot is a free data retrieval call binding the contract method 0x9030a9bb.
//
// Solidity: function verifyStateRoot(bytes32 beaconBlockRoot, (bytes32,bytes) proof) view returns()
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperSession) VerifyStateRoot(beaconBlockRoot [32]byte, proof BeaconChainProofsStateRootProof) error {
	return _BeaconChainProofsWrapper.Contract.VerifyStateRoot(&_BeaconChainProofsWrapper.CallOpts, beaconBlockRoot, proof)
}

// VerifyStateRoot is a free data retrieval call binding the contract method 0x9030a9bb.
//
// Solidity: function verifyStateRoot(bytes32 beaconBlockRoot, (bytes32,bytes) proof) view returns()
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCallerSession) VerifyStateRoot(beaconBlockRoot [32]byte, proof BeaconChainProofsStateRootProof) error {
	return _BeaconChainProofsWrapper.Contract.VerifyStateRoot(&_BeaconChainProofsWrapper.CallOpts, beaconBlockRoot, proof)
}

// VerifyValidatorBalance is a free data retrieval call binding the contract method 0x31f60d4c.
//
// Solidity: function verifyValidatorBalance(bytes32 balanceContainerRoot, uint40 validatorIndex, (bytes32,bytes32,bytes) proof) view returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCaller) VerifyValidatorBalance(opts *bind.CallOpts, balanceContainerRoot [32]byte, validatorIndex *big.Int, proof BeaconChainProofsBalanceProof) (uint64, error) {
	var out []interface{}
	err := _BeaconChainProofsWrapper.contract.Call(opts, &out, "verifyValidatorBalance", balanceContainerRoot, validatorIndex, proof)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// VerifyValidatorBalance is a free data retrieval call binding the contract method 0x31f60d4c.
//
// Solidity: function verifyValidatorBalance(bytes32 balanceContainerRoot, uint40 validatorIndex, (bytes32,bytes32,bytes) proof) view returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperSession) VerifyValidatorBalance(balanceContainerRoot [32]byte, validatorIndex *big.Int, proof BeaconChainProofsBalanceProof) (uint64, error) {
	return _BeaconChainProofsWrapper.Contract.VerifyValidatorBalance(&_BeaconChainProofsWrapper.CallOpts, balanceContainerRoot, validatorIndex, proof)
}

// VerifyValidatorBalance is a free data retrieval call binding the contract method 0x31f60d4c.
//
// Solidity: function verifyValidatorBalance(bytes32 balanceContainerRoot, uint40 validatorIndex, (bytes32,bytes32,bytes) proof) view returns(uint64)
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCallerSession) VerifyValidatorBalance(balanceContainerRoot [32]byte, validatorIndex *big.Int, proof BeaconChainProofsBalanceProof) (uint64, error) {
	return _BeaconChainProofsWrapper.Contract.VerifyValidatorBalance(&_BeaconChainProofsWrapper.CallOpts, balanceContainerRoot, validatorIndex, proof)
}

// VerifyValidatorFields is a free data retrieval call binding the contract method 0xaaa645a6.
//
// Solidity: function verifyValidatorFields(uint64 proofTimestamp, bytes32 beaconStateRoot, bytes32[] validatorFields, bytes validatorFieldsProof, uint40 validatorIndex) view returns()
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCaller) VerifyValidatorFields(opts *bind.CallOpts, proofTimestamp uint64, beaconStateRoot [32]byte, validatorFields [][32]byte, validatorFieldsProof []byte, validatorIndex *big.Int) error {
	var out []interface{}
	err := _BeaconChainProofsWrapper.contract.Call(opts, &out, "verifyValidatorFields", proofTimestamp, beaconStateRoot, validatorFields, validatorFieldsProof, validatorIndex)

	if err != nil {
		return err
	}

	return err

}

// VerifyValidatorFields is a free data retrieval call binding the contract method 0xaaa645a6.
//
// Solidity: function verifyValidatorFields(uint64 proofTimestamp, bytes32 beaconStateRoot, bytes32[] validatorFields, bytes validatorFieldsProof, uint40 validatorIndex) view returns()
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperSession) VerifyValidatorFields(proofTimestamp uint64, beaconStateRoot [32]byte, validatorFields [][32]byte, validatorFieldsProof []byte, validatorIndex *big.Int) error {
	return _BeaconChainProofsWrapper.Contract.VerifyValidatorFields(&_BeaconChainProofsWrapper.CallOpts, proofTimestamp, beaconStateRoot, validatorFields, validatorFieldsProof, validatorIndex)
}

// VerifyValidatorFields is a free data retrieval call binding the contract method 0xaaa645a6.
//
// Solidity: function verifyValidatorFields(uint64 proofTimestamp, bytes32 beaconStateRoot, bytes32[] validatorFields, bytes validatorFieldsProof, uint40 validatorIndex) view returns()
func (_BeaconChainProofsWrapper *BeaconChainProofsWrapperCallerSession) VerifyValidatorFields(proofTimestamp uint64, beaconStateRoot [32]byte, validatorFields [][32]byte, validatorFieldsProof []byte, validatorIndex *big.Int) error {
	return _BeaconChainProofsWrapper.Contract.VerifyValidatorFields(&_BeaconChainProofsWrapper.CallOpts, proofTimestamp, beaconStateRoot, validatorFields, validatorFieldsProof, validatorIndex)
}

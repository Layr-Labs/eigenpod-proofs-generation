package eigenpodproofs_test

import (
	"os"
	"testing"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	"github.com/Layr-Labs/eigenpod-proofs-generation/beacon"
	contractBeaconChainProofsWrapper "github.com/Layr-Labs/eigenpod-proofs-generation/bindings/BeaconChainProofsWrapper"
	"github.com/Layr-Labs/eigenpod-proofs-generation/common"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const RPC_URL = "https://ethereum-holesky-rpc.publicnode.com"

var BEACON_CHAIN_PROOFS_WRAPPER_ADDRESS = gethcommon.HexToAddress("0xf0B1Dd8D9476778564A515409c17c96705661e6A")

var beaconHeader *phase0.BeaconBlockHeader
var beaconState *spec.VersionedBeaconState
var beaconChainProofsWrapper *contractBeaconChainProofsWrapper.ContractBeaconChainProofsWrapper
var epp *eigenpodproofs.EigenPodProofs

// before all
func TestMain(m *testing.M) {
	var err error

	beaconHeaderFileName := "data/deneb_holesky_beacon_header_1650726.json"
	beaconHeaderBytes, err := common.ReadFile(beaconHeaderFileName)
	if err != nil {
		panic(err)
	}

	beaconStateFileName := "data/deneb_holesky_beacon_state_1650726.ssz"
	beaconStateBytes, err := common.ReadFile(beaconStateFileName)
	if err != nil {
		panic(err)

	}

	beaconHeader = &phase0.BeaconBlockHeader{}
	err = beaconHeader.UnmarshalJSON(beaconHeaderBytes)
	if err != nil {
		panic(err)
	}

	beaconState, err = beacon.UnmarshalSSZVersionedBeaconState(beaconStateBytes)
	if err != nil {
		panic(err)
	}

	ethClient, err := ethclient.Dial(RPC_URL)
	if err != nil {
		panic(err)
	}

	beaconChainProofsWrapper, err = contractBeaconChainProofsWrapper.NewContractBeaconChainProofsWrapper(BEACON_CHAIN_PROOFS_WRAPPER_ADDRESS, ethClient)
	if err != nil {
		panic(err)
	}

	epp, err = eigenpodproofs.NewEigenPodProofs(17000, 600)
	if err != nil {
		panic(err)
	}

	code := m.Run()

	// // Teardown
	// log.Println("Tearing down suite")
	// teardownSuite()

	// Exit with test result code
	os.Exit(code)
}

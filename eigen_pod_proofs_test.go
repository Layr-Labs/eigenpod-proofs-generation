package eigenpodproofs_test

import (
	"os"
	"testing"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	"github.com/Layr-Labs/eigenpod-proofs-generation/beacon"
	BeaconChainProofsWrapper "github.com/Layr-Labs/eigenpod-proofs-generation/bindings/BeaconChainProofsWrapper"
	"github.com/Layr-Labs/eigenpod-proofs-generation/common"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

var BEACON_CHAIN_PROOFS_WRAPPER_ADDRESS = gethcommon.HexToAddress("0x874Be4b0CaC8D3F6286Eee6E6196553aabA8Cb85")

var (
	beaconHeader             *phase0.BeaconBlockHeader
	beaconState              *spec.VersionedBeaconState
	beaconChainProofsWrapper *BeaconChainProofsWrapper.BeaconChainProofsWrapper
	epp                      *eigenpodproofs.EigenPodProofs
)

func loadBeaconState(headerPath, statePath string, chainID uint64) error {
	headerBytes, err := common.ReadFile(headerPath)
	if err != nil {
		return err
	}
	stateBytes, err := common.ReadFile(statePath)
	if err != nil {
		return err
	}

	beaconHeader = &phase0.BeaconBlockHeader{}
	if err := beaconHeader.UnmarshalJSON(headerBytes); err != nil {
		return err
	}

	beaconState, err = beacon.UnmarshalSSZVersionedBeaconState(stateBytes)
	if err != nil {
		return err
	}

	epp, err = eigenpodproofs.NewEigenPodProofs(chainID, 600)
	return err
}

func TestMain(m *testing.M) {
	// Load .env file
	godotenv.Load()

	rpcURL := os.Getenv("RPC_URL")
	if rpcURL == "" {
		panic("RPC_URL must be set in .env file")
	}

	ethClient, err := ethclient.Dial(rpcURL)
	if err != nil {
		panic(err)
	}

	beaconChainProofsWrapper, err = BeaconChainProofsWrapper.NewBeaconChainProofsWrapper(BEACON_CHAIN_PROOFS_WRAPPER_ADDRESS, ethClient)
	if err != nil {
		panic(err)
	}

	// Run tests for each hard fork type
	if err := loadBeaconState(
		"data/electra_mekong_beacon_headers_654719.json",
		"data/electra_mekong_beacon_state_654719.ssz",
		17000, // Use 17000 for Mekong, this check isn't relevant for Electra
	); err != nil {
		panic(err)
	}
	m.Run()

	if err := loadBeaconState(
		"data/deneb_holesky_beacon_headers_2227472.json",
		"data/deneb_holesky_beacon_state_2227472.ssz",
		17000,
	); err != nil {
		panic(err)
	}
	m.Run()
}

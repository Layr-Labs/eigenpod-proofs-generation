package beacon

import (
	"errors"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/electra"
)

func GetGenesisTime(state *spec.VersionedBeaconState) (uint64, error) {
	switch state.Version {
	case spec.DataVersionElectra:
		return state.Electra.GenesisTime, nil
	case spec.DataVersionDeneb:
		return state.Deneb.GenesisTime, nil
	default:
		return 0, errors.New("unsupported beacon state version")
	}
}

func CreateVersionedSignedBlock(block interface{}) (spec.VersionedSignedBeaconBlock, error) {
	var versionedBlock spec.VersionedSignedBeaconBlock

	switch s := block.(type) {
	case electra.BeaconBlock:
		var signedBlock electra.SignedBeaconBlock
		signedBlock.Message = &s
		versionedBlock.Electra = &signedBlock
		versionedBlock.Version = spec.DataVersionElectra
	case deneb.BeaconBlock:
		var signedBlock deneb.SignedBeaconBlock
		signedBlock.Message = &s
		versionedBlock.Deneb = &signedBlock
		versionedBlock.Version = spec.DataVersionDeneb
	default:
		return versionedBlock, errors.New("unsupported beacon block version")
	}
	return versionedBlock, nil
}

func CreateVersionedState(state interface{}) (spec.VersionedBeaconState, error) {
	var versionedState spec.VersionedBeaconState

	switch s := state.(type) {
	case *electra.BeaconState:
		versionedState.Electra = s
		versionedState.Version = spec.DataVersionElectra
	case *deneb.BeaconState:
		versionedState.Deneb = s
		versionedState.Version = spec.DataVersionDeneb
	default:
		return versionedState, errors.New("unsupported beacon state version")
	}
	return versionedState, nil
}

func UnmarshalSSZVersionedBeaconState(data []byte) (*spec.VersionedBeaconState, error) {
	beaconState := &spec.VersionedBeaconState{}
	electraBeaconState := &electra.BeaconState{}
	err := electraBeaconState.UnmarshalSSZ(data)
	if err != nil {
		// If Electra fails, try Deneb
		denebBeaconState := &deneb.BeaconState{}
		err = denebBeaconState.UnmarshalSSZ(data)
		if err != nil {
			return nil, err
		} else {
			beaconState.Deneb = denebBeaconState
			beaconState.Version = spec.DataVersionDeneb
		}
	} else {
		beaconState.Electra = electraBeaconState
		beaconState.Version = spec.DataVersionElectra
	}

	return beaconState, nil
}

func MarshalSSZVersionedBeaconState(beaconState spec.VersionedBeaconState) ([]byte, error) {
	var data []byte
	var err error
	// Try to marshal using Electra
	if beaconState.Version == spec.DataVersionElectra {
		data, err = beaconState.Electra.MarshalSSZ()
		if err != nil {
			return nil, err
		}
	} else {
		data, err = beaconState.Deneb.MarshalSSZ()
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

package common

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"math/bits"

	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ssz "github.com/ferranbt/fastssz"
)

type Bytes32 [32]byte

func (p *Bytes32) MarshalJSON() ([]byte, error) {
	return json.Marshal(hexutil.Encode(p[:]))
}

func (p *Bytes32) UnmarshalJSON(data []byte) error {
	// Unmarshal JSON string to a regular Go string
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	// Decode hex string to bytes
	decoded, err := hexutil.Decode(s)
	if err != nil {
		return err
	}

	// Length check
	if len(decoded) != 32 {
		return fmt.Errorf("expected 32 bytes but got %d", len(decoded))
	}

	// Populate the array
	copy(p[:], decoded)

	return nil
}

func BigToLittleEndian(input *big.Int) [32]byte {
	var littleEndian [32]byte
	intBytes := input.Bytes()
	intBytesLen := len(intBytes)
	for i := 0; i < intBytesLen; i++ {
		littleEndian[i] = intBytes[intBytesLen-1-i]
	}
	return littleEndian
}

func ConvertUint64ToRoot(n uint64) phase0.Root {
	hh := ssz.NewHasher()
	hh.PutUint64(uint64(n))
	return ConvertTo32ByteArray(hh.Hash())
}

func ConvertUint64ToBytes32(n uint64) Bytes32 {
	hh := ssz.NewHasher()
	hh.PutUint64(uint64(n))
	return ConvertTo32ByteArray(hh.Hash())
}

func ConvertTo32ByteArray(b []byte) [32]byte {
	var result [32]byte
	copy(result[:], b)
	return result
}

func GetDepth(d uint64) uint8 {
	if d <= 1 {
		return 0
	}
	i := NextPowerOfTwo(d)
	return 64 - uint8(bits.LeadingZeros(i)) - 1
}

func NextPowerOfTwo(v uint64) uint {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return uint(v)
}

// ceilLog2 calculates the ceiling of the base-2 logarithm of x.
func CeilLog2(x int) uint64 {
	// Using math.Log2 for logarithm base 2, and math.Ceil for ceiling
	return uint64(math.Ceil(math.Log2(float64(x))))
}

func GetSlotTimestamp(beaconState *deneb.BeaconState, blockHeader *phase0.BeaconBlockHeader) uint64 {
	return beaconState.GenesisTime + uint64(blockHeader.Slot)*12
}

func ConvertValidatorToValidatorFields(v *phase0.Validator) []Bytes32 {
	validatorFields := make([]Bytes32, 0)
	hh := ssz.NewHasher()

	hh.PutBytes(v.PublicKey[:])
	validatorFields = append(validatorFields, ConvertTo32ByteArray(hh.Hash()))
	hh.Reset()

	hh.PutBytes(v.WithdrawalCredentials)
	validatorFields = append(validatorFields, ConvertTo32ByteArray(hh.Hash()))
	hh.Reset()

	hh.PutUint64(uint64(v.EffectiveBalance))
	validatorFields = append(validatorFields, ConvertTo32ByteArray(hh.Hash()))
	hh.Reset()

	hh.PutBool(v.Slashed)
	validatorFields = append(validatorFields, ConvertTo32ByteArray(hh.Hash()))
	hh.Reset()

	hh.PutUint64(uint64(v.ActivationEligibilityEpoch))
	validatorFields = append(validatorFields, ConvertTo32ByteArray(hh.Hash()))
	hh.Reset()

	hh.PutUint64(uint64(v.ActivationEpoch))
	validatorFields = append(validatorFields, ConvertTo32ByteArray(hh.Hash()))
	hh.Reset()

	hh.PutUint64(uint64(v.ExitEpoch))
	validatorFields = append(validatorFields, ConvertTo32ByteArray(hh.Hash()))
	hh.Reset()

	hh.PutUint64(uint64(v.WithdrawableEpoch))
	validatorFields = append(validatorFields, ConvertTo32ByteArray(hh.Hash()))
	hh.Reset()

	return validatorFields
}

func ConvertWithdrawalToWithdrawalFields(w *capella.Withdrawal) []Bytes32 {
	var withdrawalFields []Bytes32
	hh := ssz.NewHasher()

	hh.PutUint64(uint64(w.Index))
	withdrawalFields = append(withdrawalFields, ConvertTo32ByteArray(hh.Hash()))
	hh.Reset()

	hh.PutUint64(uint64(w.ValidatorIndex))
	withdrawalFields = append(withdrawalFields, ConvertTo32ByteArray(hh.Hash()))
	hh.Reset()

	hh.PutBytes(w.Address[:])
	withdrawalFields = append(withdrawalFields, ConvertTo32ByteArray(hh.Hash()))
	hh.Reset()

	hh.PutUint64(uint64(w.Amount))
	withdrawalFields = append(withdrawalFields, ConvertTo32ByteArray(hh.Hash()))
	hh.Reset()

	return withdrawalFields
}

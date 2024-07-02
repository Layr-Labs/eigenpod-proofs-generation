package main

import (
	"bytes"
	sha256 "crypto/sha256"
	"errors"
	"fmt"
	"os"

	"context"

	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/onchain"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	cli "github.com/urfave/cli/v2"
)

func main() {
	var eigenpodAddress, beacon, node, owner, output string
	ctx := context.Background()

	app := &cli.App{
		Name:                   "Eigenlayer Proofs CLi",
		HelpName:               "eigenproofs",
		Usage:                  "TODO: usage",
		EnableBashCompletion:   true,
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			{
				Name:  "checkpoint",
				Usage: "Generates a proof for use with EigenPod.verifyCheckpointProofs().",
				Action: func(cctx *cli.Context) error {
					var out, owner *string = nil, nil

					if len(cctx.String("out")) > 0 {
						outProp := cctx.String("out")
						out = &outProp
					}

					if len(cctx.String("owner")) > 0 {
						ownerProp := cctx.String("owner")
						owner = &ownerProp
					}

					execute(ctx, eigenpodAddress, beacon, node, "checkpoint", out, owner)
					return nil
				},
			},
			{
				Name:  "validator",
				Usage: "Generates a proof for use with EigenPod.verifyWithdrawalCredentials()",
				Action: func(cctx *cli.Context) error {

					var out, owner *string = nil, nil

					if len(cctx.String("out")) > 0 {
						outProp := cctx.String("out")
						out = &outProp
					}

					if len(cctx.String("owner")) > 0 {
						ownerProp := cctx.String("owner")
						owner = &ownerProp
					}

					execute(ctx, eigenpodAddress, beacon, node, "validator", out, owner)
					return nil
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "eigenpodAddress",
				Aliases:     []string{"e"},
				Value:       "",
				Usage:       "[required] The onchain address of your eigenpod contract (0x123123123123)",
				Required:    true,
				Destination: &eigenpodAddress,
			},
			&cli.StringFlag{
				Name:        "beacon",
				Aliases:     []string{"b"},
				Value:       "",
				Usage:       "[required] URI to a functioning beacon node RPC (https://)",
				Required:    true,
				Destination: &beacon,
			},
			&cli.StringFlag{
				Name:        "node",
				Aliases:     []string{"n"},
				Value:       "",
				Usage:       "[required] URI to a functioning execution-layer RPC",
				Required:    true,
				Destination: &node,
			},
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Value:       "",
				Usage:       "Output path for the proof. (defaults to stdout)",
				Destination: &output,
			},
			&cli.StringFlag{
				Name:        "owner",
				Aliases:     []string{},
				Destination: &owner,
				Value:       "",
				Usage:       "Private key of the owner. If set, this will automatically submit the proofs to their corresponding onchain functions after generation. If using `checkpoint` mode, it will also begin a checkpoint if one hasn't been started already.",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err) // burn it all to the ground
	}
}

func getBeaconClient(beaconUri string) (BeaconClient, error) {
	beaconClient, _, err := NewBeaconClient(beaconUri)
	return beaconClient, err
}

func lastCheckpointedForEigenpod(eigenpodAddress string, client *ethclient.Client) uint64 {
	eigenPod, err := onchain.NewEigenPod(common.HexToAddress(eigenpodAddress), client)
	PanicOnError("failed to locate eigenpod. is your address correct?", err)

	timestamp, err := eigenPod.CurrentCheckpointTimestamp(nil)
	PanicOnError("failed to locate eigenpod. Is your address correct?", err)

	return timestamp
}

// search through beacon state for validators whose withdrawal address is set to eigenpod.
func findAllValidatorsForEigenpod(eigenpodAddress string, beaconState *spec.VersionedBeaconState) []ValidatorWithIndex {
	allValidators, err := beaconState.Validators()
	PanicOnError("failed to fetch beacon state", err)

	eigenpodAddressBytes := common.FromHex(eigenpodAddress)

	var outputValidators []ValidatorWithIndex = []ValidatorWithIndex{}
	var i uint64 = 0
	maxValidators := uint64(len(allValidators))
	for i = 0; i < maxValidators; i++ {
		validator := allValidators[i]
		if validator == nil || validator.WithdrawalCredentials[0] != 1 { // withdrawalCredentials _need_ their first byte set to 1 to withdraw to execution layer.
			continue
		}
		// we check that the last 20 bytes of expectedCredentials matches validatorCredentials.
		if bytes.Equal(
			eigenpodAddressBytes[:],
			validator.WithdrawalCredentials[12:], // first 12 bytes are not the pubKeyHash, see (https://github.com/Layr-Labs/eigenlayer-contracts/blob/d148952a2942a97a218a2ab70f9b9f1792796081/src/contracts/pods/EigenPod.sol#L663)
		) {
			outputValidators = append(outputValidators, ValidatorWithIndex{
				Validator: validator,
				Index:     i,
			})
		}
	}
	return outputValidators
}

func getOnchainValidatorInfo(client *ethclient.Client, eigenpodAddress string, allValidators []ValidatorWithIndex) []onchain.IEigenPodValidatorInfo {
	eigenPod, err := onchain.NewEigenPod(common.HexToAddress(eigenpodAddress), client)
	PanicOnError("failed to locate Eigenpod. Is your address correct?", err)

	var validatorInfo []onchain.IEigenPodValidatorInfo = []onchain.IEigenPodValidatorInfo{}

	// TODO: batch/multicall
	zeroes := [16]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	for i := 0; i < len(allValidators); i++ {
		// ssz requires values to be 32-byte aligned, which requires 16 bytes of 0's to be added
		// prior to hashing.
		pubKeyHash := sha256.Sum256(
			append(
				(allValidators[i]).Validator.PublicKey[:],
				zeroes[:]...,
			),
		)
		info, err := eigenPod.ValidatorPubkeyHashToInfo(nil, pubKeyHash)
		PanicOnError("failed to fetch validator eigeninfo.", err)
		validatorInfo = append(validatorInfo, info)
	}

	return validatorInfo
}

func getCurrentCheckpointBlockRoot(eigenpodAddress string, eth *ethclient.Client) (*[32]byte, error) {
	eigenPod, err := onchain.NewEigenPod(common.HexToAddress(eigenpodAddress), eth)
	PanicOnError("failed to locate Eigenpod. Is your address correct?", err)

	checkpoint, err := eigenPod.CurrentCheckpoint(nil)
	PanicOnError("failed to reach eigenpod.", err)

	return &checkpoint.BeaconBlockRoot, nil
}

func execute(ctx context.Context, eigenpodAddress, beacon_node_uri, node, command string, out *string, owner *string) {
	eth, err := ethclient.Dial(node)
	PanicOnError("failed to reach eth --node.", err)

	chainId, err := eth.ChainID(ctx)
	PanicOnError("failed to fetch chain id", err)

	beaconClient, err := getBeaconClient(beacon_node_uri)
	PanicOnError("failed to reach beacon chain.", err)

	if command == "checkpoint" {
		RunCheckpointProof(ctx, eigenpodAddress, eth, chainId, beaconClient, out, owner)
	} else if command == "validator" {
		RunValidatorProof(ctx, eigenpodAddress, eth, chainId, beaconClient, out, owner)
	} else {
		PanicOnError(fmt.Sprintf("invalid --prove argument. Expected 'checkpoint' or 'validator' (got `%s`)", command), errors.New("invalid command"))
	}
}

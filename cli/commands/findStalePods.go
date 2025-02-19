package commands

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/utils"
	"github.com/fatih/color"
)

type TFindStalePodsCommandArgs struct {
	EthNode    string
	BeaconNode string
	Verbose    bool
	Tolerance  float64
}

func FindStalePodsCommand(args TFindStalePodsCommandArgs) error {
	ctx := context.Background()
	eth, beacon, chainId, err := utils.GetClients(ctx, args.EthNode, args.BeaconNode /* verbose */, args.Verbose)
	utils.PanicOnError("failed to dial clients", err)

	results, err := core.FindStaleEigenpods(ctx, eth, args.EthNode, beacon, chainId, args.Verbose, args.Tolerance)
	utils.PanicOnError("failed to find stale eigenpods", err)

	if !args.Verbose {
		PrintAsJSON(results)
		return nil
	}

	if args.Verbose {
		for pod, res := range results {
			color.Red("pod %s\n", pod)
			for _, validator := range res {
				fmt.Printf("\t[#%d] (%s) - %d\n", validator.Index, func() string {
					if validator.Validator.Slashed {
						return "slashed"
					} else {
						return "not slashed"
					}
				}(), validator.Validator.EffectiveBalance)
			}
		}
	}
	return nil
}

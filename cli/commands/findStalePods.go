package commands

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core"
	"github.com/fatih/color"
)

type TFindStalePodsCommandArgs struct {
	EthNode    string
	BeaconNode string
	Verbose    bool
}

func FindStalePodsCommand(args TFindStalePodsCommandArgs) error {
	ctx := context.Background()
	eth, beacon, chainId, err := core.GetClients(ctx, args.EthNode, args.BeaconNode /* verbose */, args.Verbose)
	core.PanicOnError("failed to dial clients", err)

	results, err := core.FindStaleEigenpods(ctx, eth, args.EthNode, beacon, chainId, args.Verbose)
	core.PanicOnError("failed to find stale eigenpods", err)

	if !args.Verbose {
		// just print json and be done
		jsonStr, _ := json.MarshalIndent(results, "", "   ")
		fmt.Println(string(jsonStr))
		return nil
	}

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
	return nil
}

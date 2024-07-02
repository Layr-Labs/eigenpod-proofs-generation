//go:build !js && !wasm

package main

import (
	"context"
	"os"

	cli "github.com/urfave/cli/v2"
)

func main() {
	var eigenpodAddress, beacon, node, owner, output string
	ctx := context.Background()

	app := &cli.App{
		Name:                   "Eigenlayer Proofs CLi",
		HelpName:               "eigenproofs",
		Usage:                  "Generates proofs to (1) checkpoint your validators, or (2) verify the withdrawal credentials of an inactive validator.",
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

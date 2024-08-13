package main

import cli "github.com/urfave/cli/v2"

// Required for commands that need an EigenPod's address
var POD_ADDRESS_FLAG = &cli.StringFlag{
	Name:        "podAddress",
	Aliases:     []string{"p", "pod"},
	Value:       "",
	Usage:       "[required] The onchain `address` of your eigenpod contract (0x123123123123)",
	Required:    true,
	Destination: &eigenpodAddress,
}

// Required for commands that need a beacon chain RPC
var BEACON_NODE_FLAG = &cli.StringFlag{
	Name:        "beaconNode",
	Aliases:     []string{"b"},
	Value:       "",
	Usage:       "[required] `URL` to a functioning beacon node RPC (https://)",
	Required:    true,
	Destination: &beacon,
}

// Required for commands that need an execution layer RPC
var EXEC_NODE_FLAG = &cli.StringFlag{
	Name:        "execNode",
	Aliases:     []string{"e"},
	Value:       "",
	Usage:       "[required] `URL` to a functioning execution-layer RPC (https://)",
	Required:    true,
	Destination: &node,
}
var PRINT_CALLDATA_BUT_DO_NOT_EXECUTE_FLAG = &cli.BoolFlag{
	Name:        "print-calldata",
	Value:       false,
	Usage:       "Print the calldata for all associated transactions, but do not execute them. Note that some transactions have an order dependency (you cannot submit checkpoint proofs if you haven't started a checkpoint) so this may require you to get your pod into the correct state before usage.",
	Required:    false,
	Destination: &simulateTransaction,
}

// Optional commands:

// Optional use for commands that want direct tx submission from a specific private key
var SENDER_PK_FLAG = &cli.StringFlag{
	Name:        "sender",
	Aliases:     []string{"s"},
	Value:       "",
	Usage:       "`Private key` of the account that will send any transactions. If set, this will automatically submit the proofs to their corresponding onchain functions after generation. If using checkpoint mode, it will also begin a checkpoint if one hasn't been started already.",
	Destination: &sender,
}

// Optional use for commands that support JSON output
var PRINT_JSON_FLAG = &cli.BoolFlag{
	Name:        "json",
	Value:       false,
	Usage:       "print only plain JSON",
	Required:    false,
	Destination: &useJson,
}

var PROOF_PATH_FLAG = &cli.StringFlag{
	Name:        "proof",
	Value:       "",
	Usage:       "the `path` to a previous proof generated from this step (via -o proof.json). If provided, this proof will submitted to network via the --sender flag.",
	Destination: &proofPath,
}

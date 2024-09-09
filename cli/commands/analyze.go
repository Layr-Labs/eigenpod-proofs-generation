package commands

import (
	"context"
	"encoding/csv"
	"os"

	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fatih/color"
)

type TAnalyzeArgs struct {
	EigenpodAddress string
	DisableColor    bool
	UseJSON         bool
	Node            string
	BeaconNode      string
	Verbose         bool
}

var podDataPath = "../pod_deployed.csv"

func AnalyzeCommand(args TAnalyzeArgs) error {
	ctx := context.Background()
	if args.DisableColor {
		color.NoColor = true
	}

	isVerbose := !args.UseJSON

	eth, beaconClient, _, err := core.GetClients(ctx, args.Node, args.BeaconNode, isVerbose)
	core.PanicOnError("failed to load ethereum clients", err)

	file, err := os.Open(podDataPath)
	defer file.Close()
	core.PanicOnError("failed to open csv: %w", err)

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	core.PanicOnError("error reading records: %w", err)

	pods := make(map[string]core.PodInfo)
	for _, record := range records {
		pods[record[0]] = core.PodInfo{
			PodAddress: common.HexToAddress(record[0]),
			Owner:      common.HexToAddress(record[1]),
		}
	}

	core.AnalyzePods(ctx, pods, eth, beaconClient)

	return nil
}

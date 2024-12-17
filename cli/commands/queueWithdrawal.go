package commands

import (
	"context"

	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/IDelegationManager"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core"
	"github.com/pkg/errors"
)

type TQueueWithdrawallArgs struct {
	EthNode    string
	BeaconNode string
	EigenPod   string
}

func QueueWithdrawalCommand(args TQueueWithdrawallArgs) error {
	ctx := context.Background()
	eth, _, chainId, err := core.GetClients(ctx, args.EthNode, args.BeaconNode, false /* isVerbose */)
	core.PanicOnError("failed to dial nodes", err)

	_dm, err := IDelegationManager.NewIDelegationManager(DelegationManager(chainId), eth)
	core.PanicOnError("failed to reach delegation manager", err)

	// TODO: wait for G's conversion function from deposit[ed] shares to depositShares
	// bound the withdrawals by REG - (sum(allWithdrawalsQueued))

	/*
		struct QueuedWithdrawalParams {
			// Array of strategies that the QueuedWithdrawal contains
			IStrategy[] strategies; // native eth strategy
			// Array containing the amount of depositShares for withdrawal in each Strategy in the `strategies` array
			// Note that the actual shares received on completing withdrawal may be less than the depositShares if slashing occurred
			uint256[] depositShares;
			// The address of the withdrawer
			address withdrawer;
		}
	*/
	return errors.New("unimplemented.")
}

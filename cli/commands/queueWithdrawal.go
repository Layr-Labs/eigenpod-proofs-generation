package commands

import "github.com/pkg/errors"

func QueueWithdrawalCommand(args TComputeCheckpointableValueCommandArgs) error {
	// TODO: IDelegationManager.queueWithdrawals
	/*
		struct QueuedWithdrawalParams {
			// Array of strategies that the QueuedWithdrawal contains
			IStrategy[] strategies;
			// Array containing the amount of depositShares for withdrawal in each Strategy in the `strategies` array
			// Note that the actual shares received on completing withdrawal may be less than the depositShares if slashing occurred
			uint256[] depositShares;
			// The address of the withdrawer
			address withdrawer;
		}
	*/
	return errors.New("unimplemented.")
}

package commands

import "github.com/pkg/errors"

type TQueueWithdrawallArgs struct {
	EthNode    string
	BeaconNode string
	EigenPod   string
}

func QueueWithdrawalCommand(args TQueueWithdrawallArgs) error {
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

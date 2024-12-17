package commands

import "github.com/pkg/errors"

type TCompleteWithdrawalArgs struct {
	EthNode    string
	BeaconNode string
	EigenPod   string
}

func CompleteAllWithdrawalsCommand(args TCompleteWithdrawalArgs) error {
	/*
		TODO: IDelegationManager.completeQueuedWithdrawals(
			IERC20[][] calldata tokens,
			bool[] calldata receiveAsTokens,
			uint256 numToComplete
		)
	*/
	return errors.New("unimplemented.")
}

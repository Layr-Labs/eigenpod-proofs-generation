package commands

import "github.com/pkg/errors"

type TShowWithdrawalArgs struct {
	EthNode    string
	BeaconNode string
	EigenPod   string
}

func ShowWithdrawalsCommand(args TShowWithdrawalArgs) error {
	// IDelegationManager.getQueuedWithdrawals
	return errors.New("unimplemented.")
}

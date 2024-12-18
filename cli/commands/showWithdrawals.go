package commands

import (
	"context"
	"math/big"

	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/EigenPod"
	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/IDelegationManager"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core"
	"github.com/ethereum/go-ethereum/common"
)

type TShowWithdrawalArgs struct {
	EthNode    string
	EigenPod   string
	Strategies common.Address
}

func ShowWithdrawalsCommand(args TShowWithdrawalArgs) error {
	ctx := context.Background()
	eth, chainId, err := core.GetEthClient(ctx, args.EthNode)
	core.PanicOnError("failed to reach eth and beacon node", err)

	curBlock, err := eth.BlockByNumber(ctx, nil) /* head */
	core.PanicOnError("failed to load curBlock", err)

	dm, err := IDelegationManager.NewIDelegationManager(DelegationManager(chainId), eth)
	core.PanicOnError("failed to reach delegation manager", err)

	pod, err := EigenPod.NewEigenPod(common.HexToAddress(args.EigenPod), eth)
	core.PanicOnError("failed to reach eigenpod manager", err)

	podOwner, err := pod.PodOwner(nil)
	core.PanicOnError("failed to load podOwner", err)

	allWithdrawals, err := dm.GetQueuedWithdrawals(nil, podOwner)
	core.PanicOnError("failed to get queued withdrawals", err)

	type TWithdrawalInfo struct {
		Staker              string
		Strategy            common.Address
		CurrentBlock        uint64
		AvailableAfterBlock *big.Int
		Ready               bool
		TotalAmountETH      *big.Float
	}

	minDelay, err := dm.MinWithdrawalDelayBlocks(nil)
	core.PanicOnError("failed to get minWithdrawalDelay", err)

	withdrawalInfo := []TWithdrawalInfo{}

	for i, shares := range allWithdrawals.Shares {
		withdrawal := allWithdrawals.Withdrawals[i]

		// this cli is only for withdrawals of beaconstrategy for a single strategy
		if withdrawal.Strategies[0].Cmp(core.BeaconStrategy()) != 0 || len(withdrawal.Strategies) != 1 {
			continue
		}

		targetBlock := new(big.Int).SetUint64(uint64(allWithdrawals.Withdrawals[i].StartBlock + minDelay))

		withdrawalInfo = append(withdrawalInfo, TWithdrawalInfo{
			TotalAmountETH:      core.GweiToEther(core.WeiToGwei(shares[0])),
			Strategy:            allWithdrawals.Withdrawals[i].Strategies[0],
			Staker:              allWithdrawals.Withdrawals[i].Staker.Hex(),
			CurrentBlock:        curBlock.NumberU64(),
			AvailableAfterBlock: targetBlock,
			Ready:               targetBlock.Uint64() < curBlock.NumberU64(),
		})
	}

	printAsJSON(withdrawalInfo)
	return nil
}

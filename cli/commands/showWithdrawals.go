package commands

import (
	"context"
	"math/big"
	"time"

	lo "github.com/samber/lo"

	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/EigenPod"
	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/IDelegationManager"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core"
	"github.com/ethereum/go-ethereum/common"
)

type TShowWithdrawalArgs struct {
	EthNode    string
	BeaconNode string
	EigenPod   string
}

func ShowWithdrawalsCommand(args TShowWithdrawalArgs) error {
	ctx := context.Background()
	eth, _, chainId, err := core.GetClients(ctx, args.EthNode, args.BeaconNode, false /* isVerbose */)
	core.PanicOnError("failed to reach eth and beacon node", err)

	curBlock, err := eth.BlockByNumber(ctx, nil) /* head */
	core.PanicOnError("failed to load curBlock", err)

	genBlock, err := eth.BlockByNumber(ctx, big.NewInt(0)) /* head */
	core.PanicOnError("failed to load genesis block", err)

	timePerBlockSeconds := float64(curBlock.NumberU64()-genBlock.NumberU64()) / float64(curBlock.Time()-genBlock.Time())

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
		AvailableAfter      string
		AvailableAfterBlock *big.Int
		Ready               bool
		TotalAmountETH      *big.Float
	}

	minDelay, err := dm.MinWithdrawalDelayBlocks(nil)
	core.PanicOnError("failed to get minWithdrawalDelay", err)

	withdrawalInfo := []TWithdrawalInfo{}

	for i, shares := range allWithdrawals.Shares {
		withdrawalTotalValueWei := lo.Reduce(shares, func(accum *big.Int, item *big.Int, i int) *big.Int {
			return new(big.Int).Add(item, accum)
		}, big.NewInt(0))

		targetBlock := new(big.Int).SetUint64(uint64(allWithdrawals.Withdrawals[i].StartBlock + minDelay))

		blockDeltaSeconds := (targetBlock.Uint64() - curBlock.NumberU64()) * uint64(timePerBlockSeconds)
		availableAfter := time.Now().Add(time.Second * time.Duration(blockDeltaSeconds))

		withdrawalInfo = append(withdrawalInfo, TWithdrawalInfo{
			TotalAmountETH:      core.GweiToEther(core.WeiToGwei(withdrawalTotalValueWei)),
			Staker:              allWithdrawals.Withdrawals[i].Staker.Hex(),
			AvailableAfterBlock: targetBlock,
			AvailableAfter:      availableAfter.String(),
			Ready:               targetBlock.Uint64() <= curBlock.NumberU64(),
		})
	}
	return nil
}

package commands

import (
	"context"
	"fmt"
	"math/big"

	lo "github.com/samber/lo"

	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/EigenPod"
	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/IDelegationManager"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TCompleteWithdrawalArgs struct {
	EthNode  string
	EigenPod string
	Sender   string
}

func DelegationManager(chainId *big.Int) common.Address {
	data := map[uint64]string{
		// TODO(zeus) - make this runnable via zeus.
		1:     "0x39053D51B77DC0d36036Fc1fCc8Cb819df8Ef37A", // mainnet
		17000: "0x75dfE5B44C2E530568001400D3f704bC8AE350CC", // holesky preprod
	}
	contract, ok := data[chainId.Uint64()]
	if !ok {
		panic("no delegation manager found for chain")
	}
	addr := common.HexToAddress(contract)
	return addr
}

func CompleteAllWithdrawalsCommand(args TCompleteWithdrawalArgs) error {
	ctx := context.Background()

	eth, err := ethclient.DialContext(ctx, args.EthNode)
	core.PanicOnError("failed to reach eth node", err)

	chainId, err := eth.ChainID(nil)
	core.PanicOnError("failed to load chainId", err)

	curBlockNumber, err := eth.BlockNumber(nil)

	pod, err := EigenPod.NewEigenPod(common.HexToAddress(args.EigenPod), eth)
	core.PanicOnError("failed to reach eigenpod", err)

	reg, err := pod.WithdrawableRestakedExecutionLayerGwei(nil)
	core.PanicOnError("failed to fetch REG", err)

	podOwner, err := pod.PodOwner(nil)
	core.PanicOnError("failed to read podOwner", err)

	delegationManager, err := IDelegationManager.NewIDelegationManager(DelegationManager(chainId), eth)
	core.PanicOnError("failed to reach delegation manager", err)

	minDelay, err := delegationManager.MinWithdrawalDelayBlocks(nil)
	core.PanicOnError("failed to read MinWithdrawalDelayBlocks", err)

	queuedWithdrawals, err := delegationManager.GetQueuedWithdrawals(nil, podOwner)
	core.PanicOnError("failed to read queuedWithdrawals", err)

	beaconETHStrategy := common.HexToAddress("0xbeaC0eeEeeeeEEeEeEEEEeeEEeEeeeEeeEEBEaC0")

	eligibleWithdrawals := lo.Map(queuedWithdrawals.Withdrawals, func(withdrawal IDelegationManager.IDelegationManagerTypesWithdrawal, index int) *IDelegationManager.IDelegationManagerTypesWithdrawal {
		isBeaconWithdrawal := len(withdrawal.Strategies) == 1 && withdrawal.Strategies[0].Cmp(beaconETHStrategy) == 0
		isExecutable := curBlockNumber <= uint64(withdrawal.StartBlock+minDelay)
		if isBeaconWithdrawal && isExecutable {
			return &withdrawal
		}
		return nil
	})

	var runningSum uint64 = 0
	affordedWithdrawals := lo.Map(eligibleWithdrawals, func(withdrawal *IDelegationManager.IDelegationManagerTypesWithdrawal, index int) *IDelegationManager.IDelegationManagerTypesWithdrawal {
		if withdrawal == nil {
			return nil
		}
		withdrawalShares := queuedWithdrawals.Shares[index][0].Uint64()
		if reg < (runningSum + withdrawalShares) {
			runningSum = runningSum + withdrawalShares
			return withdrawal
		}
		return nil
	})

	// filter out any nils.
	affordedWithdrawals = lo.Filter(affordedWithdrawals, func(withdrawal *IDelegationManager.IDelegationManagerTypesWithdrawal, index int) bool {
		return withdrawal != nil
	})

	if len(affordedWithdrawals) != len(eligibleWithdrawals) {
		fmt.Printf("WARN: Your pod has %d withdrawals available, but you only have enough balance to satisfy %d of them.\n", len(eligibleWithdrawals), len(affordedWithdrawals))
		fmt.Printf("Consider checkpointing to claim beacon rewards, or depositing ETH and checkpointing to complete these withdrawals.\n\n")
	}

	fmt.Printf("Your podOwner(%s) has %d withdrawals that can be completed right now.\n", podOwner.Hex(), len(affordedWithdrawals))
	fmt.Printf("Total ETH: %sETH\n", core.GweiToEther(core.WeiToGwei(new(big.Int).SetUint64(runningSum))).String())

	core.PanicIfNoConsent("Would you like to continue?")

	withdrawals := lo.Map(affordedWithdrawals, func(w *IDelegationManager.IDelegationManagerTypesWithdrawal, i int) IDelegationManager.IDelegationManagerTypesWithdrawal {
		return *w
	})

	tokens := lo.Map(withdrawals, func(_ IDelegationManager.IDelegationManagerTypesWithdrawal, _ int) []common.Address {
		return []common.Address{common.BigToAddress(big.NewInt(0))}
	})

	receiveAsTokens := lo.Map(withdrawals, func(_ IDelegationManager.IDelegationManagerTypesWithdrawal, _ int) bool {
		return true
	})

	txn, err := delegationManager.CompleteQueuedWithdrawals0(nil, withdrawals, tokens, receiveAsTokens)
	core.PanicOnError("CompleteQueuedWithdrawals failed.", err)

	fmt.Printf("%s\n", txn.Hash())
	return nil
}

package commands

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/EigenPod"
	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/IDelegationManager"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TQueueWithdrawallArgs struct {
	EthNode  string
	EigenPod string
	Sender   string
}

func QueueWithdrawalCommand(args TQueueWithdrawallArgs) error {
	ctx := context.Background()

	eth, err := ethclient.DialContext(ctx, args.EthNode)
	core.PanicOnError("failed to reach eth node", err)

	chainId, err := eth.ChainID(nil)
	core.PanicOnError("failed to load chainId", err)

	dm, err := IDelegationManager.NewIDelegationManager(DelegationManager(chainId), eth)
	core.PanicOnError("failed to reach delegation manager", err)

	// bound the withdrawals by REG - (sum(allWithdrawalsQueued))

	pod, err := EigenPod.NewEigenPod(common.HexToAddress(args.EigenPod), eth)
	core.PanicOnError("failed to reach eigenpod", err)

	_reg, err := pod.WithdrawableRestakedExecutionLayerGwei(nil)
	core.PanicOnError("failed to load REG", err)
	reg := new(big.Int).SetUint64(_reg)

	podOwner, err := pod.PodOwner(nil)
	core.PanicOnError("failed to read podOwner", err)

	res, err := dm.GetWithdrawableShares(nil, podOwner, []common.Address{core.BeaconStrategy()})
	core.PanicOnError("failed to read beacon strategy withdrawable shares", err)

	_depositShares, err := dm.ConvertToDepositShares(nil, podOwner, []common.Address{core.BeaconStrategy()}, []*big.Int{res.WithdrawableShares[0]})
	core.PanicOnError("failed to compute deposit shares", err)
	depositShares := _depositShares[0]

	minWithdrawalDelay, err := dm.MinWithdrawalDelayBlocks(nil)
	curBlock, err := eth.BlockNumber(ctx)

	core.PanicOnError("failed to load minimum withdrawal delay", err)

	depositSharesGwei := core.IGweiToWei(depositShares)

	var amountToWithdrawDepositShares *big.Int = new(big.Int)
	*amountToWithdrawDepositShares = *depositSharesGwei

	fmt.Printf("In the Native ETH strategy, you have %sETH to be withdrawn.\n", core.GweiToEther(new(big.Float).SetInt(depositSharesGwei)))
	fmt.Printf("NOTE: If you were or become slashed on EigenLayer during the withdrawal period, the total amount received will be less any slashed amount.\n")

	if depositSharesGwei.Cmp(reg) > 0 {
		fmt.Printf("Queueing a partial withdrawal. Your pod only had %sETH available to satisfy withdrawals.")
		fmt.Printf("Checkpointing may update this balance, if you have any uncheckpointed native eth or beacon rewards.")
		*amountToWithdrawDepositShares = *reg
	}

	core.PanicIfNoConsent(fmt.Sprintf("Would you like to queue a withdrawal %sETH from the Native ETH strategy? This will be withdrawable after approximately block #%d (current block: %d)\n", amountToWithdrawDepositShares, curBlock+uint64(minWithdrawalDelay), curBlock))

	txn, err := dm.QueueWithdrawals(nil, []IDelegationManager.IDelegationManagerTypesQueuedWithdrawalParams{
		IDelegationManager.IDelegationManagerTypesQueuedWithdrawalParams{
			Strategies:    []common.Address{core.BeaconStrategy()},
			DepositShares: []*big.Int{amountToWithdrawDepositShares},
			Withdrawer:    podOwner,
		},
	})
	core.PanicOnError("failed to queue withdrawal", err)

	fmt.Printf("Txn: %s\n", txn.Hash().Hex())
	return nil
}

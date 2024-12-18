package commands

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/EigenPod"
	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/IDelegationManager"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TQueueWithdrawallArgs struct {
	EthNode     string
	EigenPod    string
	Sender      string
	EstimateGas bool
	AmountWei   uint64
}

func QueueWithdrawalCommand(args TQueueWithdrawallArgs) error {
	ctx := context.Background()

	isSimulation := args.EstimateGas

	eth, err := ethclient.DialContext(ctx, args.EthNode)
	core.PanicOnError("failed to reach eth node", err)

	chainId, err := eth.ChainID(ctx)
	core.PanicOnError("failed to load chainId", err)

	acc, err := core.PrepareAccount(&args.Sender, chainId, args.EstimateGas)
	core.PanicOnError("failed to parse private key", err)

	dm, err := IDelegationManager.NewIDelegationManager(DelegationManager(chainId), eth)
	core.PanicOnError("failed to reach delegation manager", err)

	pod, err := EigenPod.NewEigenPod(common.HexToAddress(args.EigenPod), eth)
	core.PanicOnError("failed to reach eigenpod", err)

	_reg, err := pod.WithdrawableRestakedExecutionLayerGwei(nil)
	core.PanicOnError("failed to load REG", err)

	// [withdrawable]RestakedExecutionlayerWei
	rew := core.GweiToWei(new(big.Float).SetUint64(_reg))
	if args.AmountWei > 0 && new(big.Float).SetUint64(args.AmountWei).Cmp(rew) > 0 {
		return errors.New("invalid --amountWei. must be in the range (0, pod.withdrawableRestakedExecutionLayerGwei() as wei]")
	}

	podOwner, err := pod.PodOwner(nil)
	core.PanicOnError("failed to read podOwner", err)

	res, err := dm.GetWithdrawableShares(nil, podOwner, []common.Address{core.BeaconStrategy()})
	core.PanicOnError("failed to read beacon strategy withdrawable shares", err)

	reg := new(big.Int).SetUint64(_reg)

	requestedWithdrawalSizeWei := func() *big.Int {
		if args.AmountWei == 0 {
			// if AmountWei isn't specified, we withdraw all the shares in the beacon strategy.
			return res.WithdrawableShares[0]
		}

		// if it is specified, we withdraw the specific amount.
		return new(big.Int).SetUint64(args.AmountWei)
	}()
	requestedWithdrawalSizeGwei := core.WeiToGwei(requestedWithdrawalSizeWei)

	if requestedWithdrawalSizeGwei.Cmp(new(big.Float).SetInt(res.WithdrawableShares[0])) > 0 {
		// requested to withdraw too many shares.
		return errors.New("the amount to withdraw is larger than the amount of withdrawable shares in the beacon strategy")
	}

	_depositShares, err := dm.ConvertToDepositShares(nil, podOwner, []common.Address{core.BeaconStrategy()}, []*big.Int{requestedWithdrawalSizeWei})
	core.PanicOnError("failed to compute deposit shares", err)
	depositShares := _depositShares[0]

	minWithdrawalDelay, err := dm.MinWithdrawalDelayBlocks(nil)
	core.PanicOnError("failed to load minWithdrawalDelay", err)

	curBlock, err := eth.BlockNumber(ctx)
	core.PanicOnError("failed to load current block number", err)

	depositSharesGwei := core.IGweiToWei(depositShares)

	var amountToWithdrawDepositSharesGwei *big.Int = new(big.Int)
	*amountToWithdrawDepositSharesGwei = *depositSharesGwei

	fmt.Printf("In the Native ETH strategy, you have %sETH to be withdrawn.\n", core.GweiToEther(new(big.Float).SetInt(depositSharesGwei)))
	fmt.Printf("NOTE: If you were or become slashed on EigenLayer during the withdrawal period, the total amount received will be less any slashed amount.\n")

	if depositSharesGwei.Cmp(reg) > 0 {
		fmt.Printf("Queueing a partial withdrawal. Your pod only had %sETH available to satisfy withdrawals.", core.GweiToEther(new(big.Float).SetInt(reg)).String())
		fmt.Printf("Checkpointing may update this balance, if you have any uncheckpointed native eth or beacon rewards.")
		*amountToWithdrawDepositSharesGwei = *reg
	}

	if !isSimulation {
		core.PanicIfNoConsent(fmt.Sprintf("Would you like to queue a withdrawal %sETH from the Native ETH strategy? This will be withdrawable after approximately block #%d (current block: %d)\n", core.GweiToEther(new(big.Float).SetInt(amountToWithdrawDepositSharesGwei)), curBlock+uint64(minWithdrawalDelay), curBlock))
	} else {
		fmt.Printf("THIS IS A SIMULATION. No transaction will be recorded onchain.\n")
	}
	txn, err := dm.QueueWithdrawals(acc.TransactionOptions, []IDelegationManager.IDelegationManagerTypesQueuedWithdrawalParams{
		{
			Strategies:    []common.Address{core.BeaconStrategy()},
			DepositShares: []*big.Int{core.IGweiToWei(amountToWithdrawDepositSharesGwei)},
			Withdrawer:    podOwner,
		},
	})
	core.PanicOnError("failed to queue withdrawal", err)
	if !isSimulation {
		fmt.Printf("%s\n", txn.Hash().Hex())
	} else {
		printAsJSON(Transaction{
			Type:     "queue-withdrawal",
			To:       txn.To().Hex(),
			CallData: common.Bytes2Hex(txn.Data()),
			GasEstimateGwei: func() *uint64 {
				gas := txn.Gas()
				return &gas
			}(),
		})
	}
	return nil
}

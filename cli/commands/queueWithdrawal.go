package commands

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/EigenPod"
	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/IDelegationManager"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
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
	utils.PanicOnError("failed to reach eth node", err)

	chainId, err := eth.ChainID(ctx)
	utils.PanicOnError("failed to load chainId", err)

	acc, err := utils.PrepareAccount(&args.Sender, chainId, args.EstimateGas)
	utils.PanicOnError("failed to parse private key", err)

	dm, err := IDelegationManager.NewIDelegationManager(DelegationManager(chainId), eth)
	utils.PanicOnError("failed to reach delegation manager", err)

	pod, err := EigenPod.NewEigenPod(common.HexToAddress(args.EigenPod), eth)
	utils.PanicOnError("failed to reach eigenpod", err)

	_reg, err := pod.WithdrawableRestakedExecutionLayerGwei(nil)
	utils.PanicOnError("failed to load REG", err)

	// [withdrawable]RestakedExecutionlayerWei
	rew := utils.GweiToWei(new(big.Float).SetUint64(_reg))
	if args.AmountWei > 0 && new(big.Float).SetUint64(args.AmountWei).Cmp(rew) > 0 {
		return errors.New("invalid --amountWei. must be in the range (0, pod.withdrawableRestakedExecutionLayerGwei() as wei]")
	}

	podOwner, err := pod.PodOwner(nil)
	utils.PanicOnError("failed to read podOwner", err)

	reg := new(big.Int).SetUint64(_reg)

	// TODO: maximumWithdrawalSizeGwei = reg - gwei(sumQueuedWithdrawals)
	maximumWithdrawalSizeWei := utils.IGweiToWei(reg)
	var requestedWithdrawalSizeWei *big.Int
	if args.AmountWei == 0 {
		// default to the number of withdrawable shares in the beacon strategy
		requestedWithdrawalSizeWei = maximumWithdrawalSizeWei
	} else {
		// if it is specified, we withdraw the specific amount.
		requestedWithdrawalSizeWei = new(big.Int).SetUint64(args.AmountWei)
	}

	if requestedWithdrawalSizeWei.Cmp(maximumWithdrawalSizeWei) > 0 {
		color.Red(
			"Error: the amount to withdraw from the native ETH strategy (%sETH) is larger than the total withdrawable amount from the pod (%sETH). Will attempt a smaller withdrawal.\n",
			utils.IweiToEther(requestedWithdrawalSizeWei).String(),
			utils.IweiToEther(maximumWithdrawalSizeWei).String(),
		)
		return errors.New("requested to withdraw too many shares")
	}

	_depositShares, err := dm.ConvertToDepositShares(nil, podOwner, []common.Address{core.BeaconStrategy()}, []*big.Int{requestedWithdrawalSizeWei})
	utils.PanicOnError("failed to compute deposit shares", err)
	depositShares := _depositShares[0]

	minWithdrawalDelay, err := dm.MinWithdrawalDelayBlocks(nil)
	utils.PanicOnError("failed to load minWithdrawalDelay", err)

	curBlock, err := eth.BlockNumber(ctx)
	utils.PanicOnError("failed to load current block number", err)

	requestedWithdrawalSizeEth := utils.GweiToEther(utils.WeiToGwei(requestedWithdrawalSizeWei))
	color.Blue("Withdrawing: %sETH.\n", requestedWithdrawalSizeEth.String())
	color.Yellow("NOTE: If you were or become slashed on EigenLayer during the withdrawal period, the total amount received will be less any slashed amount.\n")

	if !isSimulation {
		utils.PanicIfNoConsent(fmt.Sprintf("Would you like to queue a withdrawal %sETH from the Native ETH strategy? This will be withdrawable after approximately block #%d (current block: %d)\n", requestedWithdrawalSizeEth.String(), curBlock+uint64(minWithdrawalDelay), curBlock))
	} else {
		fmt.Printf("THIS IS A SIMULATION. No transaction will be recorded onchain.\n")
	}
	txn, err := dm.QueueWithdrawals(acc.TransactionOptions, []IDelegationManager.IDelegationManagerTypesQueuedWithdrawalParams{
		{
			Strategies:           []common.Address{core.BeaconStrategy()},
			DepositShares:        []*big.Int{depositShares},
			DeprecatedWithdrawer: podOwner,
		},
	})
	utils.PanicOnError("failed to queue withdrawal", err)
	if !isSimulation {
		txnReceipt, err := bind.WaitMined(ctx, eth, txn)
		utils.PanicOnError("failed to wait for txn", err)
		color.Green("%s\n", txnReceipt.TxHash.Hex())
	} else {
		PrintAsJSON(Transaction{
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

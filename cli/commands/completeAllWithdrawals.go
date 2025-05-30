package commands

import (
	"context"
	"fmt"
	"math/big"

	"github.com/pkg/errors"
	lo "github.com/samber/lo"

	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/EigenPod"
	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/IDelegationManager"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
)

type TCompleteWithdrawalArgs struct {
	EthNode     string
	EigenPod    string
	Sender      string
	EstimateGas bool
}

func DelegationManager(chainId *big.Int) common.Address {
	data := map[uint64]string{
		1:      "0x39053D51B77DC0d36036Fc1fCc8Cb819df8Ef37A", // mainnet
		17000:  "0xA44151489861Fe9e3055d95adC98FbD462B948e7", // holesky testnet
		560048: "0x867837a9722C512e0862d8c2E15b8bE220E8b87d", // hoodi testnet
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

	isSimulation := args.EstimateGas

	eth, err := ethclient.DialContext(ctx, args.EthNode)
	utils.PanicOnError("failed to reach eth node", err)

	chainId, err := eth.ChainID(ctx)
	utils.PanicOnError("failed to load chainId", err)

	acc, err := utils.PrepareAccount(&args.Sender, chainId, isSimulation)
	utils.PanicOnError("failed to parse private key", err)

	curBlockNumber, err := eth.BlockNumber(ctx)
	utils.PanicOnError("failed to load current block number", err)

	pod, err := EigenPod.NewEigenPod(common.HexToAddress(args.EigenPod), eth)
	utils.PanicOnError("failed to reach eigenpod", err)

	reg, err := pod.WithdrawableRestakedExecutionLayerGwei(nil)
	utils.PanicOnError("failed to fetch REG", err)
	rew := utils.GweiToWei(new(big.Float).SetUint64(reg))

	podOwner, err := pod.PodOwner(nil)
	utils.PanicOnError("failed to read podOwner", err)

	delegationManager, err := IDelegationManager.NewIDelegationManager(DelegationManager(chainId), eth)
	utils.PanicOnError("failed to reach delegation manager", err)

	minDelay, err := delegationManager.MinWithdrawalDelayBlocks(nil)
	utils.PanicOnError("failed to read MinWithdrawalDelayBlocks", err)

	queuedWithdrawals, err := delegationManager.GetQueuedWithdrawals(nil, podOwner)
	utils.PanicOnError("failed to read queuedWithdrawals", err)

	eligibleWithdrawals := lo.Map(queuedWithdrawals.Withdrawals, func(withdrawal IDelegationManager.IDelegationManagerTypesWithdrawal, index int) *IDelegationManager.IDelegationManagerTypesWithdrawal {
		isBeaconWithdrawal := len(withdrawal.Strategies) == 1 && withdrawal.Strategies[0].Cmp(core.BeaconStrategy()) == 0
		isCompletable := curBlockNumber > uint64(withdrawal.StartBlock+minDelay)
		if isBeaconWithdrawal && isCompletable {
			return &withdrawal
		}
		return nil
	})

	readyWithdrawalCount := len(lo.Filter(eligibleWithdrawals, func(withdrawal *IDelegationManager.IDelegationManagerTypesWithdrawal, index int) bool {
		return withdrawal != nil
	}))

	if readyWithdrawalCount == 0 {
		fmt.Printf("Your pod has no eligible withdrawals.\n")
		return nil
	}

	var runningSumWei *big.Float = new(big.Float)

	affordedWithdrawals := lo.Map(eligibleWithdrawals, func(withdrawal *IDelegationManager.IDelegationManagerTypesWithdrawal, index int) *IDelegationManager.IDelegationManagerTypesWithdrawal {
		if withdrawal == nil {
			return nil
		}
		withdrawalShares := queuedWithdrawals.Shares[index][0]
		// if REW > runningSumWei + withdrawalShares, we can complete with withdrawal.
		if rew.Cmp(
			new(big.Float).Add(
				runningSumWei,
				new(big.Float).SetInt(withdrawalShares),
			),
		) > 0 {
			runningSumWei = new(big.Float).Add(runningSumWei, new(big.Float).SetInt(withdrawalShares))
			return withdrawal
		}
		return nil
	})

	affordedWithdrawals = lo.Filter(affordedWithdrawals, func(withdrawal *IDelegationManager.IDelegationManagerTypesWithdrawal, index int) bool {
		return withdrawal != nil
	})

	if len(affordedWithdrawals) == 0 && readyWithdrawalCount > 0 {
		color.Yellow("WARN: Your pod has %d withdrawal(s) available, but your pod does not have enough funding to proceed.\n", readyWithdrawalCount)
		color.Yellow("Consider checkpointing to claim beacon rewards, or depositing ETH and checkpointing to complete these withdrawals.\n\n")
		return errors.New("Insufficient funds")
	}

	if len(affordedWithdrawals) != readyWithdrawalCount {
		color.Yellow("WARN: Your pod has %d withdrawal(s) available, but you only have enough balance to satisfy %d of them.\n", readyWithdrawalCount, len(affordedWithdrawals))
		color.Yellow("Consider checkpointing to claim beacon rewards, or depositing ETH and checkpointing to complete these withdrawals.\n\n")
	}

	fmt.Printf("Your podOwner(%s) has %d withdrawal(s) that can be completed right now.\n", podOwner.Hex(), len(affordedWithdrawals))
	runningSumWeiInt, _ := runningSumWei.Int(nil)
	fmt.Printf("Total ETH on all withdrawals: %sETH\n", utils.GweiToEther(utils.WeiToGwei(runningSumWeiInt)).String())

	if !isSimulation {
		utils.PanicIfNoConsent("Would you like to continue?")
	} else {
		color.Yellow("THIS IS A SIMULATION. No transaction will be recorded onchain.\n")
	}

	withdrawals := lo.Map(affordedWithdrawals, func(w *IDelegationManager.IDelegationManagerTypesWithdrawal, i int) IDelegationManager.IDelegationManagerTypesWithdrawal {
		return *w
	})

	tokens := lo.Map(withdrawals, func(_ IDelegationManager.IDelegationManagerTypesWithdrawal, _ int) []common.Address {
		return []common.Address{common.BigToAddress(big.NewInt(0))}
	})

	receiveAsTokens := lo.Map(withdrawals, func(_ IDelegationManager.IDelegationManagerTypesWithdrawal, _ int) bool {
		return true
	})

	txn, err := delegationManager.CompleteQueuedWithdrawals(acc.TransactionOptions, withdrawals, tokens, receiveAsTokens)
	utils.PanicOnError("CompleteQueuedWithdrawals failed.", err)

	if !isSimulation {
		_, err := bind.WaitMined(ctx, eth, txn)
		utils.PanicOnError("waitMined failed", err)

		color.Green("%s\n", txn.Hash().Hex())
	} else {
		PrintAsJSON(Transaction{
			Type:     "complete-withdrawals",
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

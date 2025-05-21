package utils

import (
	"fmt"
	"math/big"
)

func plural(word string, amount int) string {
	if amount > 1 {
		return word + "s"
	}
	return word
}

func StartCheckpointProofConsent() string {
	return `	This will start a new checkpoint on your eigenpod.

	Once started, you MUST complete this checkpoint. To see the full impact and transaction requirement of the checkpoint, rerun with 'status'.

	Note that if you lose the generated proofs, you'll need to recompute them. Beacon state is large, and many nodes do not retain
	long amounts of state. You can always recompute the proofs against a full archival beacon node, but you may face issues if 
	your beacon node is not archival.

	You should be comfortable with submitting all of these before beginning.
	
	PLAN: This will call EigenPod.VerifyCheckpointProofs(), with batches of proofs, to complete your checkpoint. For full details, run status.`
}

func SubmitCredentialsProofConsent(numTransactions int) string {
	return fmt.Sprintf(`	This will verify the withdrawal credentials of your validator, "restaking" your validator for the first time.
	Once submitted, future checkpoint proofs will include a balance proof against this validator. 

	PLAN: This will call EigenPod.VerifyWithdrawalCredentials() %d %s to link your %s to your eigenpod.
	`,
		numTransactions,
		plural("time", numTransactions),
		plural("validator", numTransactions),
	)
}

func SubmitSwitchRequestConsent(
	numRequests int,
	currentQueueSize *big.Int,
	currentFee string,
	totalFee string,
	totalOverestimateFee string,
	isSimulatedStr string,
) string {
	return fmt.Sprintf(`	%s This will submit switch requests to your pod, updating your validators' withdrawal credential prefixes.
	Moving from 0x01 to 0x02 withdrawal prefix can NOT be undone. 0x02 validators:
	 - Can be the 'target' of a consolidation request
	 - Have a max effective balance of 2048 ETH
	 - Will NOT be subject to withdrawal sweeps unless their balance exceeds 2048 ETH
	
	If you plan to perform partial/full withdrawals for an 0x02 validator, you will need to use EigenPod.requestWithdrawal().
	
	%s PLAN: This will call EigenPod.requestConsolidation to request prefix switches for %d validators.
	 - The EIP-7521 predeploy requires a fee sent as msg.value, depending on the number of requests in the consolidation queue.
	 - The current queue size is %d, making the current fee for each request %s.
	 - Not including gas, the total fee for your requests is currently %s. 
	 - With current overestimate settings, you will send %s along with this transaction.
	
	(Unused funds will be sent back to the caller.)

	`,
		isSimulatedStr,
		isSimulatedStr,
		numRequests,
		currentQueueSize,
		currentFee,
		totalFee,
		totalOverestimateFee,
	)
}

func SubmitSourceToTargetRequestConsent(
	numRequests int,
	currentQueueSize *big.Int,
	currentFee string,
	totalFee string,
	totalOverestimateFee string,
	targetValidator uint64,
	numSourceValidators int,
	isSimulatedStr string,
) string {
	return fmt.Sprintf(`	%s This will submit source-to-target consolidation requests to your pod, consolidating one or more source
	validators to the specified target validator. Once these requests are processed on the beacon chain, each of your source
	validators will be considered "exited", and their balances will move to the target validator. This is irreversible.

	Note that the beacon chain may reject consolidations if:
	 - the target does NOT have an 0x02 withdrawal prefix
	 - the source/target has been slashed, has initiated exit, or has been consolidated
	 - the source validator has used EigenPod.requestWithdrawal and still has a pending withdrawal request

	Before submitting this request, please make sure that none of these apply to you!
	
	%s PLAN: This will call EigenPod.requestConsolidation to submit %d source-to-target consolidation requests.
	 - The EIP-7521 predeploy requires a fee sent as msg.value, depending on the number of requests in the consolidation queue.
	 - The current queue size is %d, making the current fee for each request %s.
	 - Not including gas, the total fee for your requests is currently %s. 
	 - With current overestimate settings, you can expect to send up to %s. 
	 - Selected target validator: %d
	 - Number of source validators: %d
	
	(Unused funds will be sent back to the caller.)

	`,
		isSimulatedStr,
		isSimulatedStr,
		numRequests,
		currentQueueSize,
		currentFee,
		totalFee,
		totalOverestimateFee,
		targetValidator,
		numSourceValidators,
	)
}

package core

import (
	"fmt"
)

func plural(word string, amount int) string {
	if amount > 1 {
		return word + "s"
	}
	return word
}

func SubmitCheckpointProofConsent(numTransactions int) string {
	return fmt.Sprintf(`This will start a new checkpoint on your eigenpod.

	Once started, you MUST complete this checkpoint. This may entail submitting %d or more transactions.

	Note that if you lose the generated proofs, you'll need to recompute them. Beacon state is large, and many nodes do not retain
	long amounts of state. You can always recompute the proofs against a full archival beacon node, but you may face issues if 
	your beacon node is not archival.

	You should be comfortable with submitting all of these before beginning.
	
	PLAN: This will call EigenPod.VerifyCheckpointProofs() %d %s to complete your checkpoint.`,
		numTransactions,
		numTransactions,
		plural("time", numTransactions),
	)
}

func SubmitCredentialsProofConsent(numTransactions int) string {
	return fmt.Sprintf(`This will verify the withdrawal credentials of your validator, "restaking" your validator for the first time.
	Once submitted, future checkpoint proofs will include a balance proof against this validator. 

	PLAN: This will call EigenPod.VerifyWithdrawalCredentials() %d %s to link your %s to your eigenpod.
	`,
		numTransactions,
		plural("time", numTransactions),
		plural("validator", numTransactions),
	)
}

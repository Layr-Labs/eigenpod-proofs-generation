package generation

import (
	"flag"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// this needs to be hand crafted. If you want the root of the header at the slot x,
// then look for entry in (x)%slotsPerHistoricalRoot in the block_roots.

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Defining flags for all the parameters
	command := flag.String("command", "", "The command to execute")

	oracleBlockHeaderFile := flag.String("oracleBlockHeaderFile", "", "Oracle block header file")
	stateFile := flag.String("stateFile", "", "State file")
	validatorIndex := flag.Uint64("validatorIndex", 0, "validatorIndex")
	outputFile := flag.String("outputFile", "", "Output file")
	chainID := flag.Uint64("chainID", 0, "Chain ID")

	//WithdrawaProof specific flags
	historicalSummariesIndex := flag.Uint64("historicalSummariesIndex", 0, "Historical summaries index")
	blockHeaderIndex := flag.Uint64("blockHeaderIndex", 0, "Block header index")
	historicalSummaryStateFile := flag.String("historicalSummaryStateFile", "", "Historical summary state file")
	blockHeaderFile := flag.String("blockHeaderFile", "", "Block Header file")
	blockBodyFile := flag.String("blockBodyFile", "", "Block Body file")
	withdrawalIndex := flag.Uint64("withdrawalIndex", 0, "Withdrawal index")

	// Parse the flags
	flag.Parse()

	// Check if the required 'command' flag is provided
	if *command == "" {
		log.Debug().Msg("Error: command flag is required")
		return
	}

	var err error
	// Handling commands based on the 'command' flag
	switch *command {
	case "ValidatorFieldsProof":
		err = GenerateValidatorFieldsProof(*oracleBlockHeaderFile, *stateFile, *validatorIndex, *chainID, *outputFile)

	case "WithdrawalFieldsProof":
		err = GenerateWithdrawalFieldsProof(*oracleBlockHeaderFile, *stateFile, *historicalSummaryStateFile, *blockHeaderFile, *blockBodyFile, *validatorIndex, *withdrawalIndex, *historicalSummariesIndex, *blockHeaderIndex, *chainID, *outputFile)

	case "BalanceUpdateProof":
		err = GenerateBalanceUpdateProof(*oracleBlockHeaderFile, *stateFile, *validatorIndex, *chainID, *outputFile)

	default:
		log.Debug().Str("Unknown command:", *command)
	}
	log.Debug().AnErr("Error: ", err)
}

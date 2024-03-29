package main

import (
	"flag"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	txsubmitter "github.com/Layr-Labs/eigenpod-proofs-generation/tx_submitoor/tx_submitter"
	"github.com/caarlos0/env"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	RPC          string `json:"RPC_URL"`
	PrivateKey   string `json:"PRIVATE_KEY,required"`
	ChainID      uint64 `json:"CHAIN_ID,required"`
	CacheExpire  int    `json:"CACHE_EXPIRE" envDefault:"1000"`
	BeaconAPIURL string `json:"BEACON_API_URL,required"`
}

func main() {

	cfg := parseConfig()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Defining flags for all the parameters
	command := flag.String("command", "", "The command to execute")
	withdrawalProofConfigFile := flag.String("withdrawalProofConfig", "", "Withdrawal proof config file")
	submitTransaction := flag.Bool("submitTransaction", false, "Submit transaction to the chain")

	ethClient, err := ethclient.Dial(cfg.RPC)
	if err != nil {
		log.Panic().Msgf("failed to connect to RPC: %s", err)
	}

	chainClient, err := txsubmitter.NewChainClient(ethClient, cfg.PrivateKey)
	if err != nil {
		log.Panic().Msgf("failed to create chain client: %s", err)
	}

	eigenPodProofs, err := eigenpodproofs.NewEigenPodProofs(cfg.ChainID, cfg.CacheExpire)
	if err != nil {
		log.Panic().Msgf("failed to create eigen pod proofs: %s", err)
	}

	submitter := txsubmitter.NewEigenPodProofTxSubmitter(
		*chainClient,
		*eigenPodProofs,
	)

	// Handling commands based on the 'command' flag
	switch *command {
	case "WithdrawalFieldsProof":
		submitter.SubmitVerifyAndProcessWithdrawalsTx(*withdrawalProofConfigFile, submitTransaction)

	default:
		log.Debug().Str("Unknown command:", *command)
	}
	log.Debug().AnErr("Error: ", err)
}

func parseConfig() Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal().Msgf("failed to initialize eigenpod validator balance updater: %s", err)
	}

	return cfg
}

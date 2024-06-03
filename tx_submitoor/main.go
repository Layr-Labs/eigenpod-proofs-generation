package main

import (
	"encoding/hex"
	"flag"
	"time"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	txsubmitter "github.com/Layr-Labs/eigenpod-proofs-generation/tx_submitoor/tx_submitter"
	"github.com/caarlos0/env"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	FromAddress  string `env:"FROM_ADDRESS"`
	GasPrice     int64  `env:"GAS_PRICE" envDefault:"50"`
	GasLimit     uint64 `env:"GAS_LIMIT" envDefault:"100"`
	RPC          string `env:"RPC_URL"`
	PrivateKey   string `env:"PRIVATE_KEY"`
	ChainID      uint64 `env:"CHAIN_ID,required"`
	CacheExpire  int    `env:"CACHE_EXPIRE" envDefault:"1000"`
	BeaconAPIURL string `env:"BEACON_API_URL"`
}

func main() {

	cfg := parseConfig()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Defining flags for all the parameters
	command := flag.String("command", "", "The command to execute")
	withdrawalProofConfigFile := flag.String("withdrawalProofConfig", "", "Withdrawal proof config file")
	submitTransaction := flag.Bool("submitTransaction", false, "Submit transaction to the chain")

	flag.Parse()

	ethClient, err := ethclient.Dial(cfg.RPC)
	if err != nil {
		log.Panic().Msgf("failed to connect to RPC: %s", err)
	}

	chainClient, err := txsubmitter.NewChainClient(ethClient, cfg.PrivateKey, cfg.FromAddress, cfg.GasPrice, cfg.GasLimit)
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

	startedAt := time.Now()

	// Handling commands based on the 'command' flag
	switch *command {
	case "WithdrawalFieldsProof":
		calldata, err := submitter.SubmitVerifyAndProcessWithdrawalsTx(*withdrawalProofConfigFile, *submitTransaction)
		if err != nil {
			log.Panic().Msgf("failed to submit withdrawal proof: %s", err)
		}
		log.Info().Msgf("Withdrawal proof submitted with calldata: %s", hex.EncodeToString(calldata))

	case "WithdrawalCredentialProof":
		calldata, err := submitter.SubmitVerifyWithdrawalCredentialsTx(*withdrawalProofConfigFile, *submitTransaction)
		if err != nil {
			log.Panic().Msgf("failed to submit withdrawal credential proof: %s", err)
		}
		log.Info().Msgf("Withdrawal credential proof submitted with calldata: %s", hex.EncodeToString(calldata))

	default:
		log.Debug().Str("Unknown command:", *command)
	}
	log.Debug().AnErr("Error: ", err)
	log.Debug().Msgf("Took %f seconds", time.Since(startedAt).Seconds())

}

func parseConfig() Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal().Msgf("failed to initialize eigenpod validator balance updater: %s", err)
	}

	return cfg
}

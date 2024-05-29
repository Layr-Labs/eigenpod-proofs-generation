#!/bin/bash

cd .. && ./generation/generation \
  -command ValidatorFieldsProof \
  -oracleBlockHeaderFile "./automationScript/HEAD_FILE.json" \
  -stateFile "./automationScript/STATE_FILE.json" \
  -validatorIndex 1373594 \
  -outputFile "./automationScript/withdrawal_credential_proof.json" \
  -chainID 1
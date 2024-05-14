The current deployed test contract is at [this address](https://goerli.etherscan.io/address/0xd132dD701d3980bb5d66A21e2340f263765e4a19)

This contract takes BeaconChainProofs.sol, a library, converts it to a contract, and makes all the proof functions public.  

If needed, to generate the binding, retrieve the abi from [etherscan](https://goerli.etherscan.io/address/0xd132dD701d3980bb5d66A21e2340f263765e4a19#code) and run the following command:
```
abigen --abi build/BeaconChainProofs.abi --pkg main --type BeaconChainProofs --out BeaconChainProofs.go
```

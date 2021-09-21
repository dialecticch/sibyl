# sibyl

A more embedded version of [fxfactorial/run-evm-code](https://github.com/fxfactorial/run-evm-code/). This tool makes it 
easy to apply transactions to the current EVM state. Call it a transaction simulator or what not.

sibyl should be safe to run against your node as it does not actually commit state to disk.

## Usage

```go
// Create a new simulator using geth chaindata.
simulator, err := sybil.NewGethSimulator("geth/chaindata")
if err != nil {
    log.Panic(err)
}

// Fork to a specified block number
err = simulator.Fork(blockNumber)
if err != nil {
    log.Panic(err)
}

// Simulate a static call.
ret, err := simulator.StaticCall(sender, to, input, gas)
if err != nil {
	log.Panic(err)
}

fmt.Println(hexutil.Encode(result))
```

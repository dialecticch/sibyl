package sibyl_test

import (
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/dialecticch/sibyl"
	"github.com/dialecticch/sibyl/testdata"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestSimulator(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	backend := newBackend(privateKey)

	transactOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	if err != nil {
		t.Fatal(err)
	}

	addr, _, _, err := testdata.DeployCounter(transactOpts, backend)
	if err != nil {
		t.Fatal(err)
	}

	backend.Commit()

	s := sibyl.NewSimulator(backend.Blockchain())
}

func newBackend(key *ecdsa.PrivateKey) *backends.SimulatedBackend {
	opts, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))

	balance := new(big.Int)
	balance.SetString("10000000000000000000", 10) // 10 eth in wei

	address := opts.From
	genesisAlloc := map[common.Address]core.GenesisAccount{
		address: {
			Balance: balance,
		},
	}

	blockGasLimit := uint64(10000000)
	return backends.NewSimulatedBackend(genesisAlloc, blockGasLimit)
}

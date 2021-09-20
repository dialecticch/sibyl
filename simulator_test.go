package sibyl_test

import (
	"crypto/ecdsa"
	"math/big"
	"strings"
	"testing"

	"github.com/dialecticch/sibyl"
	"github.com/dialecticch/sibyl/testdata"
	"github.com/ethereum/go-ethereum/accounts/abi"
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

	pub := crypto.PubkeyToAddress(*privateKey.Public().(*ecdsa.PublicKey))

	backend := newBackend(privateKey)

	transactOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	if err != nil {
		t.Fatal(err)
	}

	addr, _, c, err := testdata.DeployCounter(transactOpts, backend)
	if err != nil {
		t.Fatal(err)
	}

	backend.Commit()

	s := sibyl.NewSimulator(backend.Blockchain())

	backend.Close()

	err = s.Fork(backend.Blockchain().CurrentBlock().NumberU64())
	if err != nil {
		t.Fatalf("s.Fork err: %s", err)
	}

	for i := 0; i < 10; i++ {
		opts, _ := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
		opts.NoSend = true

		tx, err := c.Tick(opts)

		_, err = s.Call(pub, addr, tx.Data(), tx.Gas(), common.Big0)
		if err != nil {
			t.Fatal(err)
		}

		cabi, err := abi.JSON(strings.NewReader(testdata.CounterMetaData.ABI))
		if err != nil {
			t.Fatal(err)
		}

		ret, err := s.StaticCall(pub, addr, cabi.Methods["count"].ID, 30_000_000)
		if err != nil {
			t.Fatal(err)
		}

		pack, err := cabi.Unpack("count", ret)
		if err != nil {
			t.Fatal(err)
		}

		parsed := *abi.ConvertType(pack[0], new(*big.Int)).(**big.Int)
		if parsed.Uint64() != uint64(i + 1) {
			t.Fatal("count did not work")
		}
	}
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

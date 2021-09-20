package sibyl

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/params"
)

// NewGethSimulator returns a new simulator that utilizes the go-ethereum LevelDB Database.
func NewGethSimulator(path string) (*Simulator, error) {
	chainDb, err := rawdb.NewLevelDBDatabase(path, 48, 48, "", true)

	if err != nil {
		return nil, err
	}

	currentHead := rawdb.ReadHeadBlockHash(chainDb)

	if currentHead == (common.Hash{}) {
		return nil, errors.New("current head is nil, check db")
	}

	cfg := vm.Config{}

	engine := ethash.New(ethash.Config{
		CachesInMem:      ethconfig.Defaults.Ethash.CachesInMem,
		CachesOnDisk:     ethconfig.Defaults.Ethash.CachesOnDisk,
		CachesLockMmap:   ethconfig.Defaults.Ethash.CachesLockMmap,
		DatasetsInMem:    ethconfig.Defaults.Ethash.DatasetsInMem,
		DatasetsOnDisk:   ethconfig.Defaults.Ethash.DatasetsOnDisk,
		DatasetsLockMmap: ethconfig.Defaults.Ethash.DatasetsLockMmap,
	}, nil, false)

	cache := &core.CacheConfig{
		TrieCleanLimit: ethconfig.Defaults.TrieCleanCache,
		TrieDirtyLimit: ethconfig.Defaults.TrieDirtyCache,
		TrieTimeLimit:  ethconfig.Defaults.TrieTimeout,
		SnapshotLimit:  0,
	}

	chain, err := core.NewBlockChain(
		chainDb, cache, params.MainnetChainConfig, engine, cfg, nil, nil,
	)

	if err != nil {
		return nil, err
	}

	return NewSimulator(chain), nil
}

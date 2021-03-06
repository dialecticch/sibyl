package sibyl

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/pkg/errors"
)

// Simulator is a transaction simulator, allowing you to run transactions against the current blockchain state
type Simulator struct {
	mux sync.Mutex

	db         *state.StateDB
	blockchain *core.BlockChain

	vm *vm.EVM
}

// NewSimulator returns a bare simulator
func NewSimulator(blockchain *core.BlockChain) *Simulator {
	return &Simulator{
		blockchain: blockchain,
	}
}

// Fork creates a new temporary context with the state for a given block number
func (s *Simulator) Fork(blockNumber uint64) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	header := s.blockchain.CurrentHeader()
	block := s.blockchain.GetBlockByNumber(blockNumber)
	db, err := s.blockchain.StateAt(block.Root())
	if err != nil {
		return errors.Wrap(err, "failed to read blockchain state for block")
	}

	s.db = db

	blockCtx := core.NewEVMBlockContext(header, s.blockchain, nil)
	txCtx := vm.TxContext{}

	s.vm = vm.NewEVM(blockCtx, txCtx, db, s.blockchain.Config(), *s.blockchain.GetVMConfig())

	return nil
}

// StaticCall executes an EVM static call on the current context
func (s *Simulator) StaticCall(sender, to common.Address, input []byte, gas uint64) ([]byte, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	ret, _, err := s.vm.StaticCall(vm.AccountRef(sender), to, input, gas)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// Call executes an EVM call on the current context
func (s *Simulator) Call(sender, to common.Address, input []byte, gas uint64, value *big.Int) ([]byte, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	ret, _, err := s.vm.Call(vm.AccountRef(sender), to, input, gas, value)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// Snapshot takes a snapshot of the current state
func (s *Simulator) Snapshot() int {
	return s.db.Snapshot()
}

// Rollback rolls the state back to a specific snapshot
func (s *Simulator) Rollback(snapshot int) {
	s.db.RevertToSnapshot(snapshot)
}

package sibyl

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/pkg/errors"
)

// Simulator is a transaction simulator, allowing you to run transactions against the current blockchain state.
type Simulator struct {
	mux sync.Mutex

	blockchain *core.BlockChain

	vm *vm.EVM
}

func NewSimulator(blockchain *core.BlockChain) *Simulator {
	return &Simulator{
		blockchain: blockchain,
	}
}

func (s *Simulator) Fork(blockNumber uint64) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	header := s.blockchain.CurrentHeader()
	block := s.blockchain.GetBlockByNumber(blockNumber)
	db, err := s.blockchain.StateAt(block.Root())
	if err != nil {
		return errors.Wrap(err, "failed to read blockchain state for block")
	}

	blockCtx := core.NewEVMBlockContext(header, s.blockchain, nil)
	txCtx := vm.TxContext{}

	s.vm = vm.NewEVM(blockCtx, txCtx, db, s.blockchain.Config(), *s.blockchain.GetVMConfig())

	return nil
}

func (s *Simulator) StaticCall(sender, to common.Address, input []byte, gas uint64) ([]byte, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	// @TODO
	s.vm.Reset(vm.TxContext{Origin: sender, GasPrice: big.NewInt(0)}, s.vm.StateDB)

	ret, _, err := s.vm.StaticCall(vm.AccountRef(sender), to, input, gas)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (s *Simulator) Call(sender, to common.Address, input []byte, gas uint64, value *big.Int) ([]byte, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.vm.Reset(vm.TxContext{Origin: sender, GasPrice: big.NewInt(0)}, s.vm.StateDB)

	ret, _, err := s.vm.Call(vm.AccountRef(sender), to, input, gas, value)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

package transactionManager

import (
	. "../alias"
	"../blockchain"
	"../utils"
	"bytes"
	"fmt"
	"sync"
)

type TransactionManager struct {
	utxoIndex        map[string]*blockchain.Input // map key is "BlockNumber:TransactionNumber:OutputNumber"
	transactionQueue []*blockchain.Transaction
	lastBlock        uint32
	lastHash         Uint256
	lastAccumulator  Uint2048
	mutex            sync.Mutex
}

func NewTransactionManager() *TransactionManager {
	result := TransactionManager{
		utxoIndex:        map[string]*blockchain.Input{},
		transactionQueue: make([]*blockchain.Transaction, 0),
		lastBlock:        0,
		lastHash:         utils.Keccak256([]byte{}), // todo define genesis hash
		lastAccumulator:  []byte{3},                 //todo define genesis accumulator
		mutex:            sync.Mutex{},
	}
	return &result
}

// ValidateInputs checks that all inputs correspond to correct unspent outputs
func (p *TransactionManager) ValidateInputs(t *blockchain.Transaction) error {
	for _, in := range t.Inputs {
		utxo := p.utxoIndex[in.GetKey()]
		if utxo == nil {
			return fmt.Errorf("no such UTXO: %s", in.GetKey())
		}
		// todo deep equal instead of explicit, or make this resistant to modification of utxo model like adding assetId
		if bytes.Compare(utxo.Owner, in.Owner) != 0 || utxo.Slice.Begin != in.Slice.Begin || utxo.Slice.End != in.Slice.End {
			return fmt.Errorf("incorrect input data for UTXO: %s", in.GetKey())
		}
		if in.BlockIndex > p.lastBlock {
			return fmt.Errorf("multiple operations on a slice within the same block are forbidden")
		}
	}

	return nil
}

func (p *TransactionManager) SubmitTransaction(t *blockchain.Transaction) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// check that transaction is fully valid
	err := t.Validate()
	if err != nil {
		return err
	}
	err = p.ValidateInputs(t)
	if err != nil {
		return err
	}

	// spend inputs, add outputs to utxo index, queue transaction for the next block
	for _, in := range t.Inputs {
		delete(p.utxoIndex, in.GetKey())
	}
	for i, out := range t.Outputs {
		in := blockchain.Input{
			Output: blockchain.Output{
				Owner: out.Owner,
				Slice: out.Slice,
			},
			BlockIndex:  p.lastBlock + 1,
			TxIndex:     uint32(len(p.transactionQueue)),
			OutputIndex: uint8(i),
		}
		p.utxoIndex[in.GetKey()] = &in
	}
	p.transactionQueue = append(p.transactionQueue, t)

	return nil
}

// todo remove this
func dereference(t []*blockchain.Transaction) []blockchain.Transaction {
	result := make([]blockchain.Transaction, 0, len(t))
	for _, v := range t {
		result = append(result, *v)
	}
	return result
}

func (p *TransactionManager) AssembleBlock() (*blockchain.Block, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	block, err := blockchain.NewBlock(p.lastBlock+1, p.lastHash, p.lastAccumulator, dereference(p.transactionQueue))
	if err != nil {
		return nil, err
	}
	p.lastBlock++
	p.lastHash = block.GetHash()
	p.lastAccumulator = block.RSAAccumulator
	p.transactionQueue = make([]*blockchain.Transaction, 0)
	return block, nil
}

// todo add utxo on deposit event, avoid double deposits
// todo spend utxo on withdraw event
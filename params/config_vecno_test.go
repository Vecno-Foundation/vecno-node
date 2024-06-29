package params

import (
	"testing"
)

// TestGenesisHashVecno tests that VecnoGenesisHash is the correct value for the genesis configuration.
func TestGenesisHashVecno(t *testing.T) {
	genesis := DefaultVecnoGenesisBlock()
	block := genesisToBlock(genesis, nil)
	if block.Hash() != VecnoGenesisHash {
		t.Errorf("want: %s, got: %s", VecnoGenesisHash.Hex(), block.Hash().Hex())
	}
}

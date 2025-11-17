package batch

import (
	"testing"

	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

func TestLinkedChain(t *testing.T) {
	chain := NewLinkedChain()

	transfer1 := types.Transfer{ID: types.ToUint128(1)}
	transfer2 := types.Transfer{ID: types.ToUint128(2)}
	transfer3 := types.Transfer{ID: types.ToUint128(3)}

	chain.AddTransfer(transfer1).
		AddTransfer(transfer2).
		AddTransfer(transfer3)

	if chain.Len() != 3 {
		t.Errorf("Expected chain length 3, got %d", chain.Len())
	}

	transfers := chain.Build()

	// First two should be linked
	if !transfers[0].Flags.Linked {
		t.Error("Expected first transfer to be linked")
	}
	if !transfers[1].Flags.Linked {
		t.Error("Expected second transfer to be linked")
	}

	// Last should not be linked
	if transfers[2].Flags.Linked {
		t.Error("Expected last transfer to not be linked")
	}
}

func TestLinkedChainEmpty(t *testing.T) {
	chain := NewLinkedChain()

	if chain.Len() != 0 {
		t.Errorf("Expected empty chain length 0, got %d", chain.Len())
	}

	transfers := chain.Build()
	if len(transfers) != 0 {
		t.Errorf("Expected empty transfers slice, got length %d", len(transfers))
	}
}

func TestMaxBatchSize(t *testing.T) {
	if MaxBatchSize != 8189 {
		t.Errorf("Expected MaxBatchSize 8189, got %d", MaxBatchSize)
	}
}

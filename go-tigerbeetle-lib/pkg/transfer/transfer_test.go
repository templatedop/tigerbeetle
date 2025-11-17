package transfer

import (
	"testing"
	"time"

	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

func TestTransferBuilder(t *testing.T) {
	id := types.ToUint128(100)
	debitID := types.ToUint128(1)
	creditID := types.ToUint128(2)
	amount := types.ToUint128(1000)

	transfer := New(id).
		DebitAccount(debitID).
		CreditAccount(creditID).
		Amount(amount).
		Ledger(1).
		Code(10).
		Build()

	if transfer.ID.ToUint128() != 100 {
		t.Errorf("Expected ID 100, got %d", transfer.ID.ToUint128())
	}

	if transfer.DebitAccountID.ToUint128() != 1 {
		t.Errorf("Expected DebitAccountID 1, got %d", transfer.DebitAccountID.ToUint128())
	}

	if transfer.CreditAccountID.ToUint128() != 2 {
		t.Errorf("Expected CreditAccountID 2, got %d", transfer.CreditAccountID.ToUint128())
	}

	if transfer.Amount.ToUint128() != 1000 {
		t.Errorf("Expected Amount 1000, got %d", transfer.Amount.ToUint128())
	}

	if transfer.Ledger != 1 {
		t.Errorf("Expected Ledger 1, got %d", transfer.Ledger)
	}

	if transfer.Code != 10 {
		t.Errorf("Expected Code 10, got %d", transfer.Code)
	}
}

func TestTransferBuilderPending(t *testing.T) {
	timeout := 1 * time.Hour

	transfer := New(types.ToUint128(1)).
		DebitAccount(types.ToUint128(1)).
		CreditAccount(types.ToUint128(2)).
		Amount(types.ToUint128(500)).
		Ledger(1).
		Code(1).
		Pending().
		Timeout(timeout).
		Build()

	if !transfer.Flags.Pending {
		t.Error("Expected Pending flag to be set")
	}

	expectedTimeout := uint64(timeout.Nanoseconds())
	if transfer.Timeout != expectedTimeout {
		t.Errorf("Expected Timeout %d, got %d", expectedTimeout, transfer.Timeout)
	}
}

func TestTransferBuilderLinked(t *testing.T) {
	transfer := New(types.ToUint128(1)).
		Linked().
		Build()

	if !transfer.Flags.Linked {
		t.Error("Expected Linked flag to be set")
	}
}

func TestTransferBuilderPostPending(t *testing.T) {
	pendingID := types.ToUint128(999)

	transfer := New(types.ToUint128(1)).
		PostPending(pendingID).
		Build()

	if !transfer.Flags.PostPendingTransfer {
		t.Error("Expected PostPendingTransfer flag to be set")
	}

	if transfer.PendingID.ToUint128() != 999 {
		t.Errorf("Expected PendingID 999, got %d", transfer.PendingID.ToUint128())
	}
}

func TestTransferBuilderVoidPending(t *testing.T) {
	pendingID := types.ToUint128(888)

	transfer := New(types.ToUint128(1)).
		VoidPending(pendingID).
		Build()

	if !transfer.Flags.VoidPendingTransfer {
		t.Error("Expected VoidPendingTransfer flag to be set")
	}

	if transfer.PendingID.ToUint128() != 888 {
		t.Errorf("Expected PendingID 888, got %d", transfer.PendingID.ToUint128())
	}
}

func TestTransferError(t *testing.T) {
	err := &TransferError{
		Index:  3,
		Result: types.CreateTransferResult(1),
	}

	errorMsg := err.Error()
	if errorMsg == "" {
		t.Error("Expected non-empty error message")
	}
}

func TestAmountUint64(t *testing.T) {
	transfer := New(types.ToUint128(1)).
		AmountUint64(12345).
		Build()

	if transfer.Amount.ToUint128() != 12345 {
		t.Errorf("Expected Amount 12345, got %d", transfer.Amount.ToUint128())
	}
}

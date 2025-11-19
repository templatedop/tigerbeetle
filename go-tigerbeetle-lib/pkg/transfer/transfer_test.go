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

	if transfer.ID != id {
		t.Errorf("Expected ID %v, got %v", id, transfer.ID)
	}

	if transfer.DebitAccountID != debitID {
		t.Errorf("Expected DebitAccountID %v, got %v", debitID, transfer.DebitAccountID)
	}

	if transfer.CreditAccountID != creditID {
		t.Errorf("Expected CreditAccountID %v, got %v", creditID, transfer.CreditAccountID)
	}

	if transfer.Amount != amount {
		t.Errorf("Expected Amount %v, got %v", amount, transfer.Amount)
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

	pendingFlag := types.TransferFlags{Pending: true}.ToUint16()
	if (transfer.Flags & pendingFlag) == 0 {
		t.Error("Expected Pending flag to be set")
	}

	expectedTimeout := uint32(timeout.Nanoseconds())
	if transfer.Timeout != expectedTimeout {
		t.Errorf("Expected Timeout %d, got %d", expectedTimeout, transfer.Timeout)
	}
}

func TestTransferBuilderLinked(t *testing.T) {
	transfer := New(types.ToUint128(1)).
		Linked().
		Build()

	linkedFlag := types.TransferFlags{Linked: true}.ToUint16()
	if (transfer.Flags & linkedFlag) == 0 {
		t.Error("Expected Linked flag to be set")
	}
}

func TestTransferBuilderPostPending(t *testing.T) {
	pendingID := types.ToUint128(999)

	transfer := New(types.ToUint128(1)).
		PostPending(pendingID).
		Build()

	postPendingFlag := types.TransferFlags{PostPendingTransfer: true}.ToUint16()
	if (transfer.Flags & postPendingFlag) == 0 {
		t.Error("Expected PostPendingTransfer flag to be set")
	}

	if transfer.PendingID != pendingID {
		t.Errorf("Expected PendingID %v, got %v", pendingID, transfer.PendingID)
	}
}

func TestTransferBuilderVoidPending(t *testing.T) {
	pendingID := types.ToUint128(888)

	transfer := New(types.ToUint128(1)).
		VoidPending(pendingID).
		Build()

	voidPendingFlag := types.TransferFlags{VoidPendingTransfer: true}.ToUint16()
	if (transfer.Flags & voidPendingFlag) == 0 {
		t.Error("Expected VoidPendingTransfer flag to be set")
	}

	if transfer.PendingID != pendingID {
		t.Errorf("Expected PendingID %v, got %v", pendingID, transfer.PendingID)
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
	expectedAmount := types.ToUint128(12345)

	transfer := New(types.ToUint128(1)).
		AmountUint64(12345).
		Build()

	if transfer.Amount != expectedAmount {
		t.Errorf("Expected Amount %v, got %v", expectedAmount, transfer.Amount)
	}
}

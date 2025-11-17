package account

import (
	"testing"

	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

func TestAccountBuilder(t *testing.T) {
	id := types.ToUint128(123)

	account := New(id).
		Ledger(1).
		Code(100).
		UserData64(456).
		HistoryEnabled().
		Build()

	if account.ID.ToUint128() != 123 {
		t.Errorf("Expected ID 123, got %d", account.ID.ToUint128())
	}

	if account.Ledger != 1 {
		t.Errorf("Expected Ledger 1, got %d", account.Ledger)
	}

	if account.Code != 100 {
		t.Errorf("Expected Code 100, got %d", account.Code)
	}

	if account.UserData64 != 456 {
		t.Errorf("Expected UserData64 456, got %d", account.UserData64)
	}

	if !account.Flags.History {
		t.Error("Expected History flag to be set")
	}
}

func TestAccountBuilderLinked(t *testing.T) {
	account := New(types.ToUint128(1)).
		LinkedAccount().
		Build()

	if !account.Flags.Linked {
		t.Error("Expected Linked flag to be set")
	}
}

func TestAccountBuilderDebitsMustNotExceedCredits(t *testing.T) {
	account := New(types.ToUint128(1)).
		DebitsMustNotExceedCredits().
		Build()

	if !account.Flags.DebitsMustNotExceedCredits {
		t.Error("Expected DebitsMustNotExceedCredits flag to be set")
	}
}

func TestAccountBuilderCreditsMustNotExceedDebits(t *testing.T) {
	account := New(types.ToUint128(1)).
		CreditsMustNotExceedDebits().
		Build()

	if !account.Flags.CreditsMustNotExceedDebits {
		t.Error("Expected CreditsMustNotExceedDebits flag to be set")
	}
}

func TestAccountError(t *testing.T) {
	err := &AccountError{
		Index:  5,
		Result: types.CreateAccountResult(1),
	}

	errorMsg := err.Error()
	if errorMsg == "" {
		t.Error("Expected non-empty error message")
	}
}

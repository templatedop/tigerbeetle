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

	if account.ID != id {
		t.Errorf("Expected ID %v, got %v", id, account.ID)
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

	historyFlag := types.AccountFlags{History: true}.ToUint16()
	if (account.Flags & historyFlag) == 0 {
		t.Error("Expected History flag to be set")
	}
}

func TestAccountBuilderLinked(t *testing.T) {
	account := New(types.ToUint128(1)).
		LinkedAccount().
		Build()

	linkedFlag := types.AccountFlags{Linked: true}.ToUint16()
	if (account.Flags & linkedFlag) == 0 {
		t.Error("Expected Linked flag to be set")
	}
}

func TestAccountBuilderDebitsMustNotExceedCredits(t *testing.T) {
	account := New(types.ToUint128(1)).
		DebitsMustNotExceedCredits().
		Build()

	debitFlag := types.AccountFlags{DebitsMustNotExceedCredits: true}.ToUint16()
	if (account.Flags & debitFlag) == 0 {
		t.Error("Expected DebitsMustNotExceedCredits flag to be set")
	}
}

func TestAccountBuilderCreditsMustNotExceedDebits(t *testing.T) {
	account := New(types.ToUint128(1)).
		CreditsMustNotExceedDebits().
		Build()

	creditFlag := types.AccountFlags{CreditsMustNotExceedDebits: true}.ToUint16()
	if (account.Flags & creditFlag) == 0 {
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

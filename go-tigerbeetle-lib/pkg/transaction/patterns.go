package transaction

import (
	"fmt"
	"time"

	"gitlab.cept.gov.in/it-2.0-common/ledgers/pkg/transfer"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

// TransactionClient defines the interface for transaction operations
type TransactionClient interface {
	CreateTransfers([]types.Transfer) ([]types.CreateTransfersError, error)
	LookupTransfers([]types.Uint128) ([]types.Transfer, error)
}

// Patterns provides high-level transaction patterns
type Patterns struct {
	client TransactionClient
}

// NewPatterns creates a new transaction patterns helper
func NewPatterns(client TransactionClient) *Patterns {
	return &Patterns{client: client}
}

// EscrowParams contains parameters for an escrow transaction
type EscrowParams struct {
	EscrowID       types.Uint128
	FromAccount    types.Uint128
	EscrowAccount  types.Uint128
	Amount         types.Uint128
	Ledger         uint32
	Code           uint16
	EscrowTimeout  time.Duration
}

// Escrow creates a pending transfer to an escrow account
// This reserves funds from the sender until the escrow is released or cancelled
func (p *Patterns) Escrow(params EscrowParams) error {
	escrowTransfer := transfer.New(params.EscrowID).
		DebitAccount(params.FromAccount).
		CreditAccount(params.EscrowAccount).
		Amount(params.Amount).
		Ledger(params.Ledger).
		Code(params.Code).
		Pending().
		Timeout(params.EscrowTimeout).
		Build()

	results, err := p.client.CreateTransfers([]types.Transfer{escrowTransfer})
	if err != nil {
		return fmt.Errorf("failed to create escrow: %w", err)
	}
	if len(results) > 0 {
		return fmt.Errorf("escrow failed: %s", results[0].Result.String())
	}
	return nil
}

// ReleaseEscrow releases escrowed funds to the final recipient
func (p *Patterns) ReleaseEscrow(
	releaseID types.Uint128,
	escrowID types.Uint128,
	escrowAccount types.Uint128,
	recipientAccount types.Uint128,
	amount types.Uint128,
	ledger uint32,
	code uint16,
) error {
	// Post the pending transfer and create a linked transfer to the recipient
	postTransfer := transfer.New(releaseID).
		PostPending(escrowID).
		Linked().
		Build()

	releaseTransfer := transfer.New(types.ToUint128(uint64(releaseID.ToUint128()) + 1)).
		DebitAccount(escrowAccount).
		CreditAccount(recipientAccount).
		Amount(amount).
		Ledger(ledger).
		Code(code).
		Build()

	results, err := p.client.CreateTransfers([]types.Transfer{postTransfer, releaseTransfer})
	if err != nil {
		return fmt.Errorf("failed to release escrow: %w", err)
	}
	if len(results) > 0 {
		return fmt.Errorf("escrow release failed: %s", results[0].Result.String())
	}
	return nil
}

// CancelEscrow cancels an escrow and returns funds to the sender
func (p *Patterns) CancelEscrow(cancelID, escrowID types.Uint128) error {
	voidTransfer := transfer.New(cancelID).
		VoidPending(escrowID).
		Build()

	results, err := p.client.CreateTransfers([]types.Transfer{voidTransfer})
	if err != nil {
		return fmt.Errorf("failed to cancel escrow: %w", err)
	}
	if len(results) > 0 {
		return fmt.Errorf("escrow cancellation failed: %s", results[0].Result.String())
	}
	return nil
}

// RefundParams contains parameters for a refund transaction
type RefundParams struct {
	RefundID        types.Uint128
	OriginalTransferID types.Uint128
	FromAccount     types.Uint128
	ToAccount       types.Uint128
	Amount          types.Uint128
	Ledger          uint32
	Code            uint16
}

// Refund creates a transfer that reverses a previous transaction
// This is typically used for returns, chargebacks, or corrections
func (p *Patterns) Refund(params RefundParams) error {
	refundTransfer := transfer.New(params.RefundID).
		DebitAccount(params.FromAccount).
		CreditAccount(params.ToAccount).
		Amount(params.Amount).
		Ledger(params.Ledger).
		Code(params.Code).
		UserData128(params.OriginalTransferID). // Track original transfer
		Build()

	results, err := p.client.CreateTransfers([]types.Transfer{refundTransfer})
	if err != nil {
		return fmt.Errorf("failed to create refund: %w", err)
	}
	if len(results) > 0 {
		return fmt.Errorf("refund failed: %s", results[0].Result.String())
	}
	return nil
}

// MultiPartyTransferParams contains parameters for a multi-party split transfer
type MultiPartyTransferParams struct {
	BaseID      types.Uint128
	FromAccount types.Uint128
	Recipients  []Recipient
	Ledger      uint32
	Code        uint16
}

// Recipient represents a recipient in a multi-party transfer
type Recipient struct {
	AccountID types.Uint128
	Amount    types.Uint128
}

// MultiPartyTransfer creates a series of linked transfers from one account to multiple recipients
// All transfers succeed or fail atomically
func (p *Patterns) MultiPartyTransfer(params MultiPartyTransferParams) error {
	if len(params.Recipients) == 0 {
		return fmt.Errorf("at least one recipient is required")
	}

	transfers := make([]types.Transfer, len(params.Recipients))
	baseIDValue := params.BaseID.ToUint128()

	for i, recipient := range params.Recipients {
		transferID := types.ToUint128(baseIDValue + uint64(i))
		builder := transfer.New(transferID).
			DebitAccount(params.FromAccount).
			CreditAccount(recipient.AccountID).
			Amount(recipient.Amount).
			Ledger(params.Ledger).
			Code(params.Code)

		// Link all transfers except the last one
		if i < len(params.Recipients)-1 {
			builder = builder.Linked()
		}

		transfers[i] = builder.Build()
	}

	results, err := p.client.CreateTransfers(transfers)
	if err != nil {
		return fmt.Errorf("failed to create multi-party transfer: %w", err)
	}
	if len(results) > 0 {
		return fmt.Errorf("multi-party transfer failed at index %d: %s",
			results[0].Index, results[0].Result.String())
	}
	return nil
}

// ExchangeParams contains parameters for a currency exchange transaction
type ExchangeParams struct {
	BaseID              types.Uint128
	FromAccount         types.Uint128
	ToAccount           types.Uint128
	IntermediaryAccount types.Uint128
	FromAmount          types.Uint128
	ToAmount            types.Uint128
	FromLedger          uint32
	ToLedger            uint32
	Code                uint16
}

// CurrencyExchange performs a two-ledger currency exchange via an intermediary account
// This atomically debits one currency and credits another
func (p *Patterns) CurrencyExchange(params ExchangeParams) error {
	baseIDValue := params.BaseID.ToUint128()

	// Transfer from source account to intermediary (source currency)
	transfer1 := transfer.New(params.BaseID).
		DebitAccount(params.FromAccount).
		CreditAccount(params.IntermediaryAccount).
		Amount(params.FromAmount).
		Ledger(params.FromLedger).
		Code(params.Code).
		Linked().
		Build()

	// Transfer from intermediary to destination (destination currency)
	transfer2 := transfer.New(types.ToUint128(baseIDValue + 1)).
		DebitAccount(params.IntermediaryAccount).
		CreditAccount(params.ToAccount).
		Amount(params.ToAmount).
		Ledger(params.ToLedger).
		Code(params.Code).
		Build()

	results, err := p.client.CreateTransfers([]types.Transfer{transfer1, transfer2})
	if err != nil {
		return fmt.Errorf("failed to create currency exchange: %w", err)
	}
	if len(results) > 0 {
		return fmt.Errorf("currency exchange failed at index %d: %s",
			results[0].Index, results[0].Result.String())
	}
	return nil
}

// ConditionalTransferParams contains parameters for a conditional transfer
type ConditionalTransferParams struct {
	TransferID     types.Uint128
	FromAccount    types.Uint128
	ToAccount      types.Uint128
	Amount         types.Uint128
	Ledger         uint32
	Code           uint16
	ConditionTimeout time.Duration
}

// ConditionalTransfer creates a pending transfer that must be explicitly confirmed or cancelled
// Useful for payment authorizations, pre-authorizations, etc.
func (p *Patterns) ConditionalTransfer(params ConditionalTransferParams) error {
	pendingTransfer := transfer.New(params.TransferID).
		DebitAccount(params.FromAccount).
		CreditAccount(params.ToAccount).
		Amount(params.Amount).
		Ledger(params.Ledger).
		Code(params.Code).
		Pending().
		Timeout(params.ConditionTimeout).
		Build()

	results, err := p.client.CreateTransfers([]types.Transfer{pendingTransfer})
	if err != nil {
		return fmt.Errorf("failed to create conditional transfer: %w", err)
	}
	if len(results) > 0 {
		return fmt.Errorf("conditional transfer failed: %s", results[0].Result.String())
	}
	return nil
}

// ConfirmConditionalTransfer confirms a conditional transfer
func (p *Patterns) ConfirmConditionalTransfer(confirmID, pendingID types.Uint128) error {
	postTransfer := transfer.New(confirmID).
		PostPending(pendingID).
		Build()

	results, err := p.client.CreateTransfers([]types.Transfer{postTransfer})
	if err != nil {
		return fmt.Errorf("failed to confirm conditional transfer: %w", err)
	}
	if len(results) > 0 {
		return fmt.Errorf("conditional transfer confirmation failed: %s", results[0].Result.String())
	}
	return nil
}

// CancelConditionalTransfer cancels a conditional transfer
func (p *Patterns) CancelConditionalTransfer(cancelID, pendingID types.Uint128) error {
	voidTransfer := transfer.New(cancelID).
		VoidPending(pendingID).
		Build()

	results, err := p.client.CreateTransfers([]types.Transfer{voidTransfer})
	if err != nil {
		return fmt.Errorf("failed to cancel conditional transfer: %w", err)
	}
	if len(results) > 0 {
		return fmt.Errorf("conditional transfer cancellation failed: %s", results[0].Result.String())
	}
	return nil
}

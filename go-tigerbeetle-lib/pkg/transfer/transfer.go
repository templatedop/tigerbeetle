package transfer

import (
	"time"

	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

// Builder provides a fluent interface for creating TigerBeetle transfers
type Builder struct {
	transfer types.Transfer
}

// New creates a new transfer builder with the given ID
func New(id types.Uint128) *Builder {
	return &Builder{
		transfer: types.Transfer{
			ID: id,
		},
	}
}

// DebitAccount sets the account to debit (sender)
func (b *Builder) DebitAccount(accountID types.Uint128) *Builder {
	b.transfer.DebitAccountID = accountID
	return b
}

// CreditAccount sets the account to credit (receiver)
func (b *Builder) CreditAccount(accountID types.Uint128) *Builder {
	b.transfer.CreditAccountID = accountID
	return b
}

// Amount sets the transfer amount
func (b *Builder) Amount(amount types.Uint128) *Builder {
	b.transfer.Amount = amount
	return b
}

// AmountUint64 sets the transfer amount from a uint64
func (b *Builder) AmountUint64(amount uint64) *Builder {
	b.transfer.Amount = types.ToUint128(amount)
	return b
}

// Ledger sets the ledger partition
// Must match the ledger of both debit and credit accounts
func (b *Builder) Ledger(ledger uint32) *Builder {
	b.transfer.Ledger = ledger
	return b
}

// Code sets the transfer reason code
// Used to categorize transfers (e.g., purchase, refund, fee, etc.)
func (b *Builder) Code(code uint16) *Builder {
	b.transfer.Code = code
	return b
}

// UserData128 sets the 128-bit user data field
func (b *Builder) UserData128(data types.Uint128) *Builder {
	b.transfer.UserData128 = data
	return b
}

// UserData64 sets the 64-bit user data field
func (b *Builder) UserData64(data uint64) *Builder {
	b.transfer.UserData64 = data
	return b
}

// UserData32 sets the 32-bit user data field
func (b *Builder) UserData32(data uint32) *Builder {
	b.transfer.UserData32 = data
	return b
}

// Timeout sets the timeout for pending transfers
// Only valid for pending transfers
func (b *Builder) Timeout(timeout time.Duration) *Builder {
	b.transfer.Timeout = uint32(timeout.Nanoseconds())
	return b
}

// Flags sets the transfer flags
func (b *Builder) Flags(flags types.TransferFlags) *Builder {
	b.transfer.Flags = flags.ToUint16()
	return b
}

// Linked marks this transfer as part of a linked chain
// Linked transfers execute atomically - all succeed or all fail
func (b *Builder) Linked() *Builder {
	b.transfer.Flags |= types.TransferFlags{Linked: true}.ToUint16()
	return b
}

// Pending marks this as a two-phase pending transfer
// Reserves funds without immediately posting them
func (b *Builder) Pending() *Builder {
	b.transfer.Flags |= types.TransferFlags{Pending: true}.ToUint16()
	return b
}

// PostPending marks this transfer as posting a pending transfer
// Moves funds from pending to posted
func (b *Builder) PostPending(pendingID types.Uint128) *Builder {
	b.transfer.Flags |= types.TransferFlags{PostPendingTransfer: true}.ToUint16()
	b.transfer.PendingID = pendingID
	return b
}

// VoidPending marks this transfer as voiding a pending transfer
// Cancels a pending transfer without posting it
func (b *Builder) VoidPending(pendingID types.Uint128) *Builder {
	b.transfer.Flags |= types.TransferFlags{VoidPendingTransfer: true}.ToUint16()
	b.transfer.PendingID = pendingID
	return b
}

// BalancingDebit allows the debit account to exceed its credits
func (b *Builder) BalancingDebit() *Builder {
	b.transfer.Flags |= types.TransferFlags{BalancingDebit: true}.ToUint16()
	return b
}

// BalancingCredit allows the credit account to exceed its debits
func (b *Builder) BalancingCredit() *Builder {
	b.transfer.Flags |= types.TransferFlags{BalancingCredit: true}.ToUint16()
	return b
}

// Build returns the constructed transfer
func (b *Builder) Build() types.Transfer {
	return b.transfer
}

// Manager provides high-level transfer management operations
type Manager struct {
	client TransferClient
}

// TransferClient defines the interface for transfer operations
type TransferClient interface {
	CreateTransfers([]types.Transfer) ([]types.TransferEventResult, error)
	LookupTransfers([]types.Uint128) ([]types.Transfer, error)
	QueryTransfers(types.QueryFilter) ([]types.Transfer, error)
	GetAccountTransfers(types.AccountFilter) ([]types.Transfer, error)
}

// NewManager creates a new transfer manager
func NewManager(client TransferClient) *Manager {
	return &Manager{client: client}
}

// Create creates a single transfer
func (m *Manager) Create(transfer types.Transfer) error {
	results, err := m.client.CreateTransfers([]types.Transfer{transfer})
	if err != nil {
		return err
	}
	if len(results) > 0 {
		return &TransferError{
			Index:  results[0].Index,
			Result: results[0].Result,
		}
	}
	return nil
}

// CreateBatch creates multiple transfers
func (m *Manager) CreateBatch(transfers []types.Transfer) ([]types.TransferEventResult, error) {
	return m.client.CreateTransfers(transfers)
}

// Lookup retrieves a single transfer by ID
func (m *Manager) Lookup(id types.Uint128) (*types.Transfer, error) {
	transfers, err := m.client.LookupTransfers([]types.Uint128{id})
	if err != nil {
		return nil, err
	}
	if len(transfers) == 0 {
		return nil, ErrTransferNotFound
	}
	return &transfers[0], nil
}

// LookupBatch retrieves multiple transfers by their IDs
func (m *Manager) LookupBatch(ids []types.Uint128) ([]types.Transfer, error) {
	return m.client.LookupTransfers(ids)
}

// Query performs a filtered query on transfers
func (m *Manager) Query(filter types.QueryFilter) ([]types.Transfer, error) {
	return m.client.QueryTransfers(filter)
}

// GetForAccount retrieves transfers for a specific account
func (m *Manager) GetForAccount(accountID types.Uint128, filter types.AccountFilter) ([]types.Transfer, error) {
	filter.AccountID = accountID
	return m.client.GetAccountTransfers(filter)
}

// SimpleTransfer creates a standard immediate transfer
func (m *Manager) SimpleTransfer(
	id types.Uint128,
	debitAccount, creditAccount types.Uint128,
	amount types.Uint128,
	ledger uint32,
	code uint16,
) error {
	transfer := New(id).
		DebitAccount(debitAccount).
		CreditAccount(creditAccount).
		Amount(amount).
		Ledger(ledger).
		Code(code).
		Build()

	return m.Create(transfer)
}

// PendingTransfer creates a two-phase pending transfer
func (m *Manager) PendingTransfer(
	id types.Uint128,
	debitAccount, creditAccount types.Uint128,
	amount types.Uint128,
	ledger uint32,
	code uint16,
	timeout time.Duration,
) error {
	transfer := New(id).
		DebitAccount(debitAccount).
		CreditAccount(creditAccount).
		Amount(amount).
		Ledger(ledger).
		Code(code).
		Pending().
		Timeout(timeout).
		Build()

	return m.Create(transfer)
}

// PostPending posts a pending transfer
func (m *Manager) PostPending(id, pendingID types.Uint128) error {
	transfer := New(id).
		PostPending(pendingID).
		Build()

	return m.Create(transfer)
}

// VoidPending voids a pending transfer
func (m *Manager) VoidPending(id, pendingID types.Uint128) error {
	transfer := New(id).
		VoidPending(pendingID).
		Build()

	return m.Create(transfer)
}

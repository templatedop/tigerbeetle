package account

import (
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

// Builder provides a fluent interface for creating TigerBeetle accounts
type Builder struct {
	account types.Account
}

// New creates a new account builder with the given ID
func New(id types.Uint128) *Builder {
	return &Builder{
		account: types.Account{
			ID: id,
		},
	}
}

// Ledger sets the ledger partition for the account
// Ledger partitions accounts that can transact together
func (b *Builder) Ledger(ledger uint32) *Builder {
	b.account.Ledger = ledger
	return b
}

// Code sets the account category code
// This is a user-defined identifier for categorizing accounts
func (b *Builder) Code(code uint16) *Builder {
	b.account.Code = code
	return b
}

// UserData128 sets the 128-bit user data field
// Can be used for external identifiers or application-specific data
func (b *Builder) UserData128(data types.Uint128) *Builder {
	b.account.UserData128 = data
	return b
}

// UserData64 sets the 64-bit user data field
func (b *Builder) UserData64(data uint64) *Builder {
	b.account.UserData64 = data
	return b
}

// UserData32 sets the 32-bit user data field
func (b *Builder) UserData32(data uint32) *Builder {
	b.account.UserData32 = data
	return b
}

// Flags sets behavioral flags for the account
func (b *Builder) Flags(flags types.AccountFlags) *Builder {
	b.account.Flags = flags
	return b
}

// LinkedAccount marks this account as part of a linked chain
// Linked accounts are created atomically - all succeed or all fail
func (b *Builder) LinkedAccount() *Builder {
	b.account.Flags |= types.AccountFlags{Linked: true}
	return b
}

// DebitsMustNotExceedCredits sets the flag to prevent overdrafts
// Ensures debits_posted + debits_pending <= credits_posted
func (b *Builder) DebitsMustNotExceedCredits() *Builder {
	b.account.Flags |= types.AccountFlags{DebitsMustNotExceedCredits: true}
	return b
}

// CreditsMustNotExceedDebits sets the flag for liability accounts
// Ensures credits_posted + credits_pending <= debits_posted
func (b *Builder) CreditsMustNotExceedDebits() *Builder {
	b.account.Flags |= types.AccountFlags{CreditsMustNotExceedDebits: true}
	return b
}

// HistoryEnabled enables account balance history tracking
// Required to use GetAccountBalances query
func (b *Builder) HistoryEnabled() *Builder {
	b.account.Flags |= types.AccountFlags{History: true}
	return b
}

// Build returns the constructed account
func (b *Builder) Build() types.Account {
	return b.account
}

// Manager provides high-level account management operations
type Manager struct {
	client AccountClient
}

// AccountClient defines the interface for account operations
type AccountClient interface {
	CreateAccounts([]types.Account) ([]types.CreateAccountsError, error)
	LookupAccounts([]types.Uint128) ([]types.Account, error)
	QueryAccounts(types.QueryFilter) ([]types.Account, error)
	GetAccountBalances(types.AccountFilter) ([]types.AccountBalance, error)
}

// NewManager creates a new account manager
func NewManager(client AccountClient) *Manager {
	return &Manager{client: client}
}

// Create creates a single account
func (m *Manager) Create(account types.Account) error {
	results, err := m.client.CreateAccounts([]types.Account{account})
	if err != nil {
		return err
	}
	if len(results) > 0 {
		return &AccountError{
			Index:  results[0].Index,
			Result: results[0].Result,
		}
	}
	return nil
}

// CreateBatch creates multiple accounts atomically
func (m *Manager) CreateBatch(accounts []types.Account) ([]types.CreateAccountsError, error) {
	return m.client.CreateAccounts(accounts)
}

// Lookup retrieves a single account by ID
func (m *Manager) Lookup(id types.Uint128) (*types.Account, error) {
	accounts, err := m.client.LookupAccounts([]types.Uint128{id})
	if err != nil {
		return nil, err
	}
	if len(accounts) == 0 {
		return nil, ErrAccountNotFound
	}
	return &accounts[0], nil
}

// LookupBatch retrieves multiple accounts by their IDs
func (m *Manager) LookupBatch(ids []types.Uint128) ([]types.Account, error) {
	return m.client.LookupAccounts(ids)
}

// Query performs a filtered query on accounts
func (m *Manager) Query(filter types.QueryFilter) ([]types.Account, error) {
	return m.client.QueryAccounts(filter)
}

// GetBalances retrieves historical balances for an account
func (m *Manager) GetBalances(accountID types.Uint128, filter types.AccountFilter) ([]types.AccountBalance, error) {
	filter.AccountID = accountID
	return m.client.GetAccountBalances(filter)
}

// GetBalance retrieves the current balance for an account
func (m *Manager) GetBalance(id types.Uint128) (*Balance, error) {
	account, err := m.Lookup(id)
	if err != nil {
		return nil, err
	}
	return &Balance{
		Account:        *account,
		DebitsPosted:   types.ToUint128(account.DebitsPosted),
		CreditsPosted:  types.ToUint128(account.CreditsPosted),
		DebitsPending:  types.ToUint128(account.DebitsPending),
		CreditsPending: types.ToUint128(account.CreditsPending),
	}, nil
}

// Balance represents an account's current balance state
type Balance struct {
	Account        types.Account
	DebitsPosted   types.Uint128
	CreditsPosted  types.Uint128
	DebitsPending  types.Uint128
	CreditsPending types.Uint128
}

// NetPosted returns the net posted balance (credits - debits)
func (b *Balance) NetPosted() types.Uint128 {
	return types.ToUint128(b.Account.CreditsPosted - b.Account.DebitsPosted)
}

// AvailableBalance returns the available balance considering pending amounts
func (b *Balance) AvailableBalance() types.Uint128 {
	// Available = credits_posted - debits_posted - debits_pending
	return types.ToUint128(
		b.Account.CreditsPosted - b.Account.DebitsPosted - b.Account.DebitsPending,
	)
}

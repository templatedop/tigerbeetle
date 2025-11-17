package query

import (
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

// AccountFilterBuilder provides a fluent interface for building account filters
type AccountFilterBuilder struct {
	filter types.AccountFilter
}

// NewAccountFilter creates a new account filter builder
func NewAccountFilter(accountID types.Uint128) *AccountFilterBuilder {
	return &AccountFilterBuilder{
		filter: types.AccountFilter{
			AccountID: accountID,
		},
	}
}

// TimestampMin sets the minimum timestamp for the filter
func (b *AccountFilterBuilder) TimestampMin(timestamp uint64) *AccountFilterBuilder {
	b.filter.TimestampMin = timestamp
	return b
}

// TimestampMax sets the maximum timestamp for the filter
func (b *AccountFilterBuilder) TimestampMax(timestamp uint64) *AccountFilterBuilder {
	b.filter.TimestampMax = timestamp
	return b
}

// Limit sets the maximum number of results
func (b *AccountFilterBuilder) Limit(limit uint32) *AccountFilterBuilder {
	b.filter.Limit = limit
	return b
}

// Debits filters for transfers where the account is debited
func (b *AccountFilterBuilder) Debits() *AccountFilterBuilder {
	b.filter.Flags = types.AccountFilterFlags{Debits: true}
	return b
}

// Credits filters for transfers where the account is credited
func (b *AccountFilterBuilder) Credits() *AccountFilterBuilder {
	b.filter.Flags = types.AccountFilterFlags{Credits: true}
	return b
}

// Reversed returns results in reverse chronological order
func (b *AccountFilterBuilder) Reversed() *AccountFilterBuilder {
	b.filter.Flags = types.AccountFilterFlags{Reversed: true}
	return b
}

// Build returns the constructed filter
func (b *AccountFilterBuilder) Build() types.AccountFilter {
	return b.filter
}

// QueryFilterBuilder provides a fluent interface for building query filters
type QueryFilterBuilder struct {
	filter types.QueryFilter
}

// NewQueryFilter creates a new query filter builder
func NewQueryFilter() *QueryFilterBuilder {
	return &QueryFilterBuilder{
		filter: types.QueryFilter{},
	}
}

// UserData128 filters by the 128-bit user data field
func (b *QueryFilterBuilder) UserData128(data types.Uint128) *QueryFilterBuilder {
	b.filter.UserData128 = data
	return b
}

// UserData64 filters by the 64-bit user data field
func (b *QueryFilterBuilder) UserData64(data uint64) *QueryFilterBuilder {
	b.filter.UserData64 = data
	return b
}

// UserData32 filters by the 32-bit user data field
func (b *QueryFilterBuilder) UserData32(data uint32) *QueryFilterBuilder {
	b.filter.UserData32 = data
	return b
}

// Ledger filters by ledger partition
func (b *QueryFilterBuilder) Ledger(ledger uint32) *QueryFilterBuilder {
	b.filter.Ledger = ledger
	return b
}

// Code filters by account or transfer code
func (b *QueryFilterBuilder) Code(code uint16) *QueryFilterBuilder {
	b.filter.Code = code
	return b
}

// TimestampMin sets the minimum timestamp for the filter
func (b *QueryFilterBuilder) TimestampMin(timestamp uint64) *QueryFilterBuilder {
	b.filter.TimestampMin = timestamp
	return b
}

// TimestampMax sets the maximum timestamp for the filter
func (b *QueryFilterBuilder) TimestampMax(timestamp uint64) *QueryFilterBuilder {
	b.filter.TimestampMax = timestamp
	return b
}

// Limit sets the maximum number of results
func (b *QueryFilterBuilder) Limit(limit uint32) *QueryFilterBuilder {
	b.filter.Limit = limit
	return b
}

// Reversed returns results in reverse chronological order
func (b *QueryFilterBuilder) Reversed() *QueryFilterBuilder {
	b.filter.Flags = types.QueryFilterFlags{Reversed: true}
	return b
}

// Build returns the constructed filter
func (b *QueryFilterBuilder) Build() types.QueryFilter {
	return b.filter
}

// Helper provides high-level query operations
type Helper struct {
	client QueryClient
}

// QueryClient defines the interface for query operations
type QueryClient interface {
	QueryAccounts(types.QueryFilter) ([]types.Account, error)
	QueryTransfers(types.QueryFilter) ([]types.Transfer, error)
	GetAccountTransfers(types.AccountFilter) ([]types.Transfer, error)
	GetAccountBalances(types.AccountFilter) ([]types.AccountBalance, error)
}

// NewHelper creates a new query helper
func NewHelper(client QueryClient) *Helper {
	return &Helper{client: client}
}

// AccountsByLedger queries all accounts in a specific ledger
func (h *Helper) AccountsByLedger(ledger uint32, limit uint32) ([]types.Account, error) {
	filter := NewQueryFilter().
		Ledger(ledger).
		Limit(limit).
		Build()
	return h.client.QueryAccounts(filter)
}

// AccountsByCode queries accounts with a specific code
func (h *Helper) AccountsByCode(code uint16, limit uint32) ([]types.Account, error) {
	filter := NewQueryFilter().
		Code(code).
		Limit(limit).
		Build()
	return h.client.QueryAccounts(filter)
}

// AccountsByUserData queries accounts by user data field
func (h *Helper) AccountsByUserData64(userData uint64, limit uint32) ([]types.Account, error) {
	filter := NewQueryFilter().
		UserData64(userData).
		Limit(limit).
		Build()
	return h.client.QueryAccounts(filter)
}

// TransfersByLedger queries all transfers in a specific ledger
func (h *Helper) TransfersByLedger(ledger uint32, limit uint32) ([]types.Transfer, error) {
	filter := NewQueryFilter().
		Ledger(ledger).
		Limit(limit).
		Build()
	return h.client.QueryTransfers(filter)
}

// TransfersByCode queries transfers with a specific code
func (h *Helper) TransfersByCode(code uint16, limit uint32) ([]types.Transfer, error) {
	filter := NewQueryFilter().
		Code(code).
		Limit(limit).
		Build()
	return h.client.QueryTransfers(filter)
}

// TransfersByTimeRange queries transfers within a time range
func (h *Helper) TransfersByTimeRange(
	minTimestamp, maxTimestamp uint64,
	limit uint32,
) ([]types.Transfer, error) {
	filter := NewQueryFilter().
		TimestampMin(minTimestamp).
		TimestampMax(maxTimestamp).
		Limit(limit).
		Build()
	return h.client.QueryTransfers(filter)
}

// AccountDebits retrieves all debit transfers for an account
func (h *Helper) AccountDebits(accountID types.Uint128, limit uint32) ([]types.Transfer, error) {
	filter := NewAccountFilter(accountID).
		Debits().
		Limit(limit).
		Build()
	return h.client.GetAccountTransfers(filter)
}

// AccountCredits retrieves all credit transfers for an account
func (h *Helper) AccountCredits(accountID types.Uint128, limit uint32) ([]types.Transfer, error) {
	filter := NewAccountFilter(accountID).
		Credits().
		Limit(limit).
		Build()
	return h.client.GetAccountTransfers(filter)
}

// AccountRecentTransfers retrieves the most recent transfers for an account
func (h *Helper) AccountRecentTransfers(accountID types.Uint128, limit uint32) ([]types.Transfer, error) {
	filter := NewAccountFilter(accountID).
		Reversed().
		Limit(limit).
		Build()
	return h.client.GetAccountTransfers(filter)
}

// AccountBalanceHistory retrieves balance history for an account
func (h *Helper) AccountBalanceHistory(accountID types.Uint128, limit uint32) ([]types.AccountBalance, error) {
	filter := NewAccountFilter(accountID).
		Limit(limit).
		Build()
	return h.client.GetAccountBalances(filter)
}

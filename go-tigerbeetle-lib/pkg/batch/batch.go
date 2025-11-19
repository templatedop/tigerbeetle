package batch

import (
	"fmt"

	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

const (
	// MaxBatchSize is the default maximum batch size for TigerBeetle
	MaxBatchSize = 8189
)

// BatchClient defines the interface for batch operations
type BatchClient interface {
	CreateAccounts([]types.Account) ([]types.AccountEventResult, error)
	CreateTransfers([]types.Transfer) ([]types.TransferEventResult, error)
}

// AccountBatcher helps batch account creation operations
type AccountBatcher struct {
	client    BatchClient
	accounts  []types.Account
	batchSize int
}

// NewAccountBatcher creates a new account batcher
func NewAccountBatcher(client BatchClient, batchSize int) *AccountBatcher {
	if batchSize <= 0 || batchSize > MaxBatchSize {
		batchSize = MaxBatchSize
	}
	return &AccountBatcher{
		client:    client,
		accounts:  make([]types.Account, 0, batchSize),
		batchSize: batchSize,
	}
}

// Add adds an account to the batch
// If the batch is full, it automatically flushes
func (b *AccountBatcher) Add(account types.Account) ([]types.AccountEventResult, error) {
	b.accounts = append(b.accounts, account)
	if len(b.accounts) >= b.batchSize {
		return b.Flush()
	}
	return nil, nil
}

// Flush sends all pending accounts to TigerBeetle
func (b *AccountBatcher) Flush() ([]types.AccountEventResult, error) {
	if len(b.accounts) == 0 {
		return nil, nil
	}

	results, err := b.client.CreateAccounts(b.accounts)
	b.accounts = b.accounts[:0] // Reset the slice
	return results, err
}

// Len returns the number of pending accounts in the batch
func (b *AccountBatcher) Len() int {
	return len(b.accounts)
}

// TransferBatcher helps batch transfer creation operations
type TransferBatcher struct {
	client    BatchClient
	transfers []types.Transfer
	batchSize int
}

// NewTransferBatcher creates a new transfer batcher
func NewTransferBatcher(client BatchClient, batchSize int) *TransferBatcher {
	if batchSize <= 0 || batchSize > MaxBatchSize {
		batchSize = MaxBatchSize
	}
	return &TransferBatcher{
		client:    client,
		transfers: make([]types.Transfer, 0, batchSize),
		batchSize: batchSize,
	}
}

// Add adds a transfer to the batch
// If the batch is full, it automatically flushes
func (b *TransferBatcher) Add(transfer types.Transfer) ([]types.TransferEventResult, error) {
	b.transfers = append(b.transfers, transfer)
	if len(b.transfers) >= b.batchSize {
		return b.Flush()
	}
	return nil, nil
}

// Flush sends all pending transfers to TigerBeetle
func (b *TransferBatcher) Flush() ([]types.TransferEventResult, error) {
	if len(b.transfers) == 0 {
		return nil, nil
	}

	results, err := b.client.CreateTransfers(b.transfers)
	b.transfers = b.transfers[:0] // Reset the slice
	return results, err
}

// Len returns the number of pending transfers in the batch
func (b *TransferBatcher) Len() int {
	return len(b.transfers)
}

// Processor provides utilities for processing large datasets in batches
type Processor struct {
	client    BatchClient
	batchSize int
}

// NewProcessor creates a new batch processor
func NewProcessor(client BatchClient, batchSize int) *Processor {
	if batchSize <= 0 || batchSize > MaxBatchSize {
		batchSize = MaxBatchSize
	}
	return &Processor{
		client:    client,
		batchSize: batchSize,
	}
}

// ProcessAccounts processes a large slice of accounts in batches
func (p *Processor) ProcessAccounts(accounts []types.Account) ([]types.AccountEventResult, error) {
	var allErrors []types.AccountEventResult

	for i := 0; i < len(accounts); i += p.batchSize {
		end := i + p.batchSize
		if end > len(accounts) {
			end = len(accounts)
		}

		batch := accounts[i:end]
		results, err := p.client.CreateAccounts(batch)
		if err != nil {
			return allErrors, fmt.Errorf("batch %d-%d failed: %w", i, end, err)
		}

		// Adjust error indices to reflect position in original slice
		for _, result := range results {
			result.Index += uint32(i)
			allErrors = append(allErrors, result)
		}
	}

	return allErrors, nil
}

// ProcessTransfers processes a large slice of transfers in batches
func (p *Processor) ProcessTransfers(transfers []types.Transfer) ([]types.TransferEventResult, error) {
	var allErrors []types.TransferEventResult

	for i := 0; i < len(transfers); i += p.batchSize {
		end := i + p.batchSize
		if end > len(transfers) {
			end = len(transfers)
		}

		batch := transfers[i:end]
		results, err := p.client.CreateTransfers(batch)
		if err != nil {
			return allErrors, fmt.Errorf("batch %d-%d failed: %w", i, end, err)
		}

		// Adjust error indices to reflect position in original slice
		for _, result := range results {
			result.Index += uint32(i)
			allErrors = append(allErrors, result)
		}
	}

	return allErrors, nil
}

// LinkedChain helps build chains of linked operations
type LinkedChain struct {
	transfers []types.Transfer
}

// NewLinkedChain creates a new linked transfer chain
func NewLinkedChain() *LinkedChain {
	return &LinkedChain{
		transfers: make([]types.Transfer, 0),
	}
}

// AddTransfer adds a transfer to the chain
// All transfers except the last will be automatically marked as linked
func (c *LinkedChain) AddTransfer(transfer types.Transfer) *LinkedChain {
	c.transfers = append(c.transfers, transfer)
	return c
}

// Build returns the chain of transfers with proper linking flags
func (c *LinkedChain) Build() []types.Transfer {
	if len(c.transfers) == 0 {
		return c.transfers
	}

	// Set linked flag on all transfers except the last
	for i := 0; i < len(c.transfers)-1; i++ {
		c.transfers[i].Flags |= types.TransferFlags{Linked: true}.ToUint16()
	}

	return c.transfers
}

// Len returns the number of transfers in the chain
func (c *LinkedChain) Len() int {
	return len(c.transfers)
}

// Execute executes the linked chain atomically
func (c *LinkedChain) Execute(client BatchClient) ([]types.TransferEventResult, error) {
	if len(c.transfers) == 0 {
		return nil, fmt.Errorf("cannot execute empty chain")
	}
	return client.CreateTransfers(c.Build())
}

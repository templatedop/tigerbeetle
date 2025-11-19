# TigerBeetle Go Library

A comprehensive, high-level Go library for interacting with TigerBeetle, the distributed financial accounting database. This library provides idiomatic Go abstractions over the official TigerBeetle client, making it easier to build financial applications with complex transaction patterns.

## Features

- **Fluent Builder API**: Easy-to-use builder patterns for creating accounts and transfers
- **High-Level Transaction Patterns**: Pre-built patterns for common scenarios:
  - Escrow transactions (two-phase transfers)
  - Multi-party transfers (payment splitting)
  - Currency exchange
  - Refunds and reversals
  - Conditional transfers
- **Efficient Batch Operations**: Tools for processing large volumes of data efficiently
- **Query Helpers**: Simplified querying and filtering of accounts and transfers
- **Type-Safe**: Leverages Go's type system for compile-time safety
- **Well-Tested**: Comprehensive test coverage
- **Production-Ready**: Built on top of the official TigerBeetle Go client

## Installation

```bash
go get gitlab.cept.gov.in/it-2.0-common/ledgers
```

## Running TigerBeetle

### Using Docker (Recommended)

The easiest way to run TigerBeetle is using Docker Compose:

```bash
# Start single node (development)
make docker-start

# Start 3-node cluster (production)
make cluster-start

# View logs
make docker-logs

# Stop and clean up
make docker-clean
```

See [DOCKER.md](DOCKER.md) for complete Docker setup documentation.

### Manual Installation

Alternatively, install TigerBeetle directly:

```bash
# Download and install TigerBeetle
curl -LO https://github.com/tigerbeetle/tigerbeetle/releases/latest/download/tigerbeetle-$(uname -s)-$(uname -m).zip
unzip tigerbeetle-*.zip
sudo mv tigerbeetle /usr/local/bin/

# Create data directory and start
tigerbeetle format --cluster=0 --replica=0 0_0.tigerbeetle
tigerbeetle start --addresses=3000 0_0.tigerbeetle
```

## Quick Start

### Creating a Client

```go
import (
    "gitlab.cept.gov.in/it-2.0-common/ledgers/pkg/client"
    "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

// Create a TigerBeetle client
c, err := client.New(client.Config{
    ClusterID:      types.ToUint128(0),
    ReplicaAddrs:   []string{"3000"},
    MaxConcurrency: 32,
})
if err != nil {
    log.Fatal(err)
}
defer c.Close()
```

### Creating Accounts

```go
import (
    "gitlab.cept.gov.in/it-2.0-common/ledgers/pkg/account"
)

// Create account manager
accountMgr := account.NewManager(c)

// Build an account with fluent API
account := account.New(types.ToUint128(1)).
    Ledger(1).
    Code(1).
    HistoryEnabled().
    DebitsMustNotExceedCredits(). // Prevent overdrafts
    Build()

// Create the account
if err := accountMgr.Create(account); err != nil {
    log.Fatal(err)
}
```

### Simple Transfer

```go
import (
    "gitlab.cept.gov.in/it-2.0-common/ledgers/pkg/transfer"
)

// Create transfer manager
transferMgr := transfer.NewManager(c)

// Perform a simple transfer
err := transferMgr.SimpleTransfer(
    types.ToUint128(1),    // Transfer ID
    types.ToUint128(1),    // From account
    types.ToUint128(2),    // To account
    types.ToUint128(1000), // Amount
    1,                     // Ledger
    1,                     // Code
)
```

### Advanced Patterns

#### Escrow Transaction

```go
import (
    "gitlab.cept.gov.in/it-2.0-common/ledgers/pkg/transaction"
    "time"
)

txPatterns := transaction.NewPatterns(c)

// Create escrow
escrowParams := transaction.EscrowParams{
    EscrowID:      types.ToUint128(1000),
    FromAccount:   types.ToUint128(100),
    EscrowAccount: types.ToUint128(102),
    Amount:        types.ToUint128(5000),
    Ledger:        1,
    Code:          1,
    EscrowTimeout: 24 * time.Hour,
}

if err := txPatterns.Escrow(escrowParams); err != nil {
    log.Fatal(err)
}

// Release escrow to recipient
err = txPatterns.ReleaseEscrow(
    types.ToUint128(1001), // Release ID
    types.ToUint128(1000), // Escrow ID
    types.ToUint128(102),  // Escrow account
    types.ToUint128(101),  // Recipient
    types.ToUint128(5000), // Amount
    1, 1,
)

// Or cancel escrow
// err = txPatterns.CancelEscrow(
//     types.ToUint128(1001),
//     types.ToUint128(1000),
// )
```

#### Multi-Party Transfer

```go
// Split payment to multiple recipients atomically
multiPartyParams := transaction.MultiPartyTransferParams{
    BaseID:      types.ToUint128(2000),
    FromAccount: types.ToUint128(200),
    Recipients: []transaction.Recipient{
        {AccountID: types.ToUint128(201), Amount: types.ToUint128(1000)},
        {AccountID: types.ToUint128(202), Amount: types.ToUint128(2000)},
        {AccountID: types.ToUint128(203), Amount: types.ToUint128(1500)},
    },
    Ledger: 1,
    Code:   1,
}

if err := txPatterns.MultiPartyTransfer(multiPartyParams); err != nil {
    log.Fatal(err)
}
```

### Batch Operations

```go
import (
    "gitlab.cept.gov.in/it-2.0-common/ledgers/pkg/batch"
)

// Create batch processor
processor := batch.NewProcessor(c, 1000)

// Process thousands of accounts efficiently
accounts := []types.Account{...} // Your accounts
errors, err := processor.ProcessAccounts(accounts)

// Or use auto-flushing batcher
batcher := batch.NewTransferBatcher(c, 100)
for _, transfer := range transfers {
    batcher.Add(transfer) // Auto-flushes when batch is full
}
batcher.Flush() // Flush remaining
```

### Query Operations

```go
import (
    "gitlab.cept.gov.in/it-2.0-common/ledgers/pkg/query"
)

queryHelper := query.NewHelper(c)

// Get all debits for an account
debits, err := queryHelper.AccountDebits(accountID, 100)

// Get recent transfers
recentTransfers, err := queryHelper.AccountRecentTransfers(accountID, 10)

// Query accounts by ledger
accounts, err := queryHelper.AccountsByLedger(1, 100)

// Query transfers by code
transfers, err := queryHelper.TransfersByCode(10, 50)
```

## Package Structure

- **pkg/client**: Core client wrapper with connection management
- **pkg/account**: Account creation and management utilities
- **pkg/transfer**: Transfer builders and operations
- **pkg/transaction**: High-level transaction patterns (escrow, refund, etc.)
- **pkg/query**: Query helpers for accounts and transfers
- **pkg/batch**: Batch processing utilities
- **pkg/insurance**: Insurance-specific operations and domain models
- **examples**: Comprehensive examples for all features
- **examples/insurance**: Complete insurance management examples

## Transaction Patterns

### Supported Patterns

1. **Simple Transfer**: Immediate one-to-one transfer
2. **Pending Transfer**: Two-phase transfer (reserve → post/void)
3. **Escrow**: Secure payment with release or cancellation
4. **Multi-Party Transfer**: Split payments to multiple recipients
5. **Currency Exchange**: Cross-ledger transfers with intermediary
6. **Refund**: Reverse previous transactions
7. **Conditional Transfer**: Authorization with confirmation

### Linked Transfers

Linked transfers execute atomically - all succeed or all fail together:

```go
chain := batch.NewLinkedChain()
chain.AddTransfer(transfer1).
      AddTransfer(transfer2).
      AddTransfer(transfer3)

// All three execute atomically
results, err := chain.Execute(c)
```

## Best Practices

### Account Design

- Use **Ledger** to partition accounts that can transact together
- Use **Code** to categorize account types (user, system, escrow, etc.)
- Enable **History** flag if you need to query historical balances
- Use **DebitsMustNotExceedCredits** to prevent overdrafts

### Transfer Design

- Always use unique IDs for transfers (they're idempotent)
- Use **UserData** fields to link to external systems
- Use **Code** to categorize transfer types (payment, refund, fee, etc.)
- Batch transfers whenever possible for better performance

### Performance

- **Batch operations**: TigerBeetle is optimized for batching
- **Max batch size**: 8,189 operations per batch
- Use the batch utilities to automatically handle batching
- Consider using linked transfers for multi-step operations

### Error Handling

```go
// Check for individual operation errors
results, err := accountMgr.CreateBatch(accounts)
if err != nil {
    // Network or system error
    log.Fatal(err)
}

// Check for business logic errors
for _, result := range results {
    log.Printf("Account %d failed: %s", result.Index, result.Result)
}
```

## Examples

See the [examples](./examples) directory for complete, runnable examples:

### General Examples
- **basic_transfer.go**: Basic account creation and transfers
- **escrow.go**: Escrow pattern with release/cancel
- **multi_party.go**: Split payments to multiple recipients
- **batch_operations.go**: Efficient batch processing
- **query_operations.go**: Querying accounts and transfers

### Insurance Examples
Complete insurance management system examples in [examples/insurance](./examples/insurance):
- **premium_collection.go**: Premium collection through cash, payment gateway, bank transfer, auto-debit, cheque
- **policy_revival.go**: Revival of lapsed policies with outstanding premiums, penalties, and interest
- **claim_processing.go**: Death claims, maturity claims, partial withdrawals, and surrender values
- **complete_insurance_workflow.go**: End-to-end insurance operations

See the [Insurance README](./examples/insurance/README.md) for detailed documentation.

## Testing

Run the test suite:

```bash
cd go-tigerbeetle-lib
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## Requirements

- Go 1.21 or higher
- TigerBeetle server running (for examples and integration tests)

## Architecture

This library wraps the official [TigerBeetle Go client](https://github.com/tigerbeetle/tigerbeetle-go) and provides:

1. **Builder patterns** for easier object construction
2. **Manager classes** for common operations
3. **Transaction patterns** for complex workflows
4. **Query helpers** for simplified data retrieval
5. **Batch utilities** for high-performance operations

The library maintains full compatibility with the underlying TigerBeetle client and doesn't hide any functionality - you can always access the native client when needed.

## Contributing

Contributions are welcome! Please ensure:

- Code follows Go conventions and passes `go fmt`
- All tests pass
- New features include tests
- Documentation is updated

## License

This library is built on top of TigerBeetle, which is licensed under the Apache License 2.0.

## Resources

- [TigerBeetle Documentation](https://docs.tigerbeetle.com/)
- [TigerBeetle Go Client](https://github.com/tigerbeetle/tigerbeetle-go)
- [TigerBeetle Repository](https://github.com/tigerbeetle/tigerbeetle)

## Support

For issues specific to this library, please open an issue on GitHub.
For TigerBeetle core issues, see the [main repository](https://github.com/tigerbeetle/tigerbeetle).

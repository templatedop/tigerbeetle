/*
Package tigerbeetle provides a comprehensive, high-level Go library for TigerBeetle,
the distributed financial accounting database.

# Overview

This library wraps the official TigerBeetle Go client with idiomatic abstractions,
making it easier to build financial applications with complex transaction patterns.

# Quick Start

Create a client and start transacting:

	import (
		"gitlab.cept.gov.in/it-2.0-common/ledgers/pkg/client"
		"gitlab.cept.gov.in/it-2.0-common/ledgers/pkg/account"
		"gitlab.cept.gov.in/it-2.0-common/ledgers/pkg/transfer"
		"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
	)

	// Connect to TigerBeetle
	c, _ := client.New(client.Config{
		ClusterID:      types.ToUint128(0),
		ReplicaAddrs:   []string{"3000"},
		MaxConcurrency: 32,
	})
	defer c.Close()

	// Create an account
	accountMgr := account.NewManager(c)
	acc := account.New(types.ToUint128(1)).
		Ledger(1).
		Code(1).
		Build()
	accountMgr.Create(acc)

	// Make a transfer
	transferMgr := transfer.NewManager(c)
	transferMgr.SimpleTransfer(
		types.ToUint128(1),    // Transfer ID
		types.ToUint128(1),    // From
		types.ToUint128(2),    // To
		types.ToUint128(1000), // Amount
		1, 1,                  // Ledger, Code
	)

# Key Packages

  - client: Core TigerBeetle client wrapper
  - account: Account creation and management
  - transfer: Transfer builders and operations
  - transaction: High-level patterns (escrow, multi-party, etc.)
  - query: Query and filtering helpers
  - batch: Batch processing utilities

# Transaction Patterns

The library provides built-in support for common financial patterns:

Escrow Transactions:

	txPatterns := transaction.NewPatterns(c)
	txPatterns.Escrow(transaction.EscrowParams{
		EscrowID:      types.ToUint128(1),
		FromAccount:   types.ToUint128(100),
		EscrowAccount: types.ToUint128(200),
		Amount:        types.ToUint128(5000),
		Ledger:        1,
		Code:          1,
		EscrowTimeout: 24 * time.Hour,
	})

Multi-Party Transfers:

	txPatterns.MultiPartyTransfer(transaction.MultiPartyTransferParams{
		BaseID:      types.ToUint128(1),
		FromAccount: types.ToUint128(100),
		Recipients: []transaction.Recipient{
			{AccountID: types.ToUint128(101), Amount: types.ToUint128(1000)},
			{AccountID: types.ToUint128(102), Amount: types.ToUint128(2000)},
		},
		Ledger: 1,
		Code:   1,
	})

# Batch Operations

For high performance, use batch utilities:

	processor := batch.NewProcessor(c, 1000)
	errors, err := processor.ProcessAccounts(accounts)

Or use auto-flushing batchers:

	batcher := batch.NewTransferBatcher(c, 100)
	for _, transfer := range transfers {
		batcher.Add(transfer) // Auto-flushes at batch size
	}
	batcher.Flush()

# Query Operations

Simplified querying and filtering:

	queryHelper := query.NewHelper(c)

	// Get account debits
	debits, _ := queryHelper.AccountDebits(accountID, 100)

	// Get recent transfers
	recent, _ := queryHelper.AccountRecentTransfers(accountID, 10)

	// Query by ledger
	accounts, _ := queryHelper.AccountsByLedger(1, 100)

# Best Practices

Account Design:
  - Use Ledger to partition transactable accounts
  - Use Code for account categorization
  - Enable History for balance tracking
  - Use DebitsMustNotExceedCredits to prevent overdrafts

Transfer Design:
  - Use unique IDs (transfers are idempotent)
  - Use UserData fields for external references
  - Batch operations for performance
  - Use linked transfers for atomic multi-step operations

Performance:
  - Batch operations whenever possible
  - Maximum batch size: 8,189 operations
  - Use provided batch utilities
  - Consider linked transfers for atomicity

# Error Handling

The library provides structured error handling:

	results, err := accountMgr.CreateBatch(accounts)
	if err != nil {
		// System/network error
		log.Fatal(err)
	}

	// Check individual operation errors
	for _, result := range results {
		log.Printf("Failed: %s", result.Result)
	}

For more information, see https://docs.tigerbeetle.com/
*/
package tigerbeetle

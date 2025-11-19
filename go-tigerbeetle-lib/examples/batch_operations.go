package main

import (
	"fmt"
	"log"

	"gitlab.cept.gov.in/it-2.0-common/ledgers/pkg/account"
	"gitlab.cept.gov.in/it-2.0-common/ledgers/pkg/batch"
	"gitlab.cept.gov.in/it-2.0-common/ledgers/pkg/client"
	"gitlab.cept.gov.in/it-2.0-common/ledgers/pkg/transfer"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

// Example demonstrating efficient batch operations
func main() {
	// Create a TigerBeetle client
	c, err := client.New(client.Config{
		ClusterID:      types.ToUint128(0),
		ReplicaAddrs:   []string{"3000"},
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	// Create a batch processor
	processor := batch.NewProcessor(c, 1000)

	// Create many accounts efficiently
	var accounts []types.Account
	for i := 0; i < 10000; i++ {
		acc := account.New(types.ToUint128(uint64(i))).
			Ledger(1).
			Code(1).
			Build()
		accounts = append(accounts, acc)
	}

	fmt.Printf("Creating %d accounts in batches...\n", len(accounts))
	errors, err := processor.ProcessAccounts(accounts)
	if err != nil {
		log.Fatalf("Failed to process accounts: %v", err)
	}
	if len(errors) > 0 {
		fmt.Printf("Warning: %d accounts failed to create\n", len(errors))
	} else {
		fmt.Println("All accounts created successfully")
	}

	// Create many transfers efficiently
	var transfers []types.Transfer
	for i := 0; i < 5000; i++ {
		// Transfer between consecutive accounts
		xfer := transfer.New(types.ToUint128(uint64(100000 + i))).
			DebitAccount(types.ToUint128(uint64(i * 2))).
			CreditAccount(types.ToUint128(uint64(i*2 + 1))).
			Amount(types.ToUint128(100)).
			Ledger(1).
			Code(1).
			Build()
		transfers = append(transfers, xfer)
	}

	fmt.Printf("Creating %d transfers in batches...\n", len(transfers))
	transferErrors, err := processor.ProcessTransfers(transfers)
	if err != nil {
		log.Fatalf("Failed to process transfers: %v", err)
	}
	if len(transferErrors) > 0 {
		fmt.Printf("Warning: %d transfers failed to create\n", len(transferErrors))
	} else {
		fmt.Println("All transfers created successfully")
	}

	// Using auto-flushing batcher
	fmt.Println("\nDemonstrating auto-flushing batcher...")
	batcher := batch.NewTransferBatcher(c, 100)

	for i := 0; i < 250; i++ {
		xfer := transfer.New(types.ToUint128(uint64(200000 + i))).
			DebitAccount(types.ToUint128(uint64(i * 2))).
			CreditAccount(types.ToUint128(uint64(i*2 + 1))).
			Amount(types.ToUint128(50)).
			Ledger(1).
			Code(1).
			Build()

		// Batcher automatically flushes when batch is full
		if _, err := batcher.Add(xfer); err != nil {
			log.Printf("Failed to add transfer: %v", err)
		}
	}

	// Don't forget to flush remaining items
	if _, err := batcher.Flush(); err != nil {
		log.Printf("Failed to flush remaining transfers: %v", err)
	}

	fmt.Printf("Processed %d transfers with auto-flushing batcher\n", 250)
}

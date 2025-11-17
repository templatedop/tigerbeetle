package main

import (
	"fmt"
	"log"

	"github.com/tigerbeetle/go-tigerbeetle-lib/pkg/account"
	"github.com/tigerbeetle/go-tigerbeetle-lib/pkg/client"
	"github.com/tigerbeetle/go-tigerbeetle-lib/pkg/query"
	"github.com/tigerbeetle/go-tigerbeetle-lib/pkg/transfer"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

// Example demonstrating query operations
func main() {
	// Create a TigerBeetle client
	c, err := client.New(client.Config{
		ClusterID:      types.ToUint128(0),
		ReplicaAddrs:   []string{"3000"},
		MaxConcurrency: 32,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	// Setup: Create some accounts and transfers
	accountMgr := account.NewManager(c)
	transferMgr := transfer.NewManager(c)

	// Create accounts
	accounts := []types.Account{
		account.New(types.ToUint128(300)).Ledger(1).Code(1).HistoryEnabled().Build(),
		account.New(types.ToUint128(301)).Ledger(1).Code(2).HistoryEnabled().Build(),
		account.New(types.ToUint128(302)).Ledger(2).Code(1).HistoryEnabled().Build(),
	}
	accountMgr.CreateBatch(accounts)

	// Create some transfers
	transferMgr.SimpleTransfer(
		types.ToUint128(3000),
		types.ToUint128(300),
		types.ToUint128(301),
		types.ToUint128(1000),
		1, 1,
	)

	// Create query helper
	queryHelper := query.NewHelper(c)

	// Query accounts by ledger
	fmt.Println("Accounts in ledger 1:")
	accountsInLedger1, err := queryHelper.AccountsByLedger(1, 100)
	if err != nil {
		log.Fatalf("Failed to query accounts: %v", err)
	}
	for _, acc := range accountsInLedger1 {
		fmt.Printf("  Account ID: %d, Code: %d\n", acc.ID.ToUint128(), acc.Code)
	}

	// Query accounts by code
	fmt.Println("\nAccounts with code 1:")
	accountsByCode, err := queryHelper.AccountsByCode(1, 100)
	if err != nil {
		log.Fatalf("Failed to query accounts by code: %v", err)
	}
	for _, acc := range accountsByCode {
		fmt.Printf("  Account ID: %d, Ledger: %d\n", acc.ID.ToUint128(), acc.Ledger)
	}

	// Get account debits and credits
	fmt.Println("\nDebits for account 300:")
	debits, err := queryHelper.AccountDebits(types.ToUint128(300), 10)
	if err != nil {
		log.Fatalf("Failed to get debits: %v", err)
	}
	for _, transfer := range debits {
		fmt.Printf("  Transfer ID: %d, Amount: %d\n",
			transfer.ID.ToUint128(), transfer.Amount.ToUint128())
	}

	fmt.Println("\nCredits for account 301:")
	credits, err := queryHelper.AccountCredits(types.ToUint128(301), 10)
	if err != nil {
		log.Fatalf("Failed to get credits: %v", err)
	}
	for _, transfer := range credits {
		fmt.Printf("  Transfer ID: %d, Amount: %d\n",
			transfer.ID.ToUint128(), transfer.Amount.ToUint128())
	}

	// Get recent transfers for an account
	fmt.Println("\nRecent transfers for account 300:")
	recentTransfers, err := queryHelper.AccountRecentTransfers(types.ToUint128(300), 5)
	if err != nil {
		log.Fatalf("Failed to get recent transfers: %v", err)
	}
	for _, transfer := range recentTransfers {
		fmt.Printf("  Transfer ID: %d, Amount: %d, From: %d, To: %d\n",
			transfer.ID.ToUint128(),
			transfer.Amount.ToUint128(),
			transfer.DebitAccountID.ToUint128(),
			transfer.CreditAccountID.ToUint128())
	}

	// Query transfers by ledger
	fmt.Println("\nTransfers in ledger 1:")
	transfersInLedger, err := queryHelper.TransfersByLedger(1, 100)
	if err != nil {
		log.Fatalf("Failed to query transfers: %v", err)
	}
	for _, transfer := range transfersInLedger {
		fmt.Printf("  Transfer ID: %d, Amount: %d\n",
			transfer.ID.ToUint128(), transfer.Amount.ToUint128())
	}
}

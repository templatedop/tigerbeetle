package main

import (
	"fmt"
	"log"

	"gitlab.cept.gov.in/it-2.0-common/ledgers/pkg/account"
	"gitlab.cept.gov.in/it-2.0-common/ledgers/pkg/client"
	"gitlab.cept.gov.in/it-2.0-common/ledgers/pkg/transfer"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

// Example demonstrating basic account creation and transfers
func main() {
	// Create a TigerBeetle client
	c, err := client.New(client.Config{
		ClusterID:    types.ToUint128(0),
		ReplicaAddrs: []string{"3000"},
		MaxConcurrency: 32,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	// Create account manager
	accountMgr := account.NewManager(c)

	// Create two accounts
	account1 := account.New(types.ToUint128(1)).
		Ledger(1).
		Code(1).
		HistoryEnabled().
		Build()

	account2 := account.New(types.ToUint128(2)).
		Ledger(1).
		Code(1).
		HistoryEnabled().
		Build()

	// Create accounts in batch
	if err := accountMgr.Create(account1); err != nil {
		log.Fatalf("Failed to create account 1: %v", err)
	}
	if err := accountMgr.Create(account2); err != nil {
		log.Fatalf("Failed to create account 2: %v", err)
	}

	fmt.Println("Created accounts successfully")

	// Create transfer manager
	transferMgr := transfer.NewManager(c)

	// Perform a simple transfer of 1000 from account 1 to account 2
	err = transferMgr.SimpleTransfer(
		types.ToUint128(1), // Transfer ID
		types.ToUint128(1), // From account
		types.ToUint128(2), // To account
		types.ToUint128(1000), // Amount
		1, // Ledger
		1, // Code
	)
	if err != nil {
		log.Fatalf("Failed to create transfer: %v", err)
	}

	fmt.Println("Transfer completed successfully")

	// Check balances
	balance1, err := accountMgr.GetBalance(types.ToUint128(1))
	if err != nil {
		log.Fatalf("Failed to get balance for account 1: %v", err)
	}

	balance2, err := accountMgr.GetBalance(types.ToUint128(2))
	if err != nil {
		log.Fatalf("Failed to get balance for account 2: %v", err)
	}

	fmt.Printf("Account 1 - Debits: %d, Credits: %d\n",
		balance1.Account.DebitsPosted, balance1.Account.CreditsPosted)
	fmt.Printf("Account 2 - Debits: %d, Credits: %d\n",
		balance2.Account.DebitsPosted, balance2.Account.CreditsPosted)
}

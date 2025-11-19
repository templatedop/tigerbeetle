package main

import (
	"fmt"
	"log"
	"time"

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
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	// Generate unique IDs based on current timestamp to avoid conflicts
	baseID := uint64(time.Now().UnixNano())

	fmt.Printf("Using base ID: %d (timestamp-based for uniqueness)\n", baseID)

	// Create account manager
	accountMgr := account.NewManager(c)

	// Create two accounts with unique IDs
	account1 := account.New(types.ToUint128(baseID)).
		Ledger(1).
		Code(1).
		HistoryEnabled().
		Build()

	account2 := account.New(types.ToUint128(baseID + 1)).
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
		types.ToUint128(baseID+1000), // Transfer ID (unique)
		types.ToUint128(baseID),      // From account
		types.ToUint128(baseID+1),    // To account
		types.ToUint128(1000),         // Amount
		1, // Ledger
		1, // Code
	)
	if err != nil {
		log.Fatalf("Failed to create transfer: %v", err)
	}

	fmt.Println("Transfer completed successfully")

	// Check balances
	balance1, err := accountMgr.GetBalance(types.ToUint128(baseID))
	if err != nil {
		log.Fatalf("Failed to get balance for account 1: %v", err)
	}

	balance2, err := accountMgr.GetBalance(types.ToUint128(baseID + 1))
	if err != nil {
		log.Fatalf("Failed to get balance for account 2: %v", err)
	}

	fmt.Printf("Account 1 - Debits: %v, Credits: %v\n",
		balance1.Account.DebitsPosted, balance1.Account.CreditsPosted)
	fmt.Printf("Account 2 - Debits: %v, Credits: %v\n",
		balance2.Account.DebitsPosted, balance2.Account.CreditsPosted)
}

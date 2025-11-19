package main

import (
	"fmt"
	"log"

	"gitlab.cept.gov.in/it-2.0-common/ledgers/pkg/account"
	"gitlab.cept.gov.in/it-2.0-common/ledgers/pkg/client"
	"gitlab.cept.gov.in/it-2.0-common/ledgers/pkg/transaction"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

// Example demonstrating multi-party transfers (splitting payments)
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

	// Create account manager
	accountMgr := account.NewManager(c)

	// Create payer and multiple recipient accounts
	accounts := []types.Account{
		account.New(types.ToUint128(200)).Ledger(1).Code(1).Build(), // Payer
		account.New(types.ToUint128(201)).Ledger(1).Code(1).Build(), // Recipient 1
		account.New(types.ToUint128(202)).Ledger(1).Code(1).Build(), // Recipient 2
		account.New(types.ToUint128(203)).Ledger(1).Code(1).Build(), // Recipient 3
	}

	if _, err := accountMgr.CreateBatch(accounts); err != nil {
		log.Fatalf("Failed to create accounts: %v", err)
	}

	fmt.Println("Created payer and recipient accounts")

	// Create transaction patterns helper
	txPatterns := transaction.NewPatterns(c)

	// Split payment: Payer sends to 3 recipients atomically
	// For example: splitting revenue, paying multiple vendors, etc.
	multiPartyParams := transaction.MultiPartyTransferParams{
		BaseID:      types.ToUint128(2000),
		FromAccount: types.ToUint128(200), // Payer
		Recipients: []transaction.Recipient{
			{AccountID: types.ToUint128(201), Amount: types.ToUint128(1000)}, // 1000 to recipient 1
			{AccountID: types.ToUint128(202), Amount: types.ToUint128(2000)}, // 2000 to recipient 2
			{AccountID: types.ToUint128(203), Amount: types.ToUint128(1500)}, // 1500 to recipient 3
		},
		Ledger: 1,
		Code:   1,
	}

	if err := txPatterns.MultiPartyTransfer(multiPartyParams); err != nil {
		log.Fatalf("Failed to create multi-party transfer: %v", err)
	}

	fmt.Println("Multi-party transfer completed successfully")

	// Check balances
	for i := 200; i <= 203; i++ {
		balance, _ := accountMgr.GetBalance(types.ToUint128(uint64(i)))
		fmt.Printf("Account %d - Debits: %d, Credits: %d\n",
			i, balance.Account.DebitsPosted, balance.Account.CreditsPosted)
	}
}

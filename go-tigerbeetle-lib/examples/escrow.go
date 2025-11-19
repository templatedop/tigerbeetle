package main

import (
	"fmt"
	"log"
	"time"

	"gitlab.cept.gov.in/it-2.0-common/ledgers/pkg/account"
	"gitlab.cept.gov.in/it-2.0-common/ledgers/pkg/client"
	"gitlab.cept.gov.in/it-2.0-common/ledgers/pkg/transaction"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

// Example demonstrating escrow pattern for secure transactions
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

	// Create account manager
	accountMgr := account.NewManager(c)

	// Create buyer, seller, and escrow accounts
	buyer := account.New(types.ToUint128(100)).
		Ledger(1).
		Code(1). // User account
		Build()

	seller := account.New(types.ToUint128(101)).
		Ledger(1).
		Code(1). // User account
		Build()

	escrowAccount := account.New(types.ToUint128(102)).
		Ledger(1).
		Code(2). // Escrow account
		Build()

	// Create all accounts
	if err := accountMgr.CreateBatch([]types.Account{buyer, seller, escrowAccount}); err != nil {
		log.Fatalf("Failed to create accounts: %v", err)
	}

	fmt.Println("Created buyer, seller, and escrow accounts")

	// Create transaction patterns helper
	txPatterns := transaction.NewPatterns(c)

	// Buyer places funds in escrow (e.g., for a purchase)
	escrowParams := transaction.EscrowParams{
		EscrowID:      types.ToUint128(1000),
		FromAccount:   types.ToUint128(100), // Buyer
		EscrowAccount: types.ToUint128(102), // Escrow
		Amount:        types.ToUint128(5000),
		Ledger:        1,
		Code:          1,
		EscrowTimeout: 24 * time.Hour, // 24 hour timeout
	}

	if err := txPatterns.Escrow(escrowParams); err != nil {
		log.Fatalf("Failed to create escrow: %v", err)
	}

	fmt.Println("Funds placed in escrow")

	// Simulate successful transaction - release funds to seller
	if err := txPatterns.ReleaseEscrow(
		types.ToUint128(1001), // Release transfer ID
		types.ToUint128(1000), // Escrow ID
		types.ToUint128(102),  // Escrow account
		types.ToUint128(101),  // Seller account
		types.ToUint128(5000), // Amount
		1, // Ledger
		1, // Code
	); err != nil {
		log.Fatalf("Failed to release escrow: %v", err)
	}

	fmt.Println("Escrow released to seller")

	// Check final balances
	buyerBalance, _ := accountMgr.GetBalance(types.ToUint128(100))
	sellerBalance, _ := accountMgr.GetBalance(types.ToUint128(101))

	fmt.Printf("Buyer - Debits: %d, Credits: %d\n",
		buyerBalance.Account.DebitsPosted, buyerBalance.Account.CreditsPosted)
	fmt.Printf("Seller - Debits: %d, Credits: %d\n",
		sellerBalance.Account.DebitsPosted, sellerBalance.Account.CreditsPosted)
}

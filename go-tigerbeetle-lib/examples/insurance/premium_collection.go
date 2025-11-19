package main

import (
	"fmt"
	"log"
	"time"

	"gitlab.cept.gov.in/it-2.0-common/ledgers/pkg/account"
	"gitlab.cept.gov.in/it-2.0-common/ledgers/pkg/client"
	"gitlab.cept.gov.in/it-2.0-common/ledgers/pkg/insurance"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

// Example demonstrating premium collection through different payment methods
// Each payment method is represented as a separate account:
// - Cash Collection Account
// - Payment Gateway Account
// - Bank Transfer Account
// - Auto-Debit Account
func main() {
	// Create TigerBeetle client
	c, err := client.New(client.Config{
		ClusterID:      types.ToUint128(0),
		ReplicaAddrs:   []string{"3000"},
		MaxConcurrency: 32,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	fmt.Println("=== Insurance Premium Collection Demo ===\n")

	// Setup: Create payment collection accounts and company accounts
	fmt.Println("Setting up payment collection accounts...")
	accountMgr := account.NewManager(c)

	paymentAccounts := []types.Account{
		// Payment collection accounts
		account.New(types.ToUint128(uint64(insurance.AccountCodeCashCollection))).
			Ledger(insurance.LedgerINR).
			Code(insurance.AccountCodeCashCollection).
			Build(),
		account.New(types.ToUint128(uint64(insurance.AccountCodePaymentGateway))).
			Ledger(insurance.LedgerINR).
			Code(insurance.AccountCodePaymentGateway).
			Build(),
		account.New(types.ToUint128(uint64(insurance.AccountCodeBankTransfer))).
			Ledger(insurance.LedgerINR).
			Code(insurance.AccountCodeBankTransfer).
			Build(),
		account.New(types.ToUint128(uint64(insurance.AccountCodeAutoDebit))).
			Ledger(insurance.LedgerINR).
			Code(insurance.AccountCodeAutoDebit).
			Build(),
		account.New(types.ToUint128(uint64(insurance.AccountCodeCheque))).
			Ledger(insurance.LedgerINR).
			Code(insurance.AccountCodeCheque).
			Build(),

		// Premium income account
		account.New(types.ToUint128(uint64(insurance.AccountCodePremiumIncome))).
			Ledger(insurance.LedgerINR).
			Code(insurance.AccountCodePremiumIncome).
			Build(),
	}

	if _, err := accountMgr.CreateBatch(paymentAccounts); err != nil {
		log.Fatalf("Failed to create payment accounts: %v", err)
	}
	fmt.Println("✓ Payment collection accounts created\n")

	// Create sample policy accounts
	fmt.Println("Creating policy accounts...")
	policyIDs := []types.Uint128{
		types.ToUint128(100001), // Policy #1
		types.ToUint128(100002), // Policy #2
		types.ToUint128(100003), // Policy #3
		types.ToUint128(100004), // Policy #4
		types.ToUint128(100005), // Policy #5
	}

	var policies []types.Account
	for i, policyID := range policyIDs {
		policy := account.New(policyID).
			Ledger(insurance.LedgerINR).
			Code(insurance.AccountCodePolicyActive).
			HistoryEnabled().
			UserData64(uint64(i + 1)). // Policy number
			Build()
		policies = append(policies, policy)
	}

	if _, err := accountMgr.CreateBatch(policies); err != nil {
		log.Fatalf("Failed to create policies: %v", err)
	}
	fmt.Printf("✓ Created %d policy accounts\n\n", len(policies))

	// Initialize insurance operations
	insuranceOps := insurance.NewOperations(c, insurance.LedgerINR)

	// Example 1: Premium collection via CASH
	fmt.Println("Example 1: Collecting premium via CASH")
	fmt.Println("----------------------------------------")
	if err := insuranceOps.CollectPremium(insurance.PremiumCollectionParams{
		TransferID:      types.ToUint128(1000),
		PolicyAccountID: policyIDs[0],
		Amount:          types.ToUint128(5000), // ₹5,000
		PaymentMethod:   insurance.PaymentMethodCash,
		IsFirstPremium:  false,
		DueDate:         time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		PaymentDate:     time.Date(2024, 1, 5, 10, 30, 0, 0, time.UTC),
	}); err != nil {
		log.Fatalf("Failed: %v", err)
	}
	fmt.Println("✓ Policy #1: ₹5,000 premium collected via CASH")
	fmt.Println("  Payment received at branch counter")
	fmt.Println("  Flow: Cash Account -> Premium Income -> Policy Account\n")

	// Example 2: Premium collection via PAYMENT GATEWAY (Online)
	fmt.Println("Example 2: Collecting premium via PAYMENT GATEWAY")
	fmt.Println("--------------------------------------------------")
	if err := insuranceOps.CollectPremium(insurance.PremiumCollectionParams{
		TransferID:      types.ToUint128(1002),
		PolicyAccountID: policyIDs[1],
		Amount:          types.ToUint128(10000), // ₹10,000
		PaymentMethod:   insurance.PaymentMethodPaymentGateway,
		IsFirstPremium:  false,
		DueDate:         time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		PaymentDate:     time.Date(2024, 1, 1, 14, 25, 0, 0, time.UTC),
	}); err != nil {
		log.Fatalf("Failed: %v", err)
	}
	fmt.Println("✓ Policy #2: ₹10,000 premium collected via PAYMENT GATEWAY")
	fmt.Println("  Online payment through website/app")
	fmt.Println("  Flow: Payment Gateway -> Premium Income -> Policy Account\n")

	// Example 3: Premium collection via BANK TRANSFER
	fmt.Println("Example 3: Collecting premium via BANK TRANSFER")
	fmt.Println("------------------------------------------------")
	if err := insuranceOps.CollectPremium(insurance.PremiumCollectionParams{
		TransferID:      types.ToUint128(1004),
		PolicyAccountID: policyIDs[2],
		Amount:          types.ToUint128(7500), // ₹7,500
		PaymentMethod:   insurance.PaymentMethodBankTransfer,
		IsFirstPremium:  false,
		DueDate:         time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		PaymentDate:     time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
	}); err != nil {
		log.Fatalf("Failed: %v", err)
	}
	fmt.Println("✓ Policy #3: ₹7,500 premium collected via BANK TRANSFER")
	fmt.Println("  Direct bank transfer/NEFT/RTGS")
	fmt.Println("  Flow: Bank Transfer Account -> Premium Income -> Policy Account\n")

	// Example 4: Premium collection via AUTO-DEBIT (Recurring)
	fmt.Println("Example 4: Collecting premium via AUTO-DEBIT")
	fmt.Println("---------------------------------------------")
	if err := insuranceOps.CollectPremium(insurance.PremiumCollectionParams{
		TransferID:      types.ToUint128(1006),
		PolicyAccountID: policyIDs[3],
		Amount:          types.ToUint128(12000), // ₹12,000
		PaymentMethod:   insurance.PaymentMethodAutoDebit,
		IsFirstPremium:  false,
		DueDate:         time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		PaymentDate:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}); err != nil {
		log.Fatalf("Failed: %v", err)
	}
	fmt.Println("✓ Policy #4: ₹12,000 premium collected via AUTO-DEBIT")
	fmt.Println("  Automatic deduction from customer's account")
	fmt.Println("  Flow: Auto-Debit Account -> Premium Income -> Policy Account\n")

	// Example 5: Premium collection via CHEQUE
	fmt.Println("Example 5: Collecting premium via CHEQUE")
	fmt.Println("-----------------------------------------")
	if err := insuranceOps.CollectPremium(insurance.PremiumCollectionParams{
		TransferID:      types.ToUint128(1008),
		PolicyAccountID: policyIDs[4],
		Amount:          types.ToUint128(8500), // ₹8,500
		PaymentMethod:   insurance.PaymentMethodCheque,
		IsFirstPremium:  false,
		DueDate:         time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		PaymentDate:     time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC),
	}); err != nil {
		log.Fatalf("Failed: %v", err)
	}
	fmt.Println("✓ Policy #5: ₹8,500 premium collected via CHEQUE")
	fmt.Println("  Cheque deposited and cleared")
	fmt.Println("  Flow: Cheque Account -> Premium Income -> Policy Account\n")

	// Example 6: Batch Premium Collection (Month-end processing)
	fmt.Println("Example 6: Batch Premium Collection (Auto-Debit for multiple policies)")
	fmt.Println("-----------------------------------------------------------------------")

	premiums := []insurance.PolicyPremium{
		{
			PolicyAccountID: policyIDs[0],
			Amount:          types.ToUint128(5000),
			PaymentMethod:   insurance.PaymentMethodAutoDebit,
			DueDate:         time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			PaymentDate:     time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			PolicyAccountID: policyIDs[1],
			Amount:          types.ToUint128(10000),
			PaymentMethod:   insurance.PaymentMethodAutoDebit,
			DueDate:         time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			PaymentDate:     time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			PolicyAccountID: policyIDs[2],
			Amount:          types.ToUint128(7500),
			PaymentMethod:   insurance.PaymentMethodAutoDebit,
			DueDate:         time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			PaymentDate:     time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			PolicyAccountID: policyIDs[3],
			Amount:          types.ToUint128(12000),
			PaymentMethod:   insurance.PaymentMethodAutoDebit,
			DueDate:         time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			PaymentDate:     time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			PolicyAccountID: policyIDs[4],
			Amount:          types.ToUint128(8500),
			PaymentMethod:   insurance.PaymentMethodAutoDebit,
			DueDate:         time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			PaymentDate:     time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	errors, err := insuranceOps.CollectMonthlyPremiums(insurance.MonthlyPremiumCollectionParams{
		BaseTransferID: types.ToUint128(2000),
		Premiums:       premiums,
		Month:          time.February,
		Year:           2024,
	})
	if err != nil {
		log.Fatalf("Failed to collect batch premiums: %v", err)
	}

	if len(errors) > 0 {
		fmt.Printf("⚠ %d premiums failed to collect\n", len(errors))
		for _, e := range errors {
			fmt.Printf("  - Failed at index %d: %s\n", e.Index, e.Result.String())
		}
	} else {
		totalAmount := uint64(5000 + 10000 + 7500 + 12000 + 8500)
		fmt.Printf("✓ Successfully collected premiums for %d policies\n", len(premiums))
		fmt.Printf("  Total Amount: ₹%d\n", totalAmount)
		fmt.Println("  All auto-debits processed atomically\n")
	}

	// Show final balances
	fmt.Println("=== Final Policy Balances ===")
	for i, policyID := range policyIDs {
		balance, err := accountMgr.GetBalance(policyID)
		if err != nil {
			fmt.Printf("Policy #%d: Error retrieving balance\n", i+1)
			continue
		}
		fmt.Printf("Policy #%d: ₹%d collected\n",
			i+1,
			balance.Account.CreditsPosted,
		)
	}

	fmt.Println("\n✓ Premium collection demo completed successfully!")
}

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

// Complete insurance workflow demonstrating:
// 1. Setting up payment collection accounts
// 2. Creating policy accounts
// 3. Premium collection through various payment methods
// 4. Policy revival after lapsing
// 5. Claim processing
func main() {
	// Create TigerBeetle client
	c, err := client.New(client.Config{
		ClusterID:      types.ToUint128(0),
		ReplicaAddrs:   []string{"3000"},
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	fmt.Println("=== Insurance Management System Demo ===\n")

	// Generate unique IDs based on timestamp to avoid conflicts with previous runs
	baseID := uint64(time.Now().UnixNano())
	fmt.Printf("Using base ID: %d (timestamp-based for uniqueness)\n\n", baseID)

	// Step 1: Setup System Accounts
	fmt.Println("Step 1: Setting up system accounts...")
	if err := setupSystemAccounts(c, baseID); err != nil {
		log.Fatalf("Failed to setup system accounts: %v", err)
	}
	fmt.Println("✓ System accounts created\n")

	// Step 2: Create Policy Accounts
	fmt.Println("Step 2: Creating policy accounts...")
	policy1ID := types.ToUint128(baseID + 100001) // Policy #1
	policy2ID := types.ToUint128(baseID + 100002) // Policy #2
	policy3ID := types.ToUint128(baseID + 100003) // Policy #3

	if err := createPolicyAccounts(c, []types.Uint128{policy1ID, policy2ID, policy3ID}); err != nil {
		log.Fatalf("Failed to create policy accounts: %v", err)
	}
	fmt.Println("✓ Created 3 policy accounts\n")

	// Step 3: Collect Premiums through Different Payment Methods
	fmt.Println("Step 3: Collecting monthly premiums...")
	insuranceOps := insurance.NewOperations(c, insurance.LedgerINR)

	// Policy 1: Premium via Cash
	if err := insuranceOps.CollectPremium(insurance.PremiumCollectionParams{
		TransferID:      1000,
		PolicyAccountID: policy1ID,
		Amount:          types.ToUint128(5000), // ₹5,000
		PaymentMethod:   insurance.PaymentMethodCash,
		IsFirstPremium:  false,
		DueDate:         time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		PaymentDate:     time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
	}); err != nil {
		log.Fatalf("Failed to collect premium via cash: %v", err)
	}
	fmt.Println("  ✓ Policy #1: ₹5,000 collected via CASH")

	// Policy 2: Premium via Payment Gateway
	if err := insuranceOps.CollectPremium(insurance.PremiumCollectionParams{
		TransferID:      1002,
		PolicyAccountID: policy2ID,
		Amount:          types.ToUint128(10000), // ₹10,000
		PaymentMethod:   insurance.PaymentMethodPaymentGateway,
		IsFirstPremium:  false,
		DueDate:         time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		PaymentDate:     time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
	}); err != nil {
		log.Fatalf("Failed to collect premium via payment gateway: %v", err)
	}
	fmt.Println("  ✓ Policy #2: ₹10,000 collected via PAYMENT GATEWAY")

	// Policy 3: Premium via Bank Transfer
	if err := insuranceOps.CollectPremium(insurance.PremiumCollectionParams{
		TransferID:      1004,
		PolicyAccountID: policy3ID,
		Amount:          types.ToUint128(7500), // ₹7,500
		PaymentMethod:   insurance.PaymentMethodBankTransfer,
		IsFirstPremium:  false,
		DueDate:         time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		PaymentDate:     time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
	}); err != nil {
		log.Fatalf("Failed to collect premium via bank transfer: %v", err)
	}
	fmt.Println("  ✓ Policy #3: ₹7,500 collected via BANK TRANSFER\n")

	// Step 4: Check Policy Balances
	fmt.Println("Step 4: Policy balances after premium collection...")
	showPolicyBalance(c, policy1ID, "Policy #1")
	showPolicyBalance(c, policy2ID, "Policy #2")
	showPolicyBalance(c, policy3ID, "Policy #3")
	fmt.Println()

	// Step 5: Policy Revival (Policy lapsed, now being revived)
	fmt.Println("Step 5: Reviving a lapsed policy...")
	lapsedPolicyID := types.ToUint128(baseID + 100004)

	// Create lapsed policy account
	accountMgr := account.NewManager(c)
	lapsedPolicy := account.New(lapsedPolicyID).
		Ledger(insurance.LedgerINR).
		Code(insurance.AccountCodePolicyLapsed).
		HistoryEnabled().
		Build()
	if err := accountMgr.Create(lapsedPolicy); err != nil {
		log.Fatalf("Failed to create lapsed policy: %v", err)
	}

	// Revive the policy with outstanding premium, penalty, and interest
	if err := insuranceOps.RevivePolicy(insurance.RevivalParams{
		BaseTransferID:     2000,
		PolicyAccountID:    lapsedPolicyID,
		OutstandingPremium: types.ToUint128(15000), // ₹15,000 outstanding
		PenaltyAmount:      types.ToUint128(1500),  // ₹1,500 penalty
		InterestAmount:     types.ToUint128(750),   // ₹750 interest
		PaymentMethod:      insurance.PaymentMethodPaymentGateway,
		LapsedDate:         time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
		RevivalDate:        time.Now(),
	}); err != nil {
		log.Fatalf("Failed to revive policy: %v", err)
	}
	fmt.Println("  ✓ Policy #4 revived")
	fmt.Printf("    - Outstanding Premium: ₹15,000\n")
	fmt.Printf("    - Penalty: ₹1,500\n")
	fmt.Printf("    - Interest: ₹750\n")
	fmt.Printf("    - Total Paid: ₹17,250\n\n")

	showPolicyBalance(c, lapsedPolicyID, "Policy #4 (Revived)")
	fmt.Println()

	// Step 6: Process a Maturity Claim
	fmt.Println("Step 6: Processing maturity claim...")
	maturedPolicyID := types.ToUint128(baseID + 100005)

	// Create matured policy with accumulated premiums
	maturedPolicy := account.New(maturedPolicyID).
		Ledger(insurance.LedgerINR).
		Code(insurance.AccountCodePolicyMatured).
		HistoryEnabled().
		Build()
	if err := accountMgr.Create(maturedPolicy); err != nil {
		log.Fatalf("Failed to create matured policy: %v", err)
	}

	// First, add some accumulated value to the policy
	if err := insuranceOps.CollectPremium(insurance.PremiumCollectionParams{
		TransferID:      2500,
		PolicyAccountID: maturedPolicyID,
		Amount:          types.ToUint128(100000), // ₹1,00,000 accumulated
		PaymentMethod:   insurance.PaymentMethodBankTransfer,
		IsFirstPremium:  true,
		DueDate:         time.Now(),
		PaymentDate:     time.Now(),
	}); err != nil {
		log.Fatalf("Failed to setup matured policy: %v", err)
	}

	// Process maturity claim
	if err := insuranceOps.ProcessClaim(insurance.ClaimPaymentParams{
		BaseTransferID:  3000,
		PolicyAccountID: maturedPolicyID,
		ClaimAmount:     types.ToUint128(100000), // ₹1,00,000 claim
		ClaimType:       insurance.TransferCodeMaturityClaim,
		PaymentMethod:   insurance.PaymentMethodBankTransfer,
		ClaimDate:       time.Now(),
	}); err != nil {
		log.Fatalf("Failed to process claim: %v", err)
	}
	fmt.Println("  ✓ Maturity claim processed")
	fmt.Printf("    - Claim Amount: ₹1,00,000\n")
	fmt.Printf("    - Payment Method: BANK TRANSFER\n\n")

	// Step 7: Batch Premium Collection (Multiple policies at once)
	fmt.Println("Step 7: Batch premium collection for the month...")

	premiums := []insurance.PolicyPremium{
		{
			PolicyAccountID: policy1ID,
			Amount:          types.ToUint128(5000),
			PaymentMethod:   insurance.PaymentMethodAutoDebit,
			DueDate:         time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			PaymentDate:     time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			PolicyAccountID: policy2ID,
			Amount:          types.ToUint128(10000),
			PaymentMethod:   insurance.PaymentMethodAutoDebit,
			DueDate:         time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			PaymentDate:     time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			PolicyAccountID: policy3ID,
			Amount:          types.ToUint128(7500),
			PaymentMethod:   insurance.PaymentMethodAutoDebit,
			DueDate:         time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			PaymentDate:     time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	errors, err := insuranceOps.CollectMonthlyPremiums(insurance.MonthlyPremiumCollectionParams{
		BaseTransferID: 4000,
		Premiums:       premiums,
		Month:          time.February,
		Year:           2024,
	})
	if err != nil {
		log.Fatalf("Failed to collect monthly premiums: %v", err)
	}
	if len(errors) > 0 {
		fmt.Printf("  ⚠ %d premiums failed to collect\n", len(errors))
	} else {
		fmt.Printf("  ✓ Successfully collected premiums for %d policies\n", len(premiums))
		fmt.Printf("    - Total Amount: ₹%d\n\n", 5000+10000+7500)
	}

	// Step 8: Pending Premium (Payment Authorization)
	fmt.Println("Step 8: Creating pending premium (payment authorization)...")
	pendingTransferID := types.ToUint128(5000)

	if err := insuranceOps.CreatePendingPremium(insurance.PendingPremiumParams{
		TransferID:      pendingTransferID,
		PolicyAccountID: policy1ID,
		Amount:          types.ToUint128(5000),
		PaymentMethod:   insurance.PaymentMethodPaymentGateway,
		Timeout:         24 * time.Hour, // 24 hour authorization
	}); err != nil {
		log.Fatalf("Failed to create pending premium: %v", err)
	}
	fmt.Println("  ✓ Payment authorization created (₹5,000 on hold)")

	// Confirm the pending premium
	if err := insuranceOps.ConfirmPendingPremium(
		types.ToUint128(5001),
		pendingTransferID,
	); err != nil {
		log.Fatalf("Failed to confirm pending premium: %v", err)
	}
	fmt.Println("  ✓ Payment authorization confirmed\n")

	// Final Summary
	fmt.Println("=== Summary ===")
	fmt.Println("✓ System accounts setup complete")
	fmt.Println("✓ Policy accounts created and managed")
	fmt.Println("✓ Premiums collected via multiple payment methods")
	fmt.Println("✓ Lapsed policy successfully revived")
	fmt.Println("✓ Maturity claim processed and paid")
	fmt.Println("✓ Batch premium collection completed")
	fmt.Println("✓ Payment authorization flow demonstrated")
	fmt.Println("\nInsurance management system demo completed successfully!")
}

// setupSystemAccounts creates all required system accounts
func setupSystemAccounts(c *client.Client, baseID uint64) error {
	accountMgr := account.NewManager(c)

	systemAccounts := []types.Account{
		// Payment collection accounts
		account.New(types.ToUint128(baseID + uint64(insurance.AccountCodeCashCollection))).
			Ledger(insurance.LedgerINR).
			Code(insurance.AccountCodeCashCollection).
			Build(),
		account.New(types.ToUint128(baseID + uint64(insurance.AccountCodePaymentGateway))).
			Ledger(insurance.LedgerINR).
			Code(insurance.AccountCodePaymentGateway).
			Build(),
		account.New(types.ToUint128(baseID + uint64(insurance.AccountCodeBankTransfer))).
			Ledger(insurance.LedgerINR).
			Code(insurance.AccountCodeBankTransfer).
			Build(),
		account.New(types.ToUint128(baseID + uint64(insurance.AccountCodeAutoDebit))).
			Ledger(insurance.LedgerINR).
			Code(insurance.AccountCodeAutoDebit).
			Build(),

		// Company accounts
		account.New(types.ToUint128(baseID + uint64(insurance.AccountCodePremiumIncome))).
			Ledger(insurance.LedgerINR).
			Code(insurance.AccountCodePremiumIncome).
			Build(),
		account.New(types.ToUint128(baseID + uint64(insurance.AccountCodeRevivalIncome))).
			Ledger(insurance.LedgerINR).
			Code(insurance.AccountCodeRevivalIncome).
			Build(),
		account.New(types.ToUint128(baseID + uint64(insurance.AccountCodePenaltyIncome))).
			Ledger(insurance.LedgerINR).
			Code(insurance.AccountCodePenaltyIncome).
			Build(),

		// Claims accounts
		account.New(types.ToUint128(baseID + uint64(insurance.AccountCodeClaimsReserve))).
			Ledger(insurance.LedgerINR).
			Code(insurance.AccountCodeClaimsReserve).
			Build(),
		account.New(types.ToUint128(baseID + uint64(insurance.AccountCodeClaimsPayable))).
			Ledger(insurance.LedgerINR).
			Code(insurance.AccountCodeClaimsPayable).
			Build(),
	}

	_, err := accountMgr.CreateBatch(systemAccounts)
	return err
}

// createPolicyAccounts creates individual policy accounts
func createPolicyAccounts(c *client.Client, policyIDs []types.Uint128) error {
	accountMgr := account.NewManager(c)

	var policyAccounts []types.Account
	for _, policyID := range policyIDs {
		acc := account.New(policyID).
			Ledger(insurance.LedgerINR).
			Code(insurance.AccountCodePolicyActive).
			HistoryEnabled().
			Build()
		policyAccounts = append(policyAccounts, acc)
	}

	_, err := accountMgr.CreateBatch(policyAccounts)
	return err
}

// showPolicyBalance displays the balance of a policy account
func showPolicyBalance(c *client.Client, policyID types.Uint128, name string) {
	accountMgr := account.NewManager(c)
	balance, err := accountMgr.GetBalance(policyID)
	if err != nil {
		fmt.Printf("  ⚠ Could not retrieve balance for %s\n", name)
		return
	}

	fmt.Printf("  %s: Credits=₹%v, Debits=₹%v\n",
		name,
		balance.Account.CreditsPosted,
		balance.Account.DebitsPosted,
	)
}

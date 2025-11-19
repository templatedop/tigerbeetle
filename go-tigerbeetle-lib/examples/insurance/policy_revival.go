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

// Example demonstrating policy revival when a policy has lapsed
// Revival requires:
// 1. Payment of outstanding premiums
// 2. Payment of penalty charges
// 3. Payment of interest on outstanding amounts
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

	fmt.Println("=== Insurance Policy Revival Demo ===\n")

	// Setup: Create necessary accounts
	fmt.Println("Setting up accounts...")
	accountMgr := account.NewManager(c)

	systemAccounts := []types.Account{
		// Payment method accounts
		account.New(types.ToUint128(uint64(insurance.AccountCodePaymentGateway))).
			Ledger(insurance.LedgerINR).
			Code(insurance.AccountCodePaymentGateway).
			Build(),
		account.New(types.ToUint128(uint64(insurance.AccountCodeCashCollection))).
			Ledger(insurance.LedgerINR).
			Code(insurance.AccountCodeCashCollection).
			Build(),

		// Company revenue accounts
		account.New(types.ToUint128(uint64(insurance.AccountCodeRevivalIncome))).
			Ledger(insurance.LedgerINR).
			Code(insurance.AccountCodeRevivalIncome).
			Build(),
		account.New(types.ToUint128(uint64(insurance.AccountCodePenaltyIncome))).
			Ledger(insurance.LedgerINR).
			Code(insurance.AccountCodePenaltyIncome).
			Build(),
	}

	if _, err := accountMgr.CreateBatch(systemAccounts); err != nil {
		log.Fatalf("Failed to create system accounts: %v", err)
	}
	fmt.Println("✓ System accounts created\n")

	// Scenario 1: Simple lapsed policy - single missed premium
	fmt.Println("Scenario 1: Policy with 1 missed premium")
	fmt.Println("=========================================")

	policy1ID := types.ToUint128(100001)
	policy1 := account.New(policy1ID).
		Ledger(insurance.LedgerINR).
		Code(insurance.AccountCodePolicyLapsed).
		HistoryEnabled().
		UserData64(12345). // Policy number
		Build()

	if err := accountMgr.Create(policy1); err != nil {
		log.Fatalf("Failed to create policy: %v", err)
	}

	fmt.Println("Policy Details:")
	fmt.Println("  Policy Number: 12345")
	fmt.Println("  Status: LAPSED")
	fmt.Println("  Lapsed Date: December 1, 2023")
	fmt.Println("  Monthly Premium: ₹5,000")
	fmt.Println("  Missed Premiums: 1 month")
	fmt.Println()

	insuranceOps := insurance.NewOperations(c, insurance.LedgerINR)

	// Calculate revival charges
	outstandingPremium := uint64(5000)
	penaltyRate := 0.10                                              // 10% penalty
	penalty := uint64(float64(outstandingPremium) * penaltyRate)
	interestRate := 0.05                                             // 5% interest
	interest := uint64(float64(outstandingPremium) * interestRate)
	totalAmount := outstandingPremium + penalty + interest

	fmt.Println("Revival Calculation:")
	fmt.Printf("  Outstanding Premium:  ₹%d\n", outstandingPremium)
	fmt.Printf("  Penalty (10%%):        ₹%d\n", penalty)
	fmt.Printf("  Interest (5%%):        ₹%d\n", interest)
	fmt.Printf("  ─────────────────────────\n")
	fmt.Printf("  Total Amount Due:     ₹%d\n\n", totalAmount)

	// Process revival
	if err := insuranceOps.RevivePolicy(insurance.RevivalParams{
		BaseTransferID:     1000,
		PolicyAccountID:    policy1ID,
		OutstandingPremium: types.ToUint128(outstandingPremium),
		PenaltyAmount:      types.ToUint128(penalty),
		InterestAmount:     types.ToUint128(interest),
		PaymentMethod:      insurance.PaymentMethodPaymentGateway,
		LapsedDate:         time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
		RevivalDate:        time.Now(),
	}); err != nil {
		log.Fatalf("Failed to revive policy: %v", err)
	}

	fmt.Println("✓ Policy #12345 successfully revived!")
	fmt.Println("  Payment Method: Online Payment Gateway")
	fmt.Println("  All charges paid atomically\n")

	// Show policy balance
	balance1, _ := accountMgr.GetBalance(policy1ID)
	fmt.Printf("  Policy Balance: ₹%d\n\n", balance1.Account.CreditsPosted)

	// Scenario 2: Multiple missed premiums
	fmt.Println("Scenario 2: Policy with 3 missed premiums")
	fmt.Println("==========================================")

	policy2ID := types.ToUint128(100002)
	policy2 := account.New(policy2ID).
		Ledger(insurance.LedgerINR).
		Code(insurance.AccountCodePolicyLapsed).
		HistoryEnabled().
		UserData64(67890).
		Build()

	if err := accountMgr.Create(policy2); err != nil {
		log.Fatalf("Failed to create policy: %v", err)
	}

	fmt.Println("Policy Details:")
	fmt.Println("  Policy Number: 67890")
	fmt.Println("  Status: LAPSED")
	fmt.Println("  Lapsed Date: October 1, 2023")
	fmt.Println("  Monthly Premium: ₹10,000")
	fmt.Println("  Missed Premiums: 3 months")
	fmt.Println()

	// Calculate revival charges for 3 months
	outstandingPremium2 := uint64(10000 * 3) // 3 months
	penalty2 := uint64(float64(outstandingPremium2) * 0.12)     // 12% penalty for longer lapse
	interest2 := uint64(float64(outstandingPremium2) * 0.08)    // 8% interest
	totalAmount2 := outstandingPremium2 + penalty2 + interest2

	fmt.Println("Revival Calculation:")
	fmt.Printf("  Outstanding Premium:  ₹%d (3 months)\n", outstandingPremium2)
	fmt.Printf("  Penalty (12%%):        ₹%d\n", penalty2)
	fmt.Printf("  Interest (8%%):        ₹%d\n", interest2)
	fmt.Printf("  ─────────────────────────\n")
	fmt.Printf("  Total Amount Due:     ₹%d\n\n", totalAmount2)

	// Process revival via cash payment
	if err := insuranceOps.RevivePolicy(insurance.RevivalParams{
		BaseTransferID:     2000,
		PolicyAccountID:    policy2ID,
		OutstandingPremium: types.ToUint128(outstandingPremium2),
		PenaltyAmount:      types.ToUint128(penalty2),
		InterestAmount:     types.ToUint128(interest2),
		PaymentMethod:      insurance.PaymentMethodCash,
		LapsedDate:         time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		RevivalDate:        time.Now(),
	}); err != nil {
		log.Fatalf("Failed to revive policy: %v", err)
	}

	fmt.Println("✓ Policy #67890 successfully revived!")
	fmt.Println("  Payment Method: Cash at Branch")
	fmt.Println("  All charges paid atomically\n")

	// Show policy balance
	balance2, _ := accountMgr.GetBalance(policy2ID)
	fmt.Printf("  Policy Balance: ₹%d\n\n", balance2.Account.CreditsPosted)

	// Scenario 3: Waived penalties (Special case)
	fmt.Println("Scenario 3: Revival with waived penalty (Special concession)")
	fmt.Println("============================================================")

	policy3ID := types.ToUint128(100003)
	policy3 := account.New(policy3ID).
		Ledger(insurance.LedgerINR).
		Code(insurance.AccountCodePolicyLapsed).
		HistoryEnabled().
		UserData64(11111).
		Build()

	if err := accountMgr.Create(policy3); err != nil {
		log.Fatalf("Failed to create policy: %v", err)
	}

	fmt.Println("Policy Details:")
	fmt.Println("  Policy Number: 11111")
	fmt.Println("  Status: LAPSED")
	fmt.Println("  Lapsed Date: November 15, 2023")
	fmt.Println("  Monthly Premium: ₹8,000")
	fmt.Println("  Special Concession: Penalty Waived")
	fmt.Println()

	outstandingPremium3 := uint64(8000 * 2) // 2 months
	interest3 := uint64(float64(outstandingPremium3) * 0.06)
	totalAmount3 := outstandingPremium3 + interest3 // No penalty

	fmt.Println("Revival Calculation:")
	fmt.Printf("  Outstanding Premium:  ₹%d (2 months)\n", outstandingPremium3)
	fmt.Printf("  Penalty:              ₹0 (WAIVED)\n")
	fmt.Printf("  Interest (6%%):        ₹%d\n", interest3)
	fmt.Printf("  ─────────────────────────\n")
	fmt.Printf("  Total Amount Due:     ₹%d\n\n", totalAmount3)

	// Process revival with zero penalty
	if err := insuranceOps.RevivePolicy(insurance.RevivalParams{
		BaseTransferID:     3000,
		PolicyAccountID:    policy3ID,
		OutstandingPremium: types.ToUint128(outstandingPremium3),
		PenaltyAmount:      types.ToUint128(0), // No penalty
		InterestAmount:     types.ToUint128(interest3),
		PaymentMethod:      insurance.PaymentMethodPaymentGateway,
		LapsedDate:         time.Date(2023, 11, 15, 0, 0, 0, 0, time.UTC),
		RevivalDate:        time.Now(),
	}); err != nil {
		log.Fatalf("Failed to revive policy: %v", err)
	}

	fmt.Println("✓ Policy #11111 successfully revived!")
	fmt.Println("  Payment Method: Online Payment")
	fmt.Println("  Special concession applied\n")

	// Show policy balance
	balance3, _ := accountMgr.GetBalance(policy3ID)
	fmt.Printf("  Policy Balance: ₹%d\n\n", balance3.Account.CreditsPosted)

	// Summary
	fmt.Println("=== Revival Summary ===")
	fmt.Printf("Total Policies Revived: 3\n")
	fmt.Printf("Policy #12345: ₹%d paid\n", totalAmount)
	fmt.Printf("Policy #67890: ₹%d paid\n", totalAmount2)
	fmt.Printf("Policy #11111: ₹%d paid\n", totalAmount3)
	fmt.Printf("───────────────────────\n")
	fmt.Printf("Total Revenue:  ₹%d\n\n", totalAmount+totalAmount2+totalAmount3)

	// Show company revenue accounts
	revivalIncome, _ := accountMgr.GetBalance(types.ToUint128(uint64(insurance.AccountCodeRevivalIncome)))
	penaltyIncome, _ := accountMgr.GetBalance(types.ToUint128(uint64(insurance.AccountCodePenaltyIncome)))

	fmt.Println("=== Company Revenue Accounts ===")
	fmt.Printf("Revival Income:  ₹%d\n", revivalIncome.Account.CreditsPosted)
	fmt.Printf("Penalty Income:  ₹%d\n\n", penaltyIncome.Account.CreditsPosted)

	fmt.Println("✓ Policy revival demo completed successfully!")
}

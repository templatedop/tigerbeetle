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

// Example demonstrating insurance claim processing
// Types of claims:
// 1. Death Claim - Paid when insured person passes away
// 2. Maturity Claim - Paid when policy matures
// 3. Partial Withdrawal - Partial amount withdrawn
// 4. Surrender Value - When policy is surrendered before maturity
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

	fmt.Println("=== Insurance Claim Processing Demo ===\n")

	// Setup: Create necessary accounts
	fmt.Println("Setting up accounts...")
	accountMgr := account.NewManager(c)

	systemAccounts := []types.Account{
		// Payment method accounts
		account.New(types.ToUint128(uint64(insurance.AccountCodeBankTransfer))).
			Ledger(insurance.LedgerINR).
			Code(insurance.AccountCodeBankTransfer).
			Build(),
		account.New(types.ToUint128(uint64(insurance.AccountCodePaymentGateway))).
			Ledger(insurance.LedgerINR).
			Code(insurance.AccountCodePaymentGateway).
			Build(),
		account.New(types.ToUint128(uint64(insurance.AccountCodeCheque))).
			Ledger(insurance.LedgerINR).
			Code(insurance.AccountCodeCheque).
			Build(),

		// Premium income (to fund policies)
		account.New(types.ToUint128(uint64(insurance.AccountCodePremiumIncome))).
			Ledger(insurance.LedgerINR).
			Code(insurance.AccountCodePremiumIncome).
			Build(),

		// Claims accounts
		account.New(types.ToUint128(uint64(insurance.AccountCodeClaimsReserve))).
			Ledger(insurance.LedgerINR).
			Code(insurance.AccountCodeClaimsReserve).
			Build(),
		account.New(types.ToUint128(uint64(insurance.AccountCodeClaimsPayable))).
			Ledger(insurance.LedgerINR).
			Code(insurance.AccountCodeClaimsPayable).
			Build(),
	}

	if _, err := accountMgr.CreateBatch(systemAccounts); err != nil {
		log.Fatalf("Failed to create system accounts: %v", err)
	}
	fmt.Println("✓ System accounts created\n")

	insuranceOps := insurance.NewOperations(c, insurance.LedgerINR)

	// Scenario 1: Death Claim
	fmt.Println("Scenario 1: Death Claim")
	fmt.Println("========================")

	deathClaimPolicyID := types.ToUint128(200001)
	deathClaimPolicy := account.New(deathClaimPolicyID).
		Ledger(insurance.LedgerINR).
		Code(insurance.AccountCodePolicyActive).
		HistoryEnabled().
		UserData64(200001). // Policy number
		Build()

	if err := accountMgr.Create(deathClaimPolicy); err != nil {
		log.Fatalf("Failed to create policy: %v", err)
	}

	// Fund the policy with accumulated premiums
	fundPolicy(c, deathClaimPolicyID, 500000) // ₹5,00,000 sum assured

	fmt.Println("Policy Details:")
	fmt.Println("  Policy Number: 200001")
	fmt.Println("  Type: Term Life Insurance")
	fmt.Println("  Sum Assured: ₹5,00,000")
	fmt.Println("  Status: Active")
	fmt.Println("  Claim Type: Death Claim")
	fmt.Println()

	// Process death claim
	deathClaimAmount := uint64(500000)
	if err := insuranceOps.ProcessClaim(insurance.ClaimPaymentParams{
		BaseTransferID:  types.ToUint128(10000),
		PolicyAccountID: deathClaimPolicyID,
		ClaimAmount:     types.ToUint128(deathClaimAmount),
		ClaimType:       insurance.TransferCodeDeathClaim,
		PaymentMethod:   insurance.PaymentMethodBankTransfer,
		ClaimDate:       time.Now(),
	}); err != nil {
		log.Fatalf("Failed to process death claim: %v", err)
	}

	fmt.Println("✓ Death claim processed successfully!")
	fmt.Printf("  Claim Amount: ₹%d\n", deathClaimAmount)
	fmt.Println("  Payment Method: Bank Transfer to Nominee")
	fmt.Println("  Flow: Policy -> Claims Reserve -> Claims Payable -> Bank Transfer\n")

	// Scenario 2: Maturity Claim
	fmt.Println("Scenario 2: Maturity Claim")
	fmt.Println("===========================")

	maturityPolicyID := types.ToUint128(200002)
	maturityPolicy := account.New(maturityPolicyID).
		Ledger(insurance.LedgerINR).
		Code(insurance.AccountCodePolicyMatured).
		HistoryEnabled().
		UserData64(200002).
		Build()

	if err := accountMgr.Create(maturityPolicy); err != nil {
		log.Fatalf("Failed to create policy: %v", err)
	}

	// Fund the policy with accumulated value
	fundPolicy(c, maturityPolicyID, 1000000) // ₹10,00,000 maturity value

	fmt.Println("Policy Details:")
	fmt.Println("  Policy Number: 200002")
	fmt.Println("  Type: Endowment Plan")
	fmt.Println("  Maturity Value: ₹10,00,000")
	fmt.Println("  Policy Term: 20 years (Completed)")
	fmt.Println("  Status: Matured")
	fmt.Println()

	// Process maturity claim
	maturityAmount := uint64(1000000)
	if err := insuranceOps.ProcessClaim(insurance.ClaimPaymentParams{
		BaseTransferID:  types.ToUint128(11000),
		PolicyAccountID: maturityPolicyID,
		ClaimAmount:     types.ToUint128(maturityAmount),
		ClaimType:       insurance.TransferCodeMaturityClaim,
		PaymentMethod:   insurance.PaymentMethodBankTransfer,
		ClaimDate:       time.Now(),
	}); err != nil {
		log.Fatalf("Failed to process maturity claim: %v", err)
	}

	fmt.Println("✓ Maturity claim processed successfully!")
	fmt.Printf("  Maturity Amount: ₹%d\n", maturityAmount)
	fmt.Println("  Payment Method: Bank Transfer to Policyholder")
	fmt.Println("  Flow: Policy -> Claims Reserve -> Claims Payable -> Bank Transfer\n")

	// Scenario 3: Partial Withdrawal
	fmt.Println("Scenario 3: Partial Withdrawal")
	fmt.Println("===============================")

	ulipPolicyID := types.ToUint128(200003)
	ulipPolicy := account.New(ulipPolicyID).
		Ledger(insurance.LedgerINR).
		Code(insurance.AccountCodePolicyActive).
		HistoryEnabled().
		UserData64(200003).
		Build()

	if err := accountMgr.Create(ulipPolicy); err != nil {
		log.Fatalf("Failed to create policy: %v", err)
	}

	// Fund the ULIP policy
	fundPolicy(c, ulipPolicyID, 750000) // ₹7,50,000 fund value

	fmt.Println("Policy Details:")
	fmt.Println("  Policy Number: 200003")
	fmt.Println("  Type: ULIP (Unit Linked Insurance Plan)")
	fmt.Println("  Current Fund Value: ₹7,50,000")
	fmt.Println("  Withdrawal Request: ₹1,00,000")
	fmt.Println("  Status: Active (Continues after withdrawal)")
	fmt.Println()

	// Process partial withdrawal
	withdrawalAmount := uint64(100000)
	if err := insuranceOps.ProcessClaim(insurance.ClaimPaymentParams{
		BaseTransferID:  types.ToUint128(12000),
		PolicyAccountID: ulipPolicyID,
		ClaimAmount:     types.ToUint128(withdrawalAmount),
		ClaimType:       insurance.TransferCodePartialWithdrawal,
		PaymentMethod:   insurance.PaymentMethodPaymentGateway,
		ClaimDate:       time.Now(),
	}); err != nil {
		log.Fatalf("Failed to process withdrawal: %v", err)
	}

	fmt.Println("✓ Partial withdrawal processed successfully!")
	fmt.Printf("  Withdrawal Amount: ₹%d\n", withdrawalAmount)
	fmt.Println("  Payment Method: Online Transfer")
	fmt.Printf("  Remaining Fund Value: ₹%d\n", 750000-100000)
	fmt.Println("  Policy continues to be active\n")

	// Scenario 4: Surrender Value
	fmt.Println("Scenario 4: Surrender Value")
	fmt.Println("============================")

	surrenderPolicyID := types.ToUint128(200004)
	surrenderPolicy := account.New(surrenderPolicyID).
		Ledger(insurance.LedgerINR).
		Code(insurance.AccountCodePolicySurrendered).
		HistoryEnabled().
		UserData64(200004).
		Build()

	if err := accountMgr.Create(surrenderPolicy); err != nil {
		log.Fatalf("Failed to create policy: %v", err)
	}

	// Fund the policy
	fundPolicy(c, surrenderPolicyID, 300000) // ₹3,00,000 accumulated

	fmt.Println("Policy Details:")
	fmt.Println("  Policy Number: 200004")
	fmt.Println("  Type: Money Back Plan")
	fmt.Println("  Years Completed: 7 out of 15")
	fmt.Println("  Premiums Paid: ₹3,00,000")
	fmt.Println("  Surrender Value: ₹2,40,000 (80% of premiums)")
	fmt.Println("  Status: Surrendered")
	fmt.Println()

	// Process surrender
	surrenderAmount := uint64(240000) // 80% of premiums paid
	if err := insuranceOps.ProcessClaim(insurance.ClaimPaymentParams{
		BaseTransferID:  types.ToUint128(13000),
		PolicyAccountID: surrenderPolicyID,
		ClaimAmount:     types.ToUint128(surrenderAmount),
		ClaimType:       insurance.TransferCodeSurrenderValue,
		PaymentMethod:   insurance.PaymentMethodCheque,
		ClaimDate:       time.Now(),
	}); err != nil {
		log.Fatalf("Failed to process surrender: %v", err)
	}

	fmt.Println("✓ Surrender value processed successfully!")
	fmt.Printf("  Surrender Amount: ₹%d\n", surrenderAmount)
	fmt.Println("  Payment Method: Cheque")
	fmt.Println("  Policy terminated\n")

	// Show claims summary
	fmt.Println("=== Claims Summary ===")
	fmt.Println()
	fmt.Printf("Death Claim (Policy #200001):       ₹%d\n", deathClaimAmount)
	fmt.Printf("Maturity Claim (Policy #200002):    ₹%d\n", maturityAmount)
	fmt.Printf("Partial Withdrawal (Policy #200003): ₹%d\n", withdrawalAmount)
	fmt.Printf("Surrender Value (Policy #200004):    ₹%d\n", surrenderAmount)
	fmt.Println("───────────────────────────────────────")
	totalClaims := deathClaimAmount + maturityAmount + withdrawalAmount + surrenderAmount
	fmt.Printf("Total Claims Paid:                   ₹%d\n\n", totalClaims)

	// Show company claims accounts
	claimsReserve, _ := accountMgr.GetBalance(types.ToUint128(uint64(insurance.AccountCodeClaimsReserve)))
	claimsPayable, _ := accountMgr.GetBalance(types.ToUint128(uint64(insurance.AccountCodeClaimsPayable)))

	fmt.Println("=== Company Claims Accounts ===")
	fmt.Printf("Claims Reserve:  Debits=₹%d, Credits=₹%d\n",
		claimsReserve.Account.DebitsPosted,
		claimsReserve.Account.CreditsPosted)
	fmt.Printf("Claims Payable:  Debits=₹%d, Credits=₹%d\n\n",
		claimsPayable.Account.DebitsPosted,
		claimsPayable.Account.CreditsPosted)

	fmt.Println("✓ Claim processing demo completed successfully!")
	fmt.Println("\nAll claims processed atomically with full audit trail!")
}

// fundPolicy is a helper function to add funds to a policy account
func fundPolicy(c *client.Client, policyID types.Uint128, amount uint64) {
	insuranceOps := insurance.NewOperations(c, insurance.LedgerINR)

	// Simulate premium collection to fund the policy
	if err := insuranceOps.CollectPremium(insurance.PremiumCollectionParams{
		TransferID:      types.ToUint128(policyID.ToUint128() + 1000000),
		PolicyAccountID: policyID,
		Amount:          types.ToUint128(amount),
		PaymentMethod:   insurance.PaymentMethodBankTransfer,
		IsFirstPremium:  true,
		DueDate:         time.Now(),
		PaymentDate:     time.Now(),
	}); err != nil {
		log.Fatalf("Failed to fund policy: %v", err)
	}
}

package insurance

import (
	"fmt"
	"time"

	"gitlab.cept.gov.in/it-2.0-common/ledgers/pkg/batch"
	"gitlab.cept.gov.in/it-2.0-common/ledgers/pkg/transfer"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

// InsuranceClient defines the interface for insurance operations
type InsuranceClient interface {
	CreateAccounts([]types.Account) ([]types.AccountEventResult, error)
	CreateTransfers([]types.Transfer) ([]types.TransferEventResult, error)
	LookupAccounts([]types.Uint128) ([]types.Account, error)
	LookupTransfers([]types.Uint128) ([]types.Transfer, error)
}

// Operations provides high-level insurance operations
type Operations struct {
	client InsuranceClient
	ledger uint32
}

// NewOperations creates a new insurance operations helper
func NewOperations(client InsuranceClient, ledger uint32) *Operations {
	return &Operations{
		client: client,
		ledger: ledger,
	}
}

// PremiumCollectionParams contains parameters for premium collection
type PremiumCollectionParams struct {
	TransferID      types.Uint128  // Unique transfer ID
	PolicyAccountID types.Uint128  // Policy account receiving premium
	Amount          types.Uint128  // Premium amount
	PaymentMethod   PaymentMethod  // How payment was received
	IsFirstPremium  bool           // Whether this is the first premium
	DueDate         time.Time      // Premium due date
	PaymentDate     time.Time      // Actual payment date
}

// CollectPremium processes a premium payment from a payment method to a policy account
// This creates a linked transfer: Payment Method -> Premium Income -> Policy Account
func (o *Operations) CollectPremium(params PremiumCollectionParams) error {
	paymentAccountCode := GetPaymentAccountCode(params.PaymentMethod)

	// Determine transfer code based on whether it's first premium
	transferCode := TransferCodePremiumPayment
	if params.IsFirstPremium {
		transferCode = TransferCodeFirstPremium
	}

	// Create linked chain: Payment source -> Premium Income -> Policy
	chain := batch.NewLinkedChain()

	// Transfer 1: From payment collection account to premium income
	transfer1 := transfer.New(params.TransferID).
		DebitAccount(types.ToUint128(uint64(paymentAccountCode))). // Payment account
		CreditAccount(types.ToUint128(uint64(AccountCodePremiumIncome))). // Premium income
		Amount(params.Amount).
		Ledger(o.ledger).
		Code(transferCode).
		UserData64(uint64(params.DueDate.Unix())). // Store due date
		Build()

	// Transfer 2: From premium income to policy account
	transfer2 := transfer.New(types.ToUint128(params.TransferID.ToUint128() + 1)).
		DebitAccount(types.ToUint128(uint64(AccountCodePremiumIncome))). // Premium income
		CreditAccount(params.PolicyAccountID). // Policy account
		Amount(params.Amount).
		Ledger(o.ledger).
		Code(transferCode).
		UserData64(uint64(params.PaymentDate.Unix())). // Store payment date
		Build()

	chain.AddTransfer(transfer1).AddTransfer(transfer2)

	results, err := chain.Execute(o.client)
	if err != nil {
		return fmt.Errorf("failed to collect premium: %w", err)
	}
	if len(results) > 0 {
		return fmt.Errorf("premium collection failed at index %d: %s",
			results[0].Index, results[0].Result.String())
	}

	return nil
}

// RevivalParams contains parameters for policy revival
type RevivalParams struct {
	BaseTransferID     types.Uint128  // Base ID for generating transfer IDs
	PolicyAccountID    types.Uint128  // Policy account to revive
	OutstandingPremium types.Uint128  // Outstanding premium amount
	PenaltyAmount      types.Uint128  // Penalty for late payment
	InterestAmount     types.Uint128  // Interest on outstanding amount
	PaymentMethod      PaymentMethod  // How payment was received
	LapsedDate         time.Time      // When policy lapsed
	RevivalDate        time.Time      // Revival date
}

// RevivePolicy processes a policy revival with outstanding premiums, penalties, and interest
// This creates an atomic linked transfer for all revival amounts
func (o *Operations) RevivePolicy(params RevivalParams) error {
	paymentAccountCode := GetPaymentAccountCode(params.PaymentMethod)
	baseID := params.BaseTransferID.ToUint128()

	chain := batch.NewLinkedChain()
	transferIndex := uint64(0)

	// Transfer outstanding premium
	if params.OutstandingPremium.ToUint128() > 0 {
		chain.AddTransfer(
			transfer.New(types.ToUint128(baseID + transferIndex)).
				DebitAccount(types.ToUint128(uint64(paymentAccountCode))).
				CreditAccount(types.ToUint128(uint64(AccountCodeRevivalIncome))).
				Amount(params.OutstandingPremium).
				Ledger(o.ledger).
				Code(TransferCodeRevivalPremium).
				UserData64(uint64(params.LapsedDate.Unix())).
				Build(),
		)
		transferIndex++

		chain.AddTransfer(
			transfer.New(types.ToUint128(baseID + transferIndex)).
				DebitAccount(types.ToUint128(uint64(AccountCodeRevivalIncome))).
				CreditAccount(params.PolicyAccountID).
				Amount(params.OutstandingPremium).
				Ledger(o.ledger).
				Code(TransferCodeRevivalPremium).
				UserData64(uint64(params.RevivalDate.Unix())).
				Build(),
		)
		transferIndex++
	}

	// Transfer penalty amount
	if params.PenaltyAmount.ToUint128() > 0 {
		chain.AddTransfer(
			transfer.New(types.ToUint128(baseID + transferIndex)).
				DebitAccount(types.ToUint128(uint64(paymentAccountCode))).
				CreditAccount(types.ToUint128(uint64(AccountCodePenaltyIncome))).
				Amount(params.PenaltyAmount).
				Ledger(o.ledger).
				Code(TransferCodeRevivalPenalty).
				Build(),
		)
		transferIndex++
	}

	// Transfer interest amount
	if params.InterestAmount.ToUint128() > 0 {
		chain.AddTransfer(
			transfer.New(types.ToUint128(baseID + transferIndex)).
				DebitAccount(types.ToUint128(uint64(paymentAccountCode))).
				CreditAccount(types.ToUint128(uint64(AccountCodePenaltyIncome))).
				Amount(params.InterestAmount).
				Ledger(o.ledger).
				Code(TransferCodeRevivalInterest).
				Build(),
		)
		transferIndex++
	}

	results, err := chain.Execute(o.client)
	if err != nil {
		return fmt.Errorf("failed to revive policy: %w", err)
	}
	if len(results) > 0 {
		return fmt.Errorf("policy revival failed at index %d: %s",
			results[0].Index, results[0].Result.String())
	}

	return nil
}

// ClaimPaymentParams contains parameters for claim payment
type ClaimPaymentParams struct {
	BaseTransferID  types.Uint128  // Base ID for generating transfer IDs
	PolicyAccountID types.Uint128  // Policy account
	ClaimAmount     types.Uint128  // Claim amount to pay
	ClaimType       uint16         // Type of claim (death, maturity, etc.)
	PaymentMethod   PaymentMethod  // How to pay the claim
	ClaimDate       time.Time      // Date of claim
}

// ProcessClaim processes an insurance claim payment
// Flow: Policy Account -> Claims Reserve -> Claims Payable -> Payment Account
func (o *Operations) ProcessClaim(params ClaimPaymentParams) error {
	paymentAccountCode := GetPaymentAccountCode(params.PaymentMethod)
	baseID := params.BaseTransferID.ToUint128()

	chain := batch.NewLinkedChain()

	// Transfer 1: Policy account to claims reserve
	chain.AddTransfer(
		transfer.New(types.ToUint128(baseID)).
			DebitAccount(params.PolicyAccountID).
			CreditAccount(types.ToUint128(uint64(AccountCodeClaimsReserve))).
			Amount(params.ClaimAmount).
			Ledger(o.ledger).
			Code(params.ClaimType).
			UserData64(uint64(params.ClaimDate.Unix())).
			Build(),
	)

	// Transfer 2: Claims reserve to claims payable
	chain.AddTransfer(
		transfer.New(types.ToUint128(baseID + 1)).
			DebitAccount(types.ToUint128(uint64(AccountCodeClaimsReserve))).
			CreditAccount(types.ToUint128(uint64(AccountCodeClaimsPayable))).
			Amount(params.ClaimAmount).
			Ledger(o.ledger).
			Code(params.ClaimType).
			Build(),
	)

	// Transfer 3: Claims payable to payment account (final settlement)
	chain.AddTransfer(
		transfer.New(types.ToUint128(baseID + 2)).
			DebitAccount(types.ToUint128(uint64(AccountCodeClaimsPayable))).
			CreditAccount(types.ToUint128(uint64(paymentAccountCode))).
			Amount(params.ClaimAmount).
			Ledger(o.ledger).
			Code(params.ClaimType).
			Build(),
	)

	results, err := chain.Execute(o.client)
	if err != nil {
		return fmt.Errorf("failed to process claim: %w", err)
	}
	if len(results) > 0 {
		return fmt.Errorf("claim processing failed at index %d: %s",
			results[0].Index, results[0].Result.String())
	}

	return nil
}

// MonthlyPremiumCollectionParams contains parameters for collecting monthly premiums
type MonthlyPremiumCollectionParams struct {
	BaseTransferID types.Uint128                    // Base ID for generating transfer IDs
	Premiums       []PolicyPremium                  // List of premiums to collect
	Month          time.Month                       // Month of collection
	Year           int                              // Year of collection
}

// PolicyPremium represents a single policy premium payment
type PolicyPremium struct {
	PolicyAccountID types.Uint128
	Amount          types.Uint128
	PaymentMethod   PaymentMethod
	DueDate         time.Time
	PaymentDate     time.Time
}

// CollectMonthlyPremiums collects premiums for multiple policies in a batch
func (o *Operations) CollectMonthlyPremiums(params MonthlyPremiumCollectionParams) ([]types.TransferEventResult, error) {
	var allTransfers []types.Transfer
	baseID := params.BaseTransferID.ToUint128()
	transferIndex := uint64(0)

	for _, premium := range params.Premiums {
		paymentAccountCode := GetPaymentAccountCode(premium.PaymentMethod)

		// Each premium requires 2 linked transfers
		transfer1 := transfer.New(types.ToUint128(baseID + transferIndex)).
			DebitAccount(types.ToUint128(uint64(paymentAccountCode))).
			CreditAccount(types.ToUint128(uint64(AccountCodePremiumIncome))).
			Amount(premium.Amount).
			Ledger(o.ledger).
			Code(TransferCodeRenewalPremium).
			UserData64(uint64(premium.DueDate.Unix())).
			Linked().
			Build()
		transferIndex++

		transfer2 := transfer.New(types.ToUint128(baseID + transferIndex)).
			DebitAccount(types.ToUint128(uint64(AccountCodePremiumIncome))).
			CreditAccount(premium.PolicyAccountID).
			Amount(premium.Amount).
			Ledger(o.ledger).
			Code(TransferCodeRenewalPremium).
			UserData64(uint64(premium.PaymentDate.Unix())).
			Build()
		transferIndex++

		allTransfers = append(allTransfers, transfer1, transfer2)
	}

	// Process all transfers in batches
	processor := batch.NewProcessor(o.client, batch.MaxBatchSize)
	return processor.ProcessTransfers(allTransfers)
}

// PendingPremiumParams contains parameters for creating a pending premium (authorization)
type PendingPremiumParams struct {
	TransferID      types.Uint128
	PolicyAccountID types.Uint128
	Amount          types.Uint128
	PaymentMethod   PaymentMethod
	Timeout         time.Duration // How long to hold the authorization
}

// CreatePendingPremium creates a pending premium payment (like payment authorization)
// This reserves the amount until confirmed or timed out
func (o *Operations) CreatePendingPremium(params PendingPremiumParams) error {
	paymentAccountCode := GetPaymentAccountCode(params.PaymentMethod)

	pendingTransfer := transfer.New(params.TransferID).
		DebitAccount(types.ToUint128(uint64(paymentAccountCode))).
		CreditAccount(params.PolicyAccountID).
		Amount(params.Amount).
		Ledger(o.ledger).
		Code(TransferCodePremiumPayment).
		Pending().
		Timeout(params.Timeout).
		Build()

	results, err := o.client.CreateTransfers([]types.Transfer{pendingTransfer})
	if err != nil {
		return fmt.Errorf("failed to create pending premium: %w", err)
	}
	if len(results) > 0 {
		return fmt.Errorf("pending premium creation failed: %s", results[0].Result.String())
	}

	return nil
}

// ConfirmPendingPremium confirms a pending premium payment
func (o *Operations) ConfirmPendingPremium(confirmID, pendingID types.Uint128) error {
	postTransfer := transfer.New(confirmID).
		PostPending(pendingID).
		Build()

	results, err := o.client.CreateTransfers([]types.Transfer{postTransfer})
	if err != nil {
		return fmt.Errorf("failed to confirm pending premium: %w", err)
	}
	if len(results) > 0 {
		return fmt.Errorf("pending premium confirmation failed: %s", results[0].Result.String())
	}

	return nil
}

// CancelPendingPremium cancels a pending premium payment
func (o *Operations) CancelPendingPremium(cancelID, pendingID types.Uint128) error {
	voidTransfer := transfer.New(cancelID).
		VoidPending(pendingID).
		Build()

	results, err := o.client.CreateTransfers([]types.Transfer{voidTransfer})
	if err != nil {
		return fmt.Errorf("failed to cancel pending premium: %w", err)
	}
	if len(results) > 0 {
		return fmt.Errorf("pending premium cancellation failed: %s", results[0].Result.String())
	}

	return nil
}

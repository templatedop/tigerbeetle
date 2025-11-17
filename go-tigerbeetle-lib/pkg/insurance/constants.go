package insurance

// Account Codes - Define different types of accounts in the insurance system
const (
	// Payment Collection Accounts (1-99)
	AccountCodeCashCollection        uint16 = 1  // Cash payment collection account
	AccountCodePaymentGateway        uint16 = 2  // Online payment gateway account
	AccountCodeBankTransfer          uint16 = 3  // Bank transfer account
	AccountCodeCheque                uint16 = 4  // Cheque payment account
	AccountCodeAutoDebit             uint16 = 5  // Auto-debit/ECS account
	AccountCodeAgentCollection       uint16 = 6  // Agent collection account

	// Policy Accounts (100-999)
	AccountCodePolicyActive          uint16 = 100 // Active policy account
	AccountCodePolicyLapsed          uint16 = 101 // Lapsed policy account
	AccountCodePolicySurrendered     uint16 = 102 // Surrendered policy account
	AccountCodePolicyMatured         uint16 = 103 // Matured policy account

	// Company Accounts (1000-1999)
	AccountCodePremiumIncome         uint16 = 1000 // Premium income account
	AccountCodeRevivalIncome         uint16 = 1001 // Revival premium income
	AccountCodePenaltyIncome         uint16 = 1002 // Late payment penalties

	// Claims Accounts (2000-2999)
	AccountCodeClaimsPayable         uint16 = 2000 // Claims payable account
	AccountCodeClaimsReserve         uint16 = 2001 // Claims reserve account
	AccountCodeClaimsPaid            uint16 = 2002 // Claims paid account

	// Surrender Accounts (3000-3999)
	AccountCodeSurrenderPayable      uint16 = 3000 // Surrender value payable
	AccountCodeSurrenderPaid         uint16 = 3001 // Surrender value paid
)

// Transfer Codes - Define different types of transactions
const (
	// Premium Collection (1-99)
	TransferCodePremiumPayment       uint16 = 1  // Regular premium payment
	TransferCodeFirstPremium         uint16 = 2  // First premium payment
	TransferCodeRenewalPremium       uint16 = 3  // Renewal premium

	// Revival Operations (100-199)
	TransferCodeRevivalPremium       uint16 = 100 // Revival premium payment
	TransferCodeRevivalPenalty       uint16 = 101 // Revival penalty charges
	TransferCodeRevivalInterest      uint16 = 102 // Revival interest charges

	// Claims (200-299)
	TransferCodeDeathClaim           uint16 = 200 // Death claim payment
	TransferCodeMaturityClaim        uint16 = 201 // Maturity claim payment
	TransferCodePartialWithdrawal    uint16 = 202 // Partial withdrawal
	TransferCodeSurrenderValue       uint16 = 203 // Surrender value payment

	// Adjustments (300-399)
	TransferCodeRefund               uint16 = 300 // Premium refund
	TransferCodeReversal             uint16 = 301 // Payment reversal
	TransferCodeAdjustment           uint16 = 302 // Manual adjustment
)

// Ledger IDs - Separate ledgers for different currencies or business units
const (
	LedgerINR  uint32 = 1  // Indian Rupee
	LedgerUSD  uint32 = 2  // US Dollar
	LedgerEUR  uint32 = 3  // Euro
	LedgerGBP  uint32 = 4  // British Pound
)

// Policy States
type PolicyState int

const (
	PolicyStateActive PolicyState = iota
	PolicyStateLapsed
	PolicyStateRevived
	PolicyStateSurrendered
	PolicyStateMatured
	PolicyStateClaimed
)

func (s PolicyState) String() string {
	switch s {
	case PolicyStateActive:
		return "Active"
	case PolicyStateLapsed:
		return "Lapsed"
	case PolicyStateRevived:
		return "Revived"
	case PolicyStateSurrendered:
		return "Surrendered"
	case PolicyStateMatured:
		return "Matured"
	case PolicyStateClaimed:
		return "Claimed"
	default:
		return "Unknown"
	}
}

// Payment Methods
type PaymentMethod string

const (
	PaymentMethodCash          PaymentMethod = "CASH"
	PaymentMethodPaymentGateway PaymentMethod = "PAYMENT_GATEWAY"
	PaymentMethodBankTransfer  PaymentMethod = "BANK_TRANSFER"
	PaymentMethodCheque        PaymentMethod = "CHEQUE"
	PaymentMethodAutoDebit     PaymentMethod = "AUTO_DEBIT"
	PaymentMethodAgent         PaymentMethod = "AGENT"
)

// GetPaymentAccountCode returns the account code for a payment method
func GetPaymentAccountCode(method PaymentMethod) uint16 {
	switch method {
	case PaymentMethodCash:
		return AccountCodeCashCollection
	case PaymentMethodPaymentGateway:
		return AccountCodePaymentGateway
	case PaymentMethodBankTransfer:
		return AccountCodeBankTransfer
	case PaymentMethodCheque:
		return AccountCodeCheque
	case PaymentMethodAutoDebit:
		return AccountCodeAutoDebit
	case PaymentMethodAgent:
		return AccountCodeAgentCollection
	default:
		return AccountCodeCashCollection
	}
}

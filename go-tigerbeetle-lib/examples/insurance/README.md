# Insurance Management Examples

This directory contains comprehensive examples demonstrating how to use TigerBeetle for insurance operations. These examples show real-world patterns for premium collection, policy management, revivals, and claims processing.

## Overview

The insurance package provides high-level abstractions for common insurance operations:

- **Premium Collection** through multiple payment methods
- **Policy Revival** with outstanding premiums, penalties, and interest
- **Claim Processing** for death, maturity, partial withdrawal, and surrender
- **Batch Operations** for month-end premium processing
- **Payment Authorization** using pending transfers

## Account Structure

### Payment Collection Accounts
Each payment method has its own account:
- Cash Collection Account (Code: 1)
- Payment Gateway Account (Code: 2)
- Bank Transfer Account (Code: 3)
- Cheque Account (Code: 4)
- Auto-Debit Account (Code: 5)
- Agent Collection Account (Code: 6)

### Policy Accounts
Each policy has its own account (Code: 100+):
- Active Policies (Code: 100)
- Lapsed Policies (Code: 101)
- Surrendered Policies (Code: 102)
- Matured Policies (Code: 103)

### Company Accounts
- Premium Income (Code: 1000)
- Revival Income (Code: 1001)
- Penalty Income (Code: 1002)
- Claims Reserve (Code: 2000)
- Claims Payable (Code: 2001)

## Examples

### 1. Premium Collection (`premium_collection.go`)

Demonstrates collecting premiums through various payment methods:

```go
insuranceOps := insurance.NewOperations(client, insurance.LedgerINR)

// Collect premium via Payment Gateway
err := insuranceOps.CollectPremium(insurance.PremiumCollectionParams{
    TransferID:      types.ToUint128(1000),
    PolicyAccountID: policyID,
    Amount:          types.ToUint128(5000), // ₹5,000
    PaymentMethod:   insurance.PaymentMethodPaymentGateway,
    IsFirstPremium:  false,
    DueDate:         time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
    PaymentDate:     time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
})
```

**Payment Flow:**
```
Payment Method Account -> Premium Income Account -> Policy Account
```

**Supported Payment Methods:**
- Cash
- Payment Gateway (Online)
- Bank Transfer
- Cheque
- Auto-Debit (Recurring)
- Agent Collection

### 2. Batch Premium Collection

Process multiple premiums atomically:

```go
premiums := []insurance.PolicyPremium{
    {
        PolicyAccountID: policy1ID,
        Amount:          types.ToUint128(5000),
        PaymentMethod:   insurance.PaymentMethodAutoDebit,
        DueDate:         time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
        PaymentDate:     time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
    },
    // ... more premiums
}

errors, err := insuranceOps.CollectMonthlyPremiums(insurance.MonthlyPremiumCollectionParams{
    BaseTransferID: types.ToUint128(4000),
    Premiums:       premiums,
    Month:          time.February,
    Year:           2024,
})
```

### 3. Policy Revival (`policy_revival.go`)

Revive lapsed policies with outstanding amounts:

```go
err := insuranceOps.RevivePolicy(insurance.RevivalParams{
    BaseTransferID:     types.ToUint128(2000),
    PolicyAccountID:    lapsedPolicyID,
    OutstandingPremium: types.ToUint128(15000), // ₹15,000
    PenaltyAmount:      types.ToUint128(1500),  // ₹1,500 penalty
    InterestAmount:     types.ToUint128(750),   // ₹750 interest
    PaymentMethod:      insurance.PaymentMethodPaymentGateway,
    LapsedDate:         time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
    RevivalDate:        time.Now(),
})
```

**Revival Components:**
1. Outstanding Premium - All missed premiums
2. Penalty - Late payment charges
3. Interest - Interest on outstanding amount

**Flow:**
All three components are processed atomically - either all succeed or all fail.

### 4. Claim Processing (`claim_processing.go`)

Process different types of insurance claims:

#### Death Claim
```go
err := insuranceOps.ProcessClaim(insurance.ClaimPaymentParams{
    BaseTransferID:  types.ToUint128(10000),
    PolicyAccountID: policyID,
    ClaimAmount:     types.ToUint128(500000), // ₹5,00,000
    ClaimType:       insurance.TransferCodeDeathClaim,
    PaymentMethod:   insurance.PaymentMethodBankTransfer,
    ClaimDate:       time.Now(),
})
```

#### Maturity Claim
```go
err := insuranceOps.ProcessClaim(insurance.ClaimPaymentParams{
    BaseTransferID:  types.ToUint128(11000),
    PolicyAccountID: maturedPolicyID,
    ClaimAmount:     types.ToUint128(1000000), // ₹10,00,000
    ClaimType:       insurance.TransferCodeMaturityClaim,
    PaymentMethod:   insurance.PaymentMethodBankTransfer,
    ClaimDate:       time.Now(),
})
```

#### Partial Withdrawal (ULIP)
```go
err := insuranceOps.ProcessClaim(insurance.ClaimPaymentParams{
    BaseTransferID:  types.ToUint128(12000),
    PolicyAccountID: ulipPolicyID,
    ClaimAmount:     types.ToUint128(100000), // ₹1,00,000
    ClaimType:       insurance.TransferCodePartialWithdrawal,
    PaymentMethod:   insurance.PaymentMethodPaymentGateway,
    ClaimDate:       time.Now(),
})
```

#### Surrender Value
```go
err := insuranceOps.ProcessClaim(insurance.ClaimPaymentParams{
    BaseTransferID:  types.ToUint128(13000),
    PolicyAccountID: surrenderPolicyID,
    ClaimAmount:     types.ToUint128(240000), // ₹2,40,000
    ClaimType:       insurance.TransferCodeSurrenderValue,
    PaymentMethod:   insurance.PaymentMethodCheque,
    ClaimDate:       time.Now(),
})
```

**Claim Flow:**
```
Policy Account -> Claims Reserve -> Claims Payable -> Payment Account
```

### 5. Payment Authorization (`complete_insurance_workflow.go`)

Use pending transfers for payment authorization:

```go
// Create authorization (reserve funds)
err := insuranceOps.CreatePendingPremium(insurance.PendingPremiumParams{
    TransferID:      types.ToUint128(5000),
    PolicyAccountID: policyID,
    Amount:          types.ToUint128(5000),
    PaymentMethod:   insurance.PaymentMethodPaymentGateway,
    Timeout:         24 * time.Hour,
})

// Confirm payment
err = insuranceOps.ConfirmPendingPremium(
    types.ToUint128(5001),
    types.ToUint128(5000),
)

// Or cancel if needed
err = insuranceOps.CancelPendingPremium(
    types.ToUint128(5001),
    types.ToUint128(5000),
)
```

### 6. Complete Workflow (`complete_insurance_workflow.go`)

A comprehensive example showing the entire insurance lifecycle:
1. System account setup
2. Policy creation
3. Premium collection (multiple methods)
4. Policy revival
5. Claim processing
6. Batch operations

## Running the Examples

### Prerequisites

1. **Start TigerBeetle Server:**
   ```bash
   # Download and install TigerBeetle
   # Then create and start a server
   ./tigerbeetle format --cluster=0 --replica=0 --replica-count=1 0_0.tigerbeetle
   ./tigerbeetle start --addresses=3000 0_0.tigerbeetle
   ```

2. **Update Module Dependencies:**
   ```bash
   cd go-tigerbeetle-lib
   go mod download
   ```

### Run Examples

```bash
# Premium collection example
go run examples/insurance/premium_collection.go

# Policy revival example
go run examples/insurance/policy_revival.go

# Claim processing example
go run examples/insurance/claim_processing.go

# Complete workflow
go run examples/insurance/complete_insurance_workflow.go
```

## Key Features

### Atomicity
All operations use linked transfers to ensure atomicity:
- Multi-step premium collection is atomic
- Revival with premium + penalty + interest is atomic
- Claim processing through multiple stages is atomic

### Audit Trail
Every transaction is immutable and traceable:
- Store due dates and payment dates in UserData fields
- Track original transaction IDs for refunds
- Full history of policy transactions

### Scalability
Designed for high-volume processing:
- Batch premium collection for thousands of policies
- Efficient monthly/yearly processing
- Support for millions of transactions per second

### Payment Method Flexibility
Support for multiple payment channels:
- Real-time (Payment Gateway, Auto-Debit)
- Delayed (Cheque, Bank Transfer)
- Offline (Cash, Agent)

## Best Practices

1. **Unique IDs**: Always use unique transfer IDs
2. **Ledger Separation**: Use different ledgers for different currencies
3. **Account Codes**: Use consistent codes for account categorization
4. **Transfer Codes**: Use transfer codes to categorize operations
5. **UserData Fields**: Use for external references and metadata
6. **Batch Operations**: Batch operations for better performance
7. **Error Handling**: Always check for operation errors
8. **Linked Transfers**: Use for multi-step atomic operations

## Account Reconciliation

TigerBeetle's immutability makes reconciliation easy:

```go
// Get all transactions for a policy
transfers, err := client.GetAccountTransfers(filter)

// Get historical balances
balances, err := client.GetAccountBalances(filter)

// Query by time period
transfers, err := client.QueryTransfers(timeRangeFilter)
```

## Support

For questions or issues:
- Check the main [README](../../README.md)
- Review [TigerBeetle Documentation](https://docs.tigerbeetle.com/)
- Open an issue on GitHub

## License

Same as the main library - Apache 2.0

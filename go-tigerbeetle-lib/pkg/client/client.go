package client

import (
	"fmt"

	tb "github.com/tigerbeetle/tigerbeetle-go"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

// Client wraps the TigerBeetle client with high-level operations
type Client struct {
	tb     tb.Client
	config Config
}

// Config holds client configuration
type Config struct {
	ClusterID      uint128
	ReplicaAddrs   []string
	MaxConcurrency uint
}

// uint128 represents a 128-bit unsigned integer
type uint128 = types.Uint128

// New creates a new TigerBeetle client with the given configuration
func New(config Config) (*Client, error) {
	if len(config.ReplicaAddrs) == 0 {
		return nil, fmt.Errorf("at least one replica address is required")
	}

	if config.MaxConcurrency == 0 {
		config.MaxConcurrency = 32 // Default concurrency
	}

	tbClient, err := tb.NewClient(
		config.ClusterID,
		config.ReplicaAddrs,
		config.MaxConcurrency,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create TigerBeetle client: %w", err)
	}

	return &Client{
		tb:     tbClient,
		config: config,
	}, nil
}

// Close closes the TigerBeetle client connection
func (c *Client) Close() {
	c.tb.Close()
}

// NativeClient returns the underlying TigerBeetle client for advanced use cases
func (c *Client) NativeClient() tb.Client {
	return c.tb
}

// CreateAccounts creates one or more accounts in TigerBeetle
func (c *Client) CreateAccounts(accounts []types.Account) ([]types.CreateAccountsError, error) {
	results, err := c.tb.CreateAccounts(accounts)
	if err != nil {
		return nil, fmt.Errorf("failed to create accounts: %w", err)
	}
	return results, nil
}

// CreateTransfers creates one or more transfers in TigerBeetle
func (c *Client) CreateTransfers(transfers []types.Transfer) ([]types.CreateTransfersError, error) {
	results, err := c.tb.CreateTransfers(transfers)
	if err != nil {
		return nil, fmt.Errorf("failed to create transfers: %w", err)
	}
	return results, nil
}

// LookupAccounts retrieves accounts by their IDs
func (c *Client) LookupAccounts(ids []types.Uint128) ([]types.Account, error) {
	accounts, err := c.tb.LookupAccounts(ids)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup accounts: %w", err)
	}
	return accounts, nil
}

// LookupTransfers retrieves transfers by their IDs
func (c *Client) LookupTransfers(ids []types.Uint128) ([]types.Transfer, error) {
	transfers, err := c.tb.LookupTransfers(ids)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup transfers: %w", err)
	}
	return transfers, nil
}

// GetAccountTransfers retrieves transfers for a specific account
func (c *Client) GetAccountTransfers(filter types.AccountFilter) ([]types.Transfer, error) {
	transfers, err := c.tb.GetAccountTransfers(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get account transfers: %w", err)
	}
	return transfers, nil
}

// GetAccountBalances retrieves historical balances for an account
func (c *Client) GetAccountBalances(filter types.AccountFilter) ([]types.AccountBalance, error) {
	balances, err := c.tb.GetAccountBalances(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get account balances: %w", err)
	}
	return balances, nil
}

// QueryAccounts queries accounts with filtering
func (c *Client) QueryAccounts(filter types.QueryFilter) ([]types.Account, error) {
	accounts, err := c.tb.QueryAccounts(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to query accounts: %w", err)
	}
	return accounts, nil
}

// QueryTransfers queries transfers with filtering
func (c *Client) QueryTransfers(filter types.QueryFilter) ([]types.Transfer, error) {
	transfers, err := c.tb.QueryTransfers(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to query transfers: %w", err)
	}
	return transfers, nil
}

package account

import (
	"errors"
	"fmt"

	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

var (
	// ErrAccountNotFound is returned when an account lookup returns no results
	ErrAccountNotFound = errors.New("account not found")
)

// AccountError represents an error that occurred during account creation
type AccountError struct {
	Index  uint32
	Result types.CreateAccountResult
}

func (e *AccountError) Error() string {
	return fmt.Sprintf("account error at index %d: %s", e.Index, e.Result.String())
}

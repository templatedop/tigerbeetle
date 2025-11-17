package transfer

import (
	"errors"
	"fmt"

	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

var (
	// ErrTransferNotFound is returned when a transfer lookup returns no results
	ErrTransferNotFound = errors.New("transfer not found")
)

// TransferError represents an error that occurred during transfer creation
type TransferError struct {
	Index  uint32
	Result types.CreateTransferResult
}

func (e *TransferError) Error() string {
	return fmt.Sprintf("transfer error at index %d: %s", e.Index, e.Result.String())
}

package model

import "fmt"

// ErrNotFound is the error for not found.
var ErrNotFound = fmt.Errorf("not found")

// ErrInsufficientFunds is the error for insufficient funds.
var ErrInsufficientFunds = fmt.Errorf("insufficient funds")

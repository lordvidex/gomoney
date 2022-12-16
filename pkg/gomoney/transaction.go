package gomoney

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	ErrTransferAccountsNotSpecified  = errors.New("transfer transaction requires two accounts")
	ErrDepositAccountNotSpecified    = errors.New("deposit transaction requires one account to deposit to")
	ErrWithdrawalAccountNotSpecified = errors.New("withdrawal transaction requires one account to withdraw from")
	ErrInvalidWithdrawalAccount      = errors.New("account is blocked or has insufficient balance")
	ErrDifferentCurrencies           = errors.New("transfer must be between accounts of the same currencies")
)

type TransactionType int

const (
	Transfer TransactionType = iota
	Deposit
	Withdrawal
)

func (t TransactionType) String() string {
	switch t {
	case Transfer:
		return "Transfer"
	case Deposit:
		return "Deposit"
	case Withdrawal:
		return "Withdrawal"
	}
	return ""
}

type Transaction struct {
	ID      uuid.UUID
	Amount  float64
	From    *Account
	To      *Account
	Created time.Time
	Type    TransactionType
}

func NewTransaction(from, to *Account, amount float64, t TransactionType) *Transaction {
	return &Transaction{
		ID:      uuid.New(),
		Amount:  amount,
		From:    from,
		To:      to,
		Created: time.Now(),
		Type:    t,
	}
}

func (t *Transaction) Validate() error {
	switch t.Type {
	case Transfer:
		if t.From == nil || t.To == nil {
			return ErrTransferAccountsNotSpecified
		}
		if t.From.Currency != t.To.Currency {
			return ErrDifferentCurrencies
		}
		if !t.From.CanTransfer(t.Amount) {
			return ErrInvalidWithdrawalAccount
		}
	case Deposit:
		if t.To == nil {
			return ErrDepositAccountNotSpecified
		}
	case Withdrawal:
		if t.From == nil {
			return ErrWithdrawalAccountNotSpecified
		}
		if !t.From.CanTransfer(t.Amount) {
			return ErrInvalidWithdrawalAccount
		}
	}
	return nil
}

type TransactionSummary struct {
	Transactions []Transaction
	Account      *Account
}

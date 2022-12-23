package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lordvidex/gomoney/api/internal/adapters/handlers/response"
	"github.com/lordvidex/gomoney/api/internal/core"
	"github.com/lordvidex/gomoney/pkg/gomoney"
)

func userFromCtx(ctx *fiber.Ctx) (*core.ApiUser, error) {
	u, ok := ctx.Locals(AuthUserPayload).(*core.ApiUser)
	if !ok {
		return nil, ctx.Status(fiber.StatusUnauthorized).
			JSON(response.ErrM("User is unauthenticated"))
	}
	return u, nil
}

func parseBody(ctx *fiber.Ctx, obj interface{}) error {
	if err := ctx.BodyParser(obj); err != nil {
		_ = ctx.Status(fiber.StatusBadRequest).
			JSON(response.ErrM(err.Error()))
		return err
	}
	return nil
}

func parseParams(ctx *fiber.Ctx, obj interface{}) error {
	if err := ctx.ParamsParser(obj); err != nil {
		_ = ctx.Status(fiber.StatusBadRequest).
			JSON(response.ErrM(err.Error()))
		return err
	}
	return nil
}

func setCtxBodyError(ctx *fiber.Ctx, err error) error {
	ge := func(ge *gomoney.Error) error {
		messages := make([]response.Error, len(ge.Messages))
		for i, m := range ge.Messages {
			messages[i] = response.Err(m, ge.Code)
		}
		return ctx.Status(response.C(ge.Code)).JSON(response.Errs(messages...))
	}
	switch err := err.(type) {
	case *gomoney.Error:
		return ge(err)
	case gomoney.Error:
		return ge(&err)
	}
	return ctx.Status(fiber.StatusInternalServerError).
		JSON(response.ErrM(err.Error()))
}

func strptr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func int64ptr(i int64) *int64 {
	return &i
}

func boolptr(b bool) *bool {
	return &b
}

func float64ptr(f float64) *float64 {
	return &f
}

func parseUser(u *gomoney.User) UserDTO {
	return UserDTO{
		ID:    strptr(u.ID.String()),
		Name:  strptr(u.Name),
		Phone: strptr(u.Phone),
	}
}

func parseAccount(a *gomoney.Account) *AccountDTO {
	if a == nil {
		return nil
	}
	return &AccountDTO{
		ID:          int64ptr(a.Id),
		Title:       strptr(a.Title),
		Description: strptr(a.Description),
		Balance:     float64ptr(a.Balance),
		Currency:    strptr(string(a.Currency)),
		IsBlocked:   boolptr(a.IsBlocked),
	}
}

func parseTransaction(t *gomoney.Transaction) TransactionDTO {
	return TransactionDTO{
		ID:        strptr(t.ID.String()),
		To:        parseAccount(t.To),
		From:      parseAccount(t.From),
		Amount:    float64ptr(t.Amount),
		CreatedAt: strptr(t.Created.Format("2006-01-02 15:04:05")),
		Type:      strptr(t.Type.String()),
	}
}

func parseTransactionSummary(s *gomoney.TransactionSummary) TransactionSummaryDTO {
	txs := make([]TransactionDTO, len(s.Transactions))
	for i, t := range s.Transactions {
		txs[i] = parseTransaction(&t)
	}
	return TransactionSummaryDTO{
		Account: parseAccount(s.Account),
		// Transaction: txs,
	}
}

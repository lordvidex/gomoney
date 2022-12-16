package redis

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/lordvidex/gomoney/api/internal/core"
	"github.com/lordvidex/gomoney/pkg/gomoney"
	"github.com/stretchr/testify/require"
)

func TestUserRepo(t *testing.T) {
	id, err := uuid.NewRandom()
	require.NoError(t, err)

	ur := NewUserRepo(testClient)

	user, password := gomoney.User{
		ID:    id,
		Name:  "John Doe",
		Phone: "+2348123456789",
	}, "password"

	err = ur.SaveUser(context.Background(), &core.ApiUser{
		User:     user,
		Password: password,
	})
	require.NoError(t, err)

	got, err := ur.GetUserFromPhone(context.Background(), user.Phone)
	require.NoError(t, err)
	require.Equal(t, user, got.User)
	require.Equal(t, password, got.Password)

	got, err = ur.GetUserFromPhone(context.Background(), "invalid")
	require.Error(t, err)
	require.Nil(t, got)
}

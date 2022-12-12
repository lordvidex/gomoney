package core

import (
	"github.com/google/uuid"
	"github.com/lordvidex/gomoney/pkg/gomoney"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestApiUser_MarshalBinary(t *testing.T) {
	id := uuid.New()
	type fields struct {
		User     gomoney.User
		Password string
	}
	tests := []struct {
		name    string
		user    fields
		wantErr bool
	}{{
		name: "OK",
		user: fields{
			User: gomoney.User{
				ID:    id,
				Name:  "Ahmad",
				Phone: "79600313041",
			},
			Password: "123123123",
		},
		wantErr: false,
	},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			au := &ApiUser{
				User:     tt.user.User,
				Password: tt.user.Password,
			}
			got, err := au.MarshalBinary()
			require.NoError(t, err, "MarshalBinary() error = %v, wantErr %v", err, tt.wantErr)

			au2 := &ApiUser{}
			_ = au2.UnmarshalBinary(got)
			require.NotEmpty(t, au2)
			require.True(t, reflect.DeepEqual(au2, au), "got = %v, want %v", au2, au)
		})
	}
}

func TestApiUserWithToken_MarshalBinary(t *testing.T) {
	id := uuid.New()
	type fields struct {
		UserWithToken ApiUserWithToken
		Password      string
	}
	tests := []struct {
		name    string
		user    fields
		wantErr bool
	}{{
		name: "OK",
		user: fields{
			UserWithToken: ApiUserWithToken{
				ApiUser: &ApiUser{
					User: gomoney.User{
						ID:    id,
						Name:  "Ahmad",
						Phone: "79600313041"},
					Password: "123123123",
				},
				Token: "123123123",
			},
			Password: "123123123",
		},
		wantErr: false,
	},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			au := &ApiUserWithToken{
				ApiUser: &ApiUser{
					User:     tt.user.UserWithToken.User,
					Password: tt.user.Password,
				},
				Token: tt.user.UserWithToken.Token,
			}
			got, err := au.MarshalBinary()
			require.NoError(t, err, "MarshalBinary() error = %v, wantErr %v", err, tt.wantErr)

			au2 := &ApiUserWithToken{}
			_ = au2.UnmarshalBinary(got)
			require.NotEmpty(t, au2)
			require.True(t, reflect.DeepEqual(au2, au), "got = %v, want %v", au2, au)
		})
	}
}

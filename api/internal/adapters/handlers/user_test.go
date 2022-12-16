package handlers

import (
	"testing"
)

func TestLogin(t *testing.T) {
	// arg := loginUserReq{Phone: "+79600313042", Password: "password"}

	// testCases := []struct {
	// 	name        string
	// 	status      int // expected status code
	// 	shouldError bool
	// 	arg         loginUserReq
	// 	stubs       func()
	// }{
	// 	{
	// 		name:        "OK",
	// 		status:      200,
	// 		shouldError: false,
	// 		arg:         arg,
	// 		stubs: func() {
	// 			mur.EXPECT().GetUserFromPhone(gomock.Any(), arg.Phone)
	// 		},
	// 	},
	// }

	// for _, tc := range testCases {
	// 	t.Run(tc.name, func(t *testing.T) {
	// 		// Set request body
	// 		body, err := encodeBody(tc.arg)
	// 		require.NoError(t, err)

	// 		// Create request
	// 		req, err := http.NewRequest("POST", "localhost/api/login", body)
	// 		require.NoError(t, err)

	// 		res, err := mr.f.Test(req, 1)
	// 		require.Equal(t, tc.status, res.StatusCode)
	// 		require.Equal(t, tc.shouldError, err != nil)
	// 	})
	// }
}

func TestRegister(t *testing.T) {

}

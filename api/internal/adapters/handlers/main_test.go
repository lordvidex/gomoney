package handlers

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	mocks "github.com/lordvidex/gomoney/api/internal/adapters/mock"
	"github.com/lordvidex/gomoney/api/internal/application"
)

var (
	mr   *router
	th   *mocks.MockTokenHelper
	ph   *mocks.MockPasswordHasher
	mur  *mocks.MockUserRepository
	msrv *mocks.MockService
)

func TestMain(m *testing.M) {
	mc := gomock.NewController(&testing.T{})

	// mocks
	th = mocks.NewMockTokenHelper(mc)
	ph = mocks.NewMockPasswordHasher(mc)
	msrv = mocks.NewMockService(mc)
	mur = mocks.NewMockUserRepository(mc)
	uc := application.New(mur, th, msrv, ph)
	mr = New(uc)

	os.Exit(m.Run())
}

func encodeBody(body interface{}) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	return &buf, err
}
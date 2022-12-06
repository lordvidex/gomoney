package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/lordvidex/gomoney/pkg/config"
	"github.com/lordvidex/gomoney/pkg/gomoney"
	app "github.com/lordvidex/gomoney/server/internal/application"
	"testing"
	"time"
)

// _conn must not be directly accessed but should be called through
// getDBConnection() to ensure that the connection is initialised.
var _conn *pgx.Conn
var _config *config.Config

func testConfig() *config.Config {
	if _config == nil {
		c := config.New()
		pw := c.Get("DB_PASSWORD")
		if pw == "" {
			panic("DB_PASSWORD not set from env")
		}
		c.Set("DATABASE_URL", fmt.Sprintf("postgres://user:%s@localhost:5432/gomoney_test?sslmode=disable", pw))
		c.Set("MIGRATION_DIRECTORY", "file://migrations/")
		_config = c
	}
	return _config
}

func getDBConnection() *pgx.Conn {
	if _conn == nil {
		con, err := NewConn(context.TODO(), testConfig())
		if err != nil {
			panic(err)
		}
		_conn = con
	}
	return _conn
}

func createRepo() app.AccountRepository {
	return NewAccount(getDBConnection())
}

func TestMustInt64(t *testing.T) {
	tests := []struct {
		name string
		n    int64
		want sql.NullInt64
	}{
		{"zero", 0, sql.NullInt64{Int64: 0, Valid: true}},
		{"positive", 1, sql.NullInt64{Int64: 1, Valid: true}},
		{"negative", -1, sql.NullInt64{Int64: -1, Valid: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mustInt64(tt.n); got != tt.want {
				t.Errorf("MustInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccountRepoTransfer(t *testing.T) {
	type args struct {
		transaction *gomoney.Transaction
	}

	repo := createRepo()
	defer cleanUp(testConfig())
	user := createTestUser(NewUser(getDBConnection()))
	acc1 := createStubAccount(repo, user)
	acc2 := createStubAccount(repo, user)

	tcs := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"transfer", args{&gomoney.Transaction{
			From:    acc1,
			To:      acc2,
			Amount:  100,
			Type:    gomoney.Transfer,
			Created: time.Now(),
			ID:      uuid.New(),
		}}, false},
		{"deposit", args{&gomoney.Transaction{
			To:      acc1,
			Amount:  100,
			Type:    gomoney.Deposit,
			Created: time.Now(),
			ID:      uuid.New(),
		}}, false},
		{"withdraw", args{&gomoney.Transaction{
			From:    acc1,
			Amount:  100,
			Type:    gomoney.Withdrawal,
			Created: time.Now(),
			ID:      uuid.New(),
		}}, false},
	}

	for _, tt := range tcs {
		t.Run(tt.name, func(t *testing.T) {
			if err := repo.Transfer(context.Background(), tt.args.transaction); (err != nil) != tt.wantErr {
				t.Errorf("AccountRepo.Transfer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func createTestUser(r app.UserRepository) uuid.UUID {
	id, e := r.CreateUser(context.Background(), app.CreateUserArg{
		Name:  "adsf",
		Phone: "12345",
	})
	if e != nil {
		return uuid.Nil
	}
	return id

}
func createStubAccount(r app.AccountRepository, userid uuid.UUID) *gomoney.Account {
	id, e := r.CreateAccount(context.Background(), app.CreateAccountArg{
		UserID:      userid,
		Title:       "test-title",
		Description: "description",
		Currency:    gomoney.USD,
	})
	if e != nil {
		return nil
	}
	acc, e := r.GetAccountByID(context.Background(), id)
	if e != nil {
		return nil
	}
	return acc
}

func cleanUp(c *config.Config) {
	dr := c.Get("MIGRATION_DIRECTORY")
	if dr == "" {
		dr = "file://migrations/"
	}
	m, err := migrate.New(dr, c.Get("DATABASE_URL"))
	if err != nil {
		return
	}
	_ = m.Down()
}

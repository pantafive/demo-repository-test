package database

import (
	"testing"

	"github.com/stretchr/testify/assert"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func Test(t *testing.T) {
	t.Parallel()

	type user struct {
		ID       int64  `db:"user_id"`
		UserName string `db:"username"`
	}

	testCases := []struct {
		name               string
		user               user
		errorAssertionFunc assert.ErrorAssertionFunc
	}{
		{
			name:               "create user: positive",
			user:               user{1, "alice"},
			errorAssertionFunc: assert.NoError,
		},
		{
			name:               "create user: false negative",
			user:               user{1, "bob"},
			errorAssertionFunc: assert.Error,
		},
	}

	query := `INSERT INTO "users" ("id", "username") VALUES ($1, $2)`

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			name := t.Name()
			db := NewTestDatabase(name)
			defer db.Close(t)

			_, err := db.Exec(query, tc.user.ID, tc.user.UserName)
			tc.errorAssertionFunc(t, err)
			t.Logf("database: %s", db)
		})
	}
}

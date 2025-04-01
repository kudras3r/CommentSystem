package migration_test

import (
	"database/sql"
	"testing"

	"github.com/kudras3r/CommentSystem/internal/storage/migration"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestMakeMigrations(t *testing.T) {
	connStr := "postgres://ozon_keker:1234@localhost:5432/comm_sys_db?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	assert.Nil(t, err)
	defer db.Close()

	err = migration.MakeMigrations(db)
	assert.NotNil(t, err)
}

func TestMakeMigrationsWithInvalidDB(t *testing.T) {
	connStr := "postgres://invalid_user:invalid_pass@localhost:5432/invalid_db?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	assert.Nil(t, err)
	defer db.Close()

	err = migration.MakeMigrations(db)
	assert.NotNil(t, err)
}

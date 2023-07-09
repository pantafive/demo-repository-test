package database

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
)

type DevDatabase struct {
	databaseName string
	// databaseWithSudo is a connection with the ability to create and drop other databases.
	// Usually, it is a connection to the default database "postgres" with superuser rights.
	databaseWithSudo *sqlx.DB
	*sqlx.DB
}

func (db *DevDatabase) String() string {
	return db.databaseName
}

func NewTestDatabase(databaseName string) *DevDatabase {
	databaseName = sanitizeDatabaseName(databaseName)

	envFile := SearchUpwardForFile("dev.env")
	if envFile == "" {
		panic("Failed to find .env file")
	}
	MustLoadEnv(envFile)

	templateName := MustNotEmptyString(os.Getenv("DATABASE_TEMPLATE"))

	superUserDatabase := sqlx.MustConnect("pgx", postgresDSN("postgres"))

	superUserDatabase.MustExec("DROP DATABASE IF EXISTS " + databaseName)
	superUserDatabase.MustExec("CREATE DATABASE " + databaseName + " TEMPLATE " + templateName)

	connectionString := postgresDSN(databaseName)

	return &DevDatabase{
		databaseName:     databaseName,
		databaseWithSudo: superUserDatabase,
		DB:               sqlx.MustConnect("pgx", connectionString),
	}
}

func sanitizeDatabaseName(databaseName string) string {
	databaseName = strings.TrimSpace(databaseName)
	databaseName = strings.ToLower(databaseName)
	databaseName = regexp.MustCompile(`[^a-zA-Z0-9_-]`).ReplaceAllString(databaseName, "_")
	databaseName = regexp.MustCompile(`_+`).ReplaceAllString(databaseName, "_")
	return databaseName
}

func postgresDSN(databaseName string) string {
	template := "postgresql://%s:%s@%s:%d/%s?sslmode=disable"

	user := MustNotEmptyString(os.Getenv("POSTGRES_USER"))
	password := MustNotEmptyString(os.Getenv("POSTGRES_PASSWORD"))
	host := MustNotEmptyString(os.Getenv("POSTGRES_HOST"))

	dsn := fmt.Sprintf(template, user, password, host, 5432, databaseName) //nolint:gomnd
	return dsn
}

type Testing interface {
	Helper()
	Failed() bool
}

func (db *DevDatabase) Close(t Testing) {
	t.Helper()
	defer db.databaseWithSudo.Close()
	if err := db.DB.Close(); err != nil {
		panic(fmt.Errorf("failed to close dev Database: %w", err))
	}

	if !t.Failed() {
		db.databaseWithSudo.MustExec("DROP DATABASE " + db.databaseName)
	}
}

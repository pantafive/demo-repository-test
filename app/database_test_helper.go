package database

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
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

func NewTestDatabase() *DevDatabase {
	envFile := SearchUpwardForFile("dev.env")
	if envFile == "" {
		panic("Failed to find .env file")
	}
	MustLoadEnv(envFile)

	templateName := MustNotEmptyString(os.Getenv("DATABASE_TEMPLATE"))

	superUserDatabase := sqlx.MustConnect("pgx", postgresDSN("postgres"))
	name := randomDatabaseNameGenerator()

	superUserDatabase.MustExec("CREATE DATABASE " + name + " TEMPLATE " + templateName)

	connectionString := postgresDSN(name)

	return &DevDatabase{
		databaseName:     name,
		databaseWithSudo: superUserDatabase,
		DB:               sqlx.MustConnect("pgx", connectionString),
	}
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
		log.Panic().Err(err).Msg("Failed to close dev Database")
	}

	if !t.Failed() {
		db.databaseWithSudo.MustExec("DROP DATABASE " + db.databaseName)
	}
}

func randomDatabaseNameGenerator() string {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec
	uid, err := ulid.New(ulid.Timestamp(time.Now()), entropy)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to generate ulid")
	}
	return "test_" + strings.ToLower(uid.String())
}

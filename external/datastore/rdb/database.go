package rdb

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rlawnsxo131/madre-server-v3/lib/env"
	"github.com/rlawnsxo131/madre-server-v3/lib/logger"
	"github.com/rs/zerolog"
)

var (
	instanceDatabase *singletonDatabase
	onceDatabase     sync.Once
)

type Database interface {
	Queryx(query string, args ...any) (*sqlx.Rows, error)
	QueryRowx(query string, args ...any) *sqlx.Row
	NamedQuery(query string, arg any) (*sqlx.Rows, error)
	PrepareNamedGet(result any, query string, arg any) error
}

func DatabaseInstance() (*singletonDatabase, error) {
	var resultError error

	onceDatabase.Do(func() {
		psqlInfo := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			env.DatabaseHost(),
			env.DatabasePort(),
			env.DatabaseUser(),
			env.DatabasePassword(),
			env.DatabaseDBName(),
			env.DatabaseSSLMode(),
		)
		logger.DefaultLogger().Info().
			Timestamp().Str("database connection info", psqlInfo).Send()

		db, err := sqlx.Connect("postgres", psqlInfo)
		if err != nil {
			resultError = errors.Wrap(err, "sqlx connect fail")
			return
		}

		l := zerolog.New(os.Stderr).With().Logger()
		instanceDatabase = &singletonDatabase{db, &l}
		initDatabase(instanceDatabase.DB)
	})

	return instanceDatabase, resultError
}

func initDatabase(db *sqlx.DB) {
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(5)
	db.SetConnMaxLifetime(time.Minute)
}

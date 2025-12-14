package sql

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBOptionFunc is a function that can be used to configure the database.
type DBOptionFunc func(*gorm.DB) error

// PostgresDB is a wrapper around gorm.DB.
type PostgresDB struct {
	*gorm.DB
}

type SQLConf struct {
	DSN             string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func NewDB(conf SQLConf, options ...DBOptionFunc) (*PostgresDB, error) {
	dialector := postgres.New(postgres.Config{
		DSN: conf.DSN,
	})

	db, err := gorm.Open(dialector, &gorm.Config{
		TranslateError:         true,
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, err
	}

	postgres, err := db.DB()
	if err != nil {
		return nil, err
	}

	postgres.SetMaxIdleConns(conf.MaxIdleConns)
	postgres.SetMaxOpenConns(conf.MaxOpenConns)
	postgres.SetConnMaxLifetime(conf.ConnMaxLifetime)
	postgres.SetConnMaxIdleTime(conf.ConnMaxIdleTime)

	for _, option := range options {
		if err := option(db); err != nil {
			return nil, err
		}
	}

	return &PostgresDB{
		db,
	}, nil
}

func (p *PostgresDB) Close() error {
	db, err := p.DB.DB()
	if err != nil {
		return err
	}

	return db.Close()
}

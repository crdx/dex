package config

import (
	"database/sql"

	"crdx.org/dex/cmd/dexd/env"
	"crdx.org/dex/db"
	"crdx.org/dex/db/schema"
	"github.com/lithammer/shortuuid/v3"
)

func GetDbConfig() *db.Config {
	return &db.Config{
		Open: func(dsn *db.DSN) (*sql.DB, error) {
			return sql.Open("mysql", dsn.Format())
		},
		DataSource: db.NewDSN().Apply(func(dsn *db.DSN) *db.DSN {
			dsn.DBName = env.DatabaseName()
			dsn.Username = env.DatabaseUsername()
			dsn.Password = env.DatabasePassword()
			dsn.Address = env.DatabaseAddress()
			dsn.Protocol = env.DatabaseProtocol()
			return dsn
		}),
		Create:       true,
		Migrations:   schema.GetMigrations(),
		EnableLogger: env.Debug(),
	}
}

func GetTestDbConfig() *db.Config {
	return &db.Config{
		Open: func(dsn *db.DSN) (*sql.DB, error) {
			return sql.Open("mysql", dsn.Format())
		},
		DataSource: db.NewDSN().Apply(func(dsn *db.DSN) *db.DSN {
			dsn.DBName = env.DatabaseName() + "_test_" + shortuuid.New()
			dsn.Username = env.DatabaseUsername()
			dsn.Password = env.DatabasePassword()
			dsn.Address = env.DatabaseAddress()
			dsn.Protocol = env.DatabaseProtocol()
			return dsn
		}),
		Migrations:   schema.GetMigrations(),
		Fresh:        true,
		EnableLogger: false,
	}
}

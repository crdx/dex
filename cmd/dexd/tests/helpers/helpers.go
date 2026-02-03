package helpers

import (
	"os"
	"testing"

	"crdx.org/dex/cmd/dexd/config"
	"crdx.org/dex/cmd/dexd/env"
	"crdx.org/dex/db"
	"github.com/samber/lo"
)

// Start begins a new db transaction and returns a func that will roll back the transaction. This
// function is NOT thread-safe, but this is fine because each test within a test package runs in
// serial.
func Start() func() {
	_ = db.BeginTransaction()

	return func() {
		_ = db.RollbackTransaction()
	}
}

// TestMain initialises the environment and database for the current test package, runs the tests,
// and then drops the test database.
func TestMain(m *testing.M) {
	env.Init()

	dbConfig := config.GetTestDbConfig()
	lo.Must0(db.Init(dbConfig))

	exitCode := m.Run()

	_, _ = db.Exec("DROP DATABASE " + dbConfig.DataSource.DBName)
	os.Exit(exitCode)
}

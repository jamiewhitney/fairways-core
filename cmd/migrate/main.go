package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jamiewhitney/fairways-core/pkg/database"
	"github.com/jamiewhitney/fairways-core/pkg/logging"
	"github.com/jamiewhitney/fairways-core/pkg/setup"
)

var (
	filePath     = flag.String("path", "migrations/", "path to migrations folder")
	databaseName = flag.String("database", "", "database name")
)

func main() {
	logger := logging.NewLoggerFromEnv()
	flag.Parse()
	if *filePath == "" {
		logger.Fatal("path to migrations folder is required")
	}

	if *databaseName == "" {
		logger.Fatal("database name is required")
	}

	var config database.Config
	_, err := setup.Setup(context.Background(), &config)
	if err != nil {
		logger.Fatalf("failed to setup database: %v", err)
	}

	dsn := fmt.Sprintf("%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true", config.UserPass, config.Host, config.Port, *databaseName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Fatalf("failed to connect to mysql: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Fatalf("failed to ping mysql: %v", err)
		return
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{
		MigrationsTable: fmt.Sprintf("migrations_%s", *databaseName),
		DatabaseName:    *databaseName,
	})
	if err != nil {
		logger.Fatal(err)
	}
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", *filePath),
		*databaseName,
		driver)
	if err != nil {
		logger.Fatal(err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Fatal(err)
	}

	logger.Info("migrations applied successfully")
}

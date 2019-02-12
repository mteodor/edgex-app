//
// Copyright (c) 2018
// Mainflux
//
// SPDX-License-Identifier: Apache-2.0
//

package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // required for SQL access
	migrate "github.com/rubenv/sql-migrate"
)

// Config defines the options that are used when connecting to a PostgreSQL instance
type Config struct {
	Host        string
	Port        string
	User        string
	Pass        string
	Name        string
	SSLMode     string
	SSLCert     string
	SSLKey      string
	SSLRootCert string
}

// Connect creates a connection to the PostgreSQL instance and applies any
// unapplied database migrations. A non-nil error is returned to indicate
// failure.
func Connect(cfg Config) (*sql.DB, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s", cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Pass, cfg.SSLMode, cfg.SSLCert, cfg.SSLKey, cfg.SSLRootCert)

	db, err := sql.Open("postgres", url)
<<<<<<< HEAD
	fmt.Println("url:" + url)
	if err != nil {
		fmt.Println("failed to connect")
=======
	if err != nil {
>>>>>>> 8b2d9977cbbd8ff37891a6bc6f99d4ffa1abe5d9
		return nil, err
	}

	if err := migrateDB(db); err != nil {
		return nil, err
	}
	return db, nil
}

func migrateDB(db *sql.DB) error {
<<<<<<< HEAD
	fmt.Println("connecting to db\n")
=======
	fmt.Printf("connecting to db\n")
>>>>>>> 8b2d9977cbbd8ff37891a6bc6f99d4ffa1abe5d9
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: "events_1",
				Up: []string{
					`CREATE TABLE IF NOT EXISTS events (
						id	 VARCHAR(254) PRIMARY KEY,
						pushed CHAR(60) NOT NULL,
						device CHAR(60) NOT NULL,
						created CHAR(20) NOT NULL,
						modified CHAR(20) NOT NULL,
						origin   CHAR(20) NOT NULL,
						event  CHAR(60) NOT NULL
					)`,
				},
				Down: []string{"DROP TABLE events"},
			},
		},
	}

	_, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
<<<<<<< HEAD
	if err != nil {
		fmt.Println("failed to create table")
	}
=======
>>>>>>> 8b2d9977cbbd8ff37891a6bc6f99d4ffa1abe5d9
	return err
}

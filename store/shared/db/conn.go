// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Drone Non-Commercial License
// that can be found in the LICENSE file.

// +build !oss

package db

import (
	"database/sql"

	"github.com/oars-sigs/drone/store/shared/migrate/mysql"
	"github.com/oars-sigs/drone/store/shared/migrate/postgres"
	"github.com/oars-sigs/drone/store/shared/migrate/sqlite"

	dblib "github.com/drone/drone/store/shared/db"
)

// Connect to a database and verify with a ping.
func Connect(driver, datasource string) (*dblib.DB, error) {
	db, err := sql.Open(driver, datasource)
	if err != nil {
		return nil, err
	}
	switch driver {
	case "mysql":
		db.SetMaxIdleConns(0)
	}
	err = setupDatabase(db, driver)
	if err != nil {
		db.Close()
		return nil, err
	}
	db.Close()
	return dblib.Connect(driver, datasource)
}

// helper function to setup the databsae by performing automated
// database migration steps.
func setupDatabase(db *sql.DB, driver string) error {
	switch driver {
	case "mysql":
		return mysql.Migrate(db)
	case "postgres":
		return postgres.Migrate(db)
	default:
		return sqlite.Migrate(db)
	}
}

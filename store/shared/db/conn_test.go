package db

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func TestConnectSqlite(t *testing.T) {
	db, err := Connect("sqlite3", "./sqlite.sqlite")
	if err != nil {
		t.Error(err)
		return
	}
	db.Close()
}

func TestConnectMysql(t *testing.T) {
	db, err := Connect("mysql", "root:@tcp(127.0.0.1:3306)/test")
	if err != nil {
		t.Error(err)
		return
	}
	db.Close()
}

package main

import (
	"database/sql"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
	"encoding/json"
	"fmt"
)

var globalMySQLDB *sql.DB

func init() {
	var dbCredentials struct{
		Username string
		Password string
	}

	file, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(file, &dbCredentials)
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/gophr", dbCredentials.Username, dbCredentials.Password)
	db, err := NewMySQLDB(dsn)
	if err != nil {
		panic(err)
	}
	globalMySQLDB = db
}

func NewMySQLDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn+"?parseTime=true")
	if err != nil {
		return nil, err
	}

	return db, db.Ping()
}
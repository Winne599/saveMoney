package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func connectDB() (*sql.DB, error) {
	// 数据库信息
	dsn := "root:root@tcp(localhost:3306)/save_money"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to the database")
	return db, nil
}

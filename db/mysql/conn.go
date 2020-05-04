package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "root:QOjM6y7471TTyNJj@tcp(127.0.0.1:3306)/fileserver?charset=utf8&parseTime=true")
	db.SetMaxOpenConns(200)
	err := db.Ping()
	if err != nil {
		fmt.Printf("Filed to connect to mysql, err: %s", err.Error())
		os.Exit(1)
	}
}

// 返回数据库连接
func DbConn() *sql.DB {
	return db
}





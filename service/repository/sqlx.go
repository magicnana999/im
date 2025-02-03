package repository

import (
	_ "github.com/go-sql-driver/mysql" // 不要忘了导入数据库驱动
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func init() {
	dns := "root:root1234@tcp(localhost:3306)/im?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai"

	_db, err := sqlx.Connect("mysql", dns)
	if err != nil {
		panic(err)
	}

	db = _db

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
}

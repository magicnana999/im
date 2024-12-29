package sql

import (
	sql "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/magicnana999/im/logger"
)

const (
	dbDriver = "mysql"
	dbUser   = "heguang"
	dbPass   = "Heguang_0315"
	dbName   = "im-platform"
	dbHost   = "rm-2ze6n61koo12t7e7coo.mysql.rds.aliyuncs.com:3306"
)

var (
	DB *sql.DB
)

func init() {
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbHost+")/"+dbName)

	// if there is an error opening the connection, handle it
	if err != nil {
		logger.Fatal(err.Error())
	}

	pingErr := db.Ping()
	if pingErr != nil {
		logger.Fatal(pingErr)
	}

	DB = db
}

type User struct {
	UserId   int64
	NickName string
}

func Select(sql string, arg ...any) error {
	rows, err := DB.Query(sql)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(arg); err != nil {
			logger.Error(err)
			return fmt.Errorf("no data found %v", err)
		}
	}

	return nil
}

package sql

import (
	sql "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
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
		log.Fatal(err.Error())
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	DB = db
}

func Select(s string) error {
	rows, err := DB.Query(s)
	if err != nil {
		return err
	}

	defer rows.Close()

	fmt.Println(rows.Columns())
	ct, err := rows.ColumnTypes()
	if err != nil {
		return err
	}

	for i := 0; i < len(ct); i++ {
		fmt.Println(ct[i].Name(), ct[i].ScanType().Name())
	}

	for rows.Next() {
		var user_id int64
		rows.Scan(&user_id)
		fmt.Println(user_id)
	}

	return nil

}

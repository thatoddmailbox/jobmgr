package data

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"

	"github.com/thatoddmailbox/jobmgr/config"
)

var DB *sql.DB

func Init() error {
	var err error

	dbConfig := mysql.NewConfig()
	dbConfig.User = config.Current.Database.Username
	dbConfig.Passwd = config.Current.Database.Password
	dbConfig.Addr = config.Current.Database.Host
	dbConfig.DBName = config.Current.Database.Database
	dbConfig.Params = map[string]string{
		"charset": "utf8mb4",
	}
	dbConfig.Collation = "utf8mb4_unicode_ci"

	DB, err = sql.Open("mysql", dbConfig.FormatDSN())
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	return nil
}

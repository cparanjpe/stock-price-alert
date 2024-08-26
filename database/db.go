package database

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql" // Import the MySQL driver
)

var DB *sql.DB

func InitDB() error {
    var err error
    DB, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/price_alert_system")
    if err != nil {
        return err
    }

    // Check if the database connection is available
    if err = DB.Ping(); err != nil {
        return err
    }

    return nil
}

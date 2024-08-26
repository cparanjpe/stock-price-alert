package database

import (
    "fmt"
    "database/sql"
    _ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() {
    var err error
    connStr := "user=username dbname=mydb sslmode=disable"
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        fmt.Println("Error connecting to the database:", err)
        return
    }

    err = db.Ping()
    if err != nil {
        fmt.Println("Database connection is not alive:", err)
        return
    }

    fmt.Println("Database connection established successfully.")
}

func GetDB() *sql.DB {
    return db
}

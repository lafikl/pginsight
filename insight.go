package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lafikl/pq"
)

// Insight is the entry point
type Insight struct {
	db *sql.DB
}

// NewInsight returns a new Insight
func NewInsight() *Insight {
	dburl := os.Getenv("PGINSIGHT_DBURL")
	if len(dburl) == 0 {
		fmt.Println("\n\tEnvironment variable $PGINSIGHT_DBURL not set")
		fmt.Println("\n\tTo set the url:\n\texport PGINSIGHT_DBURL=\"postgres://username:password@localhost/dbname?sslmode=disable\"\n")
		os.Exit(1)
	}
	db, err := sql.Open("postgres", dburl)
	if err != nil {
		fmt.Println("\n\tCouldn't connect to database\n\t Database URL: ", dburl)
		fmt.Println("\n\tTo set the url:\n\texport PGINSIGHT_DBURL=\"postgres://username:password@localhost/dbname?sslmode=disable\"\n")
		os.Exit(1)
	}
	// test connection
	var one int
	err = db.QueryRow("Select 1;").Scan(&one)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return &Insight{
		db: db,
	}
}

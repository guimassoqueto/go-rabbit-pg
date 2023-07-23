package pg

import (
	"database/sql"
	"fmt"
	"log"

	"grp/helpers"
	"grp/variables"

	_ "github.com/lib/pq"
)

func InsertMessage(message string) {
	connStr := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		variables.POSTGRES_USER, 
		variables.POSTGRES_DB, 
		variables.POSTGRES_PASSWORD, 
		variables.POSTGRES_HOST, 
		variables.POSTGRES_PORT,
	)
	db, err := sql.Open("postgres", connStr)
	helpers.FailOnError(err, "Failed to connect to Database")
	defer db.Close()

	_, err = db.Exec("INSERT INTO go_rabbit_pg(message) VALUES($1)", message)
	helpers.FailOnError(err, "Failed to insert data into Database")

	log.Print("Message successfully inserted into Database")
}
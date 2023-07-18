package initializers

import (
	"database/sql"
	"fmt"
	_"github.com/lib/pq"
	// "os"
	// "strconv"
  )
  const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "9999"
	dbname   = "Tablic"
  )

func ConnectToDb() {
	// var host     = os.Getenv("HOST")
	// var port, _  = strconv.Atoi(os.Getenv("DB_POST"))
	// var user     = os.Getenv("USER")
	// var password = os.Getenv("PASSWORD")
	// var dbname   = os.Getenv("DB_NAME")

	postgresqlDbInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	  "password=%s dbname=%s sslmode=disable",
	  host, port, user, password, dbname)

	db, err := sql.Open("postgres", postgresqlDbInfo)
	if err != nil {
	  panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
	  panic(err)
	}
	fmt.Println("Established a successful connection!")
}


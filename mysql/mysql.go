package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/tehsphinx/dbg"
)

var DB *sql.DB

func ConnectMysql(addr, database, user, password string) {

	login := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s",
		user,
		password,
		addr,
		database)

	dbg.Green(login)

	dbc, err := sql.Open("mysql", login)
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	dbc.SetConnMaxLifetime(time.Minute * 3)
	dbc.SetMaxOpenConns(10)
	dbc.SetMaxIdleConns(10)

	if err := dbc.Ping(); err != nil {
		log.Printf("Connection to database not found: %s", err.Error())
		return
	}

	DB = dbc
	initDB()
}

type counterRow struct {
	ID      int
	Counter int
}

func GetCounter() (int, error) {
	query := "SELECT counter FROM pizza WHERE id=1;"
	results, err := DB.Query(query)
	if err != nil {
		log.Printf("Error querying pizzas: %v", err)
		return -1, err
	}

	results.Next()
	var pizzaCounter counterRow
	err = results.Scan(&pizzaCounter.Counter)
	if err != nil {
		log.Printf("Error parsing pizza: %v", err)
		return -1, err
	}
	results.Close()

	return pizzaCounter.Counter, nil
}

func SetCounter(counter int) error {
	stat := fmt.Sprintf("UPDATE pizza SET counter=%d WHERE id=1;", counter)

	_, err := DB.Exec(stat)
	if err != nil {
		log.Printf("Error inserting pizza: %v", err)
		return err
	}
	return nil
}

func initDB() {
	createTab := "CREATE TABLE IF NOT EXISTS pizza (id int NOT NULL UNIQUE AUTO_INCREMENT, counter int);"

	_, err := DB.Exec(createTab)
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}

	insert := "INSERT INTO pizza (counter) VALUES (0);"
	_, err = DB.Exec(insert)
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	return
}

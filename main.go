package main

import (
	"container_example/mysql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"github.com/namsral/flag"
	"github.com/tehsphinx/dbg"
)

var (
	port        = flag.Int("apiport", 8080, "http port to listen for API requests")
	dbuser      = flag.String("dbuser", "root", "User to connect ot the MySQL database")
	dbpass      = flag.String("dbpass", "secret", "Passwort to connec to to the MySQL database")
	dbaddr      = flag.String("dbaddr", "localhost", "Address of the MySQL databse server")
	dbname      = flag.String("dbname", "pizzas", "Name of the database to be used")
	counter int = 0
)

func main() {
	flag.Parse()

	log.Printf("Listening on http port %d\n", *port)

	mysql.ConnectMysql(*dbaddr, *dbname, *dbuser, *dbpass)
	http.ListenAndServe(":"+strconv.Itoa(*port), setupRouter())
}

func setupRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/counter", GetPizzaCounter)

	static := httprouter.New()
	static.ServeFiles("/*filepath", http.Dir("static"))
	router.NotFound = static

	return router
}

func GetPizzaCounter(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	if mysql.DB == nil {
		fmt.Fprintf(w, "%d Pizzas ordered!\n", counter)
		counter++
		dbg.Red("Using local counter")
		return
	}

	dbCounter, err := mysql.GetCounter()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "%d Pizzas ordered!\n", dbCounter)
	if err := mysql.SetCounter(dbCounter + 1); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

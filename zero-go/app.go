package main

import (
	"fmt"
	"net/http"

	"database/sql"
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func hello(c echo.Context) error {
	executeSql("select id, username from person where id = ?")
	return c.String(http.StatusOK, "hello, world!")
}

// Execute sql statement from parameter, which looks like this:
// select a, b, c from some_tabble where id = ?
// Return a map
func executeSql(sqlStmt string) /*map[string]string*/ {
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	ret, err := stmt.Query(1)
	if err != nil {
		log.Fatal(err)
	}

	for ret.Next() {
		var person Person

		err = ret.Scan(&person.ID, &person.Username)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(person.ID, person.Username)
	}
}

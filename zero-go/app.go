package main

import (
	"fmt"
	"net/http"
	"time"

	"database/sql"
	"log"
	"strings"
	"context"
	"strconv"

	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"

	"github.com/graphql-go/graphql"

	"zero-go/model"
	"zero-go/schema"
	"zero-go/util"
)

func main() {
	util.ULog("start application")
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.GET("/people/all", findAllPeople)
	e.GET("/people/:id", findPeople)
	e.GET("/friends/people/:id", findFriends)

	e.GET("/query/:query", executeQuery)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func hello(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"msg": "visite `/people/:id` to get data",
	})
}

func findPeople(c echo.Context) error {
	personID := c.Param("id")

	ret, err := executeSQL("select id,password,last_login,is_superuser,username,first_name,last_name,email,is_staff,is_active,date_joined from person where id = ?", personID)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]string{"error": "Something went wrong when getting data from db"})
	}

	fmt.Println("request to find a person")
	return c.JSON(http.StatusOK, ret)
}

func findAllPeople(c echo.Context) error {
	ret, err := executeSQL("select id,password,last_login,is_superuser,username,first_name,last_name,email,is_staff,is_active,date_joined from person", "")
	if err != nil {
		return c.JSON(http.StatusOK, map[string]string{"error": "Something went wrong when getting data from db"})
	}
	fmt.Println("reqeust all people")
	return c.JSON(http.StatusOK, ret)
}

func findFriends(c echo.Context) error {
	personID := c.Param("id")
	ret, err := executeSQL("select p.id,p.password,p.last_login,p.is_superuser,p.username,p.first_name,p.last_name,p.email,p.is_staff,p.is_active,p.date_joined from person p left join person_friends pf on p.id = pf.to_person_id where pf.from_person_id = ?", personID)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]string{"error": "Something went wrong when getting data from db"})
	}

	fmt.Println("request to find all friends of a person")
	return c.JSON(http.StatusOK, ret)
}

func executeQuery(c echo.Context) error {
	query := c.Param("query")
	result := executeGraphQL(query)
	if result == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"msg": "graphql error"})
	}
	fmt.Printf("query result: %v\n", result)
	return c.JSON(http.StatusOK, result)
}

func executeGraphQL(query string) *graphql.Result {
	var client = &schema.Client{}

	ctx := context.WithValue(context.Background(), util.LoadersKey, schema.BatchedLoaders)
	ctx = context.WithValue(ctx, util.ClientKey, client)

	if schema, err := schema.CreateSchema(); err == nil {
		result := graphql.Do(graphql.Params{
			Context: ctx,
			Schema:        *schema,
			RequestString: query,
		})

		if len(result.Errors) > 0 {
			fmt.Printf("\n errors: %v\n", result.Errors)
		}

		return result
	}
	return nil
}

// Execute sql statement from parameter, which looks like this:
// select a, b, c from some_tabble where id = ?
// Return a map
func executeSQL(sqlStmt string, personID string) ([]model.Person, error) {
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer stmt.Close()

	var ret *sql.Rows
	if strings.Compare(personID, "")  == 0 {
		ret, err = stmt.Query()
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
	} else {
		ret, err = stmt.Query(personID)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
	}

	personList := make([]model.Person, 0)
	for ret.Next() {
		var person model.Person

		var personID int64
		var password sql.NullString
		var lastLogin mysql.NullTime
		var isSuperuser sql.NullBool
		var userName sql.NullString
		var firstName sql.NullString
		var lastName sql.NullString
		var email sql.NullString
		var isStaff sql.NullBool
		var isActive sql.NullBool
		var dateJoined mysql.NullTime

		err = ret.Scan(
			&personID,
			&password,
			&lastLogin,
			&isSuperuser,
			&userName,
			&firstName,
			&lastName,
			&email,
			&isStaff,
			&isActive,
			&dateJoined,
		)

		if err != nil {
			log.Fatal(err)
		}

		person.ID = personID
		person.Password = util.If(password.Valid, password.String, "").(string)
		if tempTime, ok := util.If(lastLogin.Valid, lastLogin.Time, nil).(*time.Time); ok {
			person.LastLogin = tempTime
		} else {
			person.LastLogin = nil
		}

		person.IsSuperuser = util.If(isSuperuser.Valid, isSuperuser.Bool, false).(bool)
		person.Username = util.If(userName.Valid, userName.String, "").(string)
		person.FirstName = util.If(firstName.Valid, firstName.String, "").(string)
		person.LastName = util.If(lastName.Valid, lastName.String, "").(string)
		person.Email = util.If(email.Valid, email.String, "").(string)
		person.IsStaff = util.If(isStaff.Valid, isStaff.Bool, false).(bool)
		person.IsActive = util.If(isActive.Valid, isActive.Bool, false).(bool)
		if tempTime, ok := util.If(dateJoined.Valid, dateJoined.Time, nil).(*time.Time); ok {
			person.DateJoined = tempTime
		} else {
			person.DateJoined = nil
		}

		// Get friends
		friends, err := executeSQLString("select to_person_id from person_friends where from_person_id=?", 
											strconv.FormatInt(personID, 10))
		if err == nil {
			person.Friends = friends
		}

		personList = append(personList, person)
	}

	return personList, nil
}

func executeSQLString(sqlStmt string, personID string) ([]int64, error) {
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer stmt.Close()

	var ret *sql.Rows
	if strings.Compare(personID, "")  == 0 {
		ret, err = stmt.Query()
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
	} else {
		ret, err = stmt.Query(personID)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
	}

	var friendIDs []int64
	for ret.Next() {
		var friendID sql.NullInt64

		err = ret.Scan(&friendID)

		if err != nil {
			util.ULog(err)
			return nil, err
		}

		if friendID.Valid {
			friendIDs = append(friendIDs, friendID.Int64)
		}
	}

	return friendIDs, nil
}

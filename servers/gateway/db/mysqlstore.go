package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/INFO441-19au-org/assignments-bch0ng/servers/gateway/models/users"
)

const (
	host     = "127.0.0.1"
	port     = 3306
	user     = "postgres"
	password = "your-password"
	dbname   = "psql_db"
)

type MySQLStore struct {
	MySQLDB *sql.DB
}

func scanRowsIntoUser(rows *sql.Rows, err error, store *MySQLStore) (*users.User, error) {
	concreteUser := users.User{}
	newUser := &concreteUser
	for rows.Next() {
		if err := rows.Scan(&newUser.ID, &newUser.Email, &newUser.PassHash, &newUser.UserName,
			&newUser.FirstName, &newUser.LastName, &newUser.PhotoURL); err != nil {
			fmt.Printf("error scanning row: %v\n", err)
		}
	}
	if err := rows.Err(); err != nil {
		errorStr := fmt.Sprintf("error getting next row: %v\n", err)
		return newUser, errors.New(errorStr)
	}
	return newUser, nil
}
func performGetQuery(queryStr string, store *MySQLStore) (*users.User, error) {
	rows, err := store.MySQLDB.Query(queryStr)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	} else {
		newUser, newErr := scanRowsIntoUser(rows, err, store)
		if newErr != nil {
			log.Fatal(err)
		}
		return newUser, nil
	}
	return nil, errors.New("Failed Get Query")
}

func (store *MySQLStore) GetByID(id int64) (*users.User, error) {
	query := fmt.Sprintf("select * from users where id=%d limit 1", id)
	return performGetQuery(query, store)
}

func (store *MySQLStore) GetByUserName(username string) (*users.User, error) {
	query := fmt.Sprintf("select * from users where user_name='%s' limit 1", username)
	return performGetQuery(query, store)
}

func (store *MySQLStore) Insert(user *users.User) (*users.User, error) {
	sqlStatement := `INSERT INTO users (passhash, username) VALUES (?,?);`
	res, err := store.MySQLDB.Exec(sqlStatement, user.PassHash, user.UserName)
	if err != nil {
		fmt.Printf("error inserting new row: %v\n", err)
	} else {
		//get the auto-assigned ID for the new row
		lastInsertId, err := res.LastInsertId()
		if err != nil {
			fmt.Printf("error reading new row: %v\n", err)
		}
		query := fmt.Sprintf("select * from users where id=%v limit 1;", lastInsertId)
		queryResults, err := performGetQuery(query, store)
		return queryResults, err
	}
	return nil, errors.New("Failed Insert")

}

func (store *MySQLStore) Delete(id int64) error {
	sqlStatement := `
  Delete from users where id=?`
	res, err := store.MySQLDB.Exec(sqlStatement, id)
	if err != nil {
		return err
	}
	fmt.Printf("Res found: %v", res)
	fmt.Printf("Res deleted: %v", id)
	return nil
}

func ConnectToPostgres(dsn string) (*MySQLStore, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("error opening database: %v\n", err)
		os.Exit(1)
	}
	newMySQLStore := &MySQLStore{MySQLDB: db}
	return newMySQLStore, nil
}

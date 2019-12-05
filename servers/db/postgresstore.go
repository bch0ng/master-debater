package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	. "github.com/INFO441-19au-org/assignments-bch0ng/servers/gateway/models/users"

	//_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
)

const (
	//host     = "localhost"
	host     = "127.0.0.1"
	port     = 3306
	user     = "postgres"
	password = "your-password"
	dbname   = "psql_db"
)

// PostgressStore stores the database pointer
type PostgressStore struct {
	PostgressDB *sql.DB
}

// Scan rows into the user struct.
func scanRowsIntoUser(rows *sql.Rows, err error, store *PostgressStore) (*User, error) {
	//while there are more rows
	concreteUser := User{}
	newUser := &concreteUser
	for rows.Next() {
		if err := rows.Scan(&newUser.ID, &newUser.Email, &newUser.PassHash, &newUser.UserName,
			&newUser.FirstName, &newUser.LastName, &newUser.PhotoURL); err != nil {
			fmt.Printf("error scanning row: %v\n", err)
		}
	}

	//if we got an error fetching the next row, report it
	if err := rows.Err(); err != nil {
		errorStr := fmt.Sprintf("error getting next row: %v\n", err)
		return newUser, errors.New(errorStr)
	}
	return newUser, nil
}

// SQL get request query helper
func preformGetQuery(queryStr string, store *PostgressStore) (*User, error) {
	rows, err := store.PostgressDB.Query(queryStr)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	} else {
		newUser, newErr := scanRowsIntoUser(rows, err, store)
		if newErr != nil {
			log.Fatal(err)
		}
		fmt.Printf("Get query success\n")
		return newUser, nil
	}
	return nil, errors.New("Failed Get")
}

// Get user by given ID. Will return a user struct and error (nil if no error).
func (store *PostgressStore) GetByID(id int64) (*User, error) {
	fmt.Printf("Attempting get id!\n")
	query := fmt.Sprintf("select * from contacts where id=%d limit 1", id)
	return preformGetQuery(query, store)
}

// Get user by given email. Will return a user struct and error (nil if no error).
func (store *PostgressStore) GetByEmail(email string) (*User, error) {
	query := fmt.Sprintf("select * from contacts where email='%s' limit 1", email)
	return preformGetQuery(query, store)
}

// Get user by given username. Will return a user struct and error (nil if no error).
func (store *PostgressStore) GetByUserName(username string) (*User, error) {
	query := fmt.Sprintf("select * from contacts where user_name='%s' limit 1", username)
	return preformGetQuery(query, store)
}

// Insert a contact with the given user. Will return a user struct and an error (nil if no error).
func (store *PostgressStore) Insert(user *User) (*User, error) {
	sqlStatement := `
  INSERT INTO contacts (email, passhash, user_name, first_name, last_name, photo_url)
  VALUES (?,?,?,?,?,?);`
	res, err := store.PostgressDB.Exec(sqlStatement, user.Email, user.PassHash, user.UserName,
		user.FirstName, user.LastName, user.PhotoURL)
	if err != nil {
		fmt.Printf("error inserting new row: %v\n", err)
	} else {
		//get the auto-assigned ID for the new row
		fmt.Printf("Success inserting new row:")
		lastInsertId, err := res.LastInsertId()
		if err != nil {
			fmt.Printf("error reading new row: %v\n", err)
		}
		query := fmt.Sprintf("select * from contacts where id=%v limit 1;", lastInsertId)
		queryResults, err := preformGetQuery(query, store)
		return queryResults, err
	}
	return nil, errors.New("Failed Insert")

}

// Update a contact with the given ID and update fields. Will return nil or an error.
func (store *PostgressStore) Update(id int64, updates *Updates) (*User, error) {
	sqlStatement := `
  UPDATE contacts SET first_name=?, last_name=? where id=?;`
	res, err := store.PostgressDB.Exec(sqlStatement, updates.FirstName, updates.LastName, id)
	fmt.Printf("Res:%v", res)
	if err != nil {
		return nil, errors.New("Failed Update")
	}
	query := fmt.Sprintf("select * from contacts where id=%v limit 1;", id)
	queryResults, err := preformGetQuery(query, store)
	return queryResults, err
}

// Delete a contact with the given ID. Will return nil or an error.
func (store *PostgressStore) Delete(id int64) error {
	sqlStatement := `
  Delete from contacts where id=?`
	res, err := store.PostgressDB.Exec(sqlStatement, id)
	if err != nil {
		return err
	}
	fmt.Printf("Res found: %v", res)
	fmt.Printf("Res deleted: %v", id)
	return nil
}

// AddToLog inserts a user session and returns nil or an error.
func (store *PostgressStore) AddToLog(userId string, ipAddress string) error {
	sqlStatement := `
	INSERT INTO sessions (id, sign_in_time, ip)
	VALUES (?, ?, ?);`
	ts := time.Now().Format("2006-01-02 15:04:05")
	res, err := store.PostgressDB.Exec(sqlStatement, userId, ts, ipAddress)
	query := fmt.Sprintf("Log results=%d;", res)
	fmt.Printf(query)
	if err != nil {
		fmt.Printf("error inserting new session: %v\n", err)
	}
	return nil
}

// ConnectToPostgres accepts a DSN string to initialize the
// SQL databaze connection and returns the database struct
// and an error (nil if no error).
func ConnectToPostgres(dsn string) (*PostgressStore, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		fmt.Printf("error opening database: %v\n", err)
		os.Exit(1)
	}
	newPostgressStore := &PostgressStore{PostgressDB: db}
	//ensure that the database gets closed when we are done
	//defer db.Close()

	//for now, just ping the server to ensure we have
	//a live connection to it
	err2 := db.Ping()
	if err2 != nil {
		fmt.Printf("PING ERROR:%v", err2)
		fmt.Printf("Database Not Available When Pinged")
		os.Exit(1)
	} else {
		fmt.Printf("successfully connected!\n")
	}
	return newPostgressStore, nil
}

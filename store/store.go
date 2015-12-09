package store

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3" // for database/sql driver
	"log"
)

type Users struct {
	User  string
	Pass  string
	About string
	Pic   string
}

// SqliteDB is a wrapper for the Sqlite3 database store
type SqliteDB struct{ *sql.DB }

// Init opens the database and sets up the tables if not already created
func (db *SqliteDB) Init(filename string) {
	var err error
	db.DB, err = sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatalln("open database:", err)
	}

	sqlStmt := `CREATE TABLE IF NOT EXISTS USERS( 
	USER STRING, 
	PASS STRING,
	ABOUT TEXT,
	PIC STRING
);`

	if _, err := db.Exec(sqlStmt); err != nil {
		log.Printf("Error creating table: %s", err)
	}
}

// AddUser inserts user data into the database
func AddUser(db SqliteDB, u Users) (err error) {

	sqlStmt := "INSERT INTO USERS (USER,PASS,ABOUT,PIC) VALUES(?,?,?,?)"
	_, err = db.Exec(sqlStmt, u.User, u.Pass, u.About, u.Pic)
	return err
}

// Update user data
func UpdateUser(db SqliteDB, u Users) (err error) {

	sqlStmt := "UPDATE USERS SET USER=?,PASS=?,ABOUT=?,PIC=?"
	_, err = db.Exec(sqlStmt, u.User, u.Pass, u.About, u.Pic)
	return err
}

//Query Profile info
func QueryProfile(db SqliteDB, User string) Users {

	sqlStmt := `SELECT USER,PASS,ABOUT,PIC FROM USERS WHERE USER = ?;`
	r := db.QueryRow(sqlStmt, User)
	u := Users{}
	err := r.Scan(&u.User, &u.Pass, &u.About, &u.Pic)

	if err != nil {
		return Users{
			User:  "",
			Pass:  "",
			About: "",
			Pic:   "",
		}
		log.Printf("Querying password failed: %s", err)
	}
	return u
}

//Query Password
func QueryPass(db SqliteDB, user string) string {
	var Pass string
	sqlStmt := `SELECT PASS FROM USERS WHERE USER = ?;`
	r := db.QueryRow(sqlStmt, user)

	err := r.Scan(&Pass)
	if err != nil {
		if err == sql.ErrNoRows {
			// there were no rows, but otherwise no error occurred
			log.Println("Login fail user doesn't now exist")
			return "1"
		} else {
			log.Println("Password query error:", err)
		}
	}
	return Pass
}

package irc

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"time"
)

type DBLink struct {
	key         sql.NullInt64
	link        sql.NullString
	sender_id   sql.NullInt64
	date        time.Time
	link_type   sql.NullString
	sender_name sql.NullString
}

type DBUser struct {
	key       sql.NullInt64
	user_name sql.NullString
}

// Log table action consts
const (
	Duplicate = "Duplicate"
)

type IrcDatabase struct {
	checked bool
}

func (i *IrcDatabase) Open() *sql.DB {
	config := GetConfig()
	if config.DBFile == "" {
		fmt.Println("Database path not set!")
		os.Exit(3)
	}
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?parseTime=true", config.DBFile))
	if err != nil {
		fmt.Println("Cannot open database")
		os.Exit(3)
	}
	if i.checked == false {
		i.create(db)
	}
	return db
}

func (i *IrcDatabase) create(db *sql.DB) {
	transaction, _ := db.Begin()
	db.Exec(`CREATE TABLE IF NOT EXISTS irc_link (
			id integer PRIMARY KEY,
			link varchar(2048) NOT NULL,
			sender_id integer NOT NULL,
			date datetime DEFAULT CURRENT_TIMESTAMP,
			link_type varchar(255) NOT NULL
		)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS irc_user(
			id integer PRIMARY KEY,
			user_name varchar(255) NOT NULL
		)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS irc_log(
			id integer PRIMARY KEY,
			action varchar(255) NOT NULL,
			sender_id integer NOT NULL,
			date datetime DEFAULT CURRENT_TIMESTAMP
		)`)
	transaction.Commit()
	i.checked = true
}

func (i *IrcDatabase) GetLink(url string) *DBLink {
	db := i.Open()
	defer db.Close()
	row := db.QueryRow("SELECT link.*, user.user_name FROM irc_link AS link LEFT JOIN irc_user AS user on (user.id = link.sender_id) WHERE link = ?", url)
	link := DBLink{}
	err := row.Scan(&link.key, &link.link, &link.sender_id, &link.date, &link.link_type, &link.sender_name)
	if err != nil {
		return &DBLink{}
	}
	return &link
}

func (i *IrcDatabase) GetCreateUser(userName string) *DBUser {
	db := i.Open()
	defer db.Close()
	row := db.QueryRow("SELECT * FROM irc_user WHERE user_name = ?", userName)
	user := DBUser{}
	err := row.Scan(&user.key, &user.user_name)
	if err == sql.ErrNoRows {
		stmt, err := db.Prepare("INSERT INTO irc_user(user_name) VALUES (?)")
		transaction, err := db.Begin()
		checkErr(err)
		res, err := transaction.Stmt(stmt).Exec(userName)
		checkErr(err)
		err = transaction.Commit()
		checkErr(err)
		id, err := res.LastInsertId()
		checkErr(err)
		return &DBUser{sql.NullInt64{id, true}, sql.NullString{userName, true}}
	}
	return &user
}

func (i *IrcDatabase) AddLink(link, linkType, username string) {
	user := i.GetCreateUser(username)
	db := i.Open()
	defer db.Close()
	transaction, _ := db.Begin()
	stmt, err := transaction.Prepare("INSERT INTO irc_link(link, sender_id, link_type) VALUES (?,?,?)")
	checkErr(err)
	transaction.Stmt(stmt).Exec(link, user.key.Int64, linkType)
	transaction.Commit()
}

func (i *IrcDatabase) LogDuplicate(link *DBLink, sender string) {
	user := i.GetCreateUser(sender)
	db := i.Open()
	defer db.Close()
	transaction, _ := db.Begin()
	stmt, err := transaction.Prepare("INSERT INTO irc_log(action, sender_id) VALUES (?,?)")
	checkErr(err)
	transaction.Stmt(stmt).Exec(Duplicate, user.key.Int64)
	transaction.Commit()
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

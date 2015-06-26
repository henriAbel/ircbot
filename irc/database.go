package irc

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strings"
	"time"
)

type DBLink struct {
	Key         sql.NullInt64
	Link        sql.NullString
	Sender_id   sql.NullInt64
	Date        time.Time
	Link_type   sql.NullString
	Sender_name sql.NullString
}

type DBUser struct {
	Key       sql.NullInt64
	User_name sql.NullString
}

type DBRaw struct {
	Mime_type sql.NullString
	Data      []byte
	Link_id   sql.NullInt64
	Data_type sql.NullString
}

// Those strings are used in database type fields
const (
	Duplicate = "Duplicate"
	// Gif first frame as jpeg image
	RawGifFrame       = "gif1"
	RawGif            = "gif"
	RawImage          = "image"
	RawImageThumbnail = "thumb"
	database_version  = 1
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
		i.checked = true
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
	db.Exec(`CREATE TABLE IF NOT EXISTS irc_raw(
			id integer PRIMARY KEY,
			link_id integer NOT NULL,
			data_type varchar(255) NOT NULL,
			mime_type varchar(255) NOT NULL,
			data BLOB
		)`)
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS irc_raw_link_id_type_index ON irc_raw(link_id, data_type)")
	/*current_version := i.GetDatabaseVersion()
	if current_version != database_version {
		for current_version < database_version {
			switch current_version {
			case 0:

			}

			current_version++
		}
	}*/
	transaction.Commit()
}

func (i *IrcDatabase) GetLink(url string) *DBLink {
	db := i.Open()
	defer db.Close()
	row := db.QueryRow("SELECT link.*, user.user_name FROM irc_link AS link LEFT JOIN irc_user AS user on (user.id = link.sender_id) WHERE link = ?", url)
	return rowToDBLink(row)
}

func (i *IrcDatabase) GetLinkById(id int64) *DBLink {
	db := i.Open()
	defer db.Close()
	row := db.QueryRow("SELECT link.*, user.user_name FROM irc_link AS link LEFT JOIN irc_user AS user on (user.id = link.sender_id) WHERE link.id = ?", id)
	return rowToDBLink(row)
}

func (i *IrcDatabase) GetCreateUser(userName string) *DBUser {
	db := i.Open()
	defer db.Close()
	row := db.QueryRow("SELECT * FROM irc_user WHERE user_name = ?", userName)
	user := DBUser{}
	err := row.Scan(&user.Key, &user.User_name)
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

func (i *IrcDatabase) AddLink(link, linkType, username string) *DBLink {
	user := i.GetCreateUser(username)
	db := i.Open()
	defer db.Close()
	transaction, _ := db.Begin()
	stmt, err := transaction.Prepare("INSERT INTO irc_link(link, sender_id, link_type) VALUES (?,?,?)")
	checkErr(err)
	res, _ := transaction.Stmt(stmt).Exec(link, user.Key.Int64, linkType)
	transaction.Commit()
	id, _ := res.LastInsertId()
	return &DBLink{
		sql.NullInt64{id, true},
		sql.NullString{link, true},
		sql.NullInt64{user.Key.Int64, true},
		time.Now(),
		sql.NullString{linkType, true},
		sql.NullString{username, true}}
}

func (i *IrcDatabase) LogDuplicate(link *DBLink, sender string) {
	user := i.GetCreateUser(sender)
	db := i.Open()
	defer db.Close()
	transaction, _ := db.Begin()
	stmt, err := transaction.Prepare("INSERT INTO irc_log(action, sender_id) VALUES (?,?)")
	checkErr(err)
	transaction.Stmt(stmt).Exec(Duplicate, user.Key.Int64)
	transaction.Commit()
}

func (i *IrcDatabase) AddRaw(dataType string, link *DBLink, mime_type string, data []byte) {
	db := i.Open()
	defer db.Close()
	transaction, _ := db.Begin()
	stmt, err := transaction.Prepare("INSERT INTO irc_raw(link_id, data_type, mime_type, data) VALUES (?,?,?,?)")
	checkErr(err)
	transaction.Stmt(stmt).Exec(link.Key.Int64, dataType, mime_type, data)
	transaction.Commit()
}

func (i *IrcDatabase) GetAll(limit int, offset int, link_types []string) *[]DBLink {
	db := i.Open()
	defer db.Close()
	if len(link_types) == 0 {
		link_types = []string{Image, Gif, Link, Youtube}
	}

	params := string_to_interface(&link_types)

	params[len(params)-2] = limit
	params[len(params)-1] = offset

	rows, err := db.Query(fmt.Sprintf("SELECT link.*, user.user_name FROM irc_link AS link LEFT JOIN irc_user AS user on (user.id = link.sender_id) WHERE link.link_type IN (%s) ORDER BY date DESC LIMIT ? OFFSET ?", generatePlaceHolders(&link_types)), params...)
	if nil != err {
		fmt.Println(err)
	}
	var links []DBLink
	for rows.Next() {
		link := DBLink{}
		rows.Scan(&link.Key, &link.Link, &link.Sender_id, &link.Date, &link.Link_type, &link.Sender_name)
		links = append(links, link)
	}
	return &links
}

func (i *IrcDatabase) GetCount(link_types []string) int {
	db := i.Open()
	defer db.Close()
	if len(link_types) == 0 {
		link_types = []string{Image, Gif, Link, Youtube}
	}

	params := string_to_interface(&link_types)
	row := db.QueryRow(fmt.Sprintf("SELECT count(*) FROM irc_link AS link LEFT JOIN irc_user AS user on (user.id = link.sender_id) WHERE link.link_type IN (%s)", generatePlaceHolders(&link_types)), params...)
	var count = 0
	row.Scan(&count)
	return count
}

func (i *IrcDatabase) GetRaw(linkId int64, rawType string) *DBRaw {
	db := i.Open()
	defer db.Close()
	row := db.QueryRow("SELECT mime_type, data, link_id, data_type FROM irc_raw WHERE link_id = ? AND data_type = ?", linkId, rawType)
	var raw = DBRaw{}
	row.Scan(&raw.Mime_type, &raw.Data, &raw.Link_id, &raw.Data_type)
	return &raw
}

func (i *IrcDatabase) GetDatabaseVersion() int {
	db := i.Open()
	defer db.Close()
	row := db.QueryRow("PRAGMA user_version")
	var version = 0
	row.Scan(&version)
	return version
}

func (i *IrcDatabase) setDatabaseVersion(version int) {
	db := i.Open()
	defer db.Close()
	transaction, _ := db.Begin()
	db.Exec("PRAGMA user_version=?", version)
	transaction.Commit()
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func generatePlaceHolders(list *[]string) string {
	placeholder := make([]string, len(*list))
	for i, _ := range placeholder {
		placeholder[i] = "?"
	}
	return strings.Join(placeholder, ",")
}

func rowToDBLink(row *sql.Row) *DBLink {
	link := DBLink{}
	err := row.Scan(&link.Key, &link.Link, &link.Sender_id, &link.Date, &link.Link_type, &link.Sender_name)
	if err != nil {
		return &DBLink{}
	}
	return &link
}

func string_to_interface(link_types *[]string) []interface{} {
	// Better, less hacky way??
	params := make([]interface{}, len(*link_types)+2)
	for i, s := range *link_types {
		params[i] = s
	}
	return params
}

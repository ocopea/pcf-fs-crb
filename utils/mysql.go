// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package utils

import (
	"crb/models"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// TEST is set to true when testing on DEV VM, and will cause to use a hardcoded DB url
var TEST = false

var sqlOpen = sql.Open

// DatetimeIso8601Rfc3339Format defines ISO 8601/RFC 3339 format string for MySQL
// Use golang time.RFC3339 for same time formatting
const DatetimeIso8601Rfc3339Format = "%Y-%m-%dT%TZ"

// Database is the production mysql database
type Database struct {
	db string // Pointer to a type of SQL Database
}

// OpenConnection returns a DB connection, else error.
// The close() should be handled by the callee.
func OpenConnection() (*sql.DB, error) {

	if TEST {
		return sql.Open("mysql", "root:root@tcp(localhost:3306)/test")
	}

	var err error

	creds, err := fetchDBCredentialsFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("Unable to fetch database credentials: %s", err)
	}

	//Parse the retrieved URL to break apart the individual components to assemble the connection string
	u, err := url.Parse(creds.URL)
	if err != nil {
		return nil, fmt.Errorf("Error parsing URL from credentials, URL=%s : %s", creds.URL, err)
	}
	connectString := fmt.Sprintf("%s:%s@tcp(%s)%s", creds.Username, creds.Password, u.Host, u.EscapedPath())

	sqldb, err := sqlOpen("mysql", connectString)
	if err != nil {
		return nil, fmt.Errorf("Unable to open connection to database: %s", err)
	}

	//if DB is not nil, check that connection still alive and if not alive, get connection, use sql.Ping to verify and re-establish.
	//connection not being nil doesnt ensure it's 'open' as the sql.Open function only guarantee an endpoint
	err = sqldb.Ping()
	if err != nil {
		return nil, fmt.Errorf("Unable to confirm connection is alive, db.Ping error: %s", err)
	}

	return sqldb, err
}

// CloseConnection is a no-op
func CloseConnection(sqldb *sql.DB) {
	err := sqldb.Close()
	if err != nil {
		log.Println(err)
	}
	log.Println("Open db connections :", sqldb.Stats())
}

// CreateRepositoryTable returns error
// - if the DB can't be accessed
// - if the table can't be created
func (db *Database) CreateRepositoryTable() error {
	query := "CREATE TABLE IF NOT EXISTS TARGET_CREDENTIALS (addr VARCHAR(255), user VARCHAR(255), password VARCHAR(255), PRIMARY KEY ( addr ))"
	dbc, err := OpenConnection()
	if err != nil {
		return err
	}
	defer CloseConnection(dbc)

	_, err = dbc.Exec(query)
	return err
}

// CreateCopyTable returns error
// - if the DB can't be accessed
// - if the table can't be created
func (db *Database) CreateCopyTable() error {
	query := "CREATE TABLE IF NOT EXISTS COPY_REPOSITORY (ID VARCHAR(255), date DATE, path VARCHAR(255), fileName VARCHAR(255), fileSize VARCHAR(255), PRIMARY KEY ( ID ))"
	dbc, err := OpenConnection()
	if err != nil {
		return err
	}
	defer dbc.Close()

	_, err = dbc.Exec(query)
	return err
}

// DoesRepositoryTableExist returns error
// - if the DB can't be accessed
func (db *Database) DoesRepositoryTableExist() (bool, error) {
	dbc, err := OpenConnection()
	if err != nil {
		return false, err
	}
	defer CloseConnection(dbc)

	rows, err := dbc.Query("SELECT COUNT(*) FROM TARGET_CREDENTIALS ")
	if err != nil {
		return false, err
	}
	defer CloseDbRows(rows)

	return true, nil
}

// DoesCopyTableExist returns error
// - if the DB can't be accessed
func (db *Database) DoesCopyTableExist() (bool, error) {
	return false, nil
}

// AddRepository returns error
// - if the DB can't be accessed
// - or any other error like out of space
// Note, if the table doesn't exist, it would be created and the repo added
// Currently there is only one repo supported. Adding a new one simply replaces the one in the DB.
func (db *Database) AddRepository(repoInfo *models.RepositoryInfo) error {
	query := "INSERT INTO TARGET_CREDENTIALS VALUES ('" + *repoInfo.Addr + "','" + *repoInfo.User + "','" + repoInfo.Password + "')"

	dbc, err := OpenConnection()
	if err != nil {
		return err
	}
	defer CloseConnection(dbc)

	rows, err := dbc.Query("SELECT addr FROM TARGET_CREDENTIALS")
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		deletedRows, err := dbc.Query("DELETE FROM TARGET_CREDENTIALS")
		if err != nil {
			return err
		}
		defer CloseDbRows(deletedRows)

		_, err = dbc.Exec(query)
		return err
	}
	_, err = dbc.Exec(query)
	return err
}

func handleError(err error) (*models.RepositoryInfo, error) {
	if err != nil {
		log.Println(err)
	}
	return nil, err
}

// GetRepository returns the repository entry and return an error in the following cases
// - if the DB can't be accessed
// - if the repository table doesn't exist.
func (db *Database) GetRepository() (*models.RepositoryInfo, error) {

	repoInfo := &models.RepositoryInfo{}

	sqldb, err := OpenConnection()
	if err != nil {
		return handleError(err)
	}
	defer CloseConnection(sqldb)

	exists, err := db.DoesRepositoryTableExist()

	if !exists {
		return handleError(err)
	}

	rows, err := sqldb.Query("select * from TARGET_CREDENTIALS")
	if err != nil {
		return handleError(err)
	}
	defer CloseDbRows(rows)

	for rows.Next() {
		err := rows.Scan(&repoInfo.Addr, &repoInfo.User, &repoInfo.Password)
		if err != nil {
			return handleError(err)
		}
	}

	return repoInfo, err
}

// GetCopy returns error
// - if the DB can't be accessed
// - if the id can't be found
// - if the table doesn't exist.
func (db *Database) GetCopy(id string) (*models.CopyMetaData, error) {
	dbc, err := OpenConnection()
	if err != nil {
		return nil, err
	}
	defer CloseConnection(dbc)

	query := "SELECT DATE_FORMAT(date, '" + DatetimeIso8601Rfc3339Format + "'), path, fileSize FROM COPY_REPOSITORY " + " WHERE ID = '" + id + "'"
	rows, err := dbc.Query(query)
	if err != nil {
		return nil, err
	}
	defer CloseDbRows(rows)

	metaData := &models.CopyMetaData{
		CopyID: id,
	}
	for rows.Next() {
		err := rows.Scan(&metaData.CopyTimeStamp, &metaData.CopyDataURL, &metaData.CopySize)
		if err != nil {
			return nil, err
		}
		return metaData, nil
	}

	return nil, fmt.Errorf("copyID doesn't exist %s", id)
}

// GetCopies returns error
// - if the DB can't be accessed
// - if the table doesn't exist
func (db *Database) GetCopies() ([]string, error) {
	dbc, err := OpenConnection()
	if err != nil {
		return nil, err
	}
	defer CloseConnection(dbc)

	var idList []string
	query := "SELECT ID FROM COPY_REPOSITORY"
	rows, err := dbc.Query(query)
	if err != nil {
		return nil, err
	}
	defer CloseDbRows(rows)

	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		idList = append(idList, id)
	}

	return idList, nil
}

// AddCopy returns error
// - if the DB can't be accessed
// - if the CopyInfo object is invalid
// - if the CopyInfo already exists
// - or any other error like out of space
// Note, if the table doesn't exist, it would be created and and the copy saved.
func (db *Database) AddCopy(c *models.CopyMetaData) error {
	query := fmt.Sprintf("INSERT INTO COPY_REPOSITORY VALUES ('%s', NOW(), '%s', '%s', '%d')", c.CopyID, RemotePath, c.CopyID, c.CopySize)

	dbc, err := OpenConnection()
	if err != nil {
		return err
	}
	defer CloseConnection(dbc)

	_, err = dbc.Exec(query)
	return err
}

// DoesCopyExist checks if copyId exists in the copy table or not.
// If Copy Table does not exist, this function creates one.
func (db *Database) DoesCopyExist(copyID string) (bool, error) {
	if err := db.CreateCopyTable(); err != nil {
		return false, err
	}

	if _, err := db.GetCopy(copyID); err != nil {
		if strings.Contains(err.Error(), "copyID doesn't exist") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// DeleteCopy returns error
// - if the DB can't be accessed
// - if the id can't be found
// - if the table doesn't exist,
func (db *Database) DeleteCopy(id string) error {
	query := "DELETE from COPY_REPOSITORY where ID = ('" + id + "')"

	dbc, err := OpenConnection()
	if err != nil {
		return err
	}
	defer CloseConnection(dbc)

	_, err = dbc.Exec(query)

	return err
}

// GetDatabase returns a new instance of the MySql DB interface.
func GetDatabase() DBInterface {
	db := &Database{db: "MySQL"}
	return db
}

func CloseDbRows(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		log.Println(err)
	}
}

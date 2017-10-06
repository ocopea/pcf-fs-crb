// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package utils

import (
	"database/sql"
	"errors"
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestGetDatabaseConnectionNoVCAP(t *testing.T) {
	// set oldGetEnv to real getENV
	oldEnv := getEnv

	// as we are exiting, getEnv revert back
	defer func() { getEnv = oldEnv }()

	getEnv = func(string) string {
		return ""
	}

	expected := errors.New("Unable to fetch database credentials: Unable to retrieve VCAP_SERVICES from environment")
	_, err := OpenConnection()
	if err == nil {
		t.Errorf("OpenConnection() Expected error, but did not recieve one")
	}
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("OpenConnection() mismatch: recieved error = %v, expected error = %v", err, expected)

	}
}

func TestGetDatabaseConnectionBadURL(t *testing.T) {
	// set oldGetEnv to real getENV
	oldEnv := getEnv

	// as we are exiting, getEnv revert back
	defer func() { getEnv = oldEnv }()

	getEnv = func(string) string {
		return BadURLENV
	}

	expected := errors.New(`Error parsing URL from credentials, URL=http://[fe80::%31%25en0]:8080/ : parse http://[fe80::%31%25en0]:8080/: invalid URL escape "%31"`)
	_, err := OpenConnection()
	if err == nil {
		t.Errorf("OpenConnection() Expected error, but did not recieve one")
	}
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("OpenConnection() mismatch: recieved error = %v, expected error = %v", err, expected)

	}
}

func TestGetDatabaseConnectionOpenError(t *testing.T) {
	// set oldGetEnv to real getENV
	oldEnv := getEnv
	oldOpen := sqlOpen

	// as we are exiting, revert both env and sqlopen
	defer func() { getEnv = oldEnv }()
	defer func() { sqlOpen = oldOpen }()

	getEnv = func(string) string {
		return ENV
	}
	sqlOpen = func(driver, conn string) (*sql.DB, error) {
		return nil, errors.New("failed to connect to db")
	}

	expected := errors.New("Unable to open connection to database: failed to connect to db")
	_, err := OpenConnection()
	if err == nil {
		t.Errorf("OpenConnection() Expected error, but did not recieve one")
	}
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("OpenConnection() mismatch: recieved error = %v, expected error = %v", err, expected)

	}
}

func TestGetDatabaseConnectionCantPing(t *testing.T) {
	// set oldGetEnv to real getENV
	oldEnv := getEnv

	// as we are exiting, getEnv revert back
	defer func() { getEnv = oldEnv }()

	getEnv = func(string) string {
		return BadPingENV
	}

	expected := errors.New(`Unable to confirm connection is alive, db.Ping error: invalid DSN: missing the slash separating the database name`)
	_, err := OpenConnection()
	if err == nil {
		t.Errorf("OpenConnection() Expected error, but did not recieve one")
	}
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("OpenConnection() mismatch: recieved error = %v, expected error = %v", err, expected)

	}
}

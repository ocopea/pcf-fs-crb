// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package utils

import (
	"errors"
	"reflect"
	"testing"
)

// test happy path.
// test that db.Ping throws an error
var UserProvidedEnv = "0.23"
var UserProvidedEnvEmpty = ""

var ENV = `{
       "p-mysql": [
       {
         "credentials": {
          "hostname": "mysql-broker.local.pcfdev.io",
          "jdbcUrl": "jdbc:mysql://mysql-broker.local.pcfdev.io:3306/cf_08fd31cb_3fde_4f4d_a86d_db1c04e4ccd1?user=xvQVjibAbTJoJ5Nd\u0026password=5uZ566tecLueLX2H",
          "name": "cf_08fd31cb_3fde_4f4d_a86d_db1c04e4ccd1",
          "password": "5uZ566tecLueLX2H",
          "port": 3306,
         "uri": "mysql://xvQVjibAbTJoJ5Nd:5uZ566tecLueLX2H@mysql-broker.local.pcfdev.io:3306/cf_08fd31cb_3fde_4f4d_a86d_db1c04e4ccd1?reconnect=true",
         "username": "xvQVjibAbTJoJ5Nd"
         },
      "label": "p-mysql",
       "name": "test_db",
       "plan": "1gb",
        "provider": null,
        "syslog_drain_url": null,
        "tags": [
         "mysql"
       ],
        "volume_mounts": []
        }
       ]
      }`

var BadURLENV = `{
       "p-mysql": [
       {
         "credentials": {
          "hostname": "mysql-broker.local.pcfdev.io",
          "jdbcUrl": "jdbc:mysql://mysql-broker.local.pcfdev.io:3306/cf_08fd31cb_3fde_4f4d_a86d_db1c04e4ccd1?user=xvQVjibAbTJoJ5Nd\u0026password=5uZ566tecLueLX2H",
          "name": "cf_08fd31cb_3fde_4f4d_a86d_db1c04e4ccd1",
          "password": "5uZ566tecLueLX2H",
          "port": 3306,
         "uri": "http://[fe80::%31%25en0]:8080/",
         "username": "xvQVjibAbTJoJ5Nd"
         },
      "label": "p-mysql",
       "name": "test_db",
       "plan": "1gb",
        "provider": null,
        "syslog_drain_url": null,
        "tags": [
         "mysql"
       ],
        "volume_mounts": []
        }
       ]
      }`

var BadPingENV = `{
       "p-mysql": [
       {
         "credentials": {
          "hostname": "mysql-broker.local.pcfdev.io",
          "jdbcUrl": "jdbc:mysql://mysql-broker.local.pcfdev.io:3306/cf_08fd31cb_3fde_4f4d_a86d_db1c04e4ccd1?user=xvQVjibAbTJoJ5Nd\u0026password=5uZ566tecLueLX2H",
          "name": "cf_08fd31cb_3fde_4f4d_a86d_db1c04e4ccd1",
          "password": "5uZ566tecLueLX2H",
          "port": 3306,
         "uri": "mysql://youcantpingme",
         "username": "xvQVjibAbTJoJ5Nd"
         },
      "label": "p-mysql",
       "name": "test_db",
       "plan": "1gb",
        "provider": null,
        "syslog_drain_url": null,
        "tags": [
         "mysql"
       ],
        "volume_mounts": []
        }
       ]
      }`

func Test_fetchDBCredentialsFromEnvironment(t *testing.T) {

	// set oldGetEnv to real getENV
	oldEnv := getEnv

	// as we are exiting, getEnv revert back
	defer func() { getEnv = oldEnv }()

	getEnv = func(string) string {
		return ENV
	}

	expected := &credentials{
		Password: "5uZ566tecLueLX2H",
		URL:      "mysql://xvQVjibAbTJoJ5Nd:5uZ566tecLueLX2H@mysql-broker.local.pcfdev.io:3306/cf_08fd31cb_3fde_4f4d_a86d_db1c04e4ccd1?reconnect=true",
		Username: "xvQVjibAbTJoJ5Nd",
	}
	got, err := fetchDBCredentialsFromEnvironment()
	if err != nil {
		t.Errorf("fetchDBCredentialsFromEnvironment() error = %v, expect no error", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("fetchDBCredentialsFromEnvironment() mismatch: got = %v, expected = %v", got, expected)

	}

}

func Test_fetchDBCredentialsFromEnvironmentEnvEmpty(t *testing.T) {

	// set oldGetEnv to real getENV
	oldEnv := getEnv

	// as we are exiting, getEnv revert back
	defer func() { getEnv = oldEnv }()

	getEnv = func(string) string {
		return ""
	}

	expected := errors.New("Unable to retrieve VCAP_SERVICES from environment")
	_, err := fetchDBCredentialsFromEnvironment()
	if err == nil {
		t.Errorf("fetchDBCredentialsFromEnvironment() Expected error, but did not recieve one")
	}
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("fetchDBCredentialsFromEnvironment() mismatch: recieved error = %v, expected error = %v", err, expected)

	}

}

func Test_fetchDBCredentialsFromEnvironmentEnvCorrupt(t *testing.T) {

	// set oldGetEnv to real getENV
	oldEnv := getEnv

	// as we are exiting, getEnv revert back
	defer func() { getEnv = oldEnv }()

	getEnv = func(string) string {
		return "VCAP_SERVICES"
	}

	expected := errors.New("Unable to parse VCAP_SERVICES")
	_, err := fetchDBCredentialsFromEnvironment()
	if err == nil {
		t.Errorf("fetchDBCredentialsFromEnvironment() Expected error, but did not recieve one")
	}
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("fetchDBCredentialsFromEnvironment() mismatch: recieved error = %v, expected error = %v", err, expected)

	}

}

func Test_fetchCrbVersion(t *testing.T) {

	// set oldGetEnv to real getENV
	oldEnv := getEnv

	// as we are exiting, getEnv revert back
	defer func() { getEnv = oldEnv }()

	getEnv = func(string) string {
		return UserProvidedEnv
	}

	expected := "0.23"
	got, err := FetchCrbVersion()
	if err != nil {
		t.Errorf("FetchCrbVersion() error = %v, expect no error", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("FetchCrbVersion() mismatch: got = %v, expected = %v", got, expected)

	}

}

func Test_fetchCrbVersionEmptyVersion(t *testing.T) {

	// set oldGetEnv to real getENV
	oldEnv := getEnv

	// as we are exiting, getEnv revert back
	defer func() { getEnv = oldEnv }()

	getEnv = func(string) string {
		return UserProvidedEnvEmpty
	}

	expectedErr := errors.New("Unable to retrieve CRB_VERSION from environment")
	_, err := FetchCrbVersion()
	if err == nil {
		t.Errorf("FetchCrbVersion() Expected error, but did not recieve one")
	}
	if !reflect.DeepEqual(err, expectedErr) {
		t.Errorf("FetchCrbVersion() mismatch: got = %v, expected = %v", err, expectedErr)
	}

}

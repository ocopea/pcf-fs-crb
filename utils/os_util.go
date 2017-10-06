// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package utils

import (
	"encoding/json"
	"errors"
	"os"
)

var getEnv = os.Getenv

type credentials struct {
	Password string
	URL      string
	Username string
}

func fetchDBCredentialsFromEnvironment() (*credentials, error) {
	// check if we have the DB credentials from ENV already
	// if not, go fetch the environment variables with FetchDBCredentialsFromEnvironment
	env := getEnv("VCAP_SERVICES")
	var creds credentials
	var err error
	if env != "" {
		var data map[string]interface{}
		json.Unmarshal([]byte(env), &data)
		// TODO: determine how to catch a runtime error when parsing the JSON if expected structures are not there
		// i.e. accessing array out of bounds.

		usrServices, exists := data["p-mysql"].([]interface{})
		if exists {
			creds = credentials{
				Username: usrServices[0].(map[string]interface{})["credentials"].(map[string]interface{})["username"].(string),
				Password: usrServices[0].(map[string]interface{})["credentials"].(map[string]interface{})["password"].(string),
				URL:      usrServices[0].(map[string]interface{})["credentials"].(map[string]interface{})["uri"].(string),
			}
		} else {
			err = errors.New("Unable to parse VCAP_SERVICES")
		}

	} else {
		err = errors.New("Unable to retrieve VCAP_SERVICES from environment")
	}

	return &creds, err
}

// FetchCrbVersion : This returns the CRB version information from the environment variable
func FetchCrbVersion() (string, error) {
	var err error
	crbVersion := getEnv("CRB_VERSION")
	if crbVersion == "" {
		err = errors.New("Unable to retrieve CRB_VERSION from environment")
	}
	return crbVersion, err
}

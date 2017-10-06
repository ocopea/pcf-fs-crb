// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package impl

import (
	"crb/models"
	"reflect"
	"testing"
)

var UserProvidedEnv = "0.23"

// All structure info matches
func TestCrbInfoResponse(t *testing.T) {

	// set oldGetEnv to real getENV
	oldVersion := getVersion

	// as we are exiting, getEnv revert back
	defer func() { getVersion = oldVersion }()

	getVersion = func() (string, error) {
		return UserProvidedEnv, nil
	}

	expected := &models.Info{
		Name:     "crb",
		Version:  UserProvidedEnv,
		RepoType: "file",
	}

	if got := CrbInfoResponse(); !reflect.DeepEqual(got, expected) {
		t.Errorf("%q. CrbInfoResponse() = %v, want %v", expected.Name, got, expected)
	}
}

// Version is N/A
func TestCrbInfoResponseNA(t *testing.T) {

	expected := &models.Info{
		Name:     "crb",
		Version:  "N/A",
		RepoType: "file",
	}

	if got := CrbInfoResponse(); !reflect.DeepEqual(got, expected) {
		t.Errorf("%q. CrbInfoResponse() = %v, want %v", expected.Name, got, expected)
	}
}

// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package utils

import (
	"crb/models"
)

// DBError stores user friendly message for a DB errors
type DBError struct {
	Message string
}

// This makes it a GO error
func (d *DBError) Error() string {
	return d.Message
}

// DBInterface describes DB operations performed by CRB
type DBInterface interface {
	// Return error
	// - if the DB can't be accessed
	// - if the table can't be created
	CreateRepositoryTable() error

	// Return error
	// - if the DB can't be accessed
	// - if the table can't be created
	CreateCopyTable() error

	// Return error
	// - if the DB can't be accessed
	DoesRepositoryTableExist() (bool, error)

	// Return error
	// - if the DB can't be accessed
	DoesCopyTableExist() (bool, error)

	// Return error
	// - if the copy can't be found
	// - if the table creation fails
	// Note, if the table doesn't exist, it would be created
	DoesCopyExist(copyID string) (bool, error)

	// Return error
	// - if the DB can't be accessed
	// - or any other error like out of space
	// Note, if the table doesn't exist, it would be created and the repo added
	// Currently there is only one repo supported. Adding a new one simply replaces the one in the DB.
	AddRepository(repoInfo *models.RepositoryInfo) error

	// Return error
	// - if the DB can't be accessed
	// Note, if the table doesn't exist, it would be created and empty is returned
	GetRepository() (*models.RepositoryInfo, error)

	// Return error
	// - if the DB can't be accessed
	// - if the id can't be found
	// Note, if the table doesn't exist, an error is returned
	GetCopy(id string) (*models.CopyMetaData, error)

	// Return error
	// - if the DB can't be accessed
	// Note, if the table doesn't exist, an error is returned
	GetCopies() ([]string, error)

	// Return error
	// - if the DB can't be accessed
	// - if the CopyInfo object is invalid
	// - if the CopyInfo already exists
	// - or any other error like out of space
	// Note, if the table doesn't exist, it would be created and and the copy saved.
	AddCopy(c *models.CopyMetaData) error

	// Return error
	// - if the DB can't be accessed
	// - if the id can't be found
	// Note, if the table doesn't exist, it would be created and an empty array would be returned.
	DeleteCopy(id string) error
}

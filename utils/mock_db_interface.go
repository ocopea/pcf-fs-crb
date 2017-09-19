// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package utils

import (
	"crb/models"
)

// Mockdb is a collections of mock functions for DBInterface
type Mockdb struct {
	mockCreateRepositoryTable    func() error
	mockCreateCopyTable          func() error
	mockDoesRepositoryTableExist func() (bool, error)
	mockDoesCopyTableExist       func() (bool, error)
	mockAddRepository            func() error
	mockGetRepository            func() (*models.RepositoryInfo, error)
	mockGetCopy                  func() (*models.CopyMetaData, error)
	mockGetCopies                func() ([]string, error)
	mockAddCopy                  func() error
	mockDeleteCopy               func() error
	mockDoesCopyExist            func() (bool, error)
}

type WantTableExists struct {
	Exists bool
	Err    error
}

type WantRepository struct {
	RepoInfo *models.RepositoryInfo
	Err      error
}

type WantCopy struct {
	CopyMeta *models.CopyMetaData
	Err      error
}

type WantCopies struct {
	Ids []string
	Err error
}

type WantCopyExists struct {
	Exists bool
	Err    error
}

type Wants struct {
	RepoTableExists WantTableExists
	RepoTableCreate error
	AddRepo         error
	GetRepository   WantRepository
	CopyTableCreate error
	GetCopy         WantCopy
	GetCopies       WantCopies
	AddCopy         error
	DoesCopyExist   WantCopyExists
	DeleteCopy      error
}

func (mydb *Mockdb) CreateRepositoryTable() error {
	if mydb.mockCreateRepositoryTable != nil {
		return mydb.mockCreateRepositoryTable()
	}
	return nil
}

func (mydb *Mockdb) CreateCopyTable() error {
	if mydb.mockCreateCopyTable != nil {
		return mydb.mockCreateCopyTable()
	}
	return nil
}

func (mydb *Mockdb) DoesRepositoryTableExist() (bool, error) {
	if mydb.mockDoesRepositoryTableExist != nil {
		return mydb.mockDoesRepositoryTableExist()
	}
	return false, nil
}

func (mydb *Mockdb) DoesCopyTableExist() (bool, error) {
	if mydb.mockDoesCopyTableExist != nil {
		return mydb.mockDoesCopyTableExist()
	}
	return false, nil
}

func (mydb *Mockdb) AddRepository(repoInfo *models.RepositoryInfo) error {
	if mydb.mockAddRepository != nil {
		return mydb.mockAddRepository()
	}
	return nil
}

func (mydb *Mockdb) GetRepository() (*models.RepositoryInfo, error) {
	if mydb.mockGetRepository != nil {
		return mydb.mockGetRepository()
	}
	return nil, nil
}

func (mydb *Mockdb) GetCopy(id string) (*models.CopyMetaData, error) {
	if mydb.mockGetCopy != nil {
		return mydb.mockGetCopy()
	}
	return nil, nil
}

func (mydb *Mockdb) GetCopies() ([]string, error) {
	if mydb.mockGetCopies != nil {
		return mydb.mockGetCopies()
	}
	return nil, nil
}

func (mydb *Mockdb) AddCopy(c *models.CopyMetaData) error {
	if mydb.mockAddCopy != nil {
		return mydb.mockAddCopy()
	}
	return nil
}

func (mydb *Mockdb) DeleteCopy(id string) error {
	if mydb.mockDeleteCopy != nil {
		return mydb.mockDeleteCopy()
	}
	return nil
}

func (mydb *Mockdb) DoesCopyExist(id string) (bool, error) {
	if mydb.mockDoesCopyExist != nil {
		return mydb.mockDoesCopyExist()
	}
	return false, nil
}

func SetMockDbExpectations(wants *Wants) *Mockdb {
	mockDB := &Mockdb{
		mockDoesRepositoryTableExist: func() (bool, error) {
			return wants.RepoTableExists.Exists, wants.RepoTableExists.Err
		},
		mockCreateRepositoryTable: func() error {
			return wants.RepoTableCreate
		},
		mockAddRepository: func() error {
			return wants.AddRepo
		},
		mockCreateCopyTable: func() error {
			return wants.CopyTableCreate
		},
		mockDoesCopyTableExist: func() (bool, error) {
			return false, nil
		},
		mockGetRepository: func() (*models.RepositoryInfo, error) {
			return wants.GetRepository.RepoInfo, wants.GetRepository.Err
		},
		mockGetCopy: func() (*models.CopyMetaData, error) {
			return wants.GetCopy.CopyMeta, wants.GetCopy.Err
		},
		mockGetCopies: func() ([]string, error) {
			return wants.GetCopies.Ids, wants.GetCopies.Err
		},
		mockAddCopy: func() error {
			return wants.AddCopy
		},
		mockDeleteCopy: func() error {
			return wants.DeleteCopy
		},
		mockDoesCopyExist: func() (bool, error) {
			return wants.DoesCopyExist.Exists, wants.DoesCopyExist.Err
		},
	}
	return mockDB
}

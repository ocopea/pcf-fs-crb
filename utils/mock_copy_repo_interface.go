// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package utils

import (
	"errors"
	"io"
)

// MockCopyRepo contains mock interfaces of CopyRepoInterface.
type MockCopyRepo struct {
	mockStoreCopy    func() (int64, error)
	mockRetrieveCopy func() (io.ReadCloser, error)
	mockDeleteCopy   func() error
}

// WantStoreCopy defines expected results of StoreCopy
type WantStoreCopy struct {
	CopySize int64
	Err      error
}

// WantRetrieveCopy defines expected results of RetrieveCopy.
type WantRetrieveCopy struct {
	Copydata io.ReadCloser
	Err      error
}

// WantsCopyRepoInterface defines expectations for CopyRepoInterface methods.
type WantsCopyRepoInterface struct {
	StoreCopy    WantStoreCopy
	RetrieveCopy WantRetrieveCopy
	DeleteCopy   error
}

// StoreCopy stores copy data on the copy repo at given copyID
// Returns CopyMetaData that should be persisted in the copy data DB
func (sftpClient *MockCopyRepo) StoreCopy(copyID string, copyData io.ReadCloser) (int64, error) {
	if sftpClient.mockStoreCopy != nil {
		return sftpClient.mockStoreCopy()
	}
	return 0, errors.New("not implemented")
}

// RetrieveCopy gets copy data for the given copyID from the copy repo
func (sftpClient *MockCopyRepo) RetrieveCopy(copyID string) (io.ReadCloser, error) {
	if sftpClient.mockRetrieveCopy != nil {
		return sftpClient.mockRetrieveCopy()
	}
	return nil, errors.New("not implemented")
}

// DeleteCopy deletes the copy on the copy repo
func (sftpClient *MockCopyRepo) DeleteCopy(copyID string) error {
	if sftpClient.mockDeleteCopy != nil {
		return sftpClient.mockDeleteCopy()
	}
	return errors.New("not implemented")
}

// SetMockCopyRepoExpectations sets mock expectations for CopyRepoInterface
func SetMockCopyRepoExpectations(wants *WantsCopyRepoInterface) *MockCopyRepo {
	mockCR := &MockCopyRepo{
		mockStoreCopy: func() (int64, error) {
			return wants.StoreCopy.CopySize, wants.StoreCopy.Err
		},
		mockRetrieveCopy: func() (io.ReadCloser, error) {
			return wants.RetrieveCopy.Copydata, wants.RetrieveCopy.Err
		},
		mockDeleteCopy: func() error {
			return wants.DeleteCopy
		},
	}
	return mockCR
}

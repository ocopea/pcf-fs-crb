// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package utils

import (
	"io"
)

// CopyRepoInterface defines operations specific to a copy technology required by CRB
type CopyRepoInterface interface {
	// StoreCopy stores copy data on the copy repo at given copyID
	// Returns number of bytes that have been copied i.e. size of the file
	StoreCopy(copyID string, copyData io.ReadCloser) (int64, error)

	// RetrieveCopy gets copy data for the given copyID from the copy repo
	RetrieveCopy(copyID string) (io.ReadCloser, error)

	// DeleteCopy deletes the copy on the copy repo
	DeleteCopy(copyID string) error
}

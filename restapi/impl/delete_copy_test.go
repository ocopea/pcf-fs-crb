// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package impl

import (
	"crb/models"
	"crb/restapi/operations/crb_web"
	"crb/utils"
	"errors"
	"net/http"
	"reflect"
	"testing"
)

func TestDeleteCopy(t *testing.T) {
	var mockDB *utils.Mockdb
	var mockCR *utils.MockCopyRepo

	//validReadCloser := ioutil.NopCloser(new(bytes.Buffer))

	tests := []struct {
		name    string
		copyID  string
		dbWants utils.Wants
		crWants utils.WantsCopyRepoInterface
		wantErr bool
	}{
		{
			name:    "Successful",
			copyID:  "test",
			dbWants: utils.Wants{DeleteCopy: nil},
			crWants: utils.WantsCopyRepoInterface{DeleteCopy: nil},
			wantErr: false,
		},
		{
			/* Test Case: Copy ID not found*/
			name:    "InvalidCopyId",
			copyID:  "",
			dbWants: utils.Wants{DeleteCopy: nil},
			crWants: utils.WantsCopyRepoInterface{DeleteCopy: nil},
			wantErr: true,
		},
		{
			/*
				Test Case: Couldn’t delete copy,
				           Repository credentials doesn't have write access.
			*/
			name:    "RepoFailed",
			copyID:  "testID",
			dbWants: utils.Wants{DeleteCopy: nil},
			crWants: utils.WantsCopyRepoInterface{DeleteCopy: errors.New("Unable to delete copy from repo")},
			wantErr: true,
		},
		{
			/*
				Test Case: Successful delete even if the copy is not found
				on the external repo but found in the DB.
			*/
			name:    "FileDoesntExist",
			copyID:  "testID",
			dbWants: utils.Wants{DeleteCopy: nil},
			crWants: utils.WantsCopyRepoInterface{DeleteCopy: errors.New("file does not exist")},
			wantErr: false,
		},
		{
			/*
				Test Case: Couldn’t delete DB entry after deleting copy data in repo
			*/
			name:    "DBFailed",
			copyID:  "testID",
			dbWants: utils.Wants{DeleteCopy: errors.New("Unable to delete copy metadata from DB")},
			crWants: utils.WantsCopyRepoInterface{DeleteCopy: nil},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB = utils.SetMockDbExpectations(&tt.dbWants)
			mockCR = utils.SetMockCopyRepoExpectations(&tt.crWants)
			err := DeleteCopy(tt.copyID, mockCR, mockDB)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCopyHandler() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteCopyHandler(t *testing.T) {
	testURL := "http://crb-server.local.pcfdev.io/crb/copies/testID"
	testCopyID := "testID"
	var mockDB *utils.Mockdb

	tests := []struct {
		name           string
		sftpclient     *utils.SftpClient
		dbexpectations utils.Wants
		wantCode       int
		wantErrString  string
	}{
		{
			/*
				Test Case: Invalid repo credential,
				             No repository credentials available,
							 No repo access
			*/
			name:           "InvalidSftpClient",
			sftpclient:     nil,
			dbexpectations: utils.Wants{DoesCopyExist: utils.WantCopyExists{Exists: true, Err: nil}},
			wantCode:       500,
			wantErrString:  "Connection to the repository couldn't be established",
		},
		{
			/*
			 Test Case: Copy ID not found
			*/
			name:       "NoCopyMetadata",
			sftpclient: nil,
			dbexpectations: utils.Wants{DoesCopyExist: utils.WantCopyExists{
				Exists: false, Err: nil},
			},
			wantCode:      404,
			wantErrString: "Copy ID not found",
		},
		{
			/* Test Case: No DB access */
			name:           "DBError",
			sftpclient:     nil,
			dbexpectations: utils.Wants{DoesCopyExist: utils.WantCopyExists{Exists: true, Err: errors.New("DB connection error")}},
			wantCode:       500,
			wantErrString:  "DB connection error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", testURL, nil)
			params := crb_web.DeleteCopyParams{
				HTTPRequest: req,
				CopyID:      testCopyID}
			testCode := int32(tt.wantCode)
			testErr := &models.Error{Code: &testCode, Message: &tt.wantErrString}
			want := crb_web.NewDeleteCopyDefault(tt.wantCode).WithPayload(testErr)
			mockDB = utils.SetMockDbExpectations(&tt.dbexpectations)
			if got := DeleteCopyHandler(params, mockDB, tt.sftpclient); !reflect.DeepEqual(got, want) {
				t.Errorf("DeleteCopyHandler() = %v, want %v", got, want)
			}
		})
	}
}

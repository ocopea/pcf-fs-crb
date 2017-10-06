// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package impl

import (
	"bytes"
	"crb/models"
	"crb/restapi/operations/crb_web"
	"crb/utils"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
)

func getTestCopyMeta() *models.CopyMetaData {
	return &models.CopyMetaData{
		CopyDataURL:   "testUrl",
		CopyID:        "testCopyID",
		CopySize:      128,
		CopyTimeStamp: strfmt.DateTime(time.Now()),
	}
}

func TestPostCopyHandler(t *testing.T) {
	var mockDB *utils.Mockdb
	var mockCR *utils.MockCopyRepo

	type args struct {
		copyID   string
		copyData io.ReadCloser
	}

	validReadCloser := ioutil.NopCloser(new(bytes.Buffer))

	tests := []struct {
		name    string
		args    args
		dbWants utils.Wants
		crWants utils.WantsCopyRepoInterface
		wantErr bool
	}{
		{
			name:    "TestStoreCopySuccessful",
			args:    args{copyID: "testCopyID", copyData: validReadCloser},
			dbWants: utils.Wants{AddCopy: nil},
			crWants: utils.WantsCopyRepoInterface{StoreCopy: utils.WantStoreCopy{CopySize: 128, Err: nil},
				DeleteCopy: errors.New("should not be called")},
			wantErr: false,
		},
		{
			name:    "TestStoreCopyFailed",
			args:    args{copyID: "testCopyID", copyData: validReadCloser},
			dbWants: utils.Wants{AddCopy: errors.New("should not be called")},
			crWants: utils.WantsCopyRepoInterface{StoreCopy: utils.WantStoreCopy{CopySize: 0, Err: errors.New("store copy failed")},
				DeleteCopy: errors.New("should not be called")},
			wantErr: true,
		},
		{
			name:    "TestCopyDataSizeZero",
			args:    args{copyID: "testCopyID", copyData: validReadCloser},
			dbWants: utils.Wants{AddCopy: errors.New("should not be called")},
			crWants: utils.WantsCopyRepoInterface{StoreCopy: utils.WantStoreCopy{CopySize: 0, Err: nil},
				DeleteCopy: errors.New("should not be called")},
			wantErr: true,
		},
		{
			name:    "TestPersistCopyMetaFails",
			args:    args{copyID: "testCopyID", copyData: validReadCloser},
			dbWants: utils.Wants{AddCopy: errors.New("couldn't persist copy meta")},
			crWants: utils.WantsCopyRepoInterface{StoreCopy: utils.WantStoreCopy{CopySize: 128, Err: nil}, DeleteCopy: nil},
			wantErr: true,
		},
		{
			name:    "TestPersistCopyMetaAndDeleteCopyFails",
			args:    args{copyID: "testCopyID", copyData: validReadCloser},
			dbWants: utils.Wants{AddCopy: errors.New("couldn't persist copy meta")},
			crWants: utils.WantsCopyRepoInterface{StoreCopy: utils.WantStoreCopy{CopySize: 128, Err: nil},
				DeleteCopy: errors.New("failed to delete copy data")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB = utils.SetMockDbExpectations(&tt.dbWants)
			mockCR = utils.SetMockCopyRepoExpectations(&tt.crWants)
			if err := postCopyHandler(tt.args.copyID, tt.args.copyData, mockCR, mockDB); (err != nil) != tt.wantErr {
				t.Errorf("postCopyHandler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func getFileName(nameLength int) string {
	fName := make([]rune, nameLength)

	for i := range fName {
		fName[i] = []rune("a")[0]
	}
	return string(fName)
}

func TestValidateCopyID(t *testing.T) {
	var mockDB *utils.Mockdb
	var mockCR *utils.MockCopyRepo

	type args struct {
		copyID   string
		copyData io.ReadCloser
	}

	validReadCloser := ioutil.NopCloser(new(bytes.Buffer))

	tests := []struct {
		name    string
		args    args
		dbWants utils.Wants
		crWants utils.WantsCopyRepoInterface
		wantErr bool
	}{
		{
			name:    "TestValidName",
			args:    args{copyID: "myCopyData1", copyData: validReadCloser},
			dbWants: utils.Wants{AddCopy: nil},
			crWants: utils.WantsCopyRepoInterface{StoreCopy: utils.WantStoreCopy{CopySize: 128, Err: nil},
				DeleteCopy: errors.New("should not be called")},
			wantErr: false,
		},
		{
			name:    "TestForForwardSlash",
			args:    args{copyID: "/root", copyData: validReadCloser},
			dbWants: utils.Wants{AddCopy: nil},
			crWants: utils.WantsCopyRepoInterface{StoreCopy: utils.WantStoreCopy{CopySize: 128, Err: nil},
				DeleteCopy: errors.New("should not be called")},
			wantErr: true,
		},
		{
			name:    "TestForStartingWithHyphen",
			args:    args{copyID: "-myCopyData1", copyData: validReadCloser},
			dbWants: utils.Wants{AddCopy: nil},
			crWants: utils.WantsCopyRepoInterface{StoreCopy: utils.WantStoreCopy{CopySize: 128, Err: nil},
				DeleteCopy: errors.New("should not be called")},
			wantErr: true,
		},
		{
			name:    "TestForInterleavedHyphen",
			args:    args{copyID: "myCopyData1-1", copyData: validReadCloser},
			dbWants: utils.Wants{AddCopy: nil},
			crWants: utils.WantsCopyRepoInterface{StoreCopy: utils.WantStoreCopy{CopySize: 128, Err: nil},
				DeleteCopy: errors.New("should not be called")},
			wantErr: false,
		},
		{
			name:    "TestForEmptyName",
			args:    args{copyID: "", copyData: validReadCloser},
			dbWants: utils.Wants{AddCopy: nil},
			crWants: utils.WantsCopyRepoInterface{StoreCopy: utils.WantStoreCopy{CopySize: 128, Err: nil},
				DeleteCopy: errors.New("should not be called")},
			wantErr: true,
		},
		{
			name:    "TestMaxName",
			args:    args{copyID: getFileName(255), copyData: validReadCloser},
			dbWants: utils.Wants{AddCopy: nil},
			crWants: utils.WantsCopyRepoInterface{StoreCopy: utils.WantStoreCopy{CopySize: 128, Err: nil},
				DeleteCopy: errors.New("should not be called")},
			wantErr: false,
		},
		{
			name:    "TestOversizedName",
			args:    args{copyID: getFileName(300), copyData: validReadCloser},
			dbWants: utils.Wants{AddCopy: nil},
			crWants: utils.WantsCopyRepoInterface{StoreCopy: utils.WantStoreCopy{CopySize: 128, Err: nil},
				DeleteCopy: errors.New("should not be called")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB = utils.SetMockDbExpectations(&tt.dbWants)
			mockCR = utils.SetMockCopyRepoExpectations(&tt.crWants)
			if err := postCopyHandler(tt.args.copyID, tt.args.copyData, mockCR, mockDB); (err != nil) != tt.wantErr {
				t.Errorf("postCopyHandler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getCopyInstance(t *testing.T) {
	tests := []struct {
		name            string
		postDataRequest string
		want            *models.CopyInstance
		expectFail      bool
	}{
		{
			name:            "TestPostRequestSuccess",
			postDataRequest: "http://127.0.0.1:8080/crb/copies/testId/data",
			want:            &models.CopyInstance{CopyURL: "http://127.0.0.1:8080/crb/copies/testId"},
			expectFail:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", tt.postDataRequest, nil)
			if got := getCopyInstance(req); reflect.DeepEqual(got, tt.want) {
				if tt.expectFail {
					t.Errorf("getCopyInstance() got = %v, want %v", got, tt.want)
				}
			} else {
				if !tt.expectFail {
					t.Errorf("getCopyInstance() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestCreateCopyHandler(t *testing.T) {
	testURL := "http://crb-server.local.pcfdev.io/crb/copies/test/data"
	testCopyID := "test"
	var mockDB *utils.Mockdb

	tests := []struct {
		name           string
		sftpclient     *utils.SftpClient
		dbexpectations utils.Wants
		wantCode       int
		wantErrString  string
	}{
		{
			name:           "TestInvalidSftpClient",
			sftpclient:     nil,
			dbexpectations: utils.Wants{DoesCopyExist: utils.WantCopyExists{Exists: false, Err: nil}},
			wantCode:       500,
			wantErrString:  "Sftp connection was not established",
		},
		{
			name:       "TestDBNoCopyTable",
			sftpclient: nil,
			dbexpectations: utils.Wants{DoesCopyExist: utils.WantCopyExists{
				Exists: false, Err: errors.New("Could not create CopyTable")},
			},
			wantCode:      500,
			wantErrString: "Could not create CopyTable",
		},
		{
			name:           "TestDuplicatePost",
			sftpclient:     nil,
			dbexpectations: utils.Wants{DoesCopyExist: utils.WantCopyExists{Exists: true, Err: nil}},
			wantCode:       403,
			wantErrString:  "copyID already exists",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", testURL, nil)
			testInfo := "Test Copy Info"
			params := crb_web.CreateCopyParams{
				HTTPRequest: req,
				CopyData:    ioutil.NopCloser(new(bytes.Buffer)),
				CopyID:      testCopyID,
				CopyInfo:    &testInfo}
			testCode := int32(tt.wantCode)
			testErr := &models.Error{Code: &testCode, Message: &tt.wantErrString}
			want := crb_web.NewCreateCopyDefault(tt.wantCode).WithPayload(testErr)
			mockDB = utils.SetMockDbExpectations(&tt.dbexpectations)
			if got := CreateCopyHandler(params, mockDB, tt.sftpclient); !reflect.DeepEqual(got, want) {
				t.Errorf("CreateCopyHandler() = %v, want %v", got, want)
			}
		})
	}
}

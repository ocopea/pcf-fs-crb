// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package impl

import (
	"bytes"
	"crb/models"
	"crb/restapi/operations/crb_web"
	"crb/utils"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

func TestRetrieveCopyData(t *testing.T) {
	var mockDB *utils.Mockdb
	var mockCopyClient *utils.MockCopyRepo

	tests := []struct {
		name                 string
		want                 models.OutputStream
		wantErr              bool
		dbexpectations       utils.Wants
		copyRepoExpectations utils.WantsCopyRepoInterface
	}{
		{name: "testGetCopyMetaDataUnsucccessful",
			wantErr: true,
			dbexpectations: utils.Wants{GetCopy: utils.WantCopy{
				CopyMeta: nil,
				Err:      errors.New("Error while getting copymetadata"),
			}},
			copyRepoExpectations: utils.WantsCopyRepoInterface{RetrieveCopy: utils.WantRetrieveCopy{
				Copydata: nil,
				Err:      nil,
			}},
		},
		{name: "testGetCopyUnsuccessful",
			wantErr: true,
			dbexpectations: utils.Wants{GetCopy: utils.WantCopy{
				CopyMeta: &models.CopyMetaData{
					CopyDataURL:   "/",
					CopyID:        "sample",
					CopySize:      int64(10),
					CopyTimeStamp: strfmt.NewDateTime(),
				},
				Err: nil,
			}},
			copyRepoExpectations: utils.WantsCopyRepoInterface{RetrieveCopy: utils.WantRetrieveCopy{
				Copydata: nil,
				Err:      errors.New("Error while retrieving copy"),
			}},
		},
		{name: "testGetCopySuccesful",
			wantErr: false,
			dbexpectations: utils.Wants{GetCopy: utils.WantCopy{
				CopyMeta: &models.CopyMetaData{
					CopyDataURL:   "/",
					CopyID:        "sample",
					CopySize:      int64(10),
					CopyTimeStamp: strfmt.NewDateTime(),
				},
				Err: nil,
			}},
			copyRepoExpectations: utils.WantsCopyRepoInterface{RetrieveCopy: utils.WantRetrieveCopy{
				Copydata: ioutil.NopCloser(new(bytes.Buffer)),
				Err:      nil,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			copyID := "mockid"
			mockDB = utils.SetMockDbExpectations(&tt.dbexpectations)
			mockCopyClient = utils.SetMockCopyRepoExpectations(&tt.copyRepoExpectations)
			got, err := retrieveCopyData(copyID, mockDB, mockCopyClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("retrieveCopyData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if err != nil {
					t.Errorf("retrieveCopyData() error = %v, wantErr %v", err, tt.wantErr)
				}
				if got == nil {
					t.Errorf("retrieveCopyData() got is nil")
				}

			}
		})
	}
}

func TestGetSftpConnection(t *testing.T) {
	var mockDB *utils.Mockdb
	var mockSftpClient *utils.SftpClient
	addr := "blah"
	user := "blah"
	tests := []struct {
		name           string
		args           *utils.SftpClient
		dbexpectations utils.Wants
		want           *utils.SftpClient
	}{
		{
			name:           "TestFailOnDbFailure",
			args:           mockSftpClient,
			dbexpectations: utils.Wants{GetRepository: utils.WantRepository{RepoInfo: nil, Err: errors.New("Error while getting repository")}},
			want:           nil,
		},

		{
			name: "TestFailOnSftpConnectionFailure",
			args: mockSftpClient,
			dbexpectations: utils.Wants{GetRepository: utils.WantRepository{
				RepoInfo: &models.RepositoryInfo{Addr: &addr, User: &user, Password: "blah"},
				Err:      nil},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB = utils.SetMockDbExpectations(&tt.dbexpectations)
			if got := GetSftpConnection(mockDB, mockSftpClient); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSftpConnection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRetrieveCopyHandler(t *testing.T) {
	testURL := "http://crb-server.local.pcfdev.io/crb/copies/test/data"
	testCopyID := "test"

	type args struct {
		mysqldb    *utils.Database
		sftpclient *utils.SftpClient
	}
	tests := []struct {
		name string
		args args
		want middleware.Responder
	}{
		{
			name: "TestInvalidSftpClient",
			args: args{mysqldb: nil, sftpclient: nil},
			want: crb_web.NewRetrieveCopyDefault(500),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", testURL, nil)
			params := crb_web.RetrieveCopyParams{HTTPRequest: req, CopyID: testCopyID}
			if got := RetrieveCopyHandler(params, tt.args.mysqldb, tt.args.sftpclient); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RetrieveCopyHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

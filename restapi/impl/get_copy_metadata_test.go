// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package impl

import (
	"crb/models"
	"crb/utils"
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"

	"crb/restapi/operations/crb_web"

	"github.com/go-openapi/strfmt"
)

// This is to make sure both production and expected time are fixed and same
func getTestTime() time.Time {
	t, _ := time.Parse(time.RFC3339, "2017-04-20")
	return t
}

func getTestMeta(wantCopyID string, wantCopyDataURL string) *models.CopyMetaData {
	return &models.CopyMetaData{
		CopyDataURL:   wantCopyDataURL,
		CopyID:        wantCopyID,
		CopySize:      128,
		CopyTimeStamp: strfmt.DateTime(getTestTime()),
	}
}

func TestGetCopy(t *testing.T) {
	var mockDB *utils.Mockdb

	testURL := "http://crb-server.local.pcfdev.io/crb/copies/test3"
	testCopyID := "test3"

	tests := []struct {
		name         string
		wantErr      bool
		wantCopyMeta *models.CopyMetaData
		expectations utils.Wants
	}{
		{
			name:         "testGetCopySuccessful",
			wantErr:      false,
			wantCopyMeta: getTestMeta(testCopyID, "http://crb-server.local.pcfdev.io/crb/copies/test3/data"),
			expectations: utils.Wants{GetCopy: utils.WantCopy{CopyMeta: getTestMeta(testCopyID, "/"), Err: nil}},
		},
		{
			name:         "testGetCopyUnsuccessful",
			wantErr:      true,
			wantCopyMeta: nil,
			expectations: utils.Wants{GetCopy: utils.WantCopy{CopyMeta: nil, Err: errors.New("GetCopy threw error")}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB = utils.SetMockDbExpectations(&tt.expectations)

			req, _ := http.NewRequest("GET", testURL, nil)
			params := crb_web.GetCopyMetaDataParams{HTTPRequest: req, CopyID: testCopyID}

			got, err := getCopyMetaData(params, mockDB)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCopyMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.wantCopyMeta) {
				t.Errorf("getCopyMetadata() got = %v, want %v", got, tt.wantCopyMeta)
			}
		})
	}
}

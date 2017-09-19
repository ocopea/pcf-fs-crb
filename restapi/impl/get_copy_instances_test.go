// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package impl

import (
	"crb/models"
	"crb/restapi/operations/crb_web"
	"crb/utils"
	"errors"
	"net/http"
	"testing"
)

const getURL string = "http://127.0.0.1:8080/crb/copies"

func getExpectedCopyInstanceList(ids []string) []*models.CopyInstance {
	var idList []*models.CopyInstance

	for _, id := range ids {
		var idInstance models.CopyInstance

		idInstance.CopyURL = getURL + "/" + id
		idList = append(idList, &idInstance)
	}
	return idList
}

func TestGetCopies(t *testing.T) {
	var mockDB *utils.Mockdb

	req, _ := http.NewRequest("GET", getURL, nil)
	params := crb_web.GetCopyinstancesParams{HTTPRequest: req}

	tests := []struct {
		name         string
		wantErr      bool
		expectations utils.Wants
	}{
		{name: "testGetCopiesSuccessful",
			wantErr: false,
			expectations: utils.Wants{GetCopies: utils.WantCopies{
				Ids: []string{"test3", "test6", "test7"},
				Err: nil},
			},
		},
		{name: "testGetNoCopiesSuccess",
			wantErr: false,
			expectations: utils.Wants{GetCopies: utils.WantCopies{
				Ids: []string{},
				Err: nil},
			},
		},
		{name: "testGetCopiesUnsuccessful",
			wantErr: true,
			expectations: utils.Wants{GetCopies: utils.WantCopies{
				Ids: nil,
				Err: errors.New("Error while getting repository")},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB = utils.SetMockDbExpectations(&tt.expectations)

			info, err := copyInstancesResponse(params, mockDB)
			if tt.wantErr {
				if err == nil {
					t.Errorf("copyInstancesResponse() error = %v, wantErr %v", err, tt.wantErr)
				}

				if info != nil {
					t.Errorf("copyInstancesResponse() info = %v, wanted nil", info)
				}
			} else {
				if err != nil {
					t.Errorf("copyInstancesResponse() error = %v, wantErr %v", err, tt.wantErr)
				}

				idList := getExpectedCopyInstanceList(tt.expectations.GetCopies.Ids)
				if len(info) != len(idList) {
					t.Errorf("copyInstancesResponse() got number of instances = %d, expected %d", len(info), len(idList))
				}
				for i := range info {
					if info[i].CopyURL != idList[i].CopyURL {
						t.Errorf("copyInstancesResponse() got = %v, expected %v", info[i], idList[i])
					}
				}
			}
		})
	}
}

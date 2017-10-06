// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package impl

import (
	"crb/models"
	"crb/utils"
	"errors"
	"testing"
)

func TestGetRepository(t *testing.T) {
	var mockDB *utils.Mockdb
	addr := "1.2.3.4"
	user := "user"

	tests := []struct {
		name         string
		wantErr      bool
		expectations utils.Wants
	}{
		{name: "testGetSuccessful",
			wantErr: false,
			expectations: utils.Wants{GetRepository: utils.WantRepository{
				RepoInfo: &models.RepositoryInfo{Addr: &addr, User: &user, Password: "Pass"},
				Err:      nil},
			},
		},
		{name: "testGetUnsuccessful",
			wantErr: true,
			expectations: utils.Wants{GetRepository: utils.WantRepository{
				RepoInfo: nil,
				Err:      errors.New("Error while getting repository")},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB = utils.SetMockDbExpectations(&tt.expectations)

			info, err := crbRepositoriesResponse(mockDB)
			if err != nil != tt.wantErr {
				t.Errorf("CrbRepositoriesResponse() error = %v, wantErr %v", err, tt.wantErr)
			}

			if info != tt.expectations.GetRepository.RepoInfo {
				t.Errorf("CrbRepositoriesResponse() got = %v, expected %v", info, tt.expectations.GetRepository.RepoInfo)
			}
		})
	}
}

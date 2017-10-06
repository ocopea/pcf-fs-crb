// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package impl

import (
	"crb/models"
	"crb/utils"
)

var getVersion = utils.FetchCrbVersion

// CrbInfoResponse returns an instatiated CRBInfo with hard coded info data.
func CrbInfoResponse() *models.Info {
	version := "N/A"
	versionFromEnv, err := getVersion()

	if err == nil {
		version = versionFromEnv
	}

	payload := &models.Info{
		Name:     "crb",
		Version:  version,
		RepoType: "file",
	}
	return payload
}

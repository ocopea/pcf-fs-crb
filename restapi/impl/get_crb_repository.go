// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package impl

import (
	"crb/models"
	"crb/restapi/operations/crb_web"
	"crb/utils"
	"reflect"
	"strings"

	"github.com/go-openapi/runtime/middleware"
)

// GetRepositoryInfoHandler handles getting copy repo information from DB.
// Returns failure http response if CRB is unable to retrieve information.
func GetRepositoryInfoHandler(params crb_web.GetRepositoryInfoParams) middleware.Responder {
	emptyRepoInfo := &models.RepositoryInfo{}
	response, error := crbRepositoriesResponse(utils.GetDatabase())

	if error != nil || reflect.DeepEqual(response, emptyRepoInfo) {
		var errorMsg string
		var tempCode int

		if reflect.DeepEqual(response, emptyRepoInfo) ||
			strings.Contains(error.Error(), "Error 1146") {
			// If the table doesn't exist or is empty, return 404.
			tempCode = 404
			errorMsg = "The repository information is not available."
		} else {
			tempCode = 500
			errorMsg = error.Error()
		}
		return getRepoReponseFailure(tempCode, errorMsg)
	}
	return crb_web.NewGetRepositoryInfoOK().WithPayload(response)
}

// crbRepositoriesResponse return the repository credentials
func crbRepositoriesResponse(dbUtil utils.DBInterface) (*models.RepositoryInfo, error) {
	return dbUtil.GetRepository()
}

// getRepoReponseFailure returns failure response
func getRepoReponseFailure(errorCode int, errorMsg string) *crb_web.GetRepositoryInfoDefault {
	failedResponse := crb_web.NewGetRepositoryInfoDefault(errorCode)
	modelErrorInt := int32(errorCode)
	modelsError := models.Error{Code: &modelErrorInt, Message: &errorMsg}
	failedResponse.SetPayload(&modelsError)
	return failedResponse
}

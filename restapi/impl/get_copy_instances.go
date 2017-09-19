// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package impl

import (
	"crb/models"
	"crb/restapi/operations/crb_web"
	"crb/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-openapi/runtime/middleware"
)

// GetCopyInstances fetches from the database the list of copy instances
// return error if problems found
func GetCopyInstances(params crb_web.GetCopyinstancesParams) middleware.Responder {
	var mysqldb *utils.Database
	payload, err := copyInstancesResponse(params, mysqldb)
	if err == nil {
		return crb_web.NewGetCopyinstancesOK().WithPayload(payload)
	}
	failedResponse := crb_web.NewGetCopyinstancesDefault(500)
	modelErrorInt := int32(500)
	modelErrorMsg := err.Error()
	modelsError := models.Error{Code: &modelErrorInt, Message: &modelErrorMsg}
	failedResponse.SetPayload(&modelsError)
	return failedResponse
}

// Fetch from the database the list of copy instances (urls)
// return error if problems found
func copyInstancesResponse(params crb_web.GetCopyinstancesParams, mysqldb utils.DBInterface) ([]*models.CopyInstance, error) {

	var idList []*models.CopyInstance

	// get list from database
	idArray, err := mysqldb.GetCopies()
	// loop through list and add URL
	if err != nil {
		return nil, err
	}
	for _, id := range idArray {
		var idInstance models.CopyInstance

		idInstance.CopyURL = ConstructURL(params.HTTPRequest, id)
		idList = append(idList, &idInstance)
	}

	return idList, nil
}

// ConstructURL returned by GET copies and GET Copy Metadata
func ConstructURL(request *http.Request, urlEnding string) string {
	urlPath := strings.TrimRight(request.URL.Path, "/")
	urlHost := request.Host
	scheme := "http://"

	var CopyURL string
	CopyURL = fmt.Sprintf("%s%s%s/%s", scheme, urlHost, urlPath, urlEnding)
	return CopyURL
}

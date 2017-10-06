// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package impl

import (
	"crb/models"
	"crb/restapi/operations/crb_web"
	"crb/utils"

	"github.com/go-openapi/runtime/middleware"
)

// GetCopyMetaData gets metadata for one specific ID.
func GetCopyMetaData(params crb_web.GetCopyMetaDataParams) middleware.Responder {
	var mysqldb *utils.Database
	payload, err := getCopyMetaData(params, mysqldb)
	if err == nil {
		return crb_web.NewGetCopyMetaDataOK().WithPayload(payload)
	}

	return getCopyMetaResponseFailure(500, err.Error())
}

func getCopyMetaResponseFailure(errorCode int, errorMsg string) *crb_web.GetCopyMetaDataDefault {
	failedResponse := crb_web.NewGetCopyMetaDataDefault(errorCode)
	modelErrorInt := int32(errorCode)
	modelsError := models.Error{Code: &modelErrorInt, Message: &errorMsg}
	failedResponse.SetPayload(&modelsError)
	return failedResponse
}

// Fetch from the database the metadata for the provided ID of copy
// return error if problems found
func getCopyMetaData(params crb_web.GetCopyMetaDataParams, mysqldb utils.DBInterface) (*models.CopyMetaData, error) {
	copyMetaData, err := mysqldb.GetCopy(params.CopyID)
	if err == nil {
		copyMetaData.CopyDataURL = ConstructURL(params.HTTPRequest, "data")
		return copyMetaData, nil
	}
	return nil, err
}

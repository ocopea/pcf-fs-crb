// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package impl

import (
	"crb/models"
	"crb/restapi/operations/crb_web"
	"crb/utils"
	"errors"
	"io"
	"time"

	"net/http"

	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

// CreateCopyHandler returns a success http response if copy creation succeeds else http error.
func CreateCopyHandler(params crb_web.CreateCopyParams, mysqldb utils.DBInterface, sftpClient *utils.SftpClient) middleware.Responder {
	exists, err := mysqldb.DoesCopyExist(params.CopyID)
	if err != nil {
		return createCopyReponseFailure(500, err.Error())
	}

	if exists {
		return createCopyReponseFailure(403, "copyID already exists")
	}

	// Check if  Sftp Connection is valid as it is opened before handler is called
	if sftpClient == nil {
		return createCopyReponseFailure(500, "Sftp connection was not established")
	}

	err = postCopyHandler(params.CopyID, params.CopyData, sftpClient, mysqldb)
	if err == nil {
		return crb_web.NewCreateCopyCreated().WithPayload(getCopyInstance(params.HTTPRequest))
	}
	return createCopyReponseFailure(500, err.Error())
}

// Assume you always get a valid http request
func getCopyInstance(postDataRequest *http.Request) *models.CopyInstance {
	copyURL := strings.TrimSuffix(ConstructURL(postDataRequest, ""), "/data/")
	return &models.CopyInstance{CopyURL: copyURL}
}

// postCopyHandler performs the storing of copy data and persisting the transaction in copy table.
func postCopyHandler(copyID string, copyData io.ReadCloser, mycopyrepo utils.CopyRepoInterface, mysqldb utils.DBInterface) error {
	err := utils.ValidateCopyID(copyID)
	if err != nil {
		return err
	}
	copySize, err := mycopyrepo.StoreCopy(copyID, copyData)
	if err != nil {
		return err
	}

	if copySize == 0 {
		return errors.New("copyData size cannot be 0")
	}

	copyMeta := models.CopyMetaData{
		CopyDataURL:   "", // Not stored in DB today
		CopyID:        copyID,
		CopySize:      copySize,
		CopyTimeStamp: strfmt.DateTime(time.Now()),
	}

	if err := mysqldb.AddCopy(&copyMeta); err != nil {
		// No need to check if DeleteCopy failed, as we are returning error
		// TODO: Maybe add error to application log if DeleteCopy fails
		mycopyrepo.DeleteCopy(copyID)
		return err
	}
	return nil
}

// createCopyReponseFailure returns failure response
func createCopyReponseFailure(errorCode int, errorMsg string) *crb_web.CreateCopyDefault {
	failedResponse := crb_web.NewCreateCopyDefault(errorCode)
	modelErrorInt := int32(errorCode)
	modelsError := models.Error{Code: &modelErrorInt, Message: &errorMsg}
	failedResponse.SetPayload(&modelsError)
	return failedResponse
}

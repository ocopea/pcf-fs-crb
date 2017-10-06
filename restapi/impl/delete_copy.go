// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package impl

import (
	"crb/models"
	"crb/restapi/operations/crb_web"
	"crb/utils"

	"errors"

	"strings"

	"github.com/go-openapi/runtime/middleware"
)

const errRepoDoesntExist = "Connection to the repository couldn't be established"
const errCopyID = "Copy ID not found"
const errFileMissing = "file does not exist"

// DeleteCopyHandler returns a success http response if copy deletion succeeds else http error.
func DeleteCopyHandler(params crb_web.DeleteCopyParams, mysqldb utils.DBInterface, sftpClient *utils.SftpClient) middleware.Responder {
	exists, err := mysqldb.DoesCopyExist(params.CopyID)
	if err != nil {
		return deleteCopyReponseFailure(500, err.Error())
	}

	if !exists {
		return deleteCopyReponseFailure(404, errCopyID)
	}

	// Check if  Sftp Connection is valid as it is opened before handler is called
	if sftpClient == nil {
		err = errors.New(errRepoDoesntExist)
		return deleteCopyReponseFailure(500, err.Error())
	}

	err = DeleteCopy(params.CopyID, sftpClient, mysqldb)
	if err != nil {
		return deleteCopyReponseFailure(500, err.Error())
	}
	return crb_web.NewDeleteCopyOK()
}

// DeleteCopy performs the deletion of copy data and copy metadata.
func DeleteCopy(copyID string, mycopyrepo utils.CopyRepoInterface, mysqldb utils.DBInterface) error {
	err := utils.ValidateCopyID(copyID)
	if err != nil {
		return err
	}

	err = mycopyrepo.DeleteCopy(copyID)
	if err != nil {
		// If deleting the file failed with
		// "file does not exist", then proceed to
		// deleting the DB entry.
		// If deleting file failed for some other
		// reason, then fail out here.
		if strings.Compare(err.Error(), errFileMissing) != 0 {
			return err
		}

	}

	if err = mysqldb.DeleteCopy(copyID); err != nil {
		return err
	}
	return nil
}

// deleteCopyReponseFailure returns failure response
func deleteCopyReponseFailure(errorCode int, errorMsg string) *crb_web.DeleteCopyDefault {
	failedResponse := crb_web.NewDeleteCopyDefault(errorCode)
	modelErrorInt := int32(errorCode)
	modelsError := models.Error{Code: &modelErrorInt, Message: &errorMsg}
	failedResponse.SetPayload(&modelsError)
	return failedResponse
}

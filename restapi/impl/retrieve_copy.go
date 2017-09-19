// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package impl

import (
	"crb/models"
	"crb/restapi/operations/crb_web"
	"crb/utils"
	"errors"

	"github.com/go-openapi/runtime/middleware"
)

// RetrieveCopyHandler returns a copy octet stream on success else a failure http response
func RetrieveCopyHandler(params crb_web.RetrieveCopyParams, mysqldb utils.DBInterface, sftpclient *utils.SftpClient) middleware.Responder {
	response, error := retrieveCopyData(params.CopyID, mysqldb, sftpclient)
	if error != nil {
		return sendErrorGetCopy(500, error.Error())
	}

	return crb_web.NewRetrieveCopyOK().WithPayload(response)
}

func sendErrorGetCopy(aCode int, err string) *crb_web.RetrieveCopyDefault {
	errorResponse := crb_web.NewRetrieveCopyDefault(aCode)

	// Only sending the error code instead of the models.Error because
	// we are currenlty getting this error:
	// 2017/03/17 16:25:29 http: panic serving 127.0.0.1:38054: &{0xc42040e4c8 0xc42040e4d0}
	// (*models.Error) is not supported by the ByteStreamProducer, can be resolved by
	// supporting Reader/BinaryMarshaler interface

	// code := int32(aCode)
	// newerr := models.Error{Code: &code, Message: &err}
	// errorResponse.SetPayload(&newerr)

	errorResponse.SetStatusCode(aCode)

	return errorResponse
}

// retrieveCopyData get the copy data for the given copy ID
func retrieveCopyData(copyID string, dbutil utils.DBInterface, sftpclient utils.CopyRepoInterface) (models.OutputStream, error) {
	// Check if  Sftp Connection is valid as it is opened before handler is called
	if sftpclient == nil {
		return nil, errors.New("Sftp connection was not established")
	}

	// Get the meta data for the ID
	copyInfo, error := dbutil.GetCopy(copyID)
	if error != nil {
		return nil, error
	}

	// Retrieve the copy from the extenal system.
	ioreader, error := sftpclient.RetrieveCopy(copyInfo.CopyID)
	if error != nil {
		return nil, error
	}

	return ioreader, nil
}

// GetSftpConnection opens and returns a SftpConnection for the given repo
// Returns nil on failure
// The caller for this function should close the SftpConnection
func GetSftpConnection(mysqldb utils.DBInterface, sftpClient *utils.SftpClient) *utils.SftpClient {
	repoinfo, err := mysqldb.GetRepository()
	if err != nil {
		return nil
	}

	sftpClient, err = utils.OpenSftpConnection(repoinfo)
	if err != nil {
		return nil
	}

	return sftpClient
}

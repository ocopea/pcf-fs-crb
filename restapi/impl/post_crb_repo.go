// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package impl

import (
	"crb/models"
	"crb/restapi/operations/crb_web"
	"crb/utils"
	"errors"
	"net"
	"strconv"
	"strings"

	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// StoreRepositoryInfoHandler inserts repo info into persitent DB.
// Responds with success http response else with appropriate http failure.
func StoreRepositoryInfoHandler(params crb_web.StoreRepositoryInfoParams) middleware.Responder {
	var mysqldb *utils.Database
	err := postrepohandler(params.RepositoryInfo, mysqldb)
	if err == nil {
		return crb_web.NewStoreRepositoryInfoCreated().WithPayload(getRepoInstance(params.HTTPRequest))
	}
	return createStoreRepoFailure(500, err.Error())
}

// Assume HTTP request is valid
func getRepoInstance(postRequest *http.Request) *models.RepositoryInstance {
	copyRepoURL := strings.TrimSuffix(ConstructURL(postRequest, ""), "/")
	return &models.RepositoryInstance{CopyRepoURL: copyRepoURL}
}

// createStoreRepoFailure returns failure response
func createStoreRepoFailure(errorCode int, errorMsg string) *crb_web.StoreRepositoryInfoDefault {
	failedResponse := crb_web.NewStoreRepositoryInfoDefault(errorCode)
	modelErrorInt := int32(errorCode)
	modelsError := models.Error{Code: &modelErrorInt, Message: &errorMsg}
	failedResponse.SetPayload(&modelsError)
	return failedResponse
}

// postrepohandler inserts repo info into the MySQL DB.
// If a table does not exist, then it creates the table.
// It returns nil on success or an error on failure.
func postrepohandler(repoinfo *models.RepositoryInfo, mysqldb utils.DBInterface) error {
	if err := validate(repoinfo); err != nil {
		return err
	}

	if err := mysqldb.CreateRepositoryTable(); err != nil {
		return err
	}

	if err := mysqldb.AddRepository(repoinfo); err != nil {
		return err
	}

	return nil
}

func validate(repoinfo *models.RepositoryInfo) error {
	//first see if address contains a port
	//if it has a port, split off the port numbert and see that it's an int
	//and then take the rest and validate as usual
	host, port, err := net.SplitHostPort(*repoinfo.Addr)
	if err != nil {
		// a host with no port will return an error "address <host>: missing port in address"
		if strings.Contains(err.Error(), utils.MissingPort) {
			host = *repoinfo.Addr
		} else if strings.Contains(err.Error(), utils.TooManyColons) {
			tmp := net.ParseIP(*repoinfo.Addr)
			if tmp != nil { //ipv6 address with no port
				host = *repoinfo.Addr
			}
		} else {
			return err
		}
	}
	//check that port is a uint16
	if port != "" {
		if _, err := strconv.ParseUint(port, 10, 16); err != nil {
			return errors.New("invalid port value")
		}
	}

	addresses, err := net.LookupHost(host)
	if err != nil {
		return err
	}

	if len(addresses) == 0 {
		return errors.New("invalid IP address")
	}

	if strings.TrimSpace(*(repoinfo.User)) == "" {
		return errors.New("user cannot be empty")
	}

	return nil
}

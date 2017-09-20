// Copyright (c) 2017 EMC Corporation. All Rights Reserved.
"use strict";

var httpStatusCode = require('/usr/local/lib/node_modules/http-status-codes');
var crbTestUtils = require('./crb_test_utils');
var sqlutils = require('./sqlutils');

// Test Case: Post repo and then get the same data back
exports.testPostThenGetRepo = function (callback) {
    var testcase = 'testPostThenGetRepository : '
    crbTestUtils.testPostRepoSuccess(testcase, crbTestUtils.repo_host_port, crbTestUtils.repo_password, crbTestUtils.repo_username, function () {
        console.log(testcase + "Posted Repo")
        crbTestUtils.testGetRepoSuccess(testcase, crbTestUtils.repo_host_port, crbTestUtils.repo_password, crbTestUtils.repo_username, function () {
            console.log(testcase + "Retrieved Repo")
            callback();
        });
    });
}

//test post repositories with no address - 422
exports.testPostRepoNoAddress = function (callback) {
    var testcase = 'testPostRepositoryNoAddress : '
    var expectedStatus = httpStatusCode.UNPROCESSABLE_ENTITY
    crbTestUtils.testPostRepoFailure(testcase, void 0, crbTestUtils.repo_password, crbTestUtils.repo_username, expectedStatus, function () {
        console.log(testcase + "Posted Repo - no address")
        callback();
    });
}

//test post repositories with no user - 422
exports.testPostRepoNoUser = function (callback) {
    var testcase = 'testPostRepositoryNoUser : '
    var expectedStatus = httpStatusCode.UNPROCESSABLE_ENTITY
    crbTestUtils.testPostRepoFailure(testcase, crbTestUtils.repo_host_port, crbTestUtils.repo_password, void 0, expectedStatus, function () {
        console.log(testcase + "Posted Repo - no user")
        callback();
    });
}

//test post repositories with invalid address
exports.testPostRepoInvalidAddr = function (callback) {
    var testcase = 'testPostRepositoryInvalidAddress : '
    var invalidAddrPort = 'blah:8080'
	var addrInvalidPort = '1.2.3.4:65536'
	var ipv6AddrInvalidPort = '[::1]:65536'
	var dnsAddrInvalidPort = 'localhost:blah'
	var invalidIpv6AddrFormat = '2001:0db8:85a3:0000:0000:8a2e:0370:7334:8080' // missing [] around ip literal before the port 8080
	var invalidAddr = 'blah'
    var expectedStatus = httpStatusCode.INTERNAL_SERVER_ERROR
    crbTestUtils.testPostRepoFailure(testcase, invalidAddrPort, crbTestUtils.repo_password, crbTestUtils.repo_username, expectedStatus, function () {
        console.log(testcase + "Posted Repo - invalid address with port")
        crbTestUtils.testPostRepoFailure(testcase, addrInvalidPort, crbTestUtils.repo_password, crbTestUtils.repo_username, expectedStatus, function () {
            console.log(testcase + "Posted Repo - address with invalid port")
            crbTestUtils.testPostRepoFailure(testcase, ipv6AddrInvalidPort, crbTestUtils.repo_password, crbTestUtils.repo_username, expectedStatus, function () {
                console.log(testcase + "Posted Repo - ipv6 address with invalid port")
                crbTestUtils.testPostRepoFailure(testcase, dnsAddrInvalidPort, crbTestUtils.repo_password, crbTestUtils.repo_username, expectedStatus, function () {
                    console.log(testcase + "Posted Repo - dns address with invalid port")
                    crbTestUtils.testPostRepoFailure(testcase, invalidIpv6AddrFormat, crbTestUtils.repo_password, crbTestUtils.repo_username, expectedStatus, function () {
                        console.log(testcase + "Posted Repo - invalid ipv6 address port format")
                        crbTestUtils.testPostRepoFailure(testcase, invalidAddr, crbTestUtils.repo_password, crbTestUtils.repo_username, expectedStatus, function () {
                            console.log(testcase + "Posted Repo - invalid address")
                            callback();
                        });
                    });
                });
            });
        });       
    });
}

// test get repositories with no target_credentials table in the database - 404
exports.testGetRepoNoCredTableInDB = function (callback) {
    var testcase = 'testGetRepoNoCredTableInDB :'
    var expectedStatus = httpStatusCode.NOT_FOUND
    var expectedMessage = "The repository information is not available."
    sqlutils.droptables();
    crbTestUtils.testGetRepoFailure(testcase, expectedStatus, expectedMessage, function () {
        console.log(testcase + "Get repo")
        callback();
    });
}


// test get repositories with no target credentials entry in DB table in the database - 404
exports.testGetRepoNoCredsInDB = function (callback) {
    var testcase = 'testGetRepoNoCredsInDB :'
    var expectedStatus = httpStatusCode.NOT_FOUND
    var expectedMessage = "The repository information is not available."
    crbTestUtils.testPostRepoSuccess(testcase, crbTestUtils.repo_host_port, crbTestUtils.repo_password, crbTestUtils.repo_username, function () {
        console.log(testcase + "Posted Repo")
        sqlutils.deleteAllRows('TARGET_CREDENTIALS');
        crbTestUtils.testGetRepoFailure(testcase, expectedStatus, expectedMessage, function () {
            console.log(testcase + "Get repo")
            callback();
        });
    });
}

// test get repositories with no DB Access- 500
exports.testGetRepoNoDBAccess = function (callback) {
    var testcase = 'testGetRepoNoDBAccess :'
    var expectedStatus = httpStatusCode.INTERNAL_SERVER_ERROR
    crbTestUtils.testGetRepoFailureNoMsgValidation(testcase, expectedStatus, function () {
        console.log(testcase + "Get repo")
        callback();
    });
}

// test get repositories with no DB Access- 500
exports.testPostRepoNoDBAccess = function (callback) {
    var testcase = 'testPostRepoNoDBAccess :'
    var expectedStatus = httpStatusCode.INTERNAL_SERVER_ERROR
    var expectedMessage = "Unable to confirm connection is alive, db.Ping error: Error 1045: Access denied for user '"+sqlutils.db_User+"'@'"+sqlutils.db_Ip+"' (using password: YES)"
     crbTestUtils.testPostRepoFailure(testcase, crbTestUtils.repo_host_port, crbTestUtils.repo_password,crbTestUtils.repo_username, expectedStatus, function () {
        console.log(testcase + "Post repo with No DB Access")
        callback();
    });
}

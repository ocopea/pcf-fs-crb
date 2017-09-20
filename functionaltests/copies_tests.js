// Copyright (c) 2017 EMC Corporation. All Rights Reserved.
"use strict";

var httpStatusCode = require('/usr/local/lib/node_modules/http-status-codes');
var crbTestUtils = require('./crb_test_utils');
var sqlutils = require('./sqlutils');

exports.testCopies = function (callback) {
    var testcase = 'testCopiesList : '
    var copyid = 'testCopyList'
    crbTestUtils.testPostRepoSuccess(testcase, crbTestUtils.repo_host_port, crbTestUtils.repo_password, crbTestUtils.repo_username, function () {
        console.log(testcase + "Posted Repo")
        crbTestUtils.testPostCopySuccess(testcase, copyid, function () {
            console.log(testcase + "Posted Copy")
            crbTestUtils.testGetCopiesSuccess(testcase, copyid, function() {
                console.log(testcase + "Retrieved Copy List")
                crbTestUtils.testDeleteSuccess(testcase, copyid, function () {
                    console.log(testcase + "Deleted Copy")
                    callback();
                });
            });
        });
    });
}

exports.testCopiesNoDBTable = function (callback) {
    var testcase = 'testCopiesNoDBTable : '
    var copyid = 'testCopiesNoDBTable'
    var expectedStatusCode = httpStatusCode.INTERNAL_SERVER_ERROR
    crbTestUtils.testPostRepoSuccess(testcase, crbTestUtils.repo_host_port, crbTestUtils.repo_password, crbTestUtils.repo_username, function () {
        console.log(testcase + "Posted Repo")
        crbTestUtils.testPostCopySuccess(testcase, copyid, function () {
            console.log(testcase + "Posted Copy")
            sqlutils.droptable('COPY_REPOSITORY', function () {
                crbTestUtils.testGetCopiesFailureNoMsgValidation(testcase, expectedStatusCode, function () {
                    console.log(testcase + "Retrieved Copy List")
                        callback();
                });
            });
        });
    });
}

// test get copies with no DB Access- 500
exports.testCopiesWithNoDBAccess = function (callback) {
    var testcase = 'testCopiesWithNoDBAccess :'
    var copyid = 'testCopiesWithNoDBAccess'
    var expectedStatusCode = httpStatusCode.INTERNAL_SERVER_ERROR  
    crbTestUtils.testGetCopiesFailureNoMsgValidation(testcase, expectedStatusCode, function () {
        console.log(testcase + "Get copies with no DB access")
        callback();
    });
}
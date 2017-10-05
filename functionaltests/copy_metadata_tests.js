// Copyright (c) 2017 EMC Corporation. All Rights Reserved.
"use strict";

var httpStatusCode = require('/usr/local/lib/node_modules/http-status-codes');
var crbTestUtils = require('./crb_test_utils');

exports.testCopyMetadata = function (callback) {
    var testcase = 'testCopyMetadata : '
    var copyid = 'testMetaCopy'
    crbTestUtils.testPostRepoSuccess(testcase, crbTestUtils.repo_host_port, crbTestUtils.repo_password, crbTestUtils.repo_username, function () {
        console.log(testcase + "Posted Repo")
        crbTestUtils.testPostCopySuccess(testcase, copyid, function () {
            console.log(testcase + "Posted Copy")
            crbTestUtils.testGetCopyMetadataSuccess(testcase, copyid, function() {
                console.log(testcase + "Retrieved Copy Metadata")
                crbTestUtils.testDeleteSuccess(testcase, copyid, function () {
                    console.log(testcase + "Deleted Copy")
                    callback();
                });
            });
        });
    });
}

exports.testCopyMetadataFailureBadID = function (callback) {
    var testcase = 'testCopyMetadataFailureBadID : '
    var copyid = 'idNotPosted'
    var expectedStatus = httpStatusCode.INTERNAL_SERVER_ERROR
    var expectedMessage = "copyID doesn't exist " + copyid
    crbTestUtils.testGetCopyMetadataFailure(testcase, copyid, expectedStatus, expectedMessage, function() {
        console.log(testcase + "Retrieved proper error for bad copyID")
        callback();
    });
}

// test copy metadata  with no DB Access- 500
exports.testCopyMetadataWithNoDBAccess = function (callback) {
    var testcase = 'testCopyMetadataWithNoDBAccess :'
    var copyid = 'testCopyMetadataWithNoDBAccess'
    var expectedStatus = httpStatusCode.INTERNAL_SERVER_ERROR
    crbTestUtils.testGetCopyMetadataFailureNoMsgValidation(testcase, copyid, expectedStatus, function () {
        console.log(testcase + "Get copy metadata with no DB access")
        callback();
    });
}
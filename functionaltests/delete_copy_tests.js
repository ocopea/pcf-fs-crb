// Copyright (c) 2017 EMC Corporation. All Rights Reserved.
"use strict";

var httpStatusCode = require('/usr/local/lib/node_modules/http-status-codes');
var crbTestUtils = require('./crb_test_utils');
var sqlutils = require('./sqlutils');

// Test Case: Delete copy successfully
exports.testDeleteCopy = function (callback) {
    var testcase = 'testDeleteCopy : '
    var copyid = 'testDeleteCopy'
    crbTestUtils.testPostRepoSuccess(testcase, crbTestUtils.repo_host_port, crbTestUtils.repo_password, crbTestUtils.repo_username, function () {
        console.log(testcase + "Posted Repo")
        crbTestUtils.testPostCopySuccess(testcase, copyid, function () {
            console.log(testcase + "Posted Copy")
            crbTestUtils.testDeleteSuccess(testcase, copyid, function () {
                console.log(testcase + "Deleted Copy")
                callback();
            });
        });
    });
}

// Test Case: Delete Copy with invalid repo crdentials
exports.testDeleteCopyInvalidRepoCred = function (callback) {
    var testcase = 'Test_Delete_Copy_Invalid_Repo_Cred : '
    var expectedStatusCode = httpStatusCode.INTERNAL_SERVER_ERROR
    var expectedMessage = crbTestUtils.NO_REPO_CONNECTION
    var copyid = 'Test_Delete_Copy_Invalid_Repo_Cred'
    crbTestUtils.testPostRepoSuccess(testcase, crbTestUtils.repo_host_port, crbTestUtils.repo_password, crbTestUtils.repo_username, function () {
        console.log(testcase + "Posted Repo")
        crbTestUtils.testPostCopySuccess(testcase, copyid, function () {
            console.log(testcase + "Posted Copy")
            crbTestUtils.testPostRepoSuccess(testcase, crbTestUtils.repo_host_port, 'invalid', crbTestUtils.repo_username, function () {
                console.log(testcase + "Posted Invalid Repo")
                crbTestUtils.testDeleteFailure(testcase, copyid, expectedStatusCode, expectedMessage, function () {
                    console.log(testcase + "Deleted Copy")
                    callback();
                });
            });
        });
    });
}

// Test case: Delete copy with no repo credentials
exports.testDeleteCopyNoRepoCredentials = function (callback) {
    var testcase = 'testDeleteCopyNoRepoCredentials : '
    var expectedStatusCode = httpStatusCode.INTERNAL_SERVER_ERROR
    var expectedMessage = crbTestUtils.NO_REPO_CONNECTION
    var copyid = 'testDeleteCopyNoRepoCredentials'
    crbTestUtils.testPostRepoSuccess(testcase, crbTestUtils.repo_host_port, crbTestUtils.repo_password, crbTestUtils.repo_username, function () {
        console.log(testcase + "Posted Repo")
        crbTestUtils.testPostCopySuccess(testcase, copyid, function () {
            console.log(testcase + "Posted copy")
            sqlutils.droptable('TARGET_CREDENTIALS', function () {
                console.log(testcase + "Deleted Repo credentials from DB")
                crbTestUtils.waitTime(testcase, 5000, function () {
                    console.log(testcase + "waited for 5 secs")
                    crbTestUtils.testDeleteFailure(testcase, copyid, expectedStatusCode, expectedMessage, function () {
                        console.log(testcase + "Deleted copy")
                        crbTestUtils.testPostRepoSuccess(testcase, crbTestUtils.repo_host_port, crbTestUtils.repo_password, crbTestUtils.repo_username, function () {
                            console.log(testcase + "Posted Repo")
                            callback();
                        });
                    });
                });
            });
        });
    });
}

// Test case: Delete copy with no copy data on repository
exports.testDeleteCopyNoCopyOnRepo = function (callback) {
    var sftputils = require('./sftputils');
    var testcase = 'testDeleteCopyNoCopyOnRepo : '
    var copyid = 'testDeleteCopyNoCopyOnRepo'
    sftputils.openConnection();
    crbTestUtils.testPostRepoSuccess(testcase, crbTestUtils.repo_host_port, crbTestUtils.repo_password, crbTestUtils.repo_username, function () {
        console.log(testcase + "Posted Repo")
        crbTestUtils.testPostCopySuccess(testcase, copyid, function () {
            console.log(testcase + "Posted copy")
            // remove the copy file externally
            sftputils.deleteFileonRepo(copyid, function () {
                console.log(testcase + "sftp removed file");
                crbTestUtils.testDeleteSuccess(testcase, copyid, function () {
                    console.log(testcase + "Deleted copy")
                    sftputils.closeConnection();
                    callback();
                });
            });
        });
    });
}

// Test case: Delete copy with no metedata
exports.testDeleteCopyNoMetaData = function (callback) {
    var testcase = 'testDeleteCopyNoMetaData : '
    var expectedStatusCode = httpStatusCode.NOT_FOUND
    var expectedMessage = crbTestUtils.NO_COPY_ID_FOUND
    var copyid = 'testDeleteCopyNoMetaData'
    crbTestUtils.testPostRepoSuccess(testcase, crbTestUtils.repo_host_port, crbTestUtils.repo_password, crbTestUtils.repo_username, function () {
        console.log(testcase + "Posted Repo")
        crbTestUtils.testDeleteFailure(testcase, copyid, expectedStatusCode, expectedMessage, function () {
            console.log(testcase + "Delete copy");
            callback();
        });
    });
}

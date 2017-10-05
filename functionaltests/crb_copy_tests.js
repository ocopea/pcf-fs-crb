// Copyright (c) 2017 EMC Corporation. All Rights Reserved.
"use strict";

var httpStatusCode = require('/usr/local/lib/node_modules/http-status-codes');
var crbTestUtils = require('./crb_test_utils');
var sqlutils = require('./sqlutils');

// Test Case: Post copy successfully
exports.testPostCopy = function (callback) {
    var testcase = 'testPostCopy : '
    var copyid = 'testPostCopy'
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


// Test Case: Post copy with no repo credentials - 502
exports.testPostCopyWithNoRepCred = function (callback) {
    var testcase = 'testPostCopyWithNoRepCred : '
    var copyid = 'testPostCopyWithNoRepCred'
    var expectedStatusCode = httpStatusCode.BAD_GATEWAY
    crbTestUtils.testPostRepoSuccess(testcase, crbTestUtils.repo_host_port, crbTestUtils.repo_password, crbTestUtils.repo_username, function () {
        console.log(testcase + "Posted Repo")
        sqlutils.deleteAllRows('TARGET_CREDENTIALS')
        crbTestUtils.testPostCopyFailureNoMsgValidation(testcase, copyid, expectedStatusCode, function () {
            console.log(testcase + "Posted Copy failed")
            callback();
        });
    });
}

// Test Case: Post copy with invalid repo credentials - 500
exports.testPostCopyWithInvalidRepoCred = function (callback) {
    var testcase = 'testPostCopyWithInvalidRepoCred : '
    var copyid = 'testPostCopyWithInvalidRepoCred'
    var expectedStatusCode = httpStatusCode.INTERNAL_SERVER_ERROR
    var expectedMessage = crbTestUtils.SFTP_NOT_ESTABLISHED
    crbTestUtils.testPostRepoSuccess(testcase, crbTestUtils.repo_host_port, "invalidUser", "invalidPassword", function () {
        console.log(testcase + "Posted Repo with invalid credentials")
        crbTestUtils.testPostCopyFailure(testcase, copyid, expectedStatusCode, expectedMessage, function () {
            console.log(testcase + "Posted Copy failed")
            callback();
        });
    });
}
// Test Case: Get copy with invalid repo credentials - 500
exports.testGetCopyWithInvalidRepoCred = function (callback){
    var testcase = 'testGetCopyWithInvalidRepoCred: '
    var copyid = 'testGetCopyWithInvalidRepoCred'
    var expectedStatusCode = httpStatusCode.BAD_GATEWAY
    var expectedMessage = crbTestUtils.BAD_GATEWAY
     crbTestUtils.testPostRepoSuccess(testcase, crbTestUtils.repo_host_port, crbTestUtils.repo_password, crbTestUtils.repo_username, function () {
        console.log(testcase + "Posted Repo")
        crbTestUtils.testPostCopySuccess(testcase, copyid, function () {
            console.log(testcase + "Posted Copy")
            crbTestUtils.testPostRepoSuccess(testcase, crbTestUtils.repo_host_port, "invalidUser", "invalidPassword", function () {
                console.log(testcase + "Posted Repo with invalid credentials")
                crbTestUtils.testGetCopyFailure(testcase, copyid, expectedStatusCode, expectedMessage, function () {
                    console.log(testcase + "Get copy failed")
                    callback();
                });
            });
        });
    });
}
// Test Case: Post copy with existing copy id - 403
exports.testPostCopyWithExistingCopyId = function (callback) {
    var testcase = 'testPostCopyWithExistingCopyId : '
    var copyid = 'testPostCopyWithExistingCopyId'
    var expectedStatusCode = httpStatusCode.FORBIDDEN
    var expectedMessage = crbTestUtils.COPY_ID_ALREADY_EXISTS
    crbTestUtils.testPostRepoSuccess(testcase, crbTestUtils.repo_host_port, crbTestUtils.repo_password, crbTestUtils.repo_username, function () {
        console.log(testcase + "Posted Repo")
        crbTestUtils.testPostCopySuccess(testcase, copyid, function () {
            console.log(testcase + "Posted Copy")
            crbTestUtils.testPostCopyFailure(testcase, copyid, expectedStatusCode, expectedMessage, function () {
                console.log(testcase + "Posted Copy failed with same copy id")
                callback();
            });
        });
    });
}


// Test Case: Post copy with repo credentials with invalid permisiion
exports.testPostCopyRepoCredentialsinvalidPermission = function (callback) {
    var testcase = 'testPostCopyRepoCredentialsinvalidPermission : '
    var copyid = 'testPostCopyRepoCredentialsinvalidPermission'
    var expectedStatusCode = httpStatusCode.INTERNAL_SERVER_ERROR
    var expectedMessage = crbTestUtils.SFTP_NO_PERMISSION   
    crbTestUtils.testPostRepoSuccess(testcase, crbTestUtils.repo_host_port, 'test', 'frisby', function () {
        console.log(testcase + "Posted Repo")
       crbTestUtils.testPostCopyFailure(testcase, copyid, expectedStatusCode, expectedMessage, function () {
                console.log(testcase + "Posted Copy failed with same copy id")
                callback();
            });
    });
}



// Test Case: Post copy with No db access - 500
exports.testPostCopyWithNoDBAccess = function (callback) {
    var testcase = 'testPostCopyWithNoDBAccess :'
    var copyid = 'testPostCopyWithNoDBAccess'
    var expectedStatusCode = httpStatusCode.INTERNAL_SERVER_ERROR
    crbTestUtils.testPostCopyFailureNoMsgValidation(testcase, copyid, expectedStatusCode, function () {
        console.log(testcase + "Post copy with no DB access")
        callback();
    });
}

// Test Case: Get copy successfully
exports.testGetCopy = function (callback) {
    const testcase = 'testGetCopy : '
    const copyid = 'testGetCopy'
    const orginalFile = 'testdata' // File exists in the current directory
    const copyFile = 'copytestdata' // File to be saved to the current directory
    crbTestUtils.testPostRepoSuccess(testcase, crbTestUtils.repo_host_port, crbTestUtils.repo_password, crbTestUtils.repo_username, function () {
        console.log(testcase + "Posted Repo")
        crbTestUtils.postCopyByCurl(copyid, orginalFile);
        crbTestUtils.waitTime(testcase, 5000, function () {
            console.log(testcase + "post copy using curl")
            crbTestUtils.getCopyByCurl(copyid, copyFile);
            crbTestUtils.waitTime(testcase, 5000, function () {
                console.log(testcase + "get copy using curl")
                crbTestUtils.testGetCopySuccess(testcase, copyid, function () {
                    console.log(testcase + "Get Copy")
                    crbTestUtils.validateCopy(orginalFile, copyFile);
                    crbTestUtils.waitTime(testcase, 2000, function () {
                        console.log(testcase + "validated the copy")
                        callback();
                    });
                });
            });
        });
    });
}

// Test Case: Get copy with No repo credentials - 502
exports.testGetCopyWithNoRepoCred = function (callback) {
    var testcase = 'testGetCopyWithNoRepoCred: '
    var copyid = 'testGetCopyWithNoRepoCred'
    var expectedStatusCode = httpStatusCode.BAD_GATEWAY
    var expectedMessage = crbTestUtils.BAD_GATEWAY
    crbTestUtils.testPostRepoSuccess(testcase, crbTestUtils.repo_host_port, crbTestUtils.repo_password, crbTestUtils.repo_username, function () {
        console.log(testcase + "Posted Repo")
        crbTestUtils.testPostCopySuccess(testcase, copyid, function () {
            console.log(testcase + "Posted Copy")
            sqlutils.deleteAllRows('TARGET_CREDENTIALS')
            crbTestUtils.testGetCopyFailure(testcase, copyid, expectedStatusCode, expectedMessage, function () {
                console.log(testcase + "Get copy failed")
                callback();
            });
        });
    });
}

// Test Case: Get copy with invalid copyid - 500
exports.testGetCopyWithInvalidCopyId = function (callback) {
    var testcase = 'testGetCopyWithInvalidCopyId: '
    var copyid = 'testGetCopyWithInvalidCopyId'
    var expectedStatusCode = httpStatusCode.INTERNAL_SERVER_ERROR
    crbTestUtils.testPostRepoSuccess(testcase, crbTestUtils.repo_host_port, crbTestUtils.repo_password, crbTestUtils.repo_username, function () {
        console.log(testcase + "Posted Repo")
        crbTestUtils.testPostCopySuccess(testcase, copyid, function () {
            console.log(testcase + "Posted Copy")
            crbTestUtils.testGetCopyFailureNoMsgValidation(testcase, 'invalidCopyId', expectedStatusCode, function () {
                console.log(testcase + "Get copy failed")
                callback();
            });
        });
    });
}

// Test Case: Get copy with invalid  repo IP - 500
exports.testGetCopyWithInvalidRepoIP = function (callback) {
    var testcase = 'testGetCopyWithInvalidRepoIP: '
    var copyid = 'testGetCopyWithInvalidRepoIP'
    var expectedStatusCode = httpStatusCode.INTERNAL_SERVER_ERROR
    crbTestUtils.testPostRepoSuccess(testcase, crbTestUtils.repo_host_port, crbTestUtils.repo_password, crbTestUtils.repo_username, function () {
        console.log(testcase + "Posted Repo")
        crbTestUtils.testPostCopySuccess(testcase, copyid, function () {
            console.log(testcase + "Posted Copy")
            crbTestUtils.testPostRepoSuccess(testcase, '1.2.3.4', crbTestUtils.repo_password, crbTestUtils.repo_username, function () {
                console.log(testcase + "Posted Repo")
                crbTestUtils.testGetCopyFailureNoMsgValidation(testcase, 'invalidCopyId', expectedStatusCode, function () {
                    console.log(testcase + "Get copy failed")
                    callback();
                });
            });
        });
    });
}

// Test case: Get copy with no copy data on repo system - 500
exports.testGetCopyNoCopyOnRepoSystem = function (callback) {
    var sftputils = require('./sftputils');
    var testcase = 'testGetCopyNoCopyOnRepoSystem : '
    var copyid = 'testGetCopyNoCopyOnRepoSystem'
    var expectedStatusCode = httpStatusCode.INTERNAL_SERVER_ERROR
    sftputils.openConnection();
    crbTestUtils.testPostRepoSuccess(testcase, crbTestUtils.repo_host_port, crbTestUtils.repo_password, crbTestUtils.repo_username, function () {
        console.log(testcase + "Posted Repo")
        crbTestUtils.testPostCopySuccess(testcase, copyid, function () {
            console.log(testcase + "Posted copy")
            // remove the copy file externally
            sftputils.deleteFileonRepo(copyid, function () {
                console.log(testcase + "sftp removed file");
                crbTestUtils.testGetCopyFailureNoMsgValidation(testcase, copyid, expectedStatusCode, function () {
                    console.log(testcase + "Deleted copy")
                    sftputils.closeConnection();
                    callback();
                });
            });
        });
    });
}

// Test Case: Post copy with No db access - 500
exports.testGetCopyWithNoDBAccess = function (callback) {
    var testcase = 'testPostCopyWithNoDBAccess :'
    var copyid = 'testPostCopyWithNoDBAccess'
    var expectedStatusCode = httpStatusCode.INTERNAL_SERVER_ERROR
    crbTestUtils.testGetCopyFailureNoMsgValidation(testcase, copyid, expectedStatusCode, function () {
        console.log(testcase + "Get copy with no DB access")
        callback();
    });
}


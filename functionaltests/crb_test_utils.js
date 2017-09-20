// Copyright (c) 2017 EMC Corporation. All Rights Reserved.
"use strict";
var frisby = require('/usr/local/lib/node_modules/frisby');
var fs = require('fs');
var httpStatusCode = require('/usr/local/lib/node_modules/http-status-codes');
var formData = require('/usr/local/lib/node_modules/form-data');
var path = require('/usr/local/lib/node_modules/path');
var child_process = require('child_process');

// CRB Repository information
exports.repo_host;
exports.repo_username;
exports.repo_password;
exports.repo_port;
exports.repo_host_port;

// set up CRB name
exports.crbName;

// Setup CRB urls
var crburls = { baseUrl: "", copiesUrl: "", reposUrl: "", infoUrl: "" }

exports.setCrbUrls = function (crbIP) {
    crburls.baseUrl = 'http://' + crbIP + '/crb';
    crburls.copiesUrl = crburls.baseUrl + '/copies';
    crburls.reposUrl = crburls.baseUrl + '/repositories';
    crburls.infoUrl = crburls.baseUrl + '/info';
}

// set up CRB version
var crbVersion = "";
exports.setCrbVersion = function (version) {
    crbVersion = version;
}

// Set frisby defaults
frisby.globalSetup({
    timeout: 10000
})

// CRB Error Messages
exports.NO_REPO_CONNECTION = "Connection to the repository couldn't be established";
exports.NO_COPY_ID_FOUND = "Copy ID not found";
exports.COPY_ID_ALREADY_EXISTS = "copyID already exists";
exports.SFTP_NOT_ESTABLISHED = "Sftp connection was not established";
exports.SFTP_NO_PERMISSION = 'sftp: "Permission denied" (SSH_FX_PERMISSION_DENIED)';
exports.BAD_GATEWAY = '502 Bad Gateway: Registered endpoint failed to handle the request.';

var testFilePath = path.resolve(__dirname, 'testdata')

// Test functions for each API
exports.testDeleteSuccess = function (testcase, copyID, afterFunction) {
    frisby.create(testcase + "Delete Copy")
        .delete(crburls.copiesUrl + '/' + copyID)
        .expectStatus(httpStatusCode.OK)
        .after(function () {
            afterFunction();
        })
        .toss();
}

exports.testDeleteFailure = function (testcase, copyID, statuscode, message, afterFunction) {
    frisby.create(testcase + "Delete Copy Fail")
        .delete(crburls.copiesUrl + '/' + copyID)
        .expectStatus(statuscode)
        .expectJSON({ 'message': message })
        .after(function () {
            afterFunction();
        })
        .toss();
}

exports.testPostRepoSuccess = function (testcase, addr, password, user, afterFunction) {
    frisby.create(testcase + 'Post Repo')
        .post(crburls.reposUrl,
        { 'addr': addr, 'password': password, 'user': user },
        { json: true },
        { headers: { 'Content-Type': 'application/json' } })
        //.inspectRequest()
        .expectStatus(httpStatusCode.CREATED)
        .expectJSON({
            'copyRepoURL': crburls.reposUrl
        })
        .expectJSONTypes({ 'copyRepoURL': String })
        //.inspectHeaders()
        //.inspectBody()
        .after(function () {
            afterFunction();
        })
        .toss();
}

exports.testPostRepoFailure = function (testcase, addr, password, user, status, afterFunction) {
    frisby.create(testcase + 'Post Repo')
        .post(crburls.reposUrl,
        { 'addr': addr, 'password': password, 'user': user },
        { json: true },
        { headers: { 'Content-Type': 'application/json' } })
        .expectStatus(status)
        .expectJSONTypes({ 'message': String })
        .after(function () {
            afterFunction();
        })
        .toss();
}

exports.testGetRepoSuccess = function (testcase, addr, password, user, afterFunction) {
    frisby.create(testcase + 'Get Repo')
        .get(crburls.reposUrl)
        .expectStatus(httpStatusCode.OK)
        .expectJSON({
            'addr': addr, 'password': password, 'user': user
        })
        .after(function () {
            afterFunction();
        })
        .toss();
}


exports.testGetRepoFailure = function (testcase, status, message, afterFunction) {
    frisby.create(testcase + 'Get Repo')
        .get(crburls.reposUrl)
        .expectStatus(status)
        .expectJSON({ 'message': message })
        .after(function () {
            afterFunction();
        })
        .toss();
}

exports.testGetRepoFailureNoMsgValidation = function (testcase, status, afterFunction) {
    frisby.create(testcase + 'Get Repo')
        .get(crburls.reposUrl)
        .expectStatus(status)
        .expectJSONTypes({ 'message': String })
        .after(function () {
            afterFunction();
        })
        .toss();
}

exports.testGetInfoSuccess = function (testcase, afterFunction) {
    frisby.create(testcase + 'Get Info')
        .get(crburls.infoUrl)
        .expectStatus(httpStatusCode.OK)
        .expectJSON({
            'name': 'crb',
            'version': crbVersion,
            'repoType': 'file'
        })
        .after(function () {
            afterFunction();
        })
        .toss();
}

// Function used to add delay in Frisby tests
exports.waitTime = function (testcase, time, afterFunction) {
    frisby.create(testcase + 'Get Info')
        .get(crburls.infoUrl)
        .expectStatus(httpStatusCode.OK)
        .expectJSON({
            'name': 'crb',
            'version': crbVersion,
            'repoType': 'file'
        })
        .after(function () {
            afterFunction();
        })
        .waits(time)
        .toss();
}

exports.testGetCopySuccess = function (testcase, copyID, afterFunction) {
    frisby.create(testcase + "Get Copy")
        .get(crburls.copiesUrl + '/' + copyID + '/data', {
            json: false
        })
        .expectStatus(httpStatusCode.OK)
        .after(function () {
            afterFunction();
        })
        .toss();
}

exports.testPostCopySuccess = function (testcase, copyID, afterFunction) {
    var form = new formData();
    form.append('file', fs.createReadStream(testFilePath), {
        knownLength: fs.statSync(testFilePath).size
    });

    frisby.create(testcase + "Post Copy")
        .post(crburls.copiesUrl + '/' + copyID + '/data', form, {
            json: false,
            headers: {
                'Content-Type': 'application/octet-stream',
                'Content-Length': form.getLengthSync()
            }
        })
        .expectStatus(httpStatusCode.CREATED)
        .expectJSON({
            'copyURL': crburls.copiesUrl + '/' + copyID
        })
        .after(function () {
            afterFunction();
        })
        .toss();
}

exports.testPostCopyFailure = function (testcase, copyID, expectedStatus, expectedMessage, afterFunction) {
    var form = new formData();
    form.append('file', fs.createReadStream(testFilePath), {
        knownLength: fs.statSync(testFilePath).size
    });

    frisby.create(testcase + "Post Copy")
        .post(crburls.copiesUrl + '/' + copyID + '/data', form, {
            json: false,
            headers: {
                'Content-Type': 'application/octet-stream',
                'Content-Length': form.getLengthSync()
            }
        })
        .expectStatus(expectedStatus)
        .expectJSONTypes({ 'message': String })
        .expectJSON({
            'message': expectedMessage
        })
        .after(function () {
            afterFunction();
        })
        .toss();
}

exports.testGetCopyFailure = function (testcase, copyID, expectedStatus, expectedMessage, afterFunction) {
    frisby.create(testcase + "Get Copy")
        .get(crburls.copiesUrl + '/' + copyID + '/data')
        .expectStatus(expectedStatus)
        .expectBodyContains(expectedMessage)
        .after(function () {
            afterFunction();
        })
        .toss();
}

exports.testGetCopyFailureNoMsgValidation = function (testcase, copyID, expectedStatus, afterFunction) {
    frisby.create(testcase + "Get Copy")
        .get(crburls.copiesUrl + '/' + copyID + '/data')
        .expectStatus(expectedStatus)
        .timeout(30000)
        .after(function () {
            afterFunction();
        })
        .toss();
}

exports.testPostCopyFailureNoMsgValidation = function (testcase, copyID, expectedStatus, afterFunction) {
    var form = new formData();
    form.append('file', fs.createReadStream(testFilePath), {
        knownLength: fs.statSync(testFilePath).size
    });

    frisby.create(testcase + "Post Copy")
        .post(crburls.copiesUrl + '/' + copyID + '/data', form, {
            json: false,
            headers: {
                'Content-Type': 'application/octet-stream',
                'Content-Length': form.getLengthSync()
            }
        })
        .expectStatus(expectedStatus)
        .after(function () {
            afterFunction();
        })
        .toss();
}

exports.testGetCopyMetadataSuccess = function (testcase, copyID, afterFunction) {
    var metaUrl = crburls.copiesUrl + '/' + copyID
    frisby.create(testcase + 'Get Copy Metadata')
        .get(metaUrl)
        .expectStatus(httpStatusCode.OK)
        .expectJSONTypes({
            'copySize': Number,
            'copyTimeStamp': String
        })
        .expectJSON({
            'copyId': copyID,
            'copySize': 280   //this is hardcoded to match our datafile at this time, if source file changes this line must go/change
        })
        .after(function () {
            afterFunction();
        })
        .toss();
}

exports.testGetCopyMetadataFailure = function (testcase, copyID, expectedStatus, expectedMessage, afterFunction) {
    var metaUrl = crburls.copiesUrl + '/' + copyID
    frisby.create(testcase + 'Get Copy Metadata')
        .get(metaUrl)
        .expectStatus(expectedStatus)
        .expectJSONTypes({ 'message': String })
        .expectJSON({
            'message': expectedMessage
        })
        .after(function () {
            afterFunction();
        })
        .toss();
}

exports.testGetCopyMetadataFailureNoMsgValidation = function (testcase, copyID, expectedStatus, afterFunction) {
    var metaUrl = crburls.copiesUrl + '/' + copyID
    frisby.create(testcase + 'Get Copy Metadata')
        .get(metaUrl)
        .expectStatus(expectedStatus)
        .expectJSONTypes({ 'message': String })
        .after(function () {
            afterFunction();
        })
        .toss();
}

exports.testGetCopiesSuccess = function (testcase, copyID, afterFunction) {
    var metaUrl = crburls.copiesUrl + '/' + copyID
    frisby.create(testcase + 'Get Copies')
        .get(crburls.copiesUrl)
        .expectStatus(httpStatusCode.OK)
        .expectJSON('?', {
            'copyURL': metaUrl
        })
        .after(function () {
            afterFunction();
        })
        .toss();
}

exports.testGetCopiesFailureNoMsgValidation = function (testcase, expectedStatus, afterFunction) {
    frisby.create(testcase + 'Get Copies')
        .get(crburls.copiesUrl)
        .expectStatus(expectedStatus)
        .expectJSONTypes({ 'message': String })
        .after(function () {
            afterFunction();
        })
        .toss();
}

exports.validateCopy = function (orginalFile, copyFile) {
    const originalFileBuffer = fs.readFileSync(orginalFile);
    const copyFileBuffer = fs.readFileSync(copyFile);

    expect(originalFileBuffer.length).toBe(copyFileBuffer.length)

    const crypto = require('crypto');
    const origFileHash = crypto.createHash('md5').update(originalFileBuffer).digest('hex');
    const copyFileHash = crypto.createHash('md5').update(copyFileBuffer).digest('hex');

    expect(origFileHash).toEqual(copyFileHash)
}

// function to get Copy by using curl
exports.getCopyByCurl = function (copyid, copyFile) {
    const url = crburls.copiesUrl + '/' + copyid + '/data'
    const cmd = "curl -s -o " + copyFile + " " + url;

    child_process.execSync(cmd);
}

// function to post Copy by using curl
exports.postCopyByCurl = function (copyid, orginalFile) {
    const url = crburls.copiesUrl + '/' + copyid + '/data'
    const cmd = "curl -s -T " + orginalFile + " -X POST " + url;

    child_process.execSync(cmd);
}
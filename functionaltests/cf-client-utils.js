// Copyright (c) 2017 EMC Corporation. All Rights Reserved.
"use strict";
var cfclient = require("/usr/local/lib/node_modules/cf-nodejs-client");
var crbTestUtils = require('./crb_test_utils');

const UsersUAA = new (cfclient).UsersUAA;

var cf_username;
var cf_password;
var cf_endpoint;
var token;
var crbName;
var crbAppGUID;
var serviceBindingGuid;
var CloudController;
var ServiceBindings;
var Apps;

exports.setCFcredentials = function(endpoint, user, pwd){
    cf_username = user;
    cf_password = pwd;
    cf_endpoint = "https://"+endpoint;

    CloudController = new (cfclient).CloudController(cf_endpoint);
    ServiceBindings = new (cfclient).ServiceBindings(cf_endpoint);
    Apps = new (cfclient).Apps(cf_endpoint);
}

exports.getCrbServiceBindings = function () {
    CloudController.getInfo().then((result) => {
        UsersUAA.setEndPoint(result.authorization_endpoint);
        return UsersUAA.login(cf_username, cf_password);
    }).then((result) => {
        token = result;
        Apps.setToken(token);
        return Apps.getApps();
    }).then((result) => {
        for (var i = 0; i < result["resources"].length; i++) {
            if (result["resources"][i]["entity"].name == crbTestUtils.crbName) {
                crbName = result["resources"][i]["entity"].name;
                crbAppGUID = result["resources"][i]["metadata"].guid;
                console.log("CRB : " + crbName, " CRB App guid: " + crbAppGUID);
                break;
            }
        }
        return Apps.getServiceBindings(crbAppGUID);
    }).then((result) => {
        serviceBindingGuid = result["resources"][0]["metadata"].guid;
        ServiceBindings.setToken(token);
        return ServiceBindings;
    }).then((result) => {
        console.log("Successfull getting service bindings");
    }).catch((reason) => {
        console.error("cf-client-utils getCrbServiceBindings Error: " + reason);
    });
}
// test get repositories with no DB Access- 404
exports.dbSetupAndUnbindCRB = function (callback) {
    var testcase = 'dbSetupAndUnbindCRB :'
    var copyid = 'dbSetupAndUnbindCRB'
    crbTestUtils.testPostRepoSuccess(testcase, crbTestUtils.repo_host_port, crbTestUtils.repo_password, crbTestUtils.repo_username, function () {
        console.log(testcase + "Posted Repo")
        crbTestUtils.testPostCopySuccess(testcase, copyid, function () {
            console.log(testcase + "Post Copy")
            ServiceBindings.remove(serviceBindingGuid).then((result) => {
                console.log("Unbind successful")
            }).catch((reason) => {
                console.error("cf-client-utils dbSetupAndUnbindCRB Error:" + reason);
            })
            crbTestUtils.waitTime(testcase, 5000, function () {
                console.log(testcase + "Get Info after wait")
                callback();
            });
        });
    });
}


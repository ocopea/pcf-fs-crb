// Copyright (c) 2017 EMC Corporation. All Rights Reserved.
"use strict";

var crbTestUtils = require('./crb_test_utils');


exports.testGetInfo = function (callback) {
    var testcase = 'testGetInfo : '
    crbTestUtils.testGetInfoSuccess(testcase, function () {
        console.log(testcase + "Retrieved CRB Info")
        callback();
    });
}
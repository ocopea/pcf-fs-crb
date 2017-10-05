// Copyright (c) 2017 EMC Corporation. All Rights Reserved.
"use strict";

var Client = require('/usr/local/lib/node_modules/ssh2-sftp-client');
var sftp = new Client();
var crbTestUtils = require('./crb_test_utils');
const PATH = '/var/lib/crb/'
exports.openConnection = function () {
    sftp.connect({
        host: crbTestUtils.repo_host,
        port: crbTestUtils.repo_port,
        username: crbTestUtils.repo_username,
        password: crbTestUtils.repo_password
    }).catch((err) => {
        console.log(err);
    });
}

exports.deleteFileonRepo = function (copyID, callback) {
    try {
        sftp.delete(PATH + copyID);
    }
    catch (err) {
        console.log(err);
    }
    finally {        
        callback();
    }

}

exports.closeConnection = function(){
    sftp.end();
}

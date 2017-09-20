// Copyright (c) 2017 EMC Corporation. All Rights Reserved.
"use strict";

var sqlutils = require('./sqlutils');
var crbTestUtils = require('./crb_test_utils');
var deleteCopyTests = require('./delete_copy_tests')
var crbRepoTests = require('./crb_repo_tests')
var crbInfoTests = require('./crb_info_tests')
var crbCopyTests = require('./crb_copy_tests')
var copyMetaTests = require('./copy_metadata_tests')
var copyListTests = require('./copies_tests')
var cfclientUtils = require('./cf-client-utils')
var crbip = process.env.host;
var dbIp = process.env.dbip;
var dbUser = process.env.dbuser;
var dbPassword = process.env.dbpassword;
var dbName = process.env.dbname;
var version = process.env.version;

// Settng repo credentials
crbTestUtils.repo_host = process.env.repoip;
crbTestUtils.repo_password = process.env.repopassword;
crbTestUtils.repo_username = process.env.repouser;
crbTestUtils.repo_port = process.env.repoport;
crbTestUtils.repo_host_port = crbTestUtils.repo_host + ':' + crbTestUtils.repo_port;

// Setting crb-server IP
crbTestUtils.setCrbUrls(crbip);

// Setting up CRB name
crbTestUtils.crbName = process.env.crbname;

// Setting CRB version
crbTestUtils.setCrbVersion(version);

//setting up CF credentials
cfclientUtils.setCFcredentials(process.env.cfendpoint, process.env.cfusername, process.env.cfpassword)

// Get CRB service bindings 
// This function uses CF JS client to get the CRB service bindings 
// Calling this function before the tests so that the function can run syncronusly with other tests and get the service bindings before the unbind tests are started.
cfclientUtils.getCrbServiceBindings();

var tests = [deleteCopyTests.testDeleteCopy,
deleteCopyTests.testDeleteCopyInvalidRepoCred,
deleteCopyTests.testDeleteCopyNoRepoCredentials,
deleteCopyTests.testDeleteCopyNoCopyOnRepo,
deleteCopyTests.testDeleteCopyNoMetaData,
crbRepoTests.testPostThenGetRepo,
crbRepoTests.testPostRepoNoAddress,
crbRepoTests.testPostRepoNoUser,
crbRepoTests.testPostRepoInvalidAddr,
crbRepoTests.testGetRepoNoCredTableInDB,
crbRepoTests.testGetRepoNoCredsInDB,
crbInfoTests.testGetInfo,
copyMetaTests.testCopyMetadata,
copyMetaTests.testCopyMetadataFailureBadID,
copyListTests.testCopies,
copyListTests.testCopiesNoDBTable,
crbCopyTests.testPostCopy,
crbCopyTests.testPostCopyWithNoRepCred,
crbCopyTests.testPostCopyWithInvalidRepoCred,
crbCopyTests.testPostCopyRepoCredentialsinvalidPermission,
crbCopyTests.testPostCopyWithExistingCopyId,
crbCopyTests.testGetCopy,
crbCopyTests.testGetCopyWithInvalidRepoCred,
crbCopyTests.testGetCopyWithNoRepoCred,
crbCopyTests.testGetCopyWithInvalidCopyId,
crbCopyTests.testGetCopyWithInvalidRepoIP,
crbCopyTests.testGetCopyNoCopyOnRepoSystem,
// Unbinding the CRB DB service. All the subsequent tests will be based on CRB with No DB Access
cfclientUtils.dbSetupAndUnbindCRB,
crbRepoTests.testGetRepoNoDBAccess,
crbRepoTests.testPostRepoNoDBAccess,
copyMetaTests.testCopyMetadataWithNoDBAccess,
copyListTests.testCopiesWithNoDBAccess,
crbCopyTests.testPostCopyWithNoDBAccess,
crbCopyTests.testGetCopyWithNoDBAccess
];

// Drop DB tables
sqlutils.connectMe(dbIp, dbUser, dbPassword, dbName);
// Commenting the dropTables as it is not needed since a new mysql service
// is created for each run. Also, there is a race condition between the execution
// of the drop table and the execution of the first test. If one wants to enable
// it in the future, the drop tables should be added as tests like the others to
// be executed sequentially.
//sqlutils.droptables();



function run() {
    // Run the first test from the tests array
    if (tests.length > 0) {
        tests[0](function () {            
            doNext();
        });
    }
    else {
        console.log("cleaning up db connection")
        sqlutils.endConnection();
    }
}

function doNext() {
    if (tests.length > 0) {
        tests.shift();
        run();
    }
}

run();


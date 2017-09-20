// Copyright (c) 2017 EMC Corporation. All Rights Reserved.
"use strict";

var mysql = require('/usr/local/lib/node_modules/mysql');
var connection;

exports.connectMe = function (dbIp, dbuser, dbPassword, dbName) {
  var values = { host: dbIp, user: dbuser, password: dbPassword, database: dbName }
  connection = mysql.createConnection(values);
}

exports.droptable = function (table, callback) {
  connection.query('DROP TABLE ' + table, function (error, results, fields) {
    if (error) {
      console.log(error)
    }
  });
  callback();
}

exports.endConnection = function () {
  connection.end();
}

exports.droptables = function () {
  connection.query('DROP TABLE COPY_REPOSITORY', function (error, results, fields) {
    if (error) {
      console.log(error)
    }
  });
  connection.query('DROP TABLE TARGET_CREDENTIALS', function (error, results, fields) {
    if (error) {
      console.log(error)
    }
  });
}

exports.deleteAllRows = function (table) {
  connection.query('TRUNCATE TABLE ' + table, function (error) {
    if (error) {
      console.log(error)
    }
  });
}

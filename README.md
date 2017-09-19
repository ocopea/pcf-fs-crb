PCF File System CRB
====

# Overview

FS CRB abstracts various types of persistent stores that keeps point-in-time copies of application data and provides key capability of storing and retrieving application data.
CRB service provides capabilities such as:
* Store application data on a storage system
* Store and manage copy metadata
* Provide a catalog of all application data copies and their metadata
* Allow copy data to be retrieved based on specified copy metadata
* Data transfer using SFTP

# API

The functionality above is described using OpenAPI specification by [swagger.xml](https://github.com/ocopea/pcf-fs-crb/blob/master/swagger.yaml) through [Swagger framework](https://swagger.io).

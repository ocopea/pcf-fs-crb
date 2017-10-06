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

The functionality above is described using OpenAPI specification by [swagger.yaml](swagger.yaml) through [Swagger framework](https://swagger.io).

# Functional Tests
[functionaltests](/functionaltests) contains the functional tests that validate the functionality of the CRB using [REST API](swagger.yaml)

# How to build
The crb binary is built using a [container](Dockerfile) which has the required build tools.

* git clone https://github.com/ocopea/pcf-fs-crb.git crb.

  **Note**: You must name the destination folder crb as go uses folder names as package names.
* cd crb
* Build the docker image. Run
  ```
  docker build -t ocopea-crb-build-tool .
  ```
* Build the binary. Run
  ```
  docker run -v <path-to-crb-directory>:/go/src/crb --name ocopea-crb-build-tool ocopea-crb-build-tool
  
  ```
  For example,
  ```
  docker run -v /go/src/crb:/go/src/crb --name ocopea-crb-build-tool ocopea-crb-build-tool
  ```
  
  The crb binary 'crb-server' will be saved at the cmd/crb-server/ folder.

# How to run
* Login pcf.
* Create the crb-mysql service which is needed by the crb app to store copy metadata. Run
  ```
  cf create-service p-mysql <service plan> crb-mysql
* Deploy the crb app. Run
  ```
  cf push -f <path-to-manifest>
  ```  
  For example,
  ```
  cf push -f /go/src/crb/manifest.yml
  ```  

* After the app is deployed, make sure the crb-mysql service is bound to the app successfully. Run
  ```
  cf s
  ```
  
  Here is the sample output:
  ```
  Getting services in org pcfdev-org / space pcfdev-space as admin...
  OK

  name        service   plan   bound apps   last operation
  crb-mysql   p-mysql   1gb    crb-server   create succeeded
  ```

#!/bin/bash
# Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.

CRB_SRC=$GOPATH/src/crb
CRB_BIN=cmd/crb-server/crb-server

echo 'Building pcf-fs-crb binary...'
cd $CRB_SRC
swagger generate server -f swagger.yaml -A crb
echo -e '\nInstalling packages...'
godep restore
go install ./...
cd cmd/crb-server
echo -e '\nBuilding binary...'
go build
chmod -R 777 $GOPATH/src/crb
echo -e '\nBuild completed. The binary is stored at '$CRB_BIN'.\n'


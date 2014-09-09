#!/bin/bash
#Don't change this to $home it will fail in ansible
export GOPATH=/home/capillaryDeploy/go
export GOROOT=/usr/local/go
export GOBIN=/home/capillaryDeploy/go/bin
export PATH=$PATH:$GOBIN:$GOROOT/bin
dir=$GOPATH/src/github.com/CapillarySoftware/goforward
install=/usr/local/perceptor/goforward
cd $dir
go get github.com/tools/godep
go install github.com/tools/godep
godep restore
godep go build
mkdir -p $install
cp goforward  $install/
cp seelog.xml $install/
//remove all source code after install
rm -rf $GOPATH

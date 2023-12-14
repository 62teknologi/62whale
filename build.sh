#!/bin/bash
APP_NAME="whale"
OUTPUT="binary"
#execute submodule and update
git submodule init
git submodule update
go mod tidy
#building section
ARCH="amd64"
env GOOS=darwin GOARCH=amd64 go build -x -v -o ${APP_NAME}-mac-${ARCH} main.go
env GOOS=linux GOARCH=amd64 go build -x -v -o ${APP_NAME}-linux-${ARCH} main.go
env GOOS=windows GOARCH=amd64 go build -x -v -o ${APP_NAME}-windows-${ARCH}.exe main.go
ARCH="arm64"
env GOOS=darwin GOARCH=arm64 go build -x -v -o ${APP_NAME}-mac-${ARCH} main.go
#move binary to ./build dir
rm -rf ${OUTPUT} || exit_on_error "${OUTPUT} folder didn't exist"
mkdir ${OUTPUT} || exit_on_error "${OUTPUT} folder exist"
mv ${APP_NAME}* ${OUTPUT}/
chmod +x ${OUTPUT}/*

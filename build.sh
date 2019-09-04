#!/bin/bash

cd frontend

npm run build

cd ..

mkdir -p builds

rice embed-go

env GOOS=darwin GOARCH=amd64 go build -o builds/server_darwin_amd64

env GOOS=linux GOARCH=amd64 go build -o builds/server_linux_amd64

env GOOS=windows GOARCH=amd64 go build -o builds/server_win_amd64

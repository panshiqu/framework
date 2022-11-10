#!/bin/bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o db_server ../db.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o manager_server ../manager.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o proxy_server ../proxy.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o login_server ../login.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o game_server ../game.go

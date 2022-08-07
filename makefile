SHELL := /bin/bash

tidy:
	go mod tidy 
	go mod vendor

run:
	go run app/arcaIndustria40/main.go

build:
	GOOS=windows GOARCH=amd64 go build -o bin/arca_industria_4_0_backend.exe ./app/arcaIndustria40/main.go


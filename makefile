SHELL := /bin/bash

tidy:
	go mod tidy 
	go mod vendor

run:
	go run app/arcaIndustria40/main.go

build-raspian:
	GOOS=linux GOARCH=arm go build -o bin/arca_industria_4_0_backend ./app/arcaIndustria40/main.go


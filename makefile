SHELL := /bin/bash

tidy:
	go mod tidy 
	go mod vendor

run:
	go run app/contact-api/main.go

admin:
	go run app/admin-api/main.go


FROM golang:1.17-alpine
WORKDIR /app
COPY . /app

RUN go build -o organization_ app/organization_-api/main.go

EXPOSE 9000

CMD [ "./organization" ]


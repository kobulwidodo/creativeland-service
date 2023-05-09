FROM golang:alpine

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN ls ./etc/cfg/

# install swagg
RUN go install github.com/swaggo/swag/cmd/swag@v1.6.7

#run swag init & move to path docs/swagger
RUN `go env GOPATH`/bin/swag init -g src/cmd/main.go -o docs/swagger --parseInternal

#download library or package 
RUN go mod download

RUN go build -o binary ./src/cmd

ENTRYPOINT ["/app/binary"]

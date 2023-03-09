FROM golang:1.19-alpine as builder

ENV GOPATH=/

COPY ../go.mod ../go.sum ./
RUN go mod download

COPY ./client ./

# build go app
RUN go build -o demo-docker-client .

#Build destination container
FROM alpine:latest

ENV GOPATH=/go

# copy bin
COPY --from=builder $GOPATH/demo-docker-client ./
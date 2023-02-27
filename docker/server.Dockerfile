FROM golang:1.19-alpine as builder

ENV GOPATH=/

COPY ../go.mod ../go.sum ./
RUN go mod download

COPY ./server ./

# build go app
RUN go build -o demo-docker-server ./main.go ./metrics.go

#Build destination container
FROM alpine:latest

ENV GOPATH=/go

# copy bin
COPY --from=builder $GOPATH/demo-docker-server ./

CMD [ "mkdir", "./log" ]

EXPOSE 8080
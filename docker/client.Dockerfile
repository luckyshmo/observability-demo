FROM golang:1.19-alpine as builder
WORKDIR /src/
COPY ../go.mod ../go.sum /src/
RUN go mod download
COPY ./client ./
RUN CGO_ENABLED=0 go build -o /bin/client .

FROM scratch
COPY --from=builder /bin/client /bin/client
ENTRYPOINT ["/bin/client"]
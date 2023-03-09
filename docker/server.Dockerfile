FROM golang:1.19-alpine as builder
WORKDIR /src/
COPY ../go.mod ../go.sum /src/
RUN go mod download
COPY ./server ./
RUN CGO_ENABLED=0 go build -o /bin/server .

FROM scratch
COPY --from=builder /bin/server /bin/server
ENTRYPOINT ["/bin/server"]
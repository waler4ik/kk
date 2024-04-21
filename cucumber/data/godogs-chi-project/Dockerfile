ARG GO_VERSION=1.20

FROM golang:${GO_VERSION}-alpine AS builder

RUN mkdir /user && \
  echo 'somebody:x:65534:65534:somebody:/:' > /user/passwd && \
  echo 'somebody:x:65534:' > /user/group

RUN apk add --no-cache ca-certificates

WORKDIR /src

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o /server ./cmd/server/main.go

FROM alpine:3.17.2 AS final
RUN apk update && apk upgrade

COPY --from=builder /user/group /user/passwd /etc/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /server /server

EXPOSE 6060
EXPOSE 8080

USER somebody:somebody

ENTRYPOINT ["/server","--host=0.0.0.0","--port=8080","--keep-alive=600s"]
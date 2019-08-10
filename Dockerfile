FROM alpine as go_build
RUN apk update
RUN apk upgrade
RUN apk add --update git go gcc g++

WORKDIR /go/src/github.com/simon987/ws_bucket/
ENV GOPATH /go

COPY api api
COPY *.go .
RUN go get . && CGO_ENABLED=1 GOOS=linux go build -a -o ws_bucket .

FROM alpine
WORKDIR /root

COPY --from=go_build ["/go/src/github.com/simon987/ws_bucket/ws_bucket", "./"]

ENV WS_BUCKET_DIALECT=sqlite3
ENV WS_BUCKET_CONNSTR=wsb.db
ENV WS_BUCKET_WORKDIR=/data
ENV WS_BUCKET_LOGLEVEL=info

VOLUME ["/data"]

CMD ["/root/ws_bucket"]


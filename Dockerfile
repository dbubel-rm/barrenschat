FROM golang:latest
RUN apt-get -y update
RUN apt-get -y upgrade
RUN apt-get install -y sqlite3 libsqlite3-dev

RUN mkdir -p /go/src/github.com/dbubel/bchat-api/internal/platform/web
ADD . /go/src/github.com/dbubel/bchat-api/internal/platform/web
RUN /usr/bin/sqlite3 /go/src/github.com/dbubel/bchat-api/internal/platform/web/cmd/bishopfox-api/bfscans.db < /go/src/github.com/dbubel/bchat-api/internal/platform/web/schema.sql 
WORKDIR /go/src/github.com/dbubel/bchat-api/internal/platform/web/cmd/bishopfox-api

#RUN CGO_ENABLED=0 go build -a -v -ldflags '-extldflags "-static"' main.go
RUN go build -v main.go
ENTRYPOINT [ "./main" ]
FROM golang:1

ADD app/ /go/src/github.com/globocom/secDevLabs/owasp-top10-2016-mobile/m4/note-box/server/app
WORKDIR /go/src/github.com/globocom/secDevLabs/owasp-top10-2016-mobile/m4/note-box/server/app
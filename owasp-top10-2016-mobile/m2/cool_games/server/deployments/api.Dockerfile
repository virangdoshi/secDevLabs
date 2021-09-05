FROM golang:1.17.0-bullseye

ADD app/ /go/src/github.com/globocom/secDevLabs/owasp-top10-2016-mobile/m2/cool_games/server/app
WORKDIR /go/src/github.com/globocom/secDevLabs/owasp-top10-2016-mobile/m2/cool_games/server/app
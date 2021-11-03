FROM golang:1.17.0-bullseye

ADD /api/ /app
WORKDIR /app
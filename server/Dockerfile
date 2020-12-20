# syntax = docker/dockerfile:1-experimental
FROM golang:1-alpine as build

RUN apk add --no-cache git
RUN go get github.com/markbates/pkger/cmd/pkger

ENV GOARCH=arm

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download -x

COPY . .
RUN /go/bin/pkger -o files
RUN --mount=type=cache,target=/root/.cache/go-build go build -v .

FROM scratch as export
COPY --from=build /build/spotify-screen /

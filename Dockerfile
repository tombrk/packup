# UI
FROM --platform=$BUILDPLATFORM node:14-alpine as js
COPY ./ui /ui
WORKDIR /ui
RUN yarn && yarn build

# Go Environment
FROM golang:1.19-alpine as env
RUN go env | grep -E 'GOARCH|GOOS|GOARM' > /go.env

# Go
FROM --platform=$BUILDPLATFORM golang:1.19-alpine as go
COPY --from=env /go.env /go.env
WORKDIR /packup
COPY go.mod go.sum .
RUN apk add --no-cache make git
RUN go mod download
RUN source /go.env && go env -w GOOS=$GOOS GOARCH=$GOARCH GOARM=$GOARM
COPY . .

FROM go as go-server
COPY --from=js /ui/build ui/build
RUN make server

FROM go as go-agent
RUN make agent

FROM alpine:3.17 as base
RUN apk add --no-cache restic coreutils
WORKDIR /backups

# agent
FROM base as agent
RUN apk add --no-cache sqlite postgresql13-client
COPY --from=go-agent /packup/packup-agent /usr/local/bin
COPY ./mods /mods
RUN chmod +x /mods/*
ENTRYPOINT ["packup-agent"]

# server
FROM base as server
COPY --from=go-server /packup/packup-server /usr/local/bin
ENTRYPOINT ["packup-server"]
CMD ["--config=/etc/packup/packup.yaml"]

# explicit default target
FROM server

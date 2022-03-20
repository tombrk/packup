# UI
FROM node:14-alpine as js
COPY ./ui /ui
WORKDIR /ui
RUN yarn && yarn build

# Go
FROM golang:1.18 as go
WORKDIR /packup
COPY go.mod go.sum .
RUN go mod download
COPY . .

FROM go as go-server
COPY --from=js /ui/build ui/build
RUN make server

FROM go as go-agent
RUN make agent

FROM alpine as base
RUN apk add --no-cache restic coreutils
WORKDIR /backups

# agent
FROM base as agent
RUN apk add --no-cache sqlite
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

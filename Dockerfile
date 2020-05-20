FROM node:alpine as js
COPY ./ui /ui
WORKDIR /ui
RUN yarn && yarn build

FROM golang as builder
RUN GO111MODULE=on go get github.com/markbates/pkger/cmd/pkger
COPY . /app
COPY --from=js /ui/build /app/ui/build
WORKDIR /app
RUN make static

FROM alpine
COPY --from=builder /app/packup /usr/local/bin
RUN apk add restic coreutils
WORKDIR /backups
ENTRYPOINT ["/usr/local/bin/packup"]
CMD ["--config=/etc/packup/packup.yaml"]

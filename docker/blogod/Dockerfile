FROM golang:1.16-alpine AS builder

RUN apk add --update npm make tzdata && npm install -g sass uglifycss uglify-js

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make clean && make release


FROM alpine

RUN apk add --update --no-cache bash git

WORKDIR /blogo

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

COPY --from=builder /build/build/blogod /blogo/blogod
COPY --from=builder /build/dist /blogo/dist
COPY --from=builder /build/web/template /blogo/web/template
COPY --from=builder /build/store/migration /blogo/store/migration
COPY --from=builder /build/script/wait-for-it.sh /blogo/script/

ENV GIN_MODE=release
ENV TZ Asia/Shanghai

EXPOSE 8080


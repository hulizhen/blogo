FROM golang:1.16

WORKDIR /blogo

COPY tmp/blogod /blogo/blogod
COPY dist /blogo/dist
COPY web/template /blogo/web/template
COPY store/migration /blogo/store/migration
COPY script/wait-for-it.sh /blogo/script/

ENV GIN_MODE=release

EXPOSE 8080

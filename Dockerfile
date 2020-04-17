FROM debian:stable-slim

MAINTAINER xiayesuifeng "xiayesuifeng@firerain.me"

WORKDIR /goblog

ENV GOBLOG_WEB_PATH /goblog/web

COPY goblog /goblog

RUN apt-get update -qq && apt-get install -y -qq ca-certificates

EXPOSE 80 443 8080

VOLUME /goblog/data

ENTRYPOINT ["./goblog"]
CMD ["-c","./config.json"]
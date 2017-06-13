FROM alpine:3.5

ENV CONFIG_FILEPATH /etc/push/config.toml

RUN apk --update --no-cache add ca-certificates

ADD push_linux_amd64 /usr/local/bin/push
RUN chmod +x /usr/local/bin/push

EXPOSE 3000
CMD push
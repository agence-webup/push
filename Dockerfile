FROM golang:1.12.9-alpine as builder

ENV USER root
WORKDIR /root/src

RUN apk --no-cache add make git

ADD . /root/src

RUN make build


FROM alpine:3.10
COPY --from=builder /root/src/push_linux_amd64 /usr/local/bin/pushapi
RUN chmod +x /usr/local/bin/pushapi

ENV CONFIG_FILEPATH config.toml
EXPOSE 3000
CMD ["pushapi"]
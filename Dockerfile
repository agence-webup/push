FROM golang:1.7.1-alpine

RUN apk add --no-cache git

ADD https://github.com/Masterminds/glide/releases/download/v0.12.1/glide-v0.12.1-linux-amd64.tar.gz /glide.tar.gz
RUN mkdir /glide-bin
RUN tar xzf /glide.tar.gz -C /glide-bin

RUN mkdir -p /go/src/webup/push
WORKDIR /go/src/webup/push

COPY . /go/src/webup/push

RUN /glide-bin/linux-amd64/glide install
RUN cd /go/src/webup/push/cmd/push && go install -v

CMD /go/bin/push

EXPOSE 3000
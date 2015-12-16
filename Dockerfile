FROM gliderlabs/alpine

MAINTAINER hayden@coinbase.com

RUN apk-install go git
RUN mkdir -p /go/src /go/bin && chmod -R 777 /go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

WORKDIR /go
COPY sit.go /go/
RUN go build sit.go

CMD ["/go/sit"]

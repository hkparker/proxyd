FROM quay.cbhq.net/containers/alpine:3.3

MAINTAINER hayden@hkparker.com

# ca-certificates?
RUN apk-install go git
RUN mkdir -p /go/src /go/bin && chmod -R 777 /go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN go get github.com/hkparker/TTPD
WORKDIR /go/src/github.com/hkparker/TTPD/
RUN go build

CMD ["/go/src/github.com/hkparker/TTPD/TTPD"]

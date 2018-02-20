FROM golang:1.8
ADD . /go/src/github.com/nickschuch/sherlock
WORKDIR /go/src/github.com/nickschuch/sherlock
RUN go get github.com/mitchellh/gox
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/src/github.com/nickschuch/sherlock/bin/sherlock_linux_amd64 /usr/local/bin/sherlock
CMD ["sherlock", "watson"]

FROM golang:1.8
RUN go get github.com/mitchellh/gox
ADD workspace /go
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/bin/watson_linux_amd64 /usr/local/bin/watson
CMD ["watson"]

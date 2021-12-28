FROM golang:1.17-alpine AS build-env

COPY . $GOPATH/src/Smarthome
WORKDIR $GOPATH/src/Smarthome

# go get modfile dependencies
RUN go get -d -v
RUN go build -o /go/bin/smarthome

COPY scripts /go/bin/scripts
COPY scriptfiles /go/bin/scriptfiles
COPY html /go/bin/html

WORKDIR /go/bin/
EXPOSE 8080
ENTRYPOINT ["./smarthome"]
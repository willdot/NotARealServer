FROM golang:1.12

ADD . /go/src/app
WORKDIR /go/src/app

# Go get dependancies
RUN go get -d -v ./...

ENV PORT=8080
ENV REQUESTDIRECTORY=requests

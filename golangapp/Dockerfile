FROM golang:1.11
ADD . /go/src/github.com/angadthandi/golangmongoapp/golangapp
WORKDIR /go/src/github.com/angadthandi/golangmongoapp/golangapp
RUN go build

CMD ["./golangapp"]
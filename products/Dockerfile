FROM golang:1.11
ADD . /go/src/github.com/angadthandi/golangmongoapp/products
WORKDIR /go/src/github.com/angadthandi/golangmongoapp/products
RUN go get -v
RUN go build -race

CMD ["./products"]
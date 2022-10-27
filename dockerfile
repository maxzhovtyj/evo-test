FROM golang:1.18-alpine

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go mod download github.com/ugorji/go
RUN go build -o main ./cmd/main.go

CMD ["./main"]

EXPOSE 8080
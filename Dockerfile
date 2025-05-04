FROM golang:latest

COPY ./ ./

COPY go.mod go.sum ./

RUN go mod download

RUN go build -o go-app ./main.go 

EXPOSE 9090

CMD ["./go-app"]
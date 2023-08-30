FROM golang:1.21  as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /go-docker-demo

CMD [ "/go-docker-demo" ]
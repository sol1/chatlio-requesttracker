FROM golang:1.14

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /chatlio-rt

EXPOSE 8080

CMD [ "/chatlio-rt" ]
FROM golang:1.21.6

WORKDIR /app

COPY go.mod .
COPY internal .
COPY cmd .
COPY migrations .
COPY remote .
LABEL authors="User"

RUN go build -o bin .

ENTRYPOINT ["top", "-b"]
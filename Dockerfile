FROM golang:latest

ENV APP_HOME /app

RUN mkdir $APP_HOME
WORKDIR $APP_HOME

COPY go.mod .
COPY go.sum .

RUN go mod download

EXPOSE 8080

COPY . $APP_HOME%
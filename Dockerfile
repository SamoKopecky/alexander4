FROM golang:latest

WORKDIR /alex

ENV GOPATH=
ENV ALEX_TOKEN=
ENV ALEX_PIPE=

COPY . .

RUN mkdir pipe
RUN go build -v -o main

CMD ./main 

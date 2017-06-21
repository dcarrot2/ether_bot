FROM golang:1.8.1

RUN mkdir -p /go/src/ether_bot
WORKDIR /go/src/ether_bot

COPY . /go/src/ether_bot

RUN go-wrapper download
RUN go-wrapper install

ENV PORT 80

EXPOSE 80

CMD ["go-wrapper", "run"]

FROM golang:1.15-alpine

WORKDIR /home/ubuntu/golang

COPY . .

COPY ./cs-sogang-ac-kr.crt /usr/local/share/ca-certificates

RUN apk add --no-cache ca-certificates

RUN update-ca-certificates

RUN go build

CMD ["./main"]

EXPOSE 4567

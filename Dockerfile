FROM golang:1.15-alpine
WORKDIR /home/ubuntu/golang
COPY . .
RUN go build
CMD ["./main"]
EXPOSE 4567
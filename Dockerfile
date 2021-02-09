FROM golang:1.15-alpine
WORKDIR /home/ubuntu/golang
COPY . .
RUN go build
CMD ["go", "run", "main.go"]
EXPOSE 4567
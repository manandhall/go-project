FROM golang:1.17-alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go mod tidy
COPY . ./
RUN go build -o main .
EXPOSE 5000
CMD ["/app/main"]

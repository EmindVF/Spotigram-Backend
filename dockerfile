FROM golang:alpine
WORKDIR build
ADD ./go.mod .
COPY . .
RUN go build -o ./bin/application ./cmd/app/main.go
EXPOSE 8080
CMD ["./bin/application"]
FROM golang:1.22.2
RUN apt-get update && apt-get install -y libwebp-dev
WORKDIR build
ADD ./go.mod .
COPY . .
RUN CGO_ENABLED=1 go build -o ./bin/application ./cmd/app/main.go
EXPOSE 8080
CMD ["./bin/application"]
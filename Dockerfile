FROM golang:1.24-alpine

WORKDIR /app
COPY . .

RUN go mod download 
RUN cd ./cmd/api && go build -o ./bin/api .

EXPOSE 8080

ENTRYPOINT ["/app/cmd/api/bin/api"]
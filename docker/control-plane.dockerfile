FROM golang:1.25-alpine

WORKDIR /app

# copy only control‑plane folder
COPY control-plane/ .

RUN go mod download

RUN go build -o control-plane ./cmd/server

EXPOSE 8080

CMD ["./control-plane"]

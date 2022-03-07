FROM golang:1.17.6-alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./ ./
RUN go build -o avito-test ./cmd/main.go

CMD [ "./avito-test" ]
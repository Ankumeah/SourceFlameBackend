FROM golang:alpine AS builder

WORKDIR /opt/gc/

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o ./bin/gc .

FROM scratch

COPY --from=builder /opt/gc/bin/gc .

ENTRYPOINT [ "./gc" ]

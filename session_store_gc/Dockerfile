FROM golang:alpine AS builder

WORKDIR /opt/gc/

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build go build -ldflags="-s -w" -o ./bin/gc .

FROM alpine

COPY --from=builder /opt/gc/bin/gc /opt/gc/

ENTRYPOINT [ "/opt/gc/gc" ]

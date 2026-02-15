FROM golang:alpine AS builder

WORKDIR /opt/backend/

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o ./bin/deltabase-server ./cmd/deltabase-server/

FROM alpine

RUN apk add libgit2 --no-cache

ENV GIN_MODE=release
COPY --from=builder /opt/backend/bin/deltabase-server /opt/backend/deltabase-server

ENTRYPOINT [ "/opt/backend/deltabase-server" ]

FROM golang:alpine AS builder

WORKDIR /opt/backend/

COPY go.mod go.sum ./
RUN go mod download

COPY . .
#RUN go build -ldflags="-s -w" -o ./bin/sourceflame-server ./cmd/sourceflame-server/
RUN --mount=type=cache,target=/root/.cache/go-build go build -ldflags="-s -w" -o ./bin/sourceflame-server ./cmd/sourceflame-server/

FROM alpine

RUN apk add git --no-cache

ENV GIN_MODE=release
COPY --from=builder /opt/backend/bin/sourceflame-server /opt/backend/sourceflame-server

ENTRYPOINT [ "/opt/backend/sourceflame-server" ]

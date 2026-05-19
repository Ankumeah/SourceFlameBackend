FROM golang:alpine AS builder

WORKDIR /opt/backend/

RUN apk add gcc musl-dev --no-cache

COPY go.mod go.sum ./
RUN go mod download

ARG db_type=sqlite3

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build go build -ldflags="-s -w" -tags=${db_type} -o ./bin/sourceflame-server ./cmd/sourceflame-server/

FROM alpine

RUN apk add git --no-cache

WORKDIR /opt/backend/

ENV GIN_MODE=release
COPY alpha-01.sql .
COPY --from=builder /opt/backend/bin/sourceflame-server .

ENTRYPOINT [ "./sourceflame-server" ]

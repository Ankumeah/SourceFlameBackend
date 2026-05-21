FROM golang:alpine AS builder

WORKDIR /opt/db_init/

RUN apk add gcc musl-dev --no-cache

COPY go.mod go.sum ./
RUN go mod download

ARG db_type=postgres

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build go build -ldflags="-s -w" -tags=${db_type} -o ./bin/db-init ./

FROM alpine

COPY --from=builder /opt/db_init/bin/db-init /opt/db_init/db-init

ENTRYPOINT [ "/opt/db_init/db-init" ]

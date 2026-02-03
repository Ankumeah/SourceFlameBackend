FROM golang:alpine

WORKDIR /opt/backend/

COPY ./go.mod ./go.sum ./
RUN go mod download

ENV GIN_MODE=release
RUN go build -ldflags="-s -w" -o ./bin/deltabase-server ./cmd/deltabase-server/

ENTRYPOINT [ "./bin/deltabase-server" ]

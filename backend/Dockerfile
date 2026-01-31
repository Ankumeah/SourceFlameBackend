FROM golang:alpine

WORKDIR /opt/backend/

COPY . .

ENV GIN_MODE=release

RUN go build -ldflags="-s -w" -o ./bin/deltabase-server ./cmd/deltabase-server/

ENTRYPOINT [ "./bin/deltabase-server" ]

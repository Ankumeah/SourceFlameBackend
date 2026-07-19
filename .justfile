default: vet test fmt

# Build and export dockerfiles
[linux, macos]
build *dirs:
  @echo "==> Building {{ dirs }}"
  @for dir in {{ dirs }}; do just _build ${dir}; done

  @echo
[linux, macos]
_build dir:
  #! /bin/env sh

  image_name="$(just _sanitize $(basename $(pwd))-$(basename {{ dir }})):local"

  echo "-> docker build -t ${image_name} {{ dir }}"
  docker build -t ${image_name} {{ dir }}

  echo

# Export docker images
[linux, macos]
expo *dirs:
  @echo "==> Exporting images  {{ dirs }}"
  @for dir in {{ dirs }}; do just _expo ${dir}; done

  @echo
[linux, macos]
_expo dir:
  #! /bin/env sh

  mkdir -p images/
  image_name="$(just _sanitize $(basename $(pwd))-$(basename {{ dir }})):local"

  echo "-> docker image save ${image_name} -o images/$(just _sanitize ${image_name}).tar"
  docker image save ${image_name} -o images/$(just _sanitize ${image_name}).tar

# Load docker images
[linux, macos]
load *archives = "images/*":
  @echo "==> Loading images  {{ archives }}"
  @for archive in {{ archives }}; do just _load "${archive}"; done

  @echo
[linux, macos]
_load archive:
  @echo "-> docker image load -i {{ archive }}"
  @docker image load -i {{ archive }} 1>/dev/null

# Run go test
[linux, macos]
test *options = "-tags='sqlite3'":
  @export $(cat ./env.d.example/*)

  cd ./backend/ && go test {{ options }} ./... | column -t -s $'\t'
  cd ./db_init/ && go test {{ options }} ./... | column -t -s $'\t'

# Run go vet
[linux, macos]
vet *options = "-tags='sqlite3'":
  @export $(cat ./env.d.example/*)

  cd ./backend/ && go vet {{ options }} ./...
  cd ./db_init/ && go vet {{ options }} ./...

# Run go fmt
[linux, macos]
fmt *options = "":
  @export $(cat ./env.d.example/*)

  cd ./backend/ && go fmt {{ options }} ./...

  cd ./db_init/ && go fmt {{ options }} ./...

[linux, macos]
_sanitize string:
  @echo $(basename "{{ string }}" | tr '[:upper:]' '[:lower:]' | sed "s/:/-/g" | sed "s/\//_/g" )

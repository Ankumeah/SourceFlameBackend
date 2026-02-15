# Build dockerfiles
[linux, macos]
build *dirs:
  @for dir in {{ dirs }}; do just _build "${dir}"; done
[linux, macos]
_build dir:
  @echo docker build -t "$(basename "$(pwd)" | tr '[:upper:]' '[:lower:]')-$(basename "{{ dir }}" | tr '[:upper:]' '[:lower:]'):local" "{{ dir }}"
  @docker build -t "$(basename "$(pwd)" | tr '[:upper:]' '[:lower:]')-$(basename "{{ dir }}" | tr '[:upper:]' '[:lower:]'):local" "{{ dir }}"

  @echo docker image save "$(basename "$(pwd)" | tr '[:upper:]' '[:lower:]')-$(basename "{{ dir }}" | tr '[:upper:]' '[:lower:]'):local" -o images/$(echo "$(basename "$(pwd)" | tr '[:upper:]' '[:lower:]')-$(basename "{{ dir }}" | tr '[:upper:]' '[:lower:]'):local" | sed "s/:/_/g" | sed "s/\//_/g").tar
  @docker image save "$(basename "$(pwd)" | tr '[:upper:]' '[:lower:]')-$(basename "{{ dir }}" | tr '[:upper:]' '[:lower:]'):local" -o images/$(echo "$(basename "$(pwd)" | tr '[:upper:]' '[:lower:]')-$(basename "{{ dir }}" | tr '[:upper:]' '[:lower:]'):local" | sed "s/:/_/g" | sed "s/\//_/g").tar

# Export docker images
[linux, macos]
expo *images:
  @for image in {{ images }}; do just _expo "${image}"; done
[linux, macos]
_expo image:
  @echo docker image save {{ image }} -o images/$(echo {{ image }} | sed "s/:/_/g" | sed "s/\//_/g").tar

  @mkdir -p images/
  @docker image save {{ image }} -o images/$(echo {{ image }} | sed "s/:/_/g" | sed "s/\//_/g").tar

# Load docker images
[linux, macos]
load *archives:
  @for archive in {{ archives }}; do just _load "${archive}"; done
[linux, macos]
_load archive:
  @echo docker image load -i {{ archive }}

  @docker image load -i {{ archive }}

# Run go tests on ./backend/
[linux, macos]
test *options:
  @echo go test ./... {{options}}

  @cd ./backend/ && go test ./... {{options}}

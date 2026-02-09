# Build dockerfiles
build *dirs:
  @for dir in {{ dirs }}; do just _build "${dir}"; done
_build dir:
  @echo docker build -t "$(basename "$(pwd)" | tr '[:upper:]' '[:lower:]')-$(basename "{{ dir }}" | tr '[:upper:]' '[:lower:]'):local" "{{ dir }}"

  @docker build -t "$(basename "$(pwd)" | tr '[:upper:]' '[:lower:]')-$(basename "{{ dir }}" | tr '[:upper:]' '[:lower:]'):local" "{{ dir }}"

# Export docker images
expo *images:
  @for image in {{ images }}; do just _expo "${image}"; done
_expo image:
  @echo docker image save {{ image }} -o images/$(echo {{ image }} | sed "s/:/_/g").tar

  @mkdir -p images/
  @docker image save {{ image }} -o images/$(echo {{ image }} | sed "s/:/_/g").tar

# Load docker images
load *archives:
  @for archive in {{ archives }}; do just _load "${archive}"; done
_load archive:
  @echo docker image load -i {{ archive }}

  @docker image load -i {{ archive }}

# Run go tests on ./backend/
test *options:
  @echo go test ./... {{options}}

  @cd ./backend/ && go test ./... {{options}}

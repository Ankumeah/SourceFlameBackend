# Build and export dockerfiles
[linux, macos]
build *dirs:
  @echo "==> Building {{ dirs }}"
  @for dir in {{ dirs }}; do just _build ${dir}; done

  @echo "==> Exporting images  {{ dirs }}"
  @for dir in {{ dirs }}; do just _expo ${dir}; done

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

# Run go tests on ./backend/
[linux, macos]
test *options:
  #! /bin/env sh

  echo "==> Running tests"

  echo  "-> go test ./... {{ options }}"
  cd ./backend/ && go test ./... {{ options }}

  echo

[linux, macos]
k8s base_dir = "k8s":
  #! /bin/env zsh

  echo "==> Applying all mainafests in  {{ base_dir }}"

  cd {{ base_dir }}

  for secret in secrets/*; do
    echo "-> kubectl apply -f ${secret}"
    kubectl apply -f ${secret} 1>/dev/null
  done
  for configmap in config-maps/*; do
    echo "-> kubectl apply -f ${configmap}"
    kubectl apply -f ${configmap} 1>/dev/null
  done

  echo

  for file in **/*.yaml; do
    if [[ $file == "secrets"* ]] || [[ $file == "config-maps"* ]]; then
      continue;
    fi

    echo "-> kubectl apply -f ${file}"
    kubectl apply -f ${file} 1>/dev/null
  done

  echo

  echo "==> k8s cluster status:"
  kubectl get all

  echo

[linux, macos]
_sanitize string:
  @echo $(basename "{{ string }}" | tr '[:upper:]' '[:lower:]' | sed "s/:/-/g" | sed "s/\//_/g" )

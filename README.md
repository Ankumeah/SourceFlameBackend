# DeltaBaseBackend
DeltaBase is a currently **work in progress** git hosting platform, this is the backend for it written in *golang*

![compile-check](https://github.com/Ankumeah/DeltaBaseBackend/actions/workflows/compile-check.yaml/badge.svg) ![go-tests](https://github.com/Ankumeah/DeltaBaseBackend/actions/workflows/go-tests.yaml/badge.svg)

## How to run
> WARNING: This repo is a Work in progress and has just
recently reached a barely working version, try at your own risk

> NOTE: This repo only contains the code for a backend,
while you can use this on its own if you are willing to interact
with it through raw urls, you probably want to install a
frontend for it, the frontend can be found here [Work in progress]()
### With docker compose
- Clone and cd in the repo: `git clone https://github.com/Ankumeah/DeltaBaseBackend && cd DeltaBaseBackend`
- Change the secrets in `./env.d.example/` to your likeing (**VERY IMPORTANT**)
- Start docker compose: `docker compose up -d`
### With kustomize
- Clone and cd in the repo: `git clone https://github.com/Ankumeah/DeltaBaseBackend && cd DeltaBaseBackend`
- Change the secrets in `./k8s/base/secrets.example/` to your likeing (**VERY IMPORTANT**)
- Move the files to the actual dir: `mv ./k8s/base/secrets.example/ ./k8s/base/secrets/`
- Apply the files: `kubectl apply -k ./k8s/base/`

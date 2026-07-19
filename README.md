# SourceFlameBackend
SourceFlame is a currently **work in progress** git hosting platform, this is the backend for it written in *golang*
![go-tests](https://github.com/Ankumeah/SourceFlameBackend/actions/workflows/go-tests.yaml/badge.svg)

## How to run

> [!WARNING]
> This repo is a Work in progress and has just
recently reached a barely working version, try at your own risk

> [!NOTE]
> This repo only contains the code for a backend,
while you can use this on its own if you are willing to interact
with it through raw urls, you probably want to install a
frontend for it, the frontend can be found here [Work in progress]()

### With docker compose

  > [!NOTE]
  > A new dir called data/ will appear in the project containing the data created by the project

  > [!NOTE]
  > Currently we ship support for postgres and sqlite3 out of the box,
  out of these sqlite3 is the default option when running with docker compose,
  while for any use case involving docker compose, sqlite3 is recommended,
  if you wish to switch to postgres have a look at docker-compose.yaml

- Clone and cd in the repo: `git clone https://github.com/Ankumeah/SourceFlameBackend && cd SourceFlameBackend`
- Change the secrets in `./env.d.example/` to your liking (**VERY IMPORTANT**)
- Start docker compose: `docker compose up -d`

### With kustomize

  ![Read here](https://github.com/Ankumeah/SourceFlameBackend/blob/master/k8s/base/README.md)

## How to use

> [!WARNING]
> The openapi spec may be out of date as it is harder to maintain
due to me not using an auto generator
but I try my best to make sure the quick cheatsheet is always up to date
so funnily enough the rough cheatsheet is more reliable than the formal spec

  ![openapi.yaml](https://github.com/Ankumeah/SourceFlameBackend/blob/master/openapi.yaml)

  ![Quick cheatsheet](https://github.com/Ankumeah/SourceFlameBackend/blob/master/doc.md)

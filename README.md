# DeltaBaseBackend
DeltaBase is a currently **work in progress** git hosting platform, this is the backend for it written in *golang*

![compile-check](https://github.com/Ankumeah/DeltaBaseBackend/actions/workflows/compile-check.yaml/badge.svg) ![go-tests](https://github.com/Ankumeah/DeltaBaseBackend/actions/workflows/go-tests.yaml/badge.svg)

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

- Clone and cd in the repo: `git clone https://github.com/Ankumeah/DeltaBaseBackend && cd DeltaBaseBackend`
- Change the secrets in `./env.d.example/` to your likeing (**VERY IMPORTANT**)
- Start docker compose: `docker compose up -d`

### With kustomize

  > [!NOTE]
  > This dir does not include the yaml files needed to deploy the required
  storage or the PVC used to connect the storage to the app. It is the user's
  responsibility to setup and connect the storage. The app assumes the
  presence of two PVCs by the name of `git-storage` and `database-storage`

- Clone and cd in the repo: `git clone https://github.com/Ankumeah/DeltaBaseBackend && cd DeltaBaseBackend`
- Change the secrets in `./k8s/base/secrets.example/` to your liking (**VERY IMPORTANT**)
- Move the files to the actual dir: `mv ./k8s/base/secrets.example/ ./k8s/base/secrets/`
- Add PVCs for your storage with names `git-storage` and `database-storage`
- Take a look at the yaml files and patch any values you want (Optional but recommended)
- Apply the files: `kubectl apply -k ./k8s/base/` (Adjust accordingly if patched)

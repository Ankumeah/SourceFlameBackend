# Deplyment with kustomize

> [!NOTE]
> This project uses k8s hostPath storage by default (/k8s/base/pv.yaml)
want it is recommended to change this to an actual storage provider
of your choice, the yaml resources to do so are not provided in this repo
and are assumed as the users responsibility. This app assumes the presence
of two PVCs by the names `git-storage` and `database-storage`

- Clone and cd in the repo: `git clone https://github.com/Ankumeah/SourceFlameBackend && cd SourceFlameBackend`
- Change the secrets in `./k8s/base/secrets.example/` to your liking (**VERY IMPORTANT**)
- Move the files to the actual dir: `mv ./k8s/base/secrets.example/ ./k8s/base/secrets/`
- Replace PVCs for your storage with names `git-storage` and `database-storage` (Optional but **HIGHLY** recommended)
- Take a look at the yaml files and patch any values you want (Optional but recommended)
- Apply the files: `kubectl apply -k ./k8s/base/` (Adjust accordingly if patched)

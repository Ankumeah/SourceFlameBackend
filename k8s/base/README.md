# Deplyment with kustomize

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

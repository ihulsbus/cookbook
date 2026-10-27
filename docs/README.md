# docs

Helm chart that deploys a static OpenAPI docs UI (Swagger) to Kubernetes.

The chart's `files/openapi.yaml` is a build-time copy of the repo's single
source-of-truth spec at `../openapi.yaml`. Helm can only package files that
live inside the chart directory, so run `./sync-openapi.sh` to refresh that
copy before `helm package` / `helm upgrade --install` whenever the root
`openapi.yaml` changes.

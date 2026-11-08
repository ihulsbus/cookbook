# docker-compose config files

Each microservice's `config.yaml` is mandatory - `viper.ReadInConfig()`
(see `<service>/internal/config/initFuncs.go`) fatals at startup if it
can't find one, regardless of what environment variables are set. So
`docker-compose.yml` mounts a real file into each service's container at
`/config/config.yaml` instead of trying to configure it via
`environment:`.

## First-time setup

For each service directory under `docker/config/`, copy the `.example`
file to `config.yaml` (already done if you just cloned this checkout and
these files already exist - only needed again after a fresh clone, or if
you've deleted your local copies):

```bash
for d in docker/config/*/; do
  cp "${d}config.yaml.example" "${d}config.yaml"
done
```

`config.yaml` is gitignored (`**/config.yaml` in the repo's
`.gitignore`, same rule the per-service `cmd/<service>/config.yaml` dev
files use) - your copy, including anything you put in it, never gets
committed. `config.yaml.example` is the tracked template.

Then fill in each file's `oauth.clientsecret` from your own Authentik
instance (see `authentik/README.md` and
`authentik/blueprints/cookbook-service.yaml`) - the example files ship
with it empty since it's a real credential. Nothing else needs changing
to run against this compose file's own Postgres/RabbitMQ/MinIO - the
`.example` files already point at those containers by their
docker-compose service names (`postgres`, `rabbitmq`, `minio`).

## Changing a setting

Edit the service's `config.yaml` (not `.example`) directly, then
`docker-compose restart <service>` - each service's `viper.WatchConfig()`
also picks up most changes live, but a restart is the reliable way to be
sure.

## Keeping `.example` in sync

If you add a new config key a service reads (a new `shared/httpclient`
dependency's base URL, say), update that service's `config.yaml.example`
too, with the docker-compose-network-appropriate value (container
hostnames like `postgres`/`rabbitmq`/`recipe-service`, not `localhost` -
unlike `cmd/<service>/config.yaml`, these run the service itself inside
a container on `cookbook-network`, not on your host).

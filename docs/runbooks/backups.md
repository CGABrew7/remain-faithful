# Database backups

The API app `remain-faithful-api` does not store the database. `backend/fly.toml` has no volume. The app reads `DATABASE_URL` and talks to a separate Fly Postgres app. Backups belong on that database app, not on the API machine.

Find the database before you change anything.

```bash
fly postgres list
fly mpg list
```

Use the unmanaged steps when `fly postgres list` shows the database. Use the managed steps when `fly mpg list` shows it. The API app name in the commands below is `remain-faithful-api`.

## Unmanaged Fly Postgres

Fly takes a daily snapshot of each volume and keeps it for 5 days unless you set a longer retention.

```bash
fly volumes list -a <postgres-app-name>
fly volumes update <volume-id> --snapshot-retention 14
fly volumes snapshots list <volume-id>
fly volumes snapshots create <volume-id>
```

Snapshot the primary `pg_data` volume. If the cluster has a replica, snapshot the primary.

To restore, create a new Postgres app from a snapshot. Do not overwrite the current volume in place.

```bash
fly image show -a <postgres-app-name>
fly postgres create --snapshot-id <snapshot-id> --image-ref <repository>:<tag>
```

Use the same image repository and tag the current database is running. The new volume must be at least as large as the old one.

Check the restored app, then point the API at it.

```bash
fly postgres detach <postgres-app-name> -a remain-faithful-api
fly postgres attach <restored-postgres-app> -a remain-faithful-api
```

`fly postgres attach` rewrites `DATABASE_URL` on `remain-faithful-api`. Restart or redeploy the API after the attach so every machine picks up the new address.

```bash
fly status -a remain-faithful-api
curl -fsS https://remain-faithful-api.fly.dev/health
```

Keep the old database app until you have signed in and seen a known account on the restored one.

## Managed Postgres

```bash
fly mpg backup list <cluster-id>
fly mpg backup create <cluster-id> --type full
fly mpg restore <cluster-id> --backup-id <backup-id> --name remain-faithful-db-restored
```

Restore creates a new cluster and leaves the current one in place. When the new cluster is ready, set `DATABASE_URL` on `remain-faithful-api` to that cluster and redeploy. Confirm `/health`, then sign in and look up a known account before you remove the old cluster.

## What this file does not change

`backend/fly.toml` stays a stateless API config. Adding a volume there would not back up Postgres.

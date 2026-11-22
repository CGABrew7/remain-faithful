# Rollback

Use this when a release of Remain Faithful should come off production. Roll back the surface that changed. The API, the website, and the iPhone app ship separately.

## API on Fly

App name: `remain-faithful-api`.

```bash
fly releases list -a remain-faithful-api
fly releases rollback -a remain-faithful-api
```

`fly releases rollback` with no version returns to the previous release. To pick one, pass the version from the list.

```bash
fly releases rollback <version> -a remain-faithful-api
curl -fsS https://remain-faithful-api.fly.dev/health
```

This puts the previous image back. It does not restore the database. For data, use `docs/runbooks/backups.md`.

## Website on Vercel

Project: `remain-faithful`.

1. Open the project in the Vercel dashboard.
2. Open Deployments.
3. Find the last production deployment you want to keep.
4. Use Instant Rollback on that deployment.

The rollback serves that build on remainfaithful.com. It does not change environment variables. If a bad variable shipped with the deploy, fix the variable and redeploy.

## iPhone app

App Store Connect, app Remain Faithful, version 1.0.4.

1. Open the version.
2. Use Phased Release.
3. Pause the phased release to stop the update from reaching more people, or release to all users when the build is good.

A phased release only affects people who have not installed that version yet. People who already updated keep the build until a newer version is approved. Shipping a fix means a new build (the project version in Xcode), not a rollback of the binary that is already on phones.

# Troubleshooting

Common problems and how to fix them, grouped by symptom.

---

# Noticoel exits immediately on startup

## `NOTICOEL_AUTH_TOKEN environment variable is required`

`NOTICOEL_AUTH_TOKEN` has no default and is always required. Set it — see [Configuration](../README.md#configuration).

## `NOTICOEL_TELEGRAM_BOT_TOKEN and NOTICOEL_TELEGRAM_CHAT_ID environment variables are required when Telegram is enabled`

Telegram is enabled by default (`NOTICOEL_TELEGRAM_ENABLED=true`). Either set both variables (see [Notifiers → Telegram](../README.md#telegram)), or set `NOTICOEL_TELEGRAM_ENABLED=false` if you don't want Telegram notifications.

---

# Events are accepted but no Telegram message arrives

Check the logs for a line like:

```
"msg":"notifier failed","notifier":"telegram","message":"404 Not Found"
```

The request to the events API still succeeds (`202 Accepted`) — delivery is best-effort and doesn't fail the request; only the log line tells you it failed. Common causes:

- **`404 Not Found`** — `NOTICOEL_TELEGRAM_BOT_TOKEN` is wrong, or the bot was deleted.
- **`403 Forbidden`** — the bot was removed from the chat/group, or never started (send it a direct message first).
- **Nothing logged, nothing received** — double-check `NOTICOEL_TELEGRAM_CHAT_ID`: group IDs are negative. Get it from `https://api.telegram.org/bot<TOKEN>/getUpdates` (see [Notifiers → Telegram](../README.md#telegram)).

---

# `unable to open database file` / SQLite can't write, in Docker

The `noticoel` image runs as a fixed non-root UID (`65532`) and has no Dockerfile (built via [Ko](https://ko.build)) to `chown` anything at build time. When you bind-mount a host directory (e.g. `./data:/data`), it keeps the host's ownership — usually not `65532` — so the process can't create or write the SQLite file.

Fix: `chown` the bind-mounted directory before Noticoel starts, e.g. with a one-shot init container:

```yaml
services:
  noticoel-init:
    image: busybox:latest
    command: chown -R 65532:65532 /data
    volumes:
      - ./data:/data

  noticoel:
    # ...
    depends_on:
      noticoel-init:
        condition: service_completed_successfully
    volumes:
      - ./data:/data
```

See [Configuration → Docker Compose](../README.md#docker-compose) for the full example. This only matters for bind mounts — a named Docker volume doesn't have this problem.

---

# `NOTICOEL_DATABASE_PATH` (or a mounted `config.yaml`) doesn't seem to be found

The `noticoel` image has no `WORKDIR`. A relative path (the default, `./data/noticoel.db`) resolves against `/` — the container's root, not writable by the non-root user — not against wherever you might expect (e.g. `/app`).

Fix: use an absolute path, matched to your volume mount, e.g. `NOTICOEL_DATABASE_PATH=/data/noticoel.db` with `- ./data:/data`. Same reasoning applies to a mounted `config.yaml`: mount it at `/config/config.yaml`, not `/app/config/config.yaml`.

---

# A `config.yaml` value doesn't seem to apply

Environment variables always win over the legacy YAML file, which always wins over the hardcoded default (see [Configuration](../README.md#configuration)). If `NOTICOEL_SERVER_PORT` (or any other variable) is set — including an empty default injected by your process supervisor or Compose file — it overrides whatever is in `config.yaml`. Unset the environment variable, or update it instead of the YAML file.

---

# An adapter route returns `400 Bad Request`

Each adapter (`POST /api/v1/adapters/{name}`) validates the fields it needs from the third-party payload before converting it — the response body has the specific reason (e.g. a missing field, or the wrong `object_kind` for GitLab). Common cause: the webhook is configured to send an event type the adapter doesn't handle yet (each adapter targets one representative event — see [Architecture → Adapters](architecture.md#adapters)). For GitLab specifically, make sure only "Pipeline events" is enabled for that webhook URL.

# An adapter or the Event API returns `401 Unauthorized`

Both require the same bearer token as everything else (`Authorization: Bearer $NOTICOEL_AUTH_TOKEN`). Not every webhook provider's UI lets you attach a custom header — GitHub and older GitLab versions, for example, only offer a signature secret, not an arbitrary header. Check what your provider supports (see [Adapters](../README.md#adapters)).

---

# How to check Noticoel is actually running

```bash
curl http://localhost:8080/health   # -> "OK"
curl http://localhost:8080/version  # -> the running version
```

Neither requires a bearer token.

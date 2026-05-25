# dex

**dex** is a self-hosted pastebin, file hoster, and URL shortener.

## Installation

### Client

```bash
go install crdx.org/dex/cmd/dex@latest
```

### Server

```bash
go install crdx.org/dex/cmd/dexd@latest
```

## CLI

```
Usage:
    dex [options] (paste | upload) [<paths>...] [--label NAME] [--uuid] [--force] [--type NAME] [--bare]
    dex [options] redir <url> [--label NAME] [--force]
    dex [options] cat <ref> [--force]
    dex [options] rm <refs>...
    dex [options] ls [--refs]
    dex [options] (cp | mv) <from> <to>
    dex [options] set <ref> (type | kind) <value>
    dex [options] expose <ref> [--expiry DURATION]
    dex [options] urls [rm <ids>...]
    dex [options] deploys [--all]
    dex [options] (conf [-e] | e <label> | tail | gc)

Options:
    --refs              List refs
    -u, --uuid          Use UUID instead of filenames or numbers
    -t, --type TYPE     Content type e.g. application/json
    -f, --force         Use the force
    -e, --edit          Open in $EDITOR
    -l, --label NAME    Short name for this item
    -x, --expiry DUR    Token expiry duration e.g. 1d, 12h, 30m
    -b, --bare          Create empty item (requires --label)
    -a, --all           Show all results
    -v, --verbose       Be verbose
```

## Examples

Upload a paste from stdin.

```bash
echo 'Hello, world!' | dex paste
```

Upload a file.

```bash
dex upload image.png
```

Create a URL redirect with a custom label.

```bash
dex redir https://example.com --label ex
```

Upload JSON from a pipeline with an explicit content type.

```bash
curl -s https://api.example.com/data | dex paste --label data --type application/json
```

Bulk upload files, using filenames as labels.

```bash
dex upload *.png
```

## Remote deployment

Items can be updated externally via deploy URLs. This is useful for deploying content from CI pipelines or to give your nearest LLM the ability to deploy to a single URL, securely.

Create a placeholder item and generate a deploy URL:

```bash
dex paste --bare --label dashboard
dex expose dashboard --expiry 30d
```

Or use the `e` shortcut to create a placeholder and expose it with a 1-day expiry in one command:

```bash
dex e dashboard
```

```
Public URL: https://d.example.com/dashboard
Deploy URL: https://d.example.com/deploy/abc123
Expires at: 3 Jan 2025 12:00
```

The deploy URL can then be used to deploy new content:

```bash
curl https://d.example.com/deploy/abc123 \
    -F "content=@index.html" \
    -F "change=Update metrics for December" \
    -F "deployer=Alice"
```

Manage deploy URLs with `dex urls` and view deployment history with `dex deploys`.

## Concepts

### Pastes vs uploads

Both `paste` and `upload` store content, but they differ in how the content is served:

- **Pastes** are displayed inline in the browser with an appropriate content type.
- **Uploads** are served with `Content-Disposition: attachment`, triggering a download.

### Labels

Items are identified by a label or a UUID. If no label is provided and `--uuid` is not set, an auto-incrementing numeric label is assigned (e.g. `1`, `2`, `3`).

### Placeholders

The `--bare` flag creates an empty item. This is useful for reserving a label for later use with remote deployment.

### Content-addressable storage

Blobs are stored by their SHA1 hash. Uploading the same content twice will not store duplicates. The `gc` command removes orphaned blobs that are no longer referenced by any item.

## Configuration

The client configuration file is stored at `$XDG_CONFIG_HOME/dex/config.json` (or `~/.config/dex/config.json`).

```json
{
    "base_url": "https://d.example.com",
    "api_key": "..."
}
```

Use `dex conf --edit` to open this file in `$EDITOR`.

## Server

The server requires a `.env` file and a MySQL or MariaDB database.

```bash
dexd --env .env
```

If no `--env` is specified, the server will look for a `.env` file in the current directory.

## Server configuration

dexd is configured with a `.env` file.

The first value (in **bold**) is the default.

| Variable            | Required | Values                                          | Description                                                                                                                                |
|---------------------|----------|-------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------|
| `MODE`              | no       | **`development`**, `production`                 | The execution mode. In production mode request logging is disabled.                                                                        |
| `HOST`              | no       | **`127.0.0.1`**, `0.0.0.0`                      | The interface to bind to.                                                                                                                  |
| `PORT`              | no       | **`3000`**                                      | The port to listen on.                                                                                                                     |
| `BASE_URL`          | yes      | `https://d.example.com`                         | The public base URL for generated links.                                                                                                   |
| `API_KEY`           | yes      | `hunter2`                                       | The API key for authenticating client requests. Keys are compared using a constant-time SHA256 comparison.                                 |
| `DATABASE_NAME`     | no       | **`dex`**                                       | The database name.                                                                                                                         |
| `DATABASE_USERNAME` | no       | **`dex`**                                       | The database username.                                                                                                                     |
| `DATABASE_PASSWORD` | no       | `hunter2`                                       | The database password.                                                                                                                     |
| `DATABASE_PROTOCOL` | no       | **`tcp`**, `unix`                               | The database connection protocol.                                                                                                          |
| `DATABASE_ADDRESS`  | no       | **`127.0.0.1:3306`**, `/run/mysqld/mysqld.sock` | Host and port, or a unix socket path.                                                                                                      |
| `TRUSTED_PROXIES`   | no       | `172.18.0.0/16`, `private`                      | Comma-separated IP addresses or CIDR ranges, or `private` to trust all private ranges. Ensures the correct client IP is displayed in logs. |

## Contributions

Open an [issue](https://github.com/crdx/dex/issues) or send a [pull request](https://github.com/crdx/dex/pulls).

## Licence

[GPLv3](LICENCE).

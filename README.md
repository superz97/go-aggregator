# Gator

Gator is a CLI RSS feed aggregator — follow feeds, scrape them on a schedule, and browse saved posts from the terminal.

## Prerequisites

- [Go](https://go.dev/doc/install) (1.26+)
- [PostgreSQL](https://www.postgresql.org/download/) running locally or reachable over the network

## Installation

Install the `gator` CLI with `go install`:

```bash
go install github.com/superz97/go-aggregator/cmd/gator@latest
```

This places a `gator` binary in your `$GOPATH/bin` (or `$HOME/go/bin` by default) — make sure that's on your `$PATH`.

## Configuration

Gator reads its config from `~/.gatorconfig.json`. Create it with your Postgres connection string:

```json
{
  "db_url": "postgres://user:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Run the database migrations with [goose](https://github.com/pressly/goose):

```bash
goose -dir sql/schema postgres "$DB_URL" up
```

## Usage

Register a user and log in:

```bash
gator register alice
gator login alice
```

Add a feed (this also auto-follows it):

```bash
gator addfeed "Hacker News" https://news.ycombinator.com/rss
```

Start the aggregator loop (runs forever, scrapes on an interval — Ctrl+C to stop):

```bash
gator agg 1m
```

Browse saved posts from feeds you follow. By default this shows the 2 most recent posts, but you can sort, filter, and page through results with flags:

```bash
gator browse
gator browse --limit=10
gator browse --sort=title --order=asc
gator browse --feed="Hacker News" --limit=5
```

- `--limit=<n>` — number of posts to show (default `2`)
- `--sort=published_at|title` — field to sort by (default `published_at`)
- `--order=asc|desc` — sort direction (default `desc`)
- `--feed=<name>` — only show posts from the given followed feed, matched by exact name

Other commands: `users`, `feeds`, `follow <url>`, `following`, `unfollow <url>`.

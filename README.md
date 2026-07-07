# Gator

Gator is a CLI RSS feed aggregator — follow feeds, scrape them on a schedule, and browse saved posts from the terminal.

## Prerequisites

- [Go](https://go.dev/doc/install) (1.26+)
- [PostgreSQL](https://www.postgresql.org/download/) running locally or reachable over the network, with the `pg_trgm` extension available (ships with standard Postgres/`postgresql-contrib`; used for fuzzy `search`)

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

Start the aggregator loop (runs forever, scrapes on an interval — Ctrl+C to stop). It scrapes immediately on startup, then again every interval. An optional second argument controls how many feeds are fetched concurrently per tick (default `1`):

```bash
gator agg 1m
gator agg 1m 5
```

Browse saved posts from feeds you follow. By default this shows the 2 most recent posts, but you can sort, filter, and page through results with flags:

```bash
gator browse
gator browse --limit=10
gator browse --sort=title --order=asc
gator browse --feed="Hacker News" --limit=5
gator browse --limit=10 --page=2
```

- `--limit=<n>` — number of posts to show (default `2`)
- `--page=<n>` — 1-indexed page of results to show (default `1`)
- `--sort=published_at|title` — field to sort by (default `published_at`)
- `--order=asc|desc` — sort direction (default `desc`)
- `--feed=<name>` — only show posts from the given followed feed, matched by exact name

Fuzzy search your saved posts by title (typo-tolerant, via Postgres `pg_trgm`):

```bash
gator search "murderd"
gator search --limit=5 "counterclok"
```

- `--limit=<n>` — max number of results (default `10`)
- Flags and the query can appear in any order; an unquoted multi-word query is joined with spaces.

Bookmark or like a post, identified by its URL (shown in `browse`/`search` output):

```bash
gator bookmark https://audiochuck.com
gator bookmarks
gator unbookmark https://audiochuck.com

gator like https://audiochuck.com
gator likes
gator unlike https://audiochuck.com
```

Browse posts interactively in a terminal UI — scroll the list, view a post's full detail (feed name, date, description), and open it in your browser:

```bash
gator tui
```

- `j`/`k` or arrow keys — move the selection
- `enter` — view full post detail
- `o` — open the post URL in your default browser
- `esc`/`q` — back to the list / quit

Other commands: `users`, `feeds`, `follow <url>`, `following`, `unfollow <url>`.

# Gator

Gator is a command-line RSS feed aggregator written in Go. It allows users to register, follow RSS feeds, fetch posts, and browse recent posts from their followed feeds.

## Requirements

Before installing Gator, make sure you have the following installed:

- **Go** 1.22 or later
- **PostgreSQL** (a running database server)

## Installation

Install the Gator CLI using `go install`:

```bash
go install github.com/YOUR_GITHUB_USERNAME/gator@latest
```

Make sure your Go binary directory is included in your `PATH`. By default, this is usually `$HOME/go/bin`.

Verify the installation:

```bash
gator
```

## Configuration

Gator reads its configuration from `~/.gatorconfig.json`.

Create the file with the following structure:

```json
{
  "db_url": "postgres://postgres:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Replace the database URL with your PostgreSQL connection details. Create the `gator` database if it does not already exist:

```bash
createdb gator
```

## Database setup

Apply the database migrations before running Gator. If the repository contains Goose migrations, install Goose:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Then run the migrations from the directory containing your migration files:

```bash
goose postgres "postgres://postgres:password@localhost:5432/gator?sslmode=disable" up
```

Adjust the connection string and migration directory to match your local setup.

## Usage

Gator provides commands for managing users, following RSS feeds, and browsing posts.

### Register a user

```bash
gator register mario
```

### Log in as a user

```bash
gator login mario
```

### Add an RSS feed

```bash
gator addfeed https://feeds.example.com/rss.xml
```

### List available feeds

```bash
gator feeds
```

### Follow a feed

```bash
gator follow https://feeds.example.com/rss.xml
```

### List followed feeds

```bash
gator following
```

### Fetch posts continuously

```bash
gator agg 1m
```

The aggregator periodically fetches RSS feeds. Use an appropriate duration for the aggregation interval.

### Browse recent posts

```bash
gator browse
```

By default, this displays the two most recent posts available from your followed feeds.

Specify a custom limit:

```bash
gator browse 10
```

### List users

```bash
gator users
```

## Development

Clone the repository:

```bash
git clone https://github.com/YOUR_GITHUB_USERNAME/gator.git
cd gator
```

Generate database code if required:

```bash
sqlc generate
```

Run the program locally:

```bash
go run .
```

## License

Add your chosen license here.

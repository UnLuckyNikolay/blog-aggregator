# Blog Aggregator

WIP, Readme is used as notepad for now

## Install and Run

1. Go, clone repo

2. Install Postgres v15 or later:

```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
```

To check current version:

```bash
psql --version
```

2.1. (Linux only) Update postgres password if needed:

```bash
sudo passwd postgres
```

2.2. Start Postgres server and enter the shell:

```bash
sudo service postgresql start
sudo -u postgres psql
```

3. Install Goose:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

3.1. Build database. In `sql/schema/` run:

```bash
goose postgres *database_url* up
```

?. ???

Last. Profit
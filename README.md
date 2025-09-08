# Blog Aggregator

Tool for keeping track of RSS Feeds. Allows multiple users to follow multiple feeds. Uses PostgreSQL to keep track of posts.

## Installation

1. Install [Go](https://go.dev/doc/install) 1.24.5 or higher.

2. Install Postgres v15 or later:

    ```bash
    sudo apt update
    sudo apt install postgresql postgresql-contrib
    ```

3. Start Postgres server:

    ```bash
    sudo service postgresql start
    ```

4. Get the database url. Structure: `postgres://<username>:<password>@localhost:5432/gator?sslmode=disable`. 

    Save it as an environmental variable if you wish (for later usage while building the database, not used by the app itself):

    ```bash
    export GATOR_DB_URL=<db_url>
    ```

5. Create a `.gatorconfig.json` file in yout home directory and add your url to it:

    ```json
    {"db_url":"<db_url>"}
    ```

6. Clone the repo:

    ```bash
    git clone https://github.com/UnLuckyNikolay/blog-aggregator
    cd blog-aggregator
    ```

7. Install the app:

    ```bash
    cd cmd/gator
    go install
    cd ../..
    ```

8. Build the database. Go into the `schema` folder:

    ```bash
    cd sql/schema
    ```

    and run this command for each file in order of the numbers:

    ```bash
    psql <db_url> -f <filename>
    ```

    or with an environmental variable:

    ```bash
    psql $GATOR_DB_URL -f <filename>
    ```

## Usage

Use `gator` command to interact with the app:

```bash
gator <command> <arguments>
```

To start fetching feeds run `gator agg` in a new terminal window.

Commands:
* `register <username>` - register a new user
* `login <username>` - login as the user
* `users` - list all the users
* `agg` - start fetching feeds, runs **continuously**
* `feeds` - list all the feeds
* `addfeed <feed_name> <feed_url>` - add a new feed
* `follow <feed_url>` - follow the feed
* `unfollow <feed_url>` - unfollow the feed
* `following` - list all the followed feed for the current used
* `browse <amount_of_posts>` - list N of the last posts from the followed feeds
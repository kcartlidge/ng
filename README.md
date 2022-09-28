# Near Gothic

Generate strongly-typed Go database access code directly from your Postgres database schema.

- [MIT license](./LICENSE)
- [CHANGELOG](./CHANGELOG.md)

## Prerequisites

You'll need a PostgreSQL database and valid connection details.
The connection string should be an environment value in *pgx* format.
The default environment variable name is `DB_CONNSTR`.

``` sh
export DB_CONNSTR="host=127.0.0.1 port=5432 dbname=example user=example password=example sslmode=disable"
```

*(This is an example, not a revealed secret.)*

[Here's a Postgres SQL script for a (small) example database](./postgres.sql).

## Running

There are pre-built cross-platform binaries in the `src/builds` folder.
When run they will display the command requirements (as shown below).

```
USAGE
  ng -w -folder <value> -module <value> -env <value> -schema <value>

ARGUMENTS
  -w                  overwrite any existing destination folder
  -env <value>        connection string environment variable (default `DB_CONNSTR`)
  -schema <value>     the Postgres database schema to scan (default `public`)
  -folder <value>  *  destination folder, either relative or absolute
  -module <value>  *  the Go module for created code (eg `kcartlidge/app/data`)

  * means the argument is required

EXAMPLE
  ng -w -folder ~/example/repo -module kcartlidge/app/data -env DB_CONNSTR -schema example
```

## What Near Gothic does

- If overwrite is specified (`-w`) it replaces any existing file(s)
- It uses the specified environment variable (`-env`) to connect to the database
- It scans the specified schema (`-schema`)
- It creates repository code in the output location (`-folder`)
- The Go code will use the specified module path (`-module`)

The general idea is that you'll have an API or app already created.
For the example above that would be in the `~/example` folder.
Its namespace would be `kcartlidge/app` (the module passed in, minus the end bit).

Near Gothic will scan the database and create a `repo` folder inside `~/example`.
The created repo will use the module path `kcartlidge/app/data`.

## Building cross-platform binaries

The [`src/scripts`](./src/scripts) folder has scripts to be run *on* Linux, Mac, and Windows.
Use the one for the system *you* are currently using.

Each of those scripts will produce builds for the three platforms at once.
When built they will automatically be placed in the expected `src/builds` folder.

``` shell
cd src
./scripts/macos.sh
```

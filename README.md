# Near Gothic

Generate strongly-typed Go database access code directly from your Postgres database schema.
Automatic API generation is currently in progress.

## STATUS

BETA - connections, entities, and repos are usable; API generation is ongoing

It scans the database and generates code/scripts as detailed below.
For this use case, it is stable and considered *beta* but unreleased.
The automatic generation of an API is *in progress* (see the latest commits).
It's currently an almost empty server, expected to be completed this coming week.

- [MIT license](./LICENSE)
- [CHANGELOG](./CHANGELOG.md)

# Contents

- [Prerequisites](#prerequisites)
  - [Expectations](#expectations)
- [Running](#running)
- [What Near Gothic does](#what-near-gothic-does)
- [What gets created](#what-gets-created)
  - [Sample generated folder structure](#sample-generated-folder-structure)
- [Building cross-platform binaries](#building-cross-platform-binaries)

## Prerequisites

You'll need a PostgreSQL database and valid connection details.
The connection string should be an environment value in *pgx* format.
The default environment variable name is `DB_CONNSTR`.

``` sh
export DB_CONNSTR="host=127.0.0.1 port=5432 dbname=example user=example password=example sslmode=disable"
```

*(This is an example, not a revealed secret.)*

[Here's a Postgres SQL script for a (small) example database](./postgres.sql).

Note that upon successfully running, the generated Go code will include its own `postgres.sql` file.
That file will contain scripts suitable for recreating the database entities found/used when it ran.

### Expectations

Your database *must* follow certain conventions.
This may change in the future, but currently the following applies:

- TABLES
  - Must have a primary key
    - *Only one column is allowed (no composite keys)*
      - This is not currently enforced but will be shortly
      - Their use risks data corruption
  - Names must be snaked-lowercase
    - For example `another_thing`, not `AnotherThing` or `Another_Thing`
  - May have comments, which are used when generating the API/entities
- COLUMNS
  - Names must be snaked-lowercase
    - For example `account_id`, not `AccountID` or `accountid`
  - Most common database column types are supported
    - A full list will be made available once the first release is issued
  - May have comments, which are used when generating the API/entities

## Running

There are pre-built cross-platform binaries in the `src/builds` folder.
When run they will display the command requirements (as shown below).

```
USAGE
  ng -w -env <value> -schema <value> -module <value> -folder <value> -repo <value>

ARGUMENTS
  -w                  overwrite any existing destination folder?
  -env <value>        connection string environment variable (default `DB_CONNSTR`)
  -schema <value>     the Postgres database schema to scan (default `public`)
  -module <value>  *  the Go module for created code (eg `kcartlidge/app/data`)
  -folder <value>  *  destination folder, either relative or absolute
  -repo <value>    *  repository subfolder name

  * means the argument is required

EXAMPLE
  ng -w -env DB_CONNSTR -schema example -module kcartlidge/api -folder ~/example -repo repo
```

## What Near Gothic does

- If overwrite is specified (`-w`) it replaces any existing file(s)
- It uses the specified environment variable (`-env`) to connect to the database
- It scans the specified schema (`-schema`)
- The Go code for the API will use the specified module path (`-module`)
- It creates API code in the output location (`-folder`)
- The `-repo` is appended to the `-folder` to specify where the repo is created
- The `-repo` is appended to the `-module` to form the module path for the repo

## What gets created

You get a folder structure with the following:

- Go module (with `go.mod` and `go.sum`)
- Basic JSON API with middleware etc
  - Currently just the `GET *` endpoints
  - Includes server, handler, routing, and documentation
- JSON dump file detailing what was scanned from the database
  - Useful for your own further automated processing
- A set of entities, one per database table
  - SQL comments implemented as Go comments
  - Generated property comments
    - Max length, primary key flag, sortable/filterable
  - Extra constructors
    - Construct from on a *pgx* row
    - Construct from an HTTP POST
  - Column attributes for JSON, SQL, display, and slug
  - Validation based on SQL column length
- A connection package
- A package of strongly-typed repositories
- A `README.md` detailing what the repo contains
- An emergency SQL script to recreate the entities
  - Comments, constraints, keys, defaults, and more

There's [an example of the generated code](./_example) in the `_example` folder.
And here's what it looks like when a generated API is run:

```
API v1.0.0

GET    /api/accounts
GET    /api/account-settings
GET    /api/settings

Server running at http://localhost:3000
```

The URL prefix (`/api`), port (`3000`), and version (`v1.0.0`) are fixed.
To change them update the generated `main.go` file.
Options to specify them during API generation will be provided in due course.

### Sample generated folder structure

The [`_example`](./_example) folder contains some generated code.
The `README.md` file in there shows the command used to generate it.
Here's a high-level breakdown of what the files/folders contain.

```
/_example                      // target folder
  /repo                        // generated content root
    /connection
      connection.go            // class for db connection
    /entities
      account-setting.go       // the 'account_setting' db table
      account.go               // the 'account' db table
      setting.go               // the 'setting' db table
    /support
      support.go               // support functions
    account-repo.go            // the 'account' repository
    account-setting-repo.go    // the 'account-setting' repository
    repo-base.go               // shared repository functionality
    setting-repo.go            // the 'setting' repository
  dump.json                    // JSON dump of the schema
  go.mod
  go.sum
  handlers.go                  // handlers and their registration
  main.go                      // starting point of the stub API server
  middleware.go                // some API middleware
  postgres.sql                 // SQL to recreate the entities
  README.md                    // Overview of the generated code
  server.go                    // the server functionality
```

Note that the generated `README.md` has a table in Markdown format.
This is currently not displaying in the GoLand Markdown previewer.
It does however display correct in Visual Studio Code and elsewhere.

## Building cross-platform binaries

The [`src/scripts`](./src/scripts) folder has scripts to be run *on* Linux, Mac, and Windows.
Use the one for the system *you* are currently using

Each of those scripts will produce builds for the three platforms at once.
When built they will automatically be placed in the expected `src/builds` folder.

``` shell
cd src
./scripts/macos.sh
```

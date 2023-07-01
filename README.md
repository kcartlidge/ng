# Near Gothic

Generate strongly-typed Go database access code directly from your Postgres database schema.

A simpler version of tools like *SQL Boiler*.
The created code supports filtering and sorting, comments, indexes, length restrictions, nullability, and more - with minimal setup needed.

## STATUS

*Stable. Usable. Beta.*

- [AGPL license](./LICENSE)
- [CHANGELOG](./CHANGELOG.md)

## Performance

On a 10th Gen Core i5 running Windows 11 Pro, using PostgreSQL 15 with 3 tables, it generates (and formats) a code-base in around **600ms** when averaged over a dozen runs.

Go is fast. On most modern machines even with larger collections of tables it should complete within seconds.

# Contents

- [Prerequisites](#prerequisites)
  - [Expectations](#expectations)
- [Running](#running)
- [How Near Gothic works](#how-near-gothic-works)
  - [Generated folder structure](#generated-folder-structure)
  - [Example Code](#example-code)
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

Note that upon successfully running, the generated Go code will include its own `postgres.sql` file.  That file will contain scripts suitable for recreating the database entities found/used when it ran.

### Expectations

Your database *must* follow certain conventions.
This may change in the future, but currently the following applies:

- TABLES
  - Tables must have a primary key
    - *Only one column is allowed (no composite keys)*
      - This is not currently enforced but is expected
      - Ignoring this risks data corruption
  - Table names must be snaked-lowercase
    - For example `another_thing`, not `AnotherThing` or `Another_Thing`
  - Tables may have comments, which are incorporated into the generated entities
- COLUMNS
  - Column names must also be snaked-lowercase
    - For example `account_id`, not `AccountID` or `accountid`
  - Most common database column types are supported
    - A list will be included upon first full release
  - Columns may have comments, which are incorporated into the generated entities

## Running

There are pre-built cross-platform binaries in the [`src/builds`](./src/builds) folder.
When run they will display the command requirements with details and examples (as shown below).

```
USAGE
  ng.exe [-w] [-env <value>] [-schema <value>] -module <value> -folder <value> -repo <value>

ARGUMENTS
  -w                  overwrite any existing destination folder?
  -env <value>        connection string environment variable (default `DB_CONNSTR`)
  -schema <value>     the Postgres database schema to scan (default `public`)
  -module <value>  *  the Go module for *created* code (eg `kcartlidge/app/data`)
  -folder <value>  *  *parent* module folder, either relative or absolute
  -repo <value>    *  name of the subfolder for generated code (eg `data`)

  * means the argument is required

EXAMPLE
  ng -w -env DB_CONNSTR -schema example -module kcartlidge/app/data -folder ~/example -repo data
```

The created `README.md` file will include the command you used when generating the code.

## How Near Gothic works

- It uses the named environment variable (`-env`) to connect to the database
  - It then scans the provided PostgreSQL schema (`-schema`)
- It generates code *nested within* the `-folder`
  - The created nested subfolder is named according to`-repo`
  - Generated code is within a module named as per `-module`
  - If overwrite is specified (`-w`) it replaces any existing file(s)

You get a folder structure with the following:

- Well-commented code
- A set of entities, one per database table
  - SQL table comments implemented as Go comments
  - Generated property comments
    - Max length, primary key flag, sortable/filterable
  - Extra constructors
    - Construct from a *pgx* row
    - Construct from an HTTP POST
  - Column attributes for JSON, SQL, display, and slugs
  - Validation based on SQL column length
- A connection package
- A package of strongly-typed repositories
- A `README.md` detailing what the repo contains
- A `USING.md` detailing how to use the repo
- An emergency SQL script to recreate the entities
  - Comments, constraints, keys, defaults, and more
- JSON dump file detailing what was scanned from the database
  - Useful for your own further automated processing

### Generated folder structure

Here's a high-level breakdown of what the files/folders contain.

```
/app                           // target folder (parent module)
  go.mod                       // example parent module file
  go.sum                       // example parent module file
  main.go                      // example parent module file

  /data                        // root of the generated content
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
    dump.json                  // JSON dump of the schema
    postgres.sql               // SQL to recreate the entities
    README.md                  // overview of the generated code
    repo-base.go               // shared repository functionality
    setting-repo.go            // the 'setting' repository
    USING.md                   // details of how to use the code
```

### Example Code

NearGothic does not create modules.
It generates a stand-alone package written into a subfolder of your app's module.
That package, when imported, provides strongly typed access to the database.

The example below is using a generated repo (`kcartlidge/app/data`) derived from a database containing an `account` table (amongst others).

It shows the `connection` being obtained by a connection string stored in an environment variable. It then uses that to create an instance of an `AccountRepo` via which it fetches the first 3 accounts sorted by the email address in reverse alphabetical order (for no obvious reason other than to show the somewhat-fluent query capability).

Repos are automatically created for each table found in your Postgres schema. Column types are mapped to Go types. SQL comments show as Go comments. Basic validation based on nullability and length are included. And utility methods for both filtering and sorting are added for each column that has an index (non-specific alternatives are also provided).

**Important note:** repo instances retain any filters/sorts between calls, allowing you to (for example) add a `UserId` restriction at the start of using a repo and be confident that restriction will apply to further operations.
For this same reason it is imperative that each scope creates it's own instance of any repos for use, as sharing instances can cause 'bleeding' of sorts/filters across operations leading to unexpected results.
(Repos are lightweight; the overhead is minimal and they can share a connection.)

``` go
package main

import (
	"fmt"
	"kcartlidge/app/data"
	"kcartlidge/app/data/connection"
	"log"
	"os"
)

func main() {
	// Limit queries to 100 rows for performance/scalability.
	// This can be overridden on individual calls.
	connection.MaxRows = 100

	// Connect using the connection string in the "DB_CONNSTR" environment variable.
	envName := "DB_CONNSTR"
	connectionString, hasEnv := os.LookupEnv(envName)
	if !hasEnv {
		log.Fatalf("Environment variable `%s` not found", envName)
	}
	conn := connection.NewConnection(connectionString)

	// Start a new repo and fetch the first 3 accounts in reverse email address order.
	// These lines will not build if your database tables differ (they probably do).
	ar := data.NewAccountRepo(conn)
	d, err := ar.WhereId("<", 4).ReverseByEmailAddress().List()

	// Show the result.
	fmt.Println()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("First few accounts:")
	fmt.Println(d)
	fmt.Println()

	// Deal with a specific account, using a new repo to reset the filter/sort.
	// Could also call ar.ResetConditions() and/or ar.ResetSorting() instead.
	ar = data.NewAccountRepo(conn)
	acc, err := ar.WhereEmailAddress("=", "email@example.com").List()
	if err == nil && len(acc) > 0 {
		fmt.Printf("`%s` is for account %v", acc[0].EmailAddress, acc[0].Id)
	}
	fmt.Println()
}
```

If you want to see the database operations, including the SQL that was generated, you can switch on debugging output:

``` go
connection.DebugMode = true
```

## Building cross-platform binaries

*Skip this if you are just intending to use Near Gothic rather than contribute to it.*

The [`src/scripts`](./src/scripts) folder has scripts to be run *on* Linux, Mac, and Windows. Use the one for the system *you* are currently running.
Each of those scripts will produce builds for the three supported platforms at once, and automatically place them in the expected [`src/builds`](./src/builds) folder.

Linux/Mac:

``` shell
cd src
./scripts/macos.sh
./scripts/linux.sh
```

Windows:

``` shell
cd src
scripts\windows.bat
```

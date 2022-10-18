# kcartlidge/api

## Contents

- [Running](#running)
- [Regenerating](#regenerating)
- [Entities](#entities)
- [Repository](#repository)
- [SQL Scripts](#sql-scripts)
- [Example Usage](#example-usage)

## Running

The API's main folder is the one containing this document.
Navigate to here in your command prompt, build, and run.

``` sh
cd <folder>
go build -o api
./api
```

If you are using Windows remove `./` and just run `api`.

## Regenerating

The code was generated using NearGothic via the following command:

``` sh
cd /Users/karl/source/go/ng/src
ng -w -schema example -module kcartlidge/api -folder ../_example -repo repo
```

- Please ensure the folders used if repeating these commands match your own system
- Existing code should also be copied or committed to source control first
  - This is because existing contents of the `-folder` are deleted when running
  - For extra safety ensure your initial run *does not* include `-w`

## Entities

| Struct | Table | Display | JSON | Slug |
| --- | --- | --- | --- | --- |
| [`Account`](./repo/entities/account.go) | account | *Account* | account | account |
| [`AccountSetting`](./repo/entities/account-setting.go) | account_setting | *Account Setting* | accountSetting | account-setting |
| [`Setting`](./repo/entities/setting.go) | setting | *Setting* | setting | setting |

Each entity also has methods to:
- convert a database row into an instance of the entity
- create an instance with the content of HTTP POST form variables
- perform validation checks against field lengths

## Repository

- Each entity has a repository *struct*
- They are named according to a pattern, e.g. `CustomerRepo`
- They also have a constructor, e.g. `NewCustomerRepo()`
- They have CRUD methods for `List`, `Insert`, `Update`, and `Delete`
- They have general purpose methods for maximum rows and/or paging
  - `WithLimit` adds a restriction on the number of items returned
      - Overrides the package's `MaxRows` value (for this instance only)
  - `WithOffset` enables skipping the given number of items in the result set
- They have Methods for indexed fields
  - Each indexed field gets its own set of filters/sorting
    - Multiple filters and sorts can be applied at once
    - Strongly-typed filter per field (e.g. `WhereEntryCount`)
    - Strongly-typed sorting per field
      - Ascending, e.g. `SortByEntryCount()`
      - Descending, e.g. `ReverseByEntryCount()`
  - General purpose sorting and filtering (only intended for *unindexed* column usage)
    - `Where` adds a clause to the request
    - `AddSorting` adds an ad-hoc sort by any valid column/thing

## SQL Scripts

A PostgreSQL script containing SQL (to generate the entities in this repo) is also included.

This is an emergency-use script and *does not* obviate the need for proper precautions/backups.
In particular, it only contains the information extracted to generate the entities/repos.
Any other information outside of that usage will not have been persisted.

[The script is in the `postgres.sql` file](./postgres.sql).

## Example Usage

For illustrative purposes only; does not use this code-base.

``` go
package main

import (
  "fmt"
  "os"

  "example/data/connection"
  "example/data/entities"
  "example/data/repos"
)

func main() {
  // Sample entity creation.
  customer := entities.NewCustomer()

  // Set a non-nullable field.
  customer.EmailAddress = "email@example.com"

  // Set a nullable field.
  // As it is nullable, use '&'.
  userInput := "something-that-is-too-long-for-a-16-character-field"
  customer.VerificationCode = &userInput

  // Validate.
  // This will show a collection of issues, with one entry:
  //   ["Verification Code cannot be longer than 16."]
  fmt.Println("Issues :", customer.Validate())

  // Get the connection string details.
  connStr := os.Getenv("DB_CONNSTR")
  if len(connStr) == 0 {
    fmt.Println("Need a DB_CONNSTR environment variable connection string!")
    os.Exit(1)
  }

  // Open a connection.
  // The generated code calls log.Fatal with a suitable message if this errors.
  conn := connection.NewConnection(connStr)

  // Sample repository usage.
  customerRepo := repos.NewCustomerRepo(conn)
  customers, err := customerRepo.List()
  if err != nil {
    fmt.Println("Unable to fetch the customer list")
  } else {
    fmt.Println(len(customers))
  }
}
```

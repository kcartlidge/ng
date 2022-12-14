# kcartlidge/api v1.0.0

```
Code generated by Near Gothic. DO NOT EDIT.
Generated on 2022/10/19.
Manual edits may be lost when next regenerated.

Near Gothic is (C) K Cartlidge, 2022.
No rights are asserted over generated code.
```

## Important Note About Security

- This is a naive administration API
- *Every* database table for which code was generated is supported
- **Do not expose this API publicly and unsecured as all data will be accessible!**
- I recommend that even for internal use you consider HTTPS
  - *LetsEncrypt* have a Go package to accomplish this in code

Possible options for securing the API (best to worst) include:

- Implement your organisation's existing standard methodology
- Include authentication/authorisation from a suitable professional third party
- Apply your own authentication/authorisation based on JWTs or similar
- Restrict access to your organisational network only (beware of internal staff)

It is strongly advised you follow best practices and don't create your own code.
Avoid custom implementations unless you are an expert in the field.

## Contents

- [Running](#running)
- [Regenerating](#regenerating)
- [Generated API](#generated-api)
- [Entities](#entities)
- [Repository](#repository)
- [SQL Scripts](#sql-scripts)

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

- If repeating these commands please ensure the folders used match your own system
- Existing code should also be copied or committed to source control first
  - This is because existing contents of the `-folder` are deleted when running
  - For extra safety ensure your initial run *does not* include `-w`

## Generated API

| Endpoint | Details |
| --- | --- |
| GET api/accounts | List all accounts |
| GET api/account-settings | List all account settings |
| GET api/settings | List all settings |


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

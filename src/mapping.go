package main

import (
	"fmt"
	pluralize "github.com/gertd/go-pluralize"
	"strings"
)

var plural = pluralize.NewClient()

func mapPostgresTypeToGo(postgresDataType string) string {
	switch strings.ToLower(postgresDataType) {
	case "smallint", "smallserial":
		return "int16"
	case "integer", "serial":
		return "int"
	case "bigint", "bigserial":
		return "int64"
	case "decimal", "numeric", "money":
		return "float64"
	case "real":
		return "float64"
	case "double precision":
		return "float64"
	case "bytea":
		return "[]byte"
	case "character varying", "varchar", "character", "char", "text":
		return "string"
	case "boolean":
		return "bool"
	case "bit":
		panic("Unsupported column type 'bit' - use 'boolean' instead.")
	case "timestamp", "timestamptz", "timestamp with time zone", "timestamp without time zone", "date", "time", "time with time zone", "time without time zone":
		return "*time.Time"
	case "interval":
		return "*time.Duration"
	case "uuid":
		return "Guid"
	case "json", "jsonb":
		return "string"
	case "xml":
		return "string"
	}
	panic("Unsupported Postgres data type: " + postgresDataType)
}

func isPostgresTypeCardinal(postgresDataType string) bool {
	switch strings.ToLower(postgresDataType) {
	case "smallint", "smallserial", "integer", "serial", "bigint", "bigserial":
		return true
	}
	return false
}

func mapPostgresTypeToGoWithNullable(postgresDataType string, isNullable bool) string {
	goType := mapPostgresTypeToGo(postgresDataType)
	if isNullable && !strings.HasPrefix(goType, "*") {
		return "*" + goType
	}
	return goType
}

func toJsonName(value string) string {
	if len(value) == 0 {
		return ""
	}
	s := toProper(value, false)
	return strings.ToLower(s)[:1] + s[1:]
}

func toSlug(value string) string {
	if len(value) == 0 {
		return ""
	}
	s := strings.ToLower(toProper(value, true))
	return strings.ReplaceAll(s, " ", "-")
}

func toProper(value string, forDisplay bool) string {
	if len(value) == 0 {
		return ""
	}
	s := ""
	needsUpper := true

	for _, b := range strings.ToLower(value) {
		ch := string(b)
		if ch == "_" {
			needsUpper = true
		} else {
			if needsUpper {
				needsUpper = false
				if forDisplay {
					s += " "
				}
				s += strings.ToUpper(ch)
			} else {
				s += ch
			}
		}
	}

	s = strings.TrimSpace(s)
	if strings.ToLower(s) == "id" {
		s = "Id"
	}
	return s
}

// toPlural returns a pluralised version of the given text
func toPlural(value string) string {
	return plural.Plural(value)
}

// toColumnNameListCSV returns the database column names comma-delimited
func toColumnNameListCSV(table Table) string {
	s := ""
	for _, col := range table.Columns {
		if len(s) > 0 {
			s += ","
		}
		s += col.ColumnName
	}
	return s
}

// toColumnNameListNoPrimaryKeysCSV returns the database column names comma-delimited
// Primary keys are omitted
func toColumnNameListNoPrimaryKeysCSV(table Table) string {
	s := ""
	for _, col := range table.Columns {
		if col.IsPrimaryKey {
			continue
		}
		if len(s) > 0 {
			s += ","
		}
		s += col.ColumnName
	}
	return s
}

// toPrimaryKeyParametersCSV returns the primary key fields as comma-delimited parameters
func toPrimaryKeyParametersCSV(table Table) string {
	s := ""
	for _, col := range table.Columns {
		if col.IsPrimaryKey {
			if len(s) > 0 {
				s += ","
			}
			s += fmt.Sprintf("%s %s", col.ColumnName, col.DataType)
		}
	}
	return s
}

// toCodeNameListCSV returns the code column names comma-delimited
// The prefix allows the columns to be 'attached' to something
func toCodeNameListCSV(table Table, prefix string) string {
	s := ""
	for _, col := range table.Columns {
		if len(s) > 0 {
			s += ","
		}
		s += (prefix + col.CodeName)
	}
	return s
}

// toParameterListNoPrimaryKeysCSV returns comma-delimited '$n' parameters for SQL insert statements.
// Primary keys are omitted.
func toParameterListNoPrimaryKeysCSV(table Table) string {
	s := ""
	for i, col := range table.Columns {
		if col.IsPrimaryKey {
			continue
		}
		if len(s) > 0 {
			s += ","
		}
		s += fmt.Sprintf("$%v", i)
	}
	return s
}

// toUpdateListNoPrimaryKeysCSV returns comma-delimited field='$n' parameters for SQL update statements.
// Primary keys are omitted.
func toUpdateListNoPrimaryKeysCSV(table Table) string {
	s := ""
	for i, col := range table.Columns {
		if col.IsPrimaryKey {
			continue
		}
		if len(s) > 0 {
			s += ","
		}
		s += fmt.Sprintf("%s=$%v", col.ColumnName, i)
	}
	return s
}

// columnIdxAfterPrimaryKeys returns the number of non-primary key columns.
func columnIdxAfterPrimaryKeys(table Table) int {
	c := 0
	for _, col := range table.Columns {
		if col.IsPrimaryKey {
			continue
		}
		c++
	}
	return c + 1
}

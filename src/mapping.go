package main

import (
	"strings"
)

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

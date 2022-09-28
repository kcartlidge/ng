package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	pgx "github.com/jackc/pgx/v4/pgxpool"
)

var bg = context.Background()

type scanner struct {
	Schema           Schema
	SchemaName       string
	connectionString string
}

func NewScanner(connectionString string, schemaName string) scanner {
	s := scanner{
		Schema:           Schema{},
		SchemaName:       schemaName,
		connectionString: connectionString,
	}
	return s
}

func (s *scanner) ScanPostgresDatabase() error {
	db, err := pgx.Connect(bg, s.connectionString)
	check(err)
	defer db.Close()
	check(db.Ping(bg))
	fmt.Println("Connected to Postgres")
	s.Schema = Schema{
		SchemaName:  s.SchemaName,
		CodeName:    toProper(s.SchemaName, false),
		DisplayName: toProper(s.SchemaName, true),
		JsonName:    toJsonName(s.SchemaName),
		SlugName:    toSlug(s.SchemaName),
		Owner:       s.SchemaName,
		Tables:      []Table{},
	}
	fmt.Printf("Scanning schema `%s`\n", s.SchemaName)
	s.scanTablesAndViews(db)
	return nil
}

func (s *scanner) scanTablesAndViews(db *pgx.Pool) {
	statement := "SELECT table_name, table_type, is_insertable_into, " +
		"       pg_catalog.obj_description(pgc.oid, 'pg_class') as table_description " +
		"FROM   information_schema.tables, pg_catalog.pg_class pgc " +
		"WHERE  table_schema = $1 " +
		"AND    table_name = pgc.relname " +
		"AND    table_type='BASE TABLE' " +
		"ORDER  BY table_name;"
	rows, err := db.Query(bg, statement, s.SchemaName)
	check(err)
	defer rows.Close()
	for rows.Next() {
		tableName, tableType, canInsert, comment := "", "", "", sql.NullString{}
		check(rows.Scan(&tableName, &tableType, &canInsert, &comment))
		fmt.Printf("Scanning %s `%s`\n", strings.ToLower(tableType), tableName)
		table := Table{
			SchemaName:  s.SchemaName,
			TableName:   tableName,
			CodeName:    toProper(tableName, false),
			DisplayName: toProper(tableName, true),
			JsonName:    toJsonName(tableName),
			SlugName:    toSlug(tableName),
			Owner:       s.SchemaName,
			Comment:     strings.TrimSpace(comment.String),
			TableType:   tableType,
			IsUpdatable: strings.ToLower(canInsert) == "yes",
			Columns:     s.scanColumns(db, tableName, strings.ToUpper(tableType) == "VIEW"),
			Constraints: s.scanConstraints(db, tableName),
			Indexes:     []Index{},
			CodeImports: []string{},
		}
		table.Indexes = s.scanIndexes(db, table)
		needsTime := false
		for _, col := range table.Columns {
			if col.CanFilter {
				switch col.DataType {
				case "*time.Time":
					needsTime = true
				}
			}
		}
		if needsTime {
			s.addCodeImport(&table, "time")
		}
		s.Schema.Tables = append(s.Schema.Tables, table)
	}
}

func (s *scanner) addCodeImport(table *Table, requires string) {
	for i := range table.CodeImports {
		if table.CodeImports[i] == requires {
			return
		}
	}
	table.CodeImports = append(table.CodeImports, requires)
}

func (s *scanner) scanColumns(db *pgx.Pool, tableName string, isView bool) []Column {
	result := []Column{}
	statement := "SELECT ordinal_position, column_name, is_nullable, data_type, character_maximum_length, column_default, numeric_precision, " +
		"       pg_catalog.col_description(format('%s.%s',table_schema,table_name)::regclass::oid,ordinal_position) as column_description " +
		"FROM   information_schema.columns " +
		"WHERE  table_schema = $1 " +
		"AND    table_name = $2"
	rows, err := db.Query(bg, statement, s.SchemaName, tableName)
	check(err)
	defer rows.Close()
	for rows.Next() {
		position, name, nullable, dataType, comment := 0, "", "", "", sql.NullString{}
		var maxLen *int
		var columnDefault *string
		var numericPrecision *int
		check(rows.Scan(&position, &name, &nullable, &dataType, &maxLen, &columnDefault, &numericPrecision, &comment))
		isNullable := strings.ToLower(nullable) == "yes"
		col := Column{
			Position:         position,
			ColumnName:       name,
			CodeName:         toProper(name, false),
			DisplayName:      toProper(name, true),
			JsonName:         toJsonName(name),
			SlugName:         toSlug(name),
			Comment:          strings.TrimSpace(comment.String),
			IsNullable:       isNullable,
			IsCardinal:       isPostgresTypeCardinal(dataType),
			HasMaxLen:        maxLen != nil,
			HasDefault:       columnDefault != nil,
			HasPrecision:     numericPrecision != nil,
			CanFilter:        isView,
			SqlType:          dataType,
			DataType:         mapPostgresTypeToGoWithNullable(dataType, isNullable),
			MaxLen:           maxLen,
			ColumnDefault:    columnDefault,
			NumericPrecision: numericPrecision,
		}
		result = append(result, col)
	}
	return result
}

func (s *scanner) scanConstraints(db *pgx.Pool, tableName string) []Constraint {
	result := []Constraint{}
	columnAdded := make(map[string]int)
	statement := "SELECT tc.constraint_name, kc.column_name, tc.constraint_type, " +
		"       cc.table_name as ref_table, cc.column_name as ref_column " +
		"FROM   information_schema.table_constraints tc, information_schema.key_column_usage kc, " +
		"       information_schema.constraint_column_usage cc " +
		"WHERE  kc.table_name = tc.table_name " +
		"AND    kc.table_schema = tc.table_schema " +
		"AND    kc.constraint_name = tc.constraint_name " +
		"AND    cc.constraint_name = tc.constraint_name " +
		"AND    kc.table_schema = $1 " +
		"AND    kc.table_name = $2"
	rows, err := db.Query(bg, statement, s.SchemaName, tableName)
	check(err)
	defer rows.Close()
	for rows.Next() {
		name, columnName, constraintType, refTable, refColumn := "", "", "", "", ""
		check(rows.Scan(&name, &columnName, &constraintType, &refTable, &refColumn))
		if i, ok := columnAdded[name]; ok {
			result[i].ColumnNames = append(result[i].ColumnNames, columnName)
		} else {
			constraint := Constraint{
				ConstraintName: name,
				CodeName:       toProper(name, false),
				DisplayName:    toProper(name, true),
				JsonName:       toJsonName(name),
				SlugName:       toSlug(name),
				IsPrimaryKey:   strings.ToLower(constraintType) == "primary key",
				IsForeignKey:   strings.ToLower(constraintType) == "foreign key",
				IsUniqueKey:    strings.ToLower(constraintType) == "unique",
				ColumnNames:    []string{columnName},
				ConstraintType: constraintType,
				ForeignTable:   nil,
				ForeignColumn:  nil,
			}
			if constraint.IsForeignKey {
				constraint.ForeignTable = &refTable
				constraint.ForeignColumn = &refColumn
			}
			result = append(result, constraint)
			columnAdded[name] = len(result) - 1
		}
	}
	return result
}

func (s *scanner) scanIndexes(db *pgx.Pool, table Table) []Index {
	result := []Index{}
	statement := "SELECT relname, indkey, indisprimary, indisunique " +
		"FROM   pg_class pc, pg_index pi, pg_indexes ps " +
		"WHERE  ps.indexname = relname AND ps.schemaname = $1 AND ps.tablename = $2 " +
		"AND    pc.oid = pi.indexrelid AND pc.relkind = 'i' AND pc.oid IN ( " +
		"  SELECT indexrelid " +
		"  FROM   pg_index pi2, pg_class pc2 " +
		"  WHERE  pc2.relname = $2 " +
		"  AND    pc2.oid = pi2.indrelid " +
		"); "
	rows, err := db.Query(bg, statement, s.SchemaName, table.TableName)
	check(err)
	defer rows.Close()
	for rows.Next() {
		name, indkey, isPrimary, isUnique := "", "", false, false
		check(rows.Scan(&name, &indkey, &isPrimary, &isUnique))

		idx := Index{
			IndexName:    name,
			CodeName:     toProper(name, false),
			DisplayName:  toProper(name, true),
			JsonName:     toJsonName(name),
			SlugName:     toSlug(name),
			ColumnNames:  []string{},
			IsPrimaryKey: isPrimary,
			IsUnique:     isUnique,
		}
		for _, s := range strings.Split(indkey, " ") {
			position, _ := strconv.Atoi(s)
			for _, c := range table.Columns {
				if c.Position == position {
					idx.ColumnNames = append(idx.ColumnNames, c.ColumnName)
				}
			}
		}
		for i := range table.Columns {
			if table.Columns[i].ColumnName == idx.ColumnNames[0] {
				table.Columns[i].IsPrimaryKey = idx.IsPrimaryKey
				table.Columns[i].CanFilter = true
			}
		}
		result = append(result, idx)
	}
	return result
}

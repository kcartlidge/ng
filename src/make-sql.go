package main

import (
	"fmt"
	"strings"
)

func (t Table) GetTableSQL() string {
	if t.Columns == nil || len(t.Columns) == 0 {
		return ""
	}

	schemaName := t.SchemaName
	tableName := t.TableName

	txt := "\n"
	txt += fmt.Sprintf("DROP TABLE %s.%s CASCADE;\n", schemaName, tableName)
	txt += "\n"
	txt += fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s (\n", schemaName, tableName)

	first := true
	for _, c := range t.Columns {
		if !first {
			txt += ",\n"
		}
		first = false
		txt += t.getColumnSQL(c)
	}

	txt += t.getConstraintSQL()

	txt += "\n);\n\n"
	txt += fmt.Sprintf("ALTER TABLE %s.%s OWNER TO %s;\n", schemaName, tableName, t.Owner)

	txt += t.getCommentsSQL()

	return txt
}

func (t *Table) getColumnSQL(c Column) string {
	txt := ""
	nullable := " NOT NULL"
	if c.IsNullable {
		nullable = ""
	}
	defVal := ""
	if c.HasDefault {
		defVal = fmt.Sprintf(" DEFAULT %s", *c.ColumnDefault)
	}
	name := c.ColumnName // .PadRight(table.ColumnNameWidth);
	capacity := ""
	if c.HasMaxLen && *c.MaxLen > 0 {
		capacity = fmt.Sprintf("(%v)", *c.MaxLen)
	}
	sqlType := c.SqlType + capacity
	if c.IsPrimaryKey && c.IsCardinal {
		sqlType = "BIGSERIAL"
		defVal = ""
	}
	txt += fmt.Sprintf("    %s  %s%s%s", name, sqlType, nullable, defVal)
	return txt
}

func (t *Table) getConstraintSQL() string {
	if t.Constraints == nil || len(t.Constraints) == 0 {
		return ""
	}
	txt := ""
	for _, c := range t.Constraints {
		columnNames := strings.Join(c.ColumnNames, ", ")
		txt += fmt.Sprintf(",\n    CONSTRAINT %s %s (%s)", c.ConstraintName, c.ConstraintType, columnNames)
		if c.IsForeignKey {
			txt += fmt.Sprintf("\n        REFERENCES %s.%s (%s) MATCH SIMPLE ", t.SchemaName, *c.ForeignTable, *c.ForeignColumn)
			txt += "\n        ON UPDATE NO ACTION ON DELETE NO ACTION"
		}
	}
	return txt
}

func (t *Table) getCommentsSQL() string {
	txt := ""
	if len(t.Comment) > 0 {
		comment := strings.ReplaceAll(t.Comment, "'", "''")
		txt = fmt.Sprintf("COMMENT ON TABLE %s.%s IS '%s';\n", t.SchemaName, t.TableName, comment)
	}
	for _, c := range t.Columns {
		if len(c.Comment) > 0 {
			comment := strings.ReplaceAll(c.Comment, "'", "''")
			txt += fmt.Sprintf("COMMENT ON COLUMN %s.%s.%s IS '%s';\n", t.SchemaName, t.TableName, c.ColumnName, comment)
		}
	}
	return txt
}

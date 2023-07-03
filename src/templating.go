package main

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"text/template"
	"time"
)

var (
	//go:embed templates/*.tmpl
	fsServer embed.FS
	cache    *template.Template
)

func (w *writer) getTemplatedData(data interface{}, templateName string) []byte {

	if cache == nil {
		tfs, err := fs.Sub(fsServer, "templates")
		check(err)

		cache = template.Must(template.New(templateName).Funcs(template.FuncMap{
			"lower":  strings.ToLower,
			"upper":  strings.ToUpper,
			"plural": toPlural,
			"now":    time.Now,
			"year": func() int {
				return time.Now().Year()
			},
			"sortDate": func(dtm time.Time) string {
				return dtm.Format("2006/01/02")
			},
			"sortDateTime": func(dtm time.Time) string {
				return dtm.Format("2006/01/02 15:04:05")
			},
			"longDate": func(dtm time.Time) string {
				return dtm.Format("Monday January 2, 2006 (MST)")
			},
			"longDateTime": func(dtm time.Time) string {
				return dtm.Format("Monday January 2, 2006 at 15:04 (MST)")
			},
			"TableComment": func(tbl Table) string {
				tableType := "table"
				if tbl.TableType != "BASE TABLE" {
					tableType = strings.ToLower(tbl.TableType)
				}
				txt := fmt.Sprintf("// %s is for %s `%s`", tbl.CodeName, tableType, tbl.TableName)
				if tbl.CodeName != tbl.DisplayName {
					txt += fmt.Sprintf(" (\"%s\")", tbl.DisplayName)
				}
				txt += ".\n"
				if !tbl.IsUpdatable {
					txt += "// It's READ ONLY.\n"
				}
				if len(tbl.Comment) > 0 {
					txt += fmt.Sprintf("// %s", tbl.Comment)
					if !strings.HasSuffix(tbl.Comment, ".") {
						txt += "."
					}
					txt += "\n"
				}
				return txt
			},
			"ColumnComment": func(col Column) string {
				txt := fmt.Sprintf("// %s is for column `%s`", col.CodeName, col.ColumnName)
				if col.CodeName != col.DisplayName {
					txt += fmt.Sprintf(" (\"%s\")", col.DisplayName)
				}
				txt += ".\n"
				if col.IsPrimaryKey {
					txt += "// It's a PRIMARY KEY.\n"
				}
				if col.CanFilter {
					txt += "// It's filterable/sortable.\n"
				}
				if col.HasMaxLen {
					txt += fmt.Sprintf("// It has a maximum size of %v.\n", *col.MaxLen)
				}
				if col.HasDefault {
					txt += fmt.Sprintf("//\n// Default: %s\n", *col.ColumnDefault)
				}
				if len(col.Comment) > 0 {
					txt += fmt.Sprintf("//\n// %s\n", col.Comment)
					if !strings.HasSuffix(col.Comment, ".") {
						txt += "."
					}
				}
				return txt
			},
			"CurrentFolder": func() string {
				cwd, err := os.Getwd()
				check(err)
				return cwd
			},
			"CommandLine": func() string {
				return w.commandLine
			},
			"ConnectionStringEnvArg": func() string {
				return w.connectionStringEnvArg
			},
			"SchemaName": func() string {
				return w.schema.SchemaName
			},
			"ModuleName": func() string {
				return w.module
			},
			"RepoName": func() string {
				return w.repoName
			},
			"PostgresFuncType": func(goType string) string {
				if len(goType) < 2 {
					return goType
				}
				f := goType
				if strings.HasPrefix(f, "*") {
					f = strings.TrimPrefix(f, "*") + "Nullable"
				}
				f = strings.ToUpper(f[0:1]) + f[1:]
				return f
			},
			"inc":                              func(value int) int { return value + 1 },
			"toColumnNameListCSV":              toColumnNameListCSV,
			"toColumnNameListNoPrimaryKeysCSV": toColumnNameListNoPrimaryKeysCSV,
			"toParameterListNoPrimaryKeysCSV":  toParameterListNoPrimaryKeysCSV,
			"toPrimaryKeyParametersCSV":        toPrimaryKeyParametersCSV,
			"toUpdateListNoPrimaryKeysCSV":     toUpdateListNoPrimaryKeysCSV,
			"columnIdxAfterPrimaryKeys":        columnIdxAfterPrimaryKeys,
			"toCodeNameListCSV":                toCodeNameListCSV,
		}).ParseFS(tfs, "*.tmpl"))
	}

	var wr bytes.Buffer
	check(cache.ExecuteTemplate(&wr, templateName, data))
	return wr.Bytes()
}

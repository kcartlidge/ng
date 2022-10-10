package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

// Sets whether RunCmd output is echoed to the screen.
const verbose = false

type writer struct {
	topFolder, entityFolder, repoFolder string
	connectionFolder, supportFolder     string
	schema                              Schema
	commandLine, module                 string
	connectionStringEnvArg              string
	goFilesWritten                      []string
}

func NewWriter(
	folder string,
	module string,
	commandLine string,
	connectionStringEnvArg string,
	schema Schema) writer {
	w := writer{
		topFolder:              path.Clean(folder),
		entityFolder:           path.Join(folder, "entities"),
		repoFolder:             path.Join(folder, "repos"),
		connectionFolder:       path.Join(folder, "connection"),
		supportFolder:          path.Join(folder, "support"),
		module:                 module,
		commandLine:            commandLine,
		connectionStringEnvArg: connectionStringEnvArg,
		schema:                 schema,
		goFilesWritten:         []string{},
	}
	return w
}

func (w writer) WriteStuff() {
	w.clearOutputFolder()
	w.createOutputFolders()
	w.createDumpFile()

	w.createSupportFile()
	w.createEntities()
	w.createConnection()
	w.createRepo()
	w.createEntityRepos()

	w.initModules()
	w.fetchModules()
	w.tidyModules()

	w.applyFormatting()
}

func (w *writer) clearOutputFolder() {
	fmt.Println("Clearing destination folder")
	check(os.RemoveAll(w.topFolder))
}

func (w *writer) createOutputFolders() {
	fmt.Println("Ensuring destination folders exist")
	check(os.MkdirAll(w.topFolder, 0755))
	check(os.MkdirAll(w.supportFolder, 0755))
	check(os.MkdirAll(w.entityFolder, 0755))
	check(os.MkdirAll(w.connectionFolder, 0755))
	check(os.MkdirAll(w.repoFolder, 0755))
}

func (w *writer) createDumpFile() {
	fmt.Println("Writing JSON dump file")
	filename := path.Join(w.topFolder, "dump.json")
	check(ioutil.WriteFile(filename, w.schema.ToJSON(), fs.ModePerm))
}

func (w *writer) createSupportFile() {
	fmt.Println("Creating support file")
	filename := path.Join(w.supportFolder, "support.go")
	w.writeGoFile(filename, "support", nil)
}

func (w *writer) createEntities() {
	fmt.Println("Creating entities")
	for _, table := range w.schema.Tables {
		hasPrimary := false
		for _, c := range table.Columns {
			if c.IsPrimaryKey {
				hasPrimary = true
			}
		}
		if !hasPrimary {
			check(errors.New(fmt.Sprintf("%s.%s has no primary key", w.schema.SchemaName, table.TableName)))
		}
		filename := path.Join(w.entityFolder, table.SlugName+".go")
		w.writeGoFile(filename, "entities", table)
	}
}

func (w *writer) createConnection() {
	fmt.Println("Creating connection")
	filename := path.Join(w.connectionFolder, "connection.go")
	w.writeGoFile(filename, "connection", nil)
}

func (w *writer) createRepo() {
	fmt.Println("Creating repo")
	filename := path.Join(w.repoFolder, "repo-base.go")
	w.writeGoFile(filename, "repo-base", nil)
}

func (w *writer) createEntityRepos() {
	fmt.Println("Adding entity repos")
	for _, table := range w.schema.Tables {
		filename := path.Join(w.repoFolder, table.SlugName+"-repo.go")
		w.writeGoFile(filename, "repos", table)
	}
}

func (w *writer) initModules() {
	fmt.Println("Creating Go module")
	w.runCmd(true, w.topFolder, "go", "mod", "init", w.module)
}

func (w *writer) fetchModules() {
	fmt.Println("Fetching Go modules")
	w.runCmd(true, w.topFolder, "go", "get", "github.com/jackc/pgx/v4")
	w.runCmd(true, w.topFolder, "go", "get", "github.com/jackc/pgx/v4/pgxpool")
}

func (w *writer) tidyModules() {
	fmt.Println("Tidying Go modules")
	w.runCmd(true, w.topFolder, "go", "mod", "tidy")
}

func (w *writer) runCmd(display bool, folder string, command string, args ...string) {
	var out bytes.Buffer
	cmd := exec.Command(command, args...)
	cmd.Dir = folder

	if verbose || display {
		cmd.Stdout = &out
		cmd.Stderr = &out
	}
	err := cmd.Run()
	if err != nil {
		fmt.Println()
		fmt.Println("ERROR")
		fmt.Println(err.Error())
		fmt.Println()
		fmt.Println("Folder  :", folder)
		fmt.Println("Command :", command, strings.Join(args, " "))
		fmt.Println()
		fmt.Println(out.String())
		os.Exit(cmd.ProcessState.ExitCode())
	}
	if verbose && out.Len() > 0 {
		fmt.Println(out.String())
	}
}

func (w *writer) applyFormatting() {
	fmt.Println("Formatting generated Go source")
	for _, filename := range w.goFilesWritten {
		w.runCmd(true, "", "gofmt", "-w", "-s", filename)
	}
}

func (w *writer) writeGoFile(filename string, templateName string, data interface{}) {
	w.writeFile(filename, templateName, data)
	w.goFilesWritten = append(w.goFilesWritten, filename)
}

func (w *writer) writeFile(filename string, templateName string, data interface{}) {
	b := w.getTemplatedData(data, templateName)
	check(ioutil.WriteFile(filename, b, fs.ModePerm))
}

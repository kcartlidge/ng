package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"strings"
)

// Sets whether RunCmd output is echoed to the screen.
const verbose = false

type writer struct {
	topFolder, entityFolder, repoFolder          string
	reposFolder, connectionFolder, supportFolder string
	schema                                       Schema
	commandLine, module, repoName                string
	connectionStringEnvArg                       string
	goFilesWritten                               []string
}

func NewWriter(
	folder string,
	module string,
	commandLine string,
	connectionStringEnvArg string,
	schema Schema,
	repoName string) writer {
	w := writer{
		topFolder:              path.Clean(folder),
		entityFolder:           path.Join(folder, repoName, "entities"),
		repoFolder:             path.Join(folder, repoName),
		reposFolder:            path.Join(folder, repoName, "repos"),
		connectionFolder:       path.Join(folder, repoName, "connection"),
		supportFolder:          path.Join(folder, repoName, "support"),
		module:                 module,
		commandLine:            commandLine,
		connectionStringEnvArg: connectionStringEnvArg,
		schema:                 schema,
		repoName:               repoName,
		goFilesWritten:         []string{},
	}
	return w
}

func (w writer) WriteStuff() {
	w.clearOutputFolder()
	w.createOutputFolders()
	w.createEditorConfigIfNotExists()
	w.createDumpFile()

	w.createSupportFile()
	w.createEntities()
	w.createConnection()
	w.createRepo()
	w.createEntityRepos()
	w.createReadme()
	w.createUsing()
	w.createSQL()

	w.applyFormatting()
}

// clearOutputFolder removes files and folders within the output folder.
// It doesn't remove a pre-existing output folder itself as depending
// upon the OS that can leave existing terminal/command sessions in that
// folder seeming okay but actually working against the original folder
// in the trash (which can go unnoticed).
func (w *writer) clearOutputFolder() {
	exists, err := Exists(w.repoFolder)
	check(err)
	if exists {
		fmt.Println("Clearing target folder")
		files, err := os.ReadDir(w.repoFolder)
		check(err)
		for _, f := range files {
			p := path.Join(w.repoFolder, f.Name())
			if f.IsDir() {
				check(os.RemoveAll(p))
			} else {
				check(os.Remove(p))
			}
		}
	} else {
		fmt.Println("Target folder does not exist")
	}
}

func (w *writer) createOutputFolders() {
	fmt.Println("Ensuring target folder/sub-folders exist")
	check(os.MkdirAll(w.topFolder, 0755))
	check(os.MkdirAll(w.supportFolder, 0755))
	check(os.MkdirAll(w.entityFolder, 0755))
	check(os.MkdirAll(w.connectionFolder, 0755))
	check(os.MkdirAll(w.repoFolder, 0755))
	check(os.MkdirAll(w.reposFolder, 0755))
}

func (w *writer) createEditorConfigIfNotExists() {
	filename := path.Join(w.topFolder, ".editorconfig")
	if exists, err := Exists(filename); err == nil && !exists {
		fmt.Println("Adding missing editorconfig")
		w.writeFile(filename, "editorconfig", nil)
	}
}

func (w *writer) createDumpFile() {
	fmt.Println("Writing JSON dump file")
	filename := path.Join(w.repoFolder, "dump.json")
	check(os.WriteFile(filename, w.schema.ToJSON(), fs.ModePerm))
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
		if table.IsUpdatable && !hasPrimary {
			check(fmt.Errorf("%s.%s has no primary key", w.schema.SchemaName, table.TableName))
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
	fmt.Println("Creating base repo")
	filename := path.Join(w.reposFolder, "repo-base.go")
	w.writeGoFile(filename, "repo-base", nil)
}

func (w *writer) createEntityRepos() {
	fmt.Println("Adding entity repos")
	for _, table := range w.schema.Tables {
		filename := path.Join(w.reposFolder, table.SlugName+"-repo.go")
		w.writeGoFile(filename, "repos", table)
	}
}

func (w *writer) createReadme() {
	fmt.Println("Writing README.md")
	filename := path.Join(w.repoFolder, "README.md")
	w.writeFile(filename, "readme", w.schema)
}

func (w *writer) createUsing() {
	fmt.Println("Writing USING.md")
	filename := path.Join(w.repoFolder, "USING.md")
	w.writeFile(filename, "using", w.schema)
}

func (w *writer) createSQL() {
	fmt.Println("Creating SQL script")
	filename := path.Join(w.repoFolder, "postgres.sql")
	w.writeFile(filename, "sql-scripts", w.schema)
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
	check(os.WriteFile(filename, b, fs.ModePerm))
}

package main

import (
	"bytes"
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
	topFolder string
	schema    Schema
	module    string
}

func NewWriter(
	folder string,
	module string,
	schema Schema) writer {
	w := writer{
		topFolder: path.Clean(folder),
		module:    module,
		schema:    schema,
	}
	return w
}

func (w writer) WriteStuff() {
	w.clearOutputFolder()
	w.createOutputFolders()
	w.createDumpFile()

	w.initModules()
	w.fetchModules()
	w.tidyModules()
}

func (w *writer) clearOutputFolder() {
	fmt.Println("Clearing destination folder")
	check(os.RemoveAll(w.topFolder))
}

func (w *writer) createOutputFolders() {
	fmt.Println("Ensuring destination folders exist")
	check(os.MkdirAll(w.topFolder, 0755))
}

func (w *writer) createDumpFile() {
	fmt.Println("Writing JSON dump file")
	filename := path.Join(w.topFolder, "dump.json")
	check(ioutil.WriteFile(filename, w.schema.ToJSON(), fs.ModePerm))
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

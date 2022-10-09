package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
)

type writer struct {
	topFolder string
	schema    Schema
}

func NewWriter(
	folder string,
	schema Schema) writer {
	w := writer{
		topFolder: path.Clean(folder),
		schema:    schema,
	}
	return w
}

func (w writer) WriteStuff() {
	w.clearOutputFolder()
	w.createOutputFolders()
	w.createDumpFile()
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

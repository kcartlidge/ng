package main

import (
	"fmt"
	"kcartlidge/ng/argsParser"
	"os"
)

func main() {
	// Intro.
	fmt.Println()
	fmt.Println("Near Gothic")
	fmt.Println("Generate Golang repository code directly from Postgres")
	fmt.Println()

	// Command arguments.
	var a = argsParser.New(os.Args)
	a.Example = "-w -env DB_CONNSTR -schema example -folder ~/example/repo"
	a.AddFlag("w", false, false, "overwrite any existing destination folder")
	a.AddValue("env", false, "DB_CONNSTR", "connection string environment variable")
	a.AddValue("schema", false, "public", "the Postgres database schema to scan")
	a.AddValue("folder", true, "", "destination folder, either relative or absolute")
	a.AddValue("module", true, "", "the Go module for created code (eg `kcartlidge/app/data`)")
	a.ShowUsage()
	a.Parse()
	if a.HasIssues {
		a.ShowIssues()
		os.Exit(1)
	}

	// Fetch and show config.
	overwrite := a.Flags["w"]
	env := a.Values["env"]
	schema := a.Values["schema"]
	folder := a.Values["folder"]
	module := a.Values["module"]
	fmt.Println()
	fmt.Println("Environment variable :", env)
	fmt.Println("Database schema      :", schema)
	fmt.Println("Destination folder   :", folder)
	fmt.Println("Go module name       :", module)
	fmt.Println("Overwrite existing?  :", overwrite)
	fmt.Println()
}

package main

import (
	"errors"
	"fmt"
	"kcartlidge/ng/argsParser"
	"os"
)

func main() {
	// Intro.
	fmt.Println()
	fmt.Println("NEAR GOTHIC")
	fmt.Println("Generate Golang repository code directly from Postgres")
	fmt.Println()

	// Command arguments.
	var a = argsParser.New(os.Args)
	a.Example = "-w -env DB_CONNSTR -schema example -module kcartlidge/api -folder ~/example -repo repo"
	a.AddFlag("w", false, false, "overwrite any existing destination folder?")

	a.AddValue("env", false, "DB_CONNSTR", "connection string environment variable")
	a.AddValue("schema", false, "public", "the Postgres database schema to scan")
	a.AddValue("module", true, "", "the Go module for created code (eg `kcartlidge/app/data`)")
	a.AddValue("folder", true, "", "destination folder, either relative or absolute")
	a.AddValue("repo", true, "", "repository subfolder name")
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
	module := a.Values["module"]
	folder := a.Values["folder"]
	repoName := a.Values["repo"]
	fmt.Println()
	fmt.Println("Overwrite existing?  :", overwrite)
	fmt.Println()
	fmt.Println("Environment variable :", env)
	fmt.Println("Database schema      :", schema)
	fmt.Println("Go module name       :", module)
	fmt.Println("Destination folder   :", folder)
	fmt.Println("Repo package name    :", repoName)
	fmt.Println()
	fmt.Println()

	// Fetch the connection string from the env, and test it.
	connectionString, ok := os.LookupEnv(env)
	if !ok {
		check(errors.New("environment variable missing or unreadable"))
	}
	fmt.Println("Obtained connection string from environment")

	// Scan the database to create a schema model.
	s := NewScanner(connectionString, schema)
	err := s.ScanPostgresDatabase()
	check(err)

	// Create the output.
	fmt.Println()
	w := NewWriter(folder, module, a.CommandLine, env, s.Schema, repoName)
	exists, err := Exists(w.topFolder)
	check(err)
	if exists && !overwrite {
		check(errors.New("The output folder exists (-w overwrites)"))
	}
	w.WriteStuff()

	// Done.
	fmt.Println()
	fmt.Println("Done")
	fmt.Println()
}

// If there is an error, display it and quit.
func check(err error) {
	if err != nil {
		fmt.Println()
		fmt.Println("ERROR")
		fmt.Println(err.Error())
		fmt.Println()
		os.Exit(1)
	}
}

// Exists ... Returns true if the path/filename can be found.
func Exists(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

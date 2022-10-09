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
	fmt.Println("Near Gothic")
	fmt.Println("Generate Golang repository code directly from Postgres")
	fmt.Println()

	// Command arguments.
	var a = argsParser.New(os.Args)
	a.Example = "-w -folder ~/example/repo -module kcartlidge/app/data -env DB_CONNSTR -schema example"
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
	w := NewWriter(folder, module, a.CommandLine, env, s.Schema)
	exists, err := Exists(w.topFolder)
	check(err)
	if !exists || overwrite {
		w.WriteStuff()
	} else {
		check(errors.New("Stuff exists at the output location (-w overwrites)"))
	}

	// Done.
	fmt.Println("Done")
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

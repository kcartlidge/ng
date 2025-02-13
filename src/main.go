package main

import (
	"errors"
	"fmt"
	"kcartlidge/ng/argsParser"
	"os"
	"path"
	"strings"
	"time"
)

func main() {
	started := time.Now()

	// Intro.
	fmt.Println()
	fmt.Println("NEAR GOTHIC v1.2.0")
	fmt.Println("Generate Golang repository code directly from Postgres")
	fmt.Println()

	// Command arguments.
	var a = argsParser.New(os.Args)
	a.Example = "-w -env DB_CONNSTR -schema example -module kcartlidge/app -folder ~/Source/App -repo Data"
	a.AddFlag("w", false, false, "overwrite any existing destination folder?")

	a.AddValue("env", false, "DB_CONNSTR", "connection string environment variable")
	a.AddValue("schema", false, "public", "the Postgres database schema to scan")
	a.AddValue("folder", true, "", "the *parent* module's folder (eg `~/Source/App`)")
	a.AddValue("module", true, "", "the *parent* Go module name (eg `kcartlidge/app`)")
	a.AddValue("repo", true, "", "the short folder name for generated code (eg `Data`)")

	a.AddNote("The `env` connection string should be suitable for `jackc/pgx`.")
	a.AddNote("")
	a.AddNote("Generated code is placed in a new (`repo`) folder within the")
	a.AddNote("`parent` folder. For the example above `~/Source/App` + `Data`")
	a.AddNote("means code goes into `~/Source/App/Data`.")
	a.AddNote("")
	a.AddNote("The new code will assume it is in a package named as `module`")
	a.AddNote("plus `repo` (in lower case). For the example above, that means")
	a.AddNote("`kcartlidge/app` + `Data` gives ``kcartlidge/app/data`.")

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
	parentModule := a.Values["module"]
	folder := a.Values["folder"]
	repoName := strings.ToLower(a.Values["repo"])
	module := path.Join(parentModule, repoName)
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

	// Ensure there is something to write.
	if len(s.Schema.Tables) == 0 {
		check(errors.New("no tables were found."))
	}

	// Create the output.
	fmt.Println()
	w := NewWriter(folder, module, a.CommandLine, env, s.Schema, repoName)
	exists, err := Exists(w.topFolder)
	check(err)
	if exists && !overwrite {
		check(errors.New("the output folder exists (-w overwrites)"))
	}
	w.WriteStuff()

	// Advisory.
	fmt.Println()
	fmt.Println("IMPORTANT")
	fmt.Println("This generates self-contained code, not stand-alone code.")
	fmt.Println("Read the (simple) details in the generated USING.md file.")

	// Done.
	elapsed := time.Since(started)
	fmt.Println()
	fmt.Printf("Done in %s\n", elapsed)
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

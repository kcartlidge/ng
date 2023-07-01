package argsParser

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type arguments struct {
	// AppName is the name displayed by ShowUsage.
	// It defaults to `os.Args[0]` but can be overridden.
	AppName string

	// CommandLine is the original command (including the `AppName`).
	CommandLine string

	// Example (if set) contains an example command line.
	// It should NOT include the AppName, only the example arguments.
	Example string

	// Flags holds the boolean command arguments and their value when parsed.
	Flags map[string]bool

	// Values holds the key/value command arguments and their value when parsed.
	Values map[string]string

	// Indent (defaults to 0) pushes the displayed output rightward on the console.
	Indent int

	// IsParsed tracks whether the arguments have been parsed yet.
	IsParsed bool

	// HasIssues is true if, after parsing, there were problems.
	HasIssues bool

	// Issues is a collection of all issues across all arguments when last parsed.
	// The index is the argument position, and the content is an array of messages.
	//
	// There is an index for each argument in turn. An entry will have no messages
	// if there were no issues found during parsing.
	// As the index is also the argument position, it is 1-based. However the Issues
	// collection is 0-based as index 0 holds messages that are not related to a
	// user-provided argument (eg a missing required flag).
	Issues [][]string

	flagsOrder  []string
	valuesOrder []string
	required    map[string]bool
	args        []string
	help        map[string]string
	colWidth    int
}

type foundItem struct {
	Index           int
	IsFlag, IsValue bool
	Tag, Text       string
}

// New creates a new `ArgsParser` with the given arguments.
// Normally these would be the `os.Args“ collection.
func New(args []string) arguments {
	_, fname := filepath.Split(os.Args[0])
	cmdArgs := []string{fname}
	cmdArgs = append(cmdArgs, os.Args[1:]...)
	a := arguments{
		AppName:     filepath.Base(os.Args[0]),
		CommandLine: strings.Join(cmdArgs, " "),
		Example:     "",
		Flags:       make(map[string]bool),
		Values:      make(map[string]string),
		Indent:      0,
		HasIssues:   false,
		IsParsed:    false,
		required:    make(map[string]bool),
		Issues:      [][]string{},
		help:        make(map[string]string),
		colWidth:    1,
		flagsOrder:  []string{},
		valuesOrder: []string{},
	}
	a.args = args[1:]
	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-") {
			args[i] = strings.ToLower(strings.TrimSpace(args[i]))
		}
		a.Issues = append(a.Issues, []string{})
	}
	return a
}

// AddFlag adds a boolean flag whose presence sets the flag to true.
func (a *arguments) AddFlag(name string, required bool, defaultValue bool, help string) {
	name = strings.ToLower(strings.TrimSpace(name))
	a.Flags[name] = defaultValue
	a.help[name] = help
	a.required[name] = required
	a.colWidth = int(math.Max(float64(a.colWidth), float64(len(name))))
	a.flagsOrder = append(a.flagsOrder, name)
}

// AddValue adds a value with the `name` being the key.
// Use an empty string for no `defaultValue`.
func (a *arguments) AddValue(name string, required bool, defaultValue string, help string) {
	name = strings.ToLower(strings.TrimSpace(name))
	a.Values[name] = defaultValue
	a.help[name] = help
	a.required[name] = required
	a.colWidth = int(math.Max(float64(a.colWidth), float64(len(name))))
	a.valuesOrder = append(a.valuesOrder, name)
}

// ShowUsage displays the usage details along with flags, values, and their defaults.
func (a *arguments) ShowUsage() {
	indentString := strings.Repeat(" ", a.Indent)

	// General info.
	fmt.Println()
	fmt.Printf("%sUSAGE\n", indentString)
	fmt.Printf("%s  %s", indentString, a.AppName)
	if len(a.Flags) > 0 || len(a.Values) > 0 {
		for _, k := range a.flagsOrder {
			if a.required[k] {
				fmt.Printf(" -%s", k)
			} else {
				fmt.Printf(" [-%s]", k)
			}
		}
		for _, k := range a.valuesOrder {
			if a.required[k] {
				fmt.Printf(" -%s <value>", k)
			} else {
				fmt.Printf(" [-%s <value>]", k)
			}
		}
		fmt.Println()

		fmt.Println()
		fmt.Printf("%sARGUMENTS\n", indentString)
		hasReq := false

		// Boolean flags.
		if len(a.Flags) > 0 {
			for _, k := range a.flagsOrder {
				req := " "
				if a.required[k] {
					req = "*"
					hasReq = true
				}
				pad := strings.Repeat(" ", a.colWidth-len(k))
				name := "-" + k + "        " + pad
				fmt.Printf("%s  %s  %s  %s", indentString, name, req, a.help[k])
				v := a.Flags[k]
				if v {
					fmt.Printf(" (default %v)", v)
				}
				fmt.Println()
			}
		}

		// Value arguments.
		if len(a.Values) > 0 {
			for _, k := range a.valuesOrder {
				req := " "
				if a.required[k] {
					req = "*"
					hasReq = true
				}
				pad := strings.Repeat(" ", a.colWidth-len(k))
				name := "-" + k + " <value>" + pad
				fmt.Printf("%s  %s  %s  %s", indentString, name, req, a.help[k])
				v := a.Values[k]
				if len(v) > 0 {
					fmt.Printf(" (default `%v`)", v)
				}
				fmt.Println()
			}
		}

		if hasReq {
			fmt.Println()
			fmt.Printf("%s  * means the argument is required\n", indentString)
		}
		fmt.Println()
	}

	// Any provided example.
	if len(a.Example) > 0 {
		fmt.Printf("%sEXAMPLE\n", indentString)
		fmt.Printf("%s  %s %s\n", indentString, a.AppName, a.Example)
		fmt.Println()
	}
}

// Parse extracts the arguments into flags and values.
// To look for success either check for an empty `.Issues`
// collection or look at the `.HasIssues` flag.
func (a *arguments) Parse() {
	if a.IsParsed {
		return
	}
	a.IsParsed = true

	// Gather flags and text without checking for key/value pairs.
	f := []foundItem{}
	for i, s := range a.args {
		text := strings.TrimSpace(s)
		item := foundItem{
			Index: i + 1,
			Tag:   text,
		}
		if strings.HasPrefix(item.Tag, "-") {
			item.IsFlag = true
			item.Tag = removeLeadingDashes(item.Tag)
		}
		f = append(f, item)
	}

	// Convert them into flags and key/value pairs.
	// This is done by looking at sequences of items.
	for i := range f {
		// The first entry doesn't need checking against the previous one.
		if i > 0 {
			if !f[i].IsFlag {
				if f[i-1].IsFlag {
					f[i-1].IsFlag = false
					f[i-1].IsValue = true
					f[i-1].Text = f[i].Tag
					f[i].Tag = ""
				}
			}
		}
	}

	// Add them to the collection of found entries.
	found := []foundItem{}
	for _, i := range f {
		if len(i.Tag) > 0 {
			found = append(found, i)
		} else if i.IsFlag || i.IsValue {
			a.addIssue(a.Issues, i.Index, "Cannot parse item")
		}
	}

	// Gather all the issues with what has been provided.
	for _, i := range found {
		if i.IsFlag {
			if _, found := a.Flags[i.Tag]; !found {
				a.addIssue(a.Issues, i.Index, fmt.Sprintf("Unknown flag `-%s`", i.Tag))
			}
		} else if i.IsValue {
			if _, found := a.Values[i.Tag]; !found {
				a.addIssue(a.Issues, i.Index, fmt.Sprintf("Unknown item `-%s`", i.Tag))
			}
		} else {
			a.addIssue(a.Issues, i.Index, fmt.Sprintf("Unexpected text: `%s`", i.Tag))
		}
	}

	// Gather all the issues with what has NOT been provided.
	for item := range a.Flags {
		isFound := false
		for _, rf := range found {
			if rf.IsFlag && rf.Tag == item {
				isFound = true
				a.Flags[item] = true
			}
		}
		if a.required[item] && !isFound {
			a.addIssue(a.Issues, 0, fmt.Sprintf("Missing flag `-%s`", item))
		}
	}
	for item := range a.Values {
		isFound := false
		for _, rf := range found {
			if rf.IsValue && rf.Tag == item {
				isFound = true
				a.Values[item] = rf.Text
			}
		}
		if a.required[item] && !isFound {
			a.addIssue(a.Issues, 0, fmt.Sprintf("Missing value for `-%s`", item))
		}
	}
}

// ShowIssues displays any argument issues in positional order.
// It also details missing flags and values too.
func (a *arguments) ShowIssues() {
	if !(a.IsParsed && a.HasIssues) {
		return
	}
	indentString := strings.Repeat(" ", a.Indent)
	maxPos := 0
	for position := range a.Issues {
		if position > maxPos {
			maxPos = position
		}
	}
	posWidth := len(strconv.Itoa(maxPos))
	fmt.Printf("%sISSUES\n", indentString)
	for position, ss := range a.Issues {
		for _, s := range ss {
			if position > 0 {
				pad := strings.Repeat(" ", posWidth-len(strconv.Itoa(position)))
				fmt.Printf("%s  %v%s  %s\n", indentString, position, pad, s)
			} else {
				pad := strings.Repeat(" ", posWidth-1)
				fmt.Printf("%s  -%s  %s\n", indentString, pad, s)
			}
		}
	}
	fmt.Println()
}

func (a *arguments) addIssue(issues [][]string, position int, message string) {
	issues[position] = append(issues[position], message)
	a.HasIssues = true
}

func removeLeadingDashes(text string) string {
	text = strings.TrimSpace(text)
	for {
		if strings.HasPrefix(text, "-") {
			text = strings.TrimPrefix(text, "-")
		} else {
			break
		}
	}
	text = strings.TrimSpace(text)
	return text
}

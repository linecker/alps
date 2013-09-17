package main

import "bytes"
import "fmt"
import "regexp"
import "strconv"
import "strings"

// Bundles all global variables.
type globals_t struct {
    // Points to the choosen postprocessing function.
    printer Printer

    // Holds the path to the format specification file.
    format_spec string

    // Holds the path to the fields specification file.
    fields_spec string

    // Holds all known formats (with compiled regexps).
	known_formats []Format

	// Points to the format choosen for processing.
	format *Format

    // Holds all known fields (with compiled regexps).
	known_fields []regexp.Regexp
}
var globals globals_t

// Interface definition for a postprocessing printer: A printer is a function
// that takes a list of chunks.
type Printer func(chunks []Chunk)

// Holds the specification of a format, e.g. syslog.
type Format struct {
	name string
	exp  regexp.Regexp
}

// Construct a new format.
func new_format(n string, e regexp.Regexp) *Format {
	return &Format{
		name: n,
		exp:  e,
	}
}

// Compile a regular expression for a known format.
func load_known_format(line []byte) {

	s := string(line)
	if len(strings.TrimSpace(s)) == 0 {
		return
	}

	// Split the line at the first ':', the first part beeing the name, the
	// second beeing the regexp.
	split := bytes.SplitAfterN(line, []byte(":"), 2)
	if len(split) != 2 {
		fmt.Println("Line '" + string(line) + "' doesn't have the required format 'name: regexp'")
		return
	}
	name := strings.Trim(string(split[0]), ":")
	exp := string(split[1])
	name = strings.Trim(name, " ")
	exp = strings.Trim(exp, " ")

	// Compile the regular expression.
	r, err := regexp.Compile(string(exp))
	if err != nil {
		fmt.Println("Line '"+string(line)+"' doesn't compile to a regular expression:", err.Error())
		return
	}

	// Store name and regular expression as a known format.
	f := new_format(name, *r)
	globals.known_formats = append(globals.known_formats, *f)
}

// Compile a regular expression for a known field.
func load_known_field(line []byte) {
	s := string(line)
	if len(strings.TrimSpace(s)) == 0 {
		return
	}
	r, err := regexp.Compile(s)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	globals.known_fields = append(globals.known_fields, *r)
}

// Searches the known formats for a format name.
func search_known_formats(name string) *Format {
	for i := range globals.known_formats {
		if globals.known_formats[i].name == name {
			return &globals.known_formats[i]
		}
	}
	return nil
}

// Dumps a description of all known formats to the console.
func dump_known_formats() {
	fmt.Printf("\nKnown formats:\n")
	for i := range globals.known_formats {
		fmt.Printf("  '%s'\n", globals.known_formats[i].name)
	}
	total := strconv.Itoa(len(globals.known_formats))
	fmt.Println("  Total " + total)
}

// Dumps a description of all known fields to the console.
func dump_known_fields() {
	fmt.Printf("\nKnown fields:\n")
	for i := range globals.known_fields {
		fmt.Printf("  '%s' containing ", globals.known_fields[i].String())
		names := globals.known_fields[i].SubexpNames()
		for j := range names {
			if j != 0 {
				fmt.Printf("<%s>", names[j])
			}
		}
		fmt.Printf("\n")
	}
	total := strconv.Itoa(len(globals.known_fields))
	fmt.Println("  Total " + total)
}

// Load default values.
func load_config_defaults() {
	globals.printer = print_color
	globals.format_spec = "../data/known-formats.txt"
	globals.fields_spec = "../data/known-fields.txt"
}

// Read the config files.
func read_config_files() {
	for_each_line(globals.format_spec, load_known_format)
	for_each_line(globals.fields_spec, load_known_field)
}

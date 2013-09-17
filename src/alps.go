// Entry point and command line handling.

package main

import "bufio"
import "flag"
import "fmt"
import "os"
import "io"
import "strings"

// Print usage message.
func usage() {

	fmt.Println("Usage: alps [flags] filename.")
	flag.PrintDefaults()
	os.Exit(2)
}

// Linewise tagging from stdin, e.g. "tail -f /var/log/messages | alps".
func read_from_stdin() {
	bio := bufio.NewReader(os.Stdin)
	for {
		line, err := bio.ReadString('\n')
		if err == nil {
			clean := strings.TrimSpace(string(line))
			tag_line([]byte(clean))
		} else if err == io.EOF {
			return
		} else {
			panic(err)
		}
	}
}

func main() {
	var dump bool
	var format string
	var output string
    var format_spec string
    var fields_spec string

    load_config_defaults()

	flag.Usage = usage
	flag.BoolVar(&dump, "dump", false, "Dump config and exit")
	flag.StringVar(&format, "format", "", "Use specific input format")
	flag.StringVar(&output, "output", "", "Use specific output format (color|table|json)")
	flag.StringVar(&format_spec, "format_file", "", "Use a custom format specification file")
	flag.StringVar(&fields_spec, "fields_file", "", "Use a custom field specification file")
	flag.Parse()

    if format_spec != "" {
        globals.format_spec = format_spec
    }
    if fields_spec != "" {
        globals.fields_spec = fields_spec
    }

	// Config files and debug.
	read_config_files()
	if dump {
		dump_known_formats()
		dump_known_fields()
		os.Exit(0)
	}

	// Handle custom postprocessing function (default is color).
	if output != "" {
		if output == "table" {
			globals.printer = print_table
		} else if output == "json" {
			globals.printer = print_json
		}
	}

	// Handle custom input format (default is syslog).
	if format != "" {
		f := search_known_formats(format)
		if f == nil {
			fmt.Println("Unknown format " + format)
			os.Exit(1)
		}
		globals.format = f
	} else {
		globals.format = &globals.known_formats[0]
	}

	// If a file was given, we read the file. Otherwise we use stdin.
	args := flag.Args()
	if len(args) == 1 {
		input_file := args[0]
		for_each_line(input_file, tag_line)
	} else {
		read_from_stdin()
	}
}

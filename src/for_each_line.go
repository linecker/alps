package main

import "bufio"
import "fmt"
import "io"
import "os"
import "strings"

// Remove comments, which are indicated by "<<<<<<", from a line.
func uncomment(line string) string {
    pos := strings.Index(string(line), " <<<<<<")
    if pos == -1 {
        return line
    }
    trimmed := line[0:pos]

    //fmt.Println("input  =", line)
    //fmt.Println("output =", trimmed)
    return trimmed
}

// Opens the given file and calls the given function for each line in the file.
func for_each_line(filename string, function func([]byte)) {
    f, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    defer f.Close()
    r := bufio.NewReader(f)
    for {
        line, err := r.ReadBytes('\n')
        if err == nil {
            clean := strings.Trim(string(line), "\r\n")
            clean = uncomment(clean)
            function([]byte(clean))
        } else if err == io.EOF {
            clean := strings.Trim(string(line), "\r\n")
            clean = uncomment(clean)
            function([]byte(clean))
            break
        } else {
            fmt.Printf(err.Error())
        }
    }
}

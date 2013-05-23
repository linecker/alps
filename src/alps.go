// Next:
// - Better output format(?)
// - Proper offset

package main

import "bufio"
import "fmt"
import "io"
import "os"
import "regexp"
import "strings"
import "strconv"

type Chunk struct {
    pos1 int
    pos2 int
    payload []byte
    tag string
}

var known_formats []regexp.Regexp
var known_fields []regexp.Regexp

func print(chunks []Chunk) {
    c := 1
    for i := range chunks {
        payload := chunks[i].payload
        tag := chunks[i].tag
        var x int
        if tag == "" {
            x = 7 // white
        } else {
            x = c
            c++
            if c > 6 {
                c = 1
            }
        }
        color := "\x1b[3" + strconv.Itoa(x) + "m"
        fmt.Printf("%s%s", color, payload)
    }
    fmt.Println()
}

func print_table(chunks []Chunk) {
    fmt.Println()
    fmt.Println(" Tag              | Payload                                  | Pos1     | Pos2")
    fmt.Println("------------------|----------------------------------------------------------------")
    for i := range chunks {
        payload := chunks[i].payload
        n := 40
        if len(payload) < 40 {
            n = len(payload)
        }
        if len(payload) > 1 {
            tag := chunks[i].tag
            if tag == "" {
                tag = "..."
            }
            fmt.Printf(" %-16s | %-40s | %-8d | %-8d\n", tag, payload[:n], 
                chunks[i].pos1, chunks[i].pos2)
        }
    }
    fmt.Println()
}

func chop_add(chunks []Chunk, bytes []byte, tag string, pos1 int, pos2 int, offset int) []Chunk {
    if len(bytes) == 0 {
        return chunks;
    }
    var c Chunk
    c.pos1 = offset + pos1
    c.pos2 = offset + pos2
    c.payload = bytes
    c.tag = tag
    return append(chunks, c)
}

// Try to match the src Chunk for the regular expression re
// and chop it into pieces...
func chop(src Chunk, re regexp.Regexp, offset int) []Chunk {
    if src.tag != "" {
        return []Chunk{src}
    }
    payload := src.payload
    t := re.FindSubmatchIndex(payload)
    if t == nil { // No match
        return []Chunk{src}
    }
    names := re.SubexpNames()
    var chunked []Chunk
    last := 0
    tag := "";
    for i := 2; i < len(t); i++ {
        cur := t[i]
        if i % 2 != 0 {
            tag = names[i/2]
        } else {
            tag = ""
        }
        chunked = chop_add(chunked, payload[last:cur], tag, last, cur, offset)
        last = cur
    }
    chunked = chop_add(chunked, payload[last:], "", last, len(payload), offset)
    return chunked
}

func preprocess(line []byte) []Chunk {
    var c Chunk
    c.pos1 = 0
    c.pos2 = len(line)
    c.payload = line
    return []Chunk{c}
}

func rec(chunk Chunk) []Chunk{
    var result []Chunk
    for j := 0; j < len(known_fields); j++ {
        r := chop(chunk, known_fields[j], chunk.pos1)
        if len(r) > 1 {
        //if ( r[0].tag != "" ) {
            for k := range r {
                r2 := rec(r[k])
                for l := range r2 {
                    result = append(result, r2[l])
                }
            }
            return result
        }
    }
    return []Chunk{chunk}
} 

func process_line(line []byte) {
    c := preprocess(line)
    chunked := chop(c[0], known_formats[0], 0)
    var result []Chunk
    for i := range chunked {
        r := rec(chunked[i])
        for k := range r {
            result = append(result, r[k])
        }
    }
    print(result)
}

func prepare_regexps(line []byte, list *[]regexp.Regexp) {
    
    r, err := regexp.Compile(string(line))
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    *list = append(*(list), *r)
}

func prepare_known_formats(line []byte) {
    prepare_regexps(line, &known_formats)
}

func prepare_known_fields(line []byte) {
    prepare_regexps(line, &known_fields)
}

// Call function for each line in a file.
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
            function([]byte(clean))
        } else if err == io.EOF {
            break
        } else {
            fmt.Printf(err.Error())
        }
    }
}

func main() {
    for_each_line("../misc/known-formats.txt", prepare_known_formats)
    for_each_line("../misc/known-fields.txt", prepare_known_fields)
    for_each_line("../misc/syslog-sample-01.txt", process_line)
}

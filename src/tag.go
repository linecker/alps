// All things tagging.

package main

import "bufio"
import "fmt"
import "io"
import "os"
import "regexp"
import "strings"

type Chunk struct {
	pos1    int
	pos2    int
	payload []byte
	tag     string
}

// Opens named file and calls named func for each line within the file. Calling
// func happens with a trimmed version of line.
func for_each_line(filename string, function func([]byte)) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	bio := bufio.NewReader(file)
	for {
		line, err := bio.ReadBytes('\n')
		if err != nil && err != io.EOF {
			fmt.Printf(err.Error())
			continue
		}
		clean := strings.TrimSpace(string(line))
		clean = uncomment(clean)
		function([]byte(clean))
		if err == io.EOF {
			break
		}
	}
}

// TODO: Description
func chop_add(chunks []Chunk, bytes []byte, tag string, pos1 int, pos2 int,
	offset int) []Chunk {
	if len(bytes) == 0 {
		return chunks
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

//
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
	tag := ""
	for i := 2; i < len(t); i++ {
		cur := t[i]
		if i%2 != 0 {
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

// Recursively walk to all
func recursive_xxx(chunk Chunk) []Chunk {
	var result []Chunk
	for j := 0; j < len(globals.known_fields); j++ {
		r := chop(chunk, globals.known_fields[j], chunk.pos1)
		if len(r) > 1 {
			//if ( r[0].tag != "" ) {
			for k := range r {
				r2 := recursive_xxx(r[k])
				for l := range r2 {
					result = append(result, r2[l])
				}
			}
			return result
		}
	}
	return []Chunk{chunk}
}

func tag_line(line []byte) {
	initial_chunk := preprocess(line)
	chunked := chop(initial_chunk[0], globals.format.exp, 0)
	var result []Chunk
	for i := range chunked {
		r := recursive_xxx(chunked[i])
		for k := range r {
			result = append(result, r[k])
		}
	}
	globals.printer(result)
}

// Preprocessing functions.

package main

import "strings"

// Converts a slice of bytes into a single Chunk.
func preprocess(payload []byte) []Chunk {
	var c Chunk
	c.pos1 = 0
	c.pos2 = len(payload) - 1
	c.payload = payload
	return []Chunk{c}
}

// Remove trailing comment, indicated by "<<<<<<", from a line.
func uncomment(line string) string {
	pos := strings.Index(string(line), "<<<<<<")
	if pos == -1 {
		return line
	}
	trimmed := line[0:pos]
	return trimmed
}

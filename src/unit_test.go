// Unit tests.

package main

import "testing"
import "fmt"

// See file pre.go.
func Test_preprocess(t *testing.T) {
	input := "test"
	chunk := preprocess([]byte(input))
	if chunk[0].pos1 != 0 {
		t.Error("Invalid pos1")
	}
	if chunk[0].pos2 != 3 {
		t.Error("Invalid pos2")
	}
	if string(chunk[0].payload) != "test" {
		t.Error("Invalid payload")
	}
	if chunk[0].tag != "" {
		t.Error("Invalid tag")
	}
}

// See file pre.go.
func Test_uncomment(t *testing.T) {
	if uncomment("test") != "test" {
		t.Error()
	}
	if uncomment("") != "" {
		t.Error()
	}
	if uncomment("<<<<<<") != "" {
		t.Error()
	}
	if uncomment("test <<<<<") != "test <<<<<" {
		t.Error()
	}
	if uncomment("test <<<<<<") != "test " {
		t.Error()
	}
	if uncomment("test <<<<<<<<") != "test " {
		t.Error()
	}
}

// See file tag.go.
func Test_chop_add(t *testing.T) {

	var c1 Chunk
	c1.pos1 = 0
	c1.pos2 = 8
	c1.payload = []byte("p1")
	c1.tag = "t1"

	var c2 Chunk
	c2.pos1 = 9
	c2.pos2 = 17
	c2.payload = []byte("p2")
	c2.tag = "t2"

	var chunks []Chunk
	chunks = append(chunks, c1)
	chunks = append(chunks, c2)

	result := chop_add(chunks, []byte("p3"), "t3", 0, 8, 17)

	fmt.Print(result)
}

// See file tag.go
func Test_tag_line(t *testing.T) {

}

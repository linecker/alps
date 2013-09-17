// Postprocessing functions.

package main

import "fmt"
import "strconv"

// Print the given chunks in different colors. Colored means tagged, white means
// untagged.
func print_color(chunks []Chunk) {
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

// Print the given chunks as table.
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
		if len(payload) > 0 {
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

// Print the given chunks as json.
func print_json(chunks []Chunk) {
	fmt.Println("{")
	for i := range chunks {
		payload := chunks[i].payload
		tag := chunks[i].tag
		if i > 0 {
			fmt.Println(",")
		}
		fmt.Printf("   \"%s\": \"%s\"", tag, payload)
	}
	fmt.Printf("\n}\n")
}

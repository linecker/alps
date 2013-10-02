// Try to lookup an IPV4 address (convert it to a hostname)
package main

import "fmt"
import "net"
import "os"

func main() {
    input := os.Args[1]
    output := input // fallback
    name, err := net.LookupAddr(input)
    if err == nil {
    	pretty := string(name[0])
    	pretty = pretty[0:len(pretty)-1]
    	output = pretty
    }
    fmt.Print(output)
}

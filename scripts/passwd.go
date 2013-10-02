// run using "go run passwd.go"
package main

import "fmt"
import "os"
import "io"
import "bufio"
import "strings"

func main() {
  f, err := os.Open("/etc/passwd")
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  defer f.Close()
  r := bufio.NewReader(f)
  for {
    line, err := r.ReadString(10) // 0x0A separator = newline
    if err == io.EOF {
      break
    } else if err != nil {
    	fmt.Println(err)
    	os.Exit(1)
    }
    user := strings.SplitAfter(line, ":")[0]
    if user[0] != '#' {
    	user = user[0:len(user)-1]
    	user = "(?P<entity>" + user + ")"
    	fmt.Println(user)
	}
  }
  os.Exit(0)
}


package main

import "os"
import "os/exec"

// Try to compile a plugin.
func try_compile_plugin(tag string) {
    filename := "../plugins/" + tag + ".go"
    if file_exists(filename) {
        compile_plugin(filename, tag)
    }
}

// Compile a plugin to the tmp directory.
func compile_plugin(filename string, tag string) {
    out := globals.tmp_directory + "/" + tag
    cmd := "go"
    arg1 := "build"
    arg2 := "-o"
    arg3 := out
    arg4 := filename
    o, err := exec.Command(cmd, arg1, arg2, arg3, arg4).Output()
    if err != nil {
        println(err.Error())
        println(o)
        return
    }
}

// Try to execute a plugin on named tag.
func try_apply_plugin(tag string, payload []byte) []byte {
    if tag == "" {
        return nil
    }
    filename := globals.tmp_directory + "/" + tag
    if file_exists(filename) {
        return execute_plugin(filename, payload)
    }
    return nil
}

// Checks whether a file exists.
func file_exists(filename string) bool {
    if _, err := os.Stat(filename); err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return true
}

// Execute a compiled plugin.
func execute_plugin(filename string, input []byte) []byte {
    cmd := filename
    arg := string(input)
    out, err := exec.Command(cmd, arg).Output()
    if err != nil {
        println(err.Error())
        return nil
    }
    return out
}

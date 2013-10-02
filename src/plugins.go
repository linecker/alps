
package main

import "os"
import "os/exec"

// TODO: doc
func try_apply_plugin(tag string, payload []byte) []byte {
    filename := "../plugins/" + tag + ".go"
    if plugin_exists(filename) {
        return execute_plugin(filename, payload)
    }
    return nil
}

// Checks whether the plugin file exists.
func plugin_exists(filename string) bool {
    if _, err := os.Stat(filename); err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return true
}

// TODO: doc
func execute_plugin(filename string, input []byte) []byte {
    app := "go"
    arg1 := "run"
    arg2 := filename
    arg3 := string(input)
    out, err := exec.Command(app, arg1, arg2, arg3).Output()
    if err != nil {
        println(err.Error())
        return nil
    }
    return out
}

package main
import (
    "fmt"
    "net/http"
    "syscall"
    "os"
)
func handler(w http.ResponseWriter, r *http.Request) {
    var stat syscall.Statfs_t
    wd, err := os.Getwd()
    if err != nil {
        panic(err)
    }
    syscall.Statfs(wd, &stat)
    fmt.Fprintf(w, "WebServer tests... ")
    fmt.Fprintf(w, "%d", stat.Bavail * uint64(stat.Bsize))
}
func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}

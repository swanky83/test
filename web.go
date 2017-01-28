package main
 
import (
    "net/http"
    "fmt"
    "strings"
    "encoding/json"
    "os"
    "github.com/guillermo/go.procmeminfo"
    "syscall"
    "io/ioutil"
    "strconv"
    "time"
)
type DiskStatus struct {
    All  uint64 `json:"all"`
    Used uint64 `json:"used"`
    Free uint64 `json:"free"`
}
 
// disk usage of path/disk
func DiskUsage(path string) (disk DiskStatus) {
    fs := syscall.Statfs_t{}
    err := syscall.Statfs(path, &fs)
    if err != nil {
        return
    }
    disk.All = fs.Blocks * uint64(fs.Bsize)
    disk.Free = fs.Bfree * uint64(fs.Bsize)
    disk.Used = disk.All - disk.Free
    return
}
 
const (
    B  = 1
    KB = 1024 * B
    MB = 1024 * KB
    GB = 1024 * MB
)
 
func getCPUavg() (idle, total uint64) {
    contents, err := ioutil.ReadFile("/proc/stat")
    if err != nil {
        return
    }
    lines := strings.Split(string(contents), "\n")
    for _, line := range(lines) {
        fields := strings.Fields(line)
        if fields[0] == "cpu" {
            numFields := len(fields)
            for i := 1; i < numFields; i++ {
                val, err := strconv.ParseUint(fields[i], 10, 64)
                if err != nil {
                    fmt.Println("Error: ", i, fields[i], err)
                }
                total += val // tally up all the numbers to get total ticks
                if i == 4 {  // idle is the 5th field in the cpu line
                    idle = val
                }
            }
            return
        }
    }
    return
}
 
type ServerStatus struct {
    Hostname  string `json:"hostname"`
    CPUAvg    float64 `json:"cpu_avg"`
    DiskUsage float64 `json:"disk_usage"`
    MemUsage  uint64 `json:"mem_usage"`
    CheckTime string `json:"check_time"`
}
 
func statusHandler(w http.ResponseWriter, r *http.Request) {
    name, err := os.Hostname()
 
    if err != nil {
        panic(err)
    }
 
    meminfo := &procmeminfo.MemInfo{}
    meminfo.Update()
    usage_memory := meminfo.Used()
    if usage_memory > 0 {
        usage_memory = usage_memory / 1024 /1024
    }
    disk := DiskUsage("/")
    free_space := float64(disk.Free)/float64(GB)
 
 
    idle0, total0 := getCPUavg()
    time.Sleep(3 * time.Second)
    idle1, total1 := getCPUavg()
 
    idleTicks := float64(idle1 - idle0)
    totalTicks := float64(total1 - total0)
    cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks
 
 
    sstatus := ServerStatus{name, cpuUsage, free_space, usage_memory, "Alex5"}
 
    js, err := json.Marshal(sstatus)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
 
    w.Header().Set("Content-Type", "application/json")
    w.Write(js)
}
 
func helloHandler(w http.ResponseWriter, r *http.Request) {
    remPartOfURL := r.URL.Path[len("/hello/"):] //get everything after the /hello/ part of the URL
    fmt.Fprintf(w, "Hello %s!", remPartOfURL)
}
 
func shouthelloHandler(w http.ResponseWriter, r *http.Request) {
    remPartOfURL := r.URL.Path[len("/shouthello/"):] //get everything after the /shouthello/ part of the URL
    fmt.Fprintf(w, "Hello %s!", strings.ToUpper(remPartOfURL))
}
 
func main() {
    http.HandleFunc("/status", statusHandler)
    http.ListenAndServe("0.0.0.0:9999", nil)
}
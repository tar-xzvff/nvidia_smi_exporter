package main

import (
    "bytes"
    "encoding/csv"
    "fmt"
    "net/http"
    "log"
    "os"
    "os/exec"
    "strings"
)


// name, index, temperature.gpu, utilization.gpu,
// utilization.memory, memory.total, memory.free, memory.used

func metrics(response http.ResponseWriter, request *http.Request) {
    out, err := exec.Command(
        "nvidia-smi",
        "--query-gpu=name,index,timestamp,temperature.gpu,utilization.gpu,utilization.memory,memory.total,memory.free,memory.used,clocks_throttle_reasons.sw_power_cap,clocks_throttle_reasons.sw_thermal_slowdown,clocks.current.graphics,temperature.memory,power.draw",
        "--format=csv,noheader,nounits").Output()

    if err != nil {
        fmt.Printf("%s\n", err)
        return
    }

    csvReader := csv.NewReader(bytes.NewReader(out))
    csvReader.TrimLeadingSpace = true
    records, err := csvReader.ReadAll()

    if err != nil {
        fmt.Printf("%s\n", err)
        return
    }

    metricList := []string {
    	"timestamp",
        "temperature.gpu", "utilization.gpu",
        "utilization.memory", "memory.total", "memory.free", "memory.used",
        "clocks_throttle_reasons.sw_power_cap", "clocks_throttle_reasons.sw_thermal_slowdown", "clocks.current.graphics",
        "temperature.memory", "power.draw"}

    result := ""
    for _, row := range records {
        name := fmt.Sprintf("%s[%s]", row[0], row[1])
        for idx, value := range row[2:] {
            result = fmt.Sprintf(
                "%s%s{gpu=\"%s\"} %s\n", result,
                strings.Replace(metricList[idx], ".", "_", -1), name, value)
        }
    }
    fmt.Fprintf(response, strings.Replace(result, "", "", -1))
}

func main() {
    addr := ":9101"
    if len(os.Args) > 1 {
        addr = ":" + os.Args[1]
    }

    http.HandleFunc("/metrics/", metrics)
    err := http.ListenAndServe(addr, nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

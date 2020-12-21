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
        "--query-gpu=name,index,timestamp,temperature.gpu,utilization.gpu,utilization.memory,memory.total,memory.free,memory.used,clocks_throttle_reasons.sw_power_cap,clocks_throttle_reasons.sw_thermal_slowdown,clocks.current.graphics,temperature.memory,power.draw,driver_version,count,gpu_serial,pcie.link.gen.current,pcie.link.gen.max,pcie.link.width.current,pcie.link.width.max,fan.speed,pstate,clocks_throttle_reasons.supported,clocks_throttle_reasons.active,clocks_throttle_reasons.gpu_idle,clocks_throttle_reasons.applications_clocks_setting,clocks_throttle_reasons.hw_slowdown,clocks_throttle_reasons.hw_thermal_slowdown,clocks_throttle_reasons.hw_power_brake_slowdown,clocks_throttle_reasons.sync_boost,power.limit,enforced.power.limit,power.default_limit,power.min_limit,power.max_limit,clocks.current.sm,clocks.current.memory,clocks.current.video",
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
        "temperature.memory", "power.draw", "driver_version","count","gpu_serial","pcie.link.gen.current",
        "pcie.link.gen.max", "pcie.link.width.current",
        "pcie.link.width.max", "fan.speed",
        "pstate", "clocks_throttle_reasons.supported",
        "clocks_throttle_reasons.active",
        "clocks_throttle_reasons.gpu_idle",
        "clocks_throttle_reasons.applications_clocks_setting",
        "clocks_throttle_reasons.hw_slowdown", "clocks_throttle_reasons.hw_thermal_slowdown",
        "clocks_throttle_reasons.hw_power_brake_slowdown",
        "clocks_throttle_reasons.sync_boost", "power.limit",
        "enforced.power.limit", "power.default_limit",
        "power.min_limit", "power.max_limit",
        "clocks.current.sm",
        "clocks.current.memory", "clocks.current.video"}

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

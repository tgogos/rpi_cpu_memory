package main

import (
    "log"
    "fmt"
    "time"
    linuxproc "github.com/c9s/goprocinfo/linux"
)


// source: http://stackoverflow.com/questions/23367857/accurate-calculation-of-cpu-usage-given-in-percentage-in-linux

// PrevIdle = previdle + previowait
// Idle = idle + iowait
//
// PrevNonIdle = prevuser + prevnice + prevsystem + previrq + prevsoftirq + prevsteal
// NonIdle = user + nice + system + irq + softirq + steal
//
// PrevTotal = PrevIdle + PrevNonIdle
// Total = Idle + NonIdle
//
// # differentiate: actual value minus the previous one
// totald = Total - PrevTotal
// idled = Idle - PrevIdle
//
// CPU_Percentage = (totald - idled)/totald


// type CPUStat struct {
//     Id        string `json:"id"`
//     User      uint64 `json:"user"`
//     Nice      uint64 `json:"nice"`
//     System    uint64 `json:"system"`
//     Idle      uint64 `json:"idle"`
//     IOWait    uint64 `json:"iowait"`
//     IRQ       uint64 `json:"irq"`
//     SoftIRQ   uint64 `json:"softirq"`
//     Steal     uint64 `json:"steal"`
//     Guest     uint64 `json:"guest"`
//     GuestNice uint64 `json:"guest_nice"`
// }

type MyCPUStats struct {
  Idle      int
  NonIdle   int
  Total     int
}


func main() {

  var prevCPUStats *linuxproc.CPUStat
  var currCPUStats *linuxproc.CPUStat

  currCPUStats := ReadCPUStats()
  prevCPUStats := ReadCPUStats()

  for {
    time.Sleep(time.Second * 10)

    currCPUStats := ReadCPUStats()


    time.Sleep(time.Second * 10)
  }

}


func ReadCPUStats() *linuxproc.Stat {
  stat, err := linuxproc.ReadStat("/proc/stat")
  if err != nil {
      log.Fatal("stat read fail")
  }
  //fmt.Println(stat)
  return stat

  // for _, s := range stat.CPUStats {
  //     // s.User
  //     // s.Nice
  //     // s.System
  //     // s.Idle
  //     // s.IOWait
  //     fmt.Println(s)
  // }
}


func calcMyCPUStats() *MyCPUStats {
  stats := ReadCPUStats()
  
}

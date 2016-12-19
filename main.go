package main

import (
	"fmt"
	"log"
	"time"

	linuxproc "github.com/c9s/goprocinfo/linux"
)

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
	Cpu0 float32
	Cpu1 float32
	Cpu2 float32
	Cpu3 float32
}

func main() {

	//var prevCPUStats linuxproc.CPUStat
	//var currCPUStats linuxproc.CPUStat

	currCPUStats := ReadCPUStats()
	prevCPUStats := ReadCPUStats()

	for {
		time.Sleep(time.Second * 4)

		currCPUStats = ReadCPUStats()
		calcMyCPUStats(currCPUStats, prevCPUStats)
		prevCPUStats = currCPUStats

	}

}

func ReadCPUStats() *linuxproc.Stat {
	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		log.Fatal("stat read fail")
	}
	// fmt.Println(stat)
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

func calcMyCPUStats(curr, prev *linuxproc.Stat) *MyCPUStats {

	var stats MyCPUStats

	// fmt.Println("cpu0: ", calcSingleCoreUsage(curr.CPUStats[0], prev.CPUStats[0]))
	// fmt.Println("cpu1: ", calcSingleCoreUsage(curr.CPUStats[1], prev.CPUStats[1]))
	// fmt.Println("cpu2: ", calcSingleCoreUsage(curr.CPUStats[2], prev.CPUStats[2]))
	// fmt.Println("cpu3: ", calcSingleCoreUsage(curr.CPUStats[3], prev.CPUStats[3]))

	stats.Cpu0 = calcSingleCoreUsage(curr.CPUStats[0], prev.CPUStats[0])
	stats.Cpu1 = calcSingleCoreUsage(curr.CPUStats[1], prev.CPUStats[1])
	stats.Cpu2 = calcSingleCoreUsage(curr.CPUStats[2], prev.CPUStats[2])
	stats.Cpu3 = calcSingleCoreUsage(curr.CPUStats[3], prev.CPUStats[3])

	fmt.Printf("Stats:\n%+v", stats)

	return &stats
}

func calcSingleCoreUsage(curr, prev linuxproc.CPUStat) float32 {

	/*
	 *    source: http://stackoverflow.com/questions/23367857/accurate-calculation-of-cpu-usage-given-in-percentage-in-linux
	 *
	 *    PrevIdle = previdle + previowait
	 *    Idle = idle + iowait
	 *
	 *    PrevNonIdle = prevuser + prevnice + prevsystem + previrq + prevsoftirq + prevsteal
	 *    NonIdle = user + nice + system + irq + softirq + steal
	 *
	 *    PrevTotal = PrevIdle + PrevNonIdle
	 *    Total = Idle + NonIdle
	 *
	 *    # differentiate: actual value minus the previous one
	 *    totald = Total - PrevTotal
	 *    idled = Idle - PrevIdle
	 *
	 *    CPU_Percentage = (totald - idled)/totald
	 */

	// fmt.Printf("curr:\n %+v\n", curr)
	// fmt.Printf("prev:\n %+v\n", prev)

	PrevIdle := prev.Idle + prev.IOWait
	Idle := curr.Idle + curr.IOWait

	PrevNonIdle := prev.User + prev.Nice + prev.System + prev.IRQ + prev.SoftIRQ + prev.Steal
	NonIdle := curr.User + curr.Nice + curr.System + curr.IRQ + curr.SoftIRQ + curr.Steal

	PrevTotal := PrevIdle + PrevNonIdle
	Total := Idle + NonIdle
	// fmt.Println(PrevIdle, Idle, PrevNonIdle, NonIdle, PrevTotal, Total)

	//  differentiate: actual value minus the previous one
	totald := Total - PrevTotal
	idled := Idle - PrevIdle

	CPU_Percentage := (float32(totald) - float32(idled)) / float32(totald)

	return CPU_Percentage
}

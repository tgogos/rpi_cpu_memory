package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	linuxproc "github.com/c9s/goprocinfo/linux"
)

type MyCPUStats struct {
	Cpu0 float32
	Cpu1 float32
	Cpu2 float32
	Cpu3 float32
}

func main() {

	time_interval := 1 // this number represents seconds
	push_to_influx := true
	influxUrl := "http://10.143.0.218:8086"
	cpuDBname := "pi_cpu"
	// memoDBname := "pi_memo"

	currCPUStats := ReadCPUStats()
	prevCPUStats := ReadCPUStats()
	info := ReadMemoInfo()

	for {
		time.Sleep(time.Second * time.Duration(time_interval))

		currCPUStats = ReadCPUStats()
		coreStats := calcMyCPUStats(currCPUStats, prevCPUStats)
		fmt.Println(coreStats)
		if push_to_influx {
			url := influxUrl + "/write?db=" + cpuDBname
			body := []byte("cpu0,coreID=0 value=" + strconv.FormatFloat(float64(coreStats.Cpu0), 'f', -1, 32) + "\n" +
				"cpu1,coreID=1 value=" + strconv.FormatFloat(float64(coreStats.Cpu0), 'f', -1, 32) + "\n" +
				"cpu2,coreID=2 value=" + strconv.FormatFloat(float64(coreStats.Cpu0), 'f', -1, 32) + "\n" +
				"cpu3,coreID=3 value=" + strconv.FormatFloat(float64(coreStats.Cpu0), 'f', -1, 32))
			// fmt.Printf("%s", body)
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
			hc := http.Client{}
			resp, err := hc.Do(req)
			if err != nil {
				log.Fatal("could not send POST")
			}
			fmt.Println(resp)
		}
		prevCPUStats = currCPUStats

		info = ReadMemoInfo()
		memUsed := info.MemTotal - info.MemAvailable
		fmt.Println("total: ", info.MemTotal, " free: ", info.MemFree)
		fmt.Printf("Memory used: %d\n\n\n", memUsed)

	}

}

func ReadCPUStats() *linuxproc.Stat {
	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		log.Fatal("stat read fail")
	}
	// fmt.Println(stat)
	return stat
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

	fmt.Printf("CPU stats:\n%+v\n", stats)

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

//
//  Memory
//
//

func ReadMemoInfo() *linuxproc.MemInfo {
	info, err := linuxproc.ReadMemInfo("/proc/meminfo")
	if err != nil {
		log.Fatal("info read fail")
	}
	// fmt.Printf("Memory info struct:\n%+v", info)
	return info
}

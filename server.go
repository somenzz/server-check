package main

import (
	"runtime"
	"sort"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

type Server struct {
	Os   Os   `json:"os"`
	Cpu  Cpu  `json:"cpu"`
	Ram  Ram  `json:"ram"`
	Disk Disk `json:"disk"`
}

type Os struct {
	GOOS         string `json:"goos"`
	NumCPU       int    `json:"numCpu"`
	Compiler     string `json:"compiler"`
	GoVersion    string `json:"goVersion"`
	NumGoroutine int    `json:"numGoroutine"`
}

type Cpu struct {
	Cpus  []float64 `json:"cpus"`
	Cores int       `json:"cores"`
}

type Ram struct {
	UsedMB      int     `json:"usedMb"`
	TotalMB     int     `json:"totalMb"`
	UsedPercent float64 `json:"usedPercent"`
}

type Disk struct {
	UsedMB      int     `json:"usedMb"`
	UsedGB      int     `json:"usedGb"`
	TotalMB     int     `json:"totalMb"`
	TotalGB     int     `json:"totalGb"`
	UsedPercent float64 `json:"usedPercent"`
}

type ProcessInfo struct {
	Name       string
	Exe        string
	CPUPercent float64
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: InitCPU
//@description: OS信息
//@return: o Os, err error

func InitOS() (o Os) {
	o.GOOS = runtime.GOOS
	o.NumCPU = runtime.NumCPU()
	o.Compiler = runtime.Compiler
	o.GoVersion = runtime.Version()
	o.NumGoroutine = runtime.NumGoroutine()
	return o
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: InitCPU
//@description: CPU信息
//@return: c Cpu, err error

func InitCPU() (c Cpu, err error) {
	if cores, err := cpu.Counts(false); err != nil {
		return c, err
	} else {
		c.Cores = cores
	}
	if cpus, err := cpu.Percent(time.Duration(1)*time.Second, false); err != nil {
		return c, err
	} else {
		c.Cpus = cpus
	}
	return c, nil
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: InitRAM
//@description: RAM信息
//@return: r Ram, err error

func InitRAM() (r Ram, err error) {
	if u, err := mem.VirtualMemory(); err != nil {
		return r, err
	} else {
		r.UsedMB = int(u.Used) / MB
		r.TotalMB = int(u.Total) / MB
		r.UsedPercent = u.UsedPercent
	}
	return r, nil
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: InitDisk
//@description: 硬盘信息
//@return: d Disk, err error

func InitDisk() (d Disk, err error) {
	if u, err := disk.Usage("/"); err != nil {
		return d, err
	} else {
		d.UsedMB = int(u.Used) / MB
		d.UsedGB = int(u.Used) / GB
		d.TotalMB = int(u.Total) / MB
		d.TotalGB = int(u.Total) / GB
		d.UsedPercent = u.UsedPercent
	}
	return d, nil
}

func InitProcess() (processInfos []ProcessInfo, err error) {
	processes, err := process.Processes()
	if err != nil {
		return processInfos, err
	}

	for _, p := range processes {
		cpuPercent, err := p.CPUPercent()
		if err == nil {
			cpuPercent = float64(int(cpuPercent*100)) / 100
		}
		if err != nil {
			continue
		}

		name, _ := p.Name()
		exe, _ := p.Exe()

		processInfos = append(processInfos, ProcessInfo{
			Name:       name,
			Exe:        exe,
			CPUPercent: cpuPercent,
		})
	}

	// Sort processes by CPU usage in descending order
	sort.Slice(processInfos, func(i, j int) bool {
		return processInfos[i].CPUPercent > processInfos[j].CPUPercent
	})

	return processInfos[:5], nil
}

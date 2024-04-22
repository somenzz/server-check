package main

import (
	"fmt"
	"syscall"

	"github.com/somenzz/ewechat"
)

func checkDiskSpace(path string, threshold float64) {
	var ewechat = ewechat.EWechat{
		CorpID:     CFG.EWeChat.CorpID,
		CorpSecret: CFG.EWeChat.CorpSecret,
		AgentID:    CFG.EWeChat.AgentID,
	}
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		fmt.Printf("Error getting filesystem info: %s\n", err)
		return
	}

	// Available blocks * size per block = available space
	available := float64(fs.Bavail * uint64(fs.Bsize))
	// Total blocks * size per block = total space
	total := float64(fs.Blocks * uint64(fs.Bsize))

	// Calculate the percentage of free disk space
	percentFree := (available / total) * 100
	var msg string
	fmt.Println(threshold)
	if percentFree < threshold {
		msg = fmt.Sprintf("Warning: Disk space is %.2f%% and below %.2f%% on %s\n", percentFree, threshold, path)
		ewechat.SendMessage(msg, CFG.EWeChat.Receivers)
	} else {
		msg = fmt.Sprintf("Disk space is sufficient (%.2f%% free) on %s\n", percentFree, path)

	}
	fmt.Println(msg)

}

func main() {
	var ewechat = ewechat.EWechat{
		CorpID:     CFG.EWeChat.CorpID,
		CorpSecret: CFG.EWeChat.CorpSecret,
		AgentID:    CFG.EWeChat.AgentID,
	}

	disk, err := InitDisk()
	if err != nil {
		ewechat.SendMessage(fmt.Sprintf("disk read error: %s", err.Error()), CFG.EWeChat.Receivers)
	}

	if disk.UsedPercent > CFG.DiskUsageRate {

		msg := fmt.Sprintf("Warning: Disk usage rate is %.2f%% and over DiskUsageRate %.2f%%", disk.UsedPercent, CFG.DiskUsageRate)
		// fmt.Println(msg)
		ewechat.SendMessage(msg, CFG.EWeChat.Receivers)

	}

	cpu, err := InitCPU()
	if err != nil {
		ewechat.SendMessage(fmt.Sprintf("cpu read error: %s", err.Error()), CFG.EWeChat.Receivers)
	}

	if cpu.Cpus[0] > CFG.DiskUsageRate {

		msg := fmt.Sprintf("Warning: CPU usage rate is %.2f%% and over DiskUsageRate %.2f%%", cpu.Cpus[0], CFG.DiskUsageRate)
		// fmt.Println(msg)
		ewechat.SendMessage(msg, CFG.EWeChat.Receivers)

	}

	ram, err := InitRAM()
	if err != nil {
		ewechat.SendMessage(fmt.Sprintf("ram read error: %s", err.Error()), CFG.EWeChat.Receivers)
	}

	if ram.UsedPercent > CFG.MemUsageRate {

		msg := fmt.Sprintf("Warning: Ram usage rate is %.2f%% and over DiskUsageRate %.2f%%", ram.UsedPercent, CFG.MemUsageRate)
		// fmt.Println(msg)
		ewechat.SendMessage(msg, CFG.EWeChat.Receivers)

	}

}

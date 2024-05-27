package main

import (
	"fmt"
	"net"

	"github.com/somenzz/ewechat"
)

func getLocalIP() ([]string, error) {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips, nil
}

func main() {

	ips, err := getLocalIP()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, ip := range ips {
		fmt.Println("Local machine IP address:", ip)
	}

	var ewechat = ewechat.EWechat{
		CorpID:     CFG.EWeChat.CorpID,
		CorpSecret: CFG.EWeChat.CorpSecret,
		AgentID:    CFG.EWeChat.AgentID,
	}
	msg_prefix := fmt.Sprintf("IP address: %s", ips[0])
	disk, err := InitDisk()

	if err != nil {
		ewechat.SendMessage(fmt.Sprintf("%s disk read error: %s", msg_prefix, err.Error()), CFG.EWeChat.Receivers)
	}

	if disk.UsedPercent > CFG.DiskUsageRate {

		msg := fmt.Sprintf("%s Warning: Disk usage rate is %.2f%% and over DiskUsageRate %.2f%%", msg_prefix, disk.UsedPercent, CFG.DiskUsageRate)
		// fmt.Println(msg)
		ewechat.SendMessage(msg, CFG.EWeChat.Receivers)

	}

	cpu, err := InitCPU()
	if err != nil {
		ewechat.SendMessage(fmt.Sprintf("%s cpu read error: %s", msg_prefix, err.Error()), CFG.EWeChat.Receivers)
	}

	if cpu.Cpus[0] > CFG.CpuUsageRate {

		msg := fmt.Sprintf("%s Warning: CPU usage rate is %.2f%% and over CpuUsageRate %.2f%%", msg_prefix, cpu.Cpus[0], CFG.CpuUsageRate)
		// fmt.Println(msg)
		ewechat.SendMessage(msg, CFG.EWeChat.Receivers)

	}

	ram, err := InitRAM()
	if err != nil {
		ewechat.SendMessage(fmt.Sprintf("%s ram read error: %s", msg_prefix, err.Error()), CFG.EWeChat.Receivers)
	}

	if ram.UsedPercent > CFG.MemUsageRate {

		msg := fmt.Sprintf("%s Warning: Ram usage rate is %.2f%% and over MemUsageRate %.2f%%", msg_prefix, ram.UsedPercent, CFG.MemUsageRate)
		// fmt.Println(msg)
		ewechat.SendMessage(msg, CFG.EWeChat.Receivers)

	}

}

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/somenzz/ewechat"
	"github.com/somenzz/server-check/http_check"
)

var CFG = GetConfig()

var ewechatSender = ewechat.EWechat{
	CorpID:     CFG.EWeChat.CorpID,
	CorpSecret: CFG.EWeChat.CorpSecret,
	AgentID:    CFG.EWeChat.AgentID,
}

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

func CheckUrlIsHealth(url, method string, expectStatusCode int, expectBody string) {

	maxRetries := 3
	retryDelay := time.Second * 5

	for i := 0; i < maxRetries; i++ {
		if http_check.CheckHealth(url, method, expectStatusCode, expectBody) {
			log.Printf("Service at %s is healthy\n", url)
			return
		}

		if i < maxRetries-1 {
			log.Printf("Service unhealthy. Retrying in %v...\n", retryDelay)
			time.Sleep(retryDelay)
		}
	}

	ewechatSender.SendMessage(fmt.Sprintf("Service at %s is unhealthy after %d attempts\n", url, maxRetries), CFG.EWeChat.Receivers)

}

func main() {

	// Get the path to the executable.
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	// Resolve the directory of the executable.
	exePath := filepath.Dir(exe)
	log.Println("working path:", exePath)
	// Change the working directory to the executable's directory.
	err = os.Chdir(exePath)
	if err != nil {
		log.Fatal(err)
	}

	ips, err := getLocalIP()
	if err != nil {
		log.Fatal(err)
	}
	for _, ip := range ips {
		log.Println("Local machine IP address:", ip)
	}

	msg_prefix := fmt.Sprintf("IP address: %s", ips[0])
	disk, err := InitDisk()

	if err != nil {
		ewechatSender.SendMessage(fmt.Sprintf("%s disk read error: %s", msg_prefix, err.Error()), CFG.EWeChat.Receivers)
	}

	if disk.UsedPercent > CFG.DiskUsageRate {

		msg := fmt.Sprintf("%s Warning: Disk usage rate is %.2f%% and over DiskUsageRate %.2f%%", msg_prefix, disk.UsedPercent, CFG.DiskUsageRate)
		// fmt.Println(msg)
		ewechatSender.SendMessage(msg, CFG.EWeChat.Receivers)

	}

	cpu, err := InitCPU()
	if err != nil {
		ewechatSender.SendMessage(fmt.Sprintf("%s cpu read error: %s", msg_prefix, err.Error()), CFG.EWeChat.Receivers)
	}

	if cpu.Cpus[0] > CFG.CpuUsageRate {

		msg := fmt.Sprintf("%s Warning: CPU usage rate is %.2f%% and over CpuUsageRate %.2f%%", msg_prefix, cpu.Cpus[0], CFG.CpuUsageRate)
		// fmt.Println(msg)
		ewechatSender.SendMessage(msg, CFG.EWeChat.Receivers)

	}

	ram, err := InitRAM()
	if err != nil {
		ewechatSender.SendMessage(fmt.Sprintf("%s ram read error: %s", msg_prefix, err.Error()), CFG.EWeChat.Receivers)
	}

	if ram.UsedPercent > CFG.MemUsageRate {

		msg := fmt.Sprintf("%s Warning: Ram usage rate is %.2f%% and over MemUsageRate %.2f%%", msg_prefix, ram.UsedPercent, CFG.MemUsageRate)
		// fmt.Println(msg)
		ewechatSender.SendMessage(msg, CFG.EWeChat.Receivers)

	}

	//url 健康检查

	for _, url := range CFG.CheckUrl {
		CheckUrlIsHealth(url.Url, url.Method, url.ExpectStatusCode, url.ExpectBody)
	}

}

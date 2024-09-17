package main

import (
	"fmt"
	"github.com/somenzz/ewechat"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

func healthCheck(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("无法访问 %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("API 返回状态码不正常，状态码: %d，内容: %s", resp.StatusCode, string(body))
	}

	body, _ := ioutil.ReadAll(resp.Body)
	body_str := string(body)
	if body_str != `"ok"` {
		fmt.Println(body_str)
		return fmt.Errorf("API 返回结果不正确，期望 'ok'，实际: %s", body_str)
	}

	return nil
	//return fmt.Errorf("API health") //for test
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

	err = healthCheck(CFG.ApiUrl)
	if err != nil {
		ewechat.SendMessage(err.Error(), CFG.EWeChat.Receivers)
	}
}

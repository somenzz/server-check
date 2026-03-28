package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/somenzz/ewechat"
	"github.com/somenzz/server-check/check_pid"
	"github.com/somenzz/server-check/http_check"
)

var CFG = GetConfig()

var ewechatSender = ewechat.EWechat{
	CorpID:     CFG.EWeChat.CorpID,
	CorpSecret: CFG.EWeChat.CorpSecret,
	AgentID:    CFG.EWeChat.AgentID,
}

func sendWechatMessage(message string, receivers string) {
	msg, err := ewechatSender.SendMessage(message, receivers)
	if err != nil {
		log.Printf("send message error: %s", err.Error())
	}
	log.Printf("send message: %s", msg)
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

func CheckUrlIsHealth(ctx context.Context, url, method string, expectStatusCode int, expectBody string) {

	maxRetries := 3
	retryDelay := time.Second * 5

	for i := 0; i < maxRetries; i++ {
		if http_check.CheckHealth(ctx, url, method, expectStatusCode, expectBody) {
			log.Printf("Service at %s is healthy\n", url)
			return
		}

		if i < maxRetries-1 {
			log.Printf("Service unhealthy. Retrying in %v...\n", retryDelay)
			select {
			case <-time.After(retryDelay):
			case <-ctx.Done():
				log.Printf("Timeout checking URL %s", url)
				return
			}
		}
	}

	sendWechatMessage(fmt.Sprintf("Service at %s is unhealthy after %d attempts\n", url, maxRetries), CFG.EWeChat.Receivers)

}

func CheckTcpIsHealth(ctx context.Context, host string, port int) {
	maxRetries := 3
	retryDelay := time.Second * 5
	address := net.JoinHostPort(host, strconv.Itoa(port))

	for i := 0; i < maxRetries; i++ {
		var dialer net.Dialer
		conn, err := dialer.DialContext(ctx, "tcp", address)
		if err == nil {
			log.Printf("TCP Service at %s is healthy\n", address)
			conn.Close()
			return
		}

		if i < maxRetries-1 {
			log.Printf("TCP Service %s unhealthy. Retrying in %v...\n", address, retryDelay)
			select {
			case <-time.After(retryDelay):
			case <-ctx.Done():
				log.Printf("Timeout checking TCP %s", address)
				return
			}
		}
	}

	sendWechatMessage(fmt.Sprintf("TCP Service at %s is unhealthy after %d attempts\n", address, maxRetries), CFG.EWeChat.Receivers)
}

func CheckPidIsHealth(ctx context.Context, pidFile string) {
	maxRetries := 3
	retryDelay := time.Second * 5

	for i := 0; i < maxRetries; i++ {
		running, err := check_pid.CheckHealth(ctx, pidFile)
		if err == nil && running {
			log.Printf("Process with PID file %s is running\n", pidFile)
			return
		}

		if i < maxRetries-1 {
			if err != nil {
				log.Printf("Failed to check PID file %s: %v. Retrying in %v...\n", pidFile, err, retryDelay)
			} else {
				log.Printf("Process with PID file %s is not running. Retrying in %v...\n", pidFile, retryDelay)
			}
			select {
			case <-time.After(retryDelay):
			case <-ctx.Done():
				log.Printf("Timeout checking PID %s", pidFile)
				return
			}
		}
	}

	sendWechatMessage(fmt.Sprintf("Process with PID file %s is not running or accessible after %d attempts\n", pidFile, maxRetries), CFG.EWeChat.Receivers)

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
		sendWechatMessage(fmt.Sprintf("%s disk read error: %s", msg_prefix, err.Error()), CFG.EWeChat.Receivers)
	}

	if disk.UsedPercent > CFG.DiskUsageRate {

		msg := fmt.Sprintf("%s Warning: Disk usage rate is %.2f%% and over DiskUsageRate %.2f%%", msg_prefix, disk.UsedPercent, CFG.DiskUsageRate)
		log.Println(msg)
		sendWechatMessage(msg, CFG.EWeChat.Receivers)

	}

	cpu, err := InitCPU()
	if err != nil {
		sendWechatMessage(fmt.Sprintf("%s cpu read error: %s", msg_prefix, err.Error()), CFG.EWeChat.Receivers)
	}

	if cpu.Cpus[0] > CFG.CpuUsageRate {

		msg := fmt.Sprintf("%s Warning: CPU usage rate is %.2f%% and over CpuUsageRate %.2f%%\n", msg_prefix, cpu.Cpus[0], CFG.CpuUsageRate)

		processInfos, err := InitProcess()
		if err != nil {
			sendWechatMessage(fmt.Sprintf("%s process read error: %s", msg_prefix, err.Error()), CFG.EWeChat.Receivers)
		}
		for _, p := range processInfos {
			msg += fmt.Sprintf("ProcessInfo: %s - %.2f%%\n", p.Exe, p.CPUPercent)
		}
		log.Println(msg)

		sendWechatMessage(msg, CFG.EWeChat.Receivers)
	}

	ram, err := InitRAM()
	if err != nil {
		sendWechatMessage(fmt.Sprintf("%s ram read error: %s", msg_prefix, err.Error()), CFG.EWeChat.Receivers)
	}

	if ram.UsedPercent > CFG.MemUsageRate {

		msg := fmt.Sprintf("%s Warning: Ram usage rate is %.2f%% and over MemUsageRate %.2f%%", msg_prefix, ram.UsedPercent, CFG.MemUsageRate)
		log.Println(msg)
		sendWechatMessage(msg, CFG.EWeChat.Receivers)

	}

	//url 健康检查
	var wg sync.WaitGroup
	sem := make(chan struct{}, 10)

	for _, url := range CFG.CheckUrl {
		u := url
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()
			CheckUrlIsHealth(ctx, u.Url, u.Method, u.ExpectStatusCode, u.ExpectBody)
		}()
	}

	//tcp 健康检查
	for _, tcpConfig := range CFG.CheckTcp {
		t := tcpConfig
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			CheckTcpIsHealth(ctx, t.Host, t.Port)
		}()
	}

	//pid 功能检查
	for _, pidConfig := range CFG.CheckPid {
		p := pidConfig
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			CheckPidIsHealth(ctx, p.Pid)
		}()
	}

	wg.Wait()

}

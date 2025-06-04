package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type ProcessInfo struct {
	PID        string `json:"pid"`
	Executable string `json:"executable"`
	Protocol   string `json:"protocol"`
	Port       int    `json:"port"`
}

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
)

func FindPortUnix(ports []int, jsonOutput bool, useColor bool) {
	found := false
	var results []ProcessInfo

	for _, protoFile := range []string{"/proc/net/tcp", "/proc/net/tcp6", "/proc/net/udp", "/proc/net/udp6"} {
		procsFile, err := os.ReadFile(protoFile)
		if err != nil {
			continue
		}

		protocol := "TCP"
		if strings.Contains(protoFile, "udp") {
			protocol = "UDP"
		}

		lines := strings.Split(string(procsFile), "\n")
		for i, line := range lines {
			if i == 0 {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) < 10 {
				continue
			}
			localAddrPort := fields[1]
			tcpInode := fields[9]

			parts := strings.Split(localAddrPort, ":")
			if len(parts) != 2 {
				continue
			}
			hexPort := parts[1]

			num, err := strconv.ParseInt(hexPort, 16, 64)
			if err != nil {
				continue
			}
			port := int(num)

			for _, targetPort := range ports {
				if port == targetPort {
					pid := findPIDBySocketInode(tcpInode)
					if pid == "" {
						continue
					}

					exePath := filepath.Join("/proc", pid, "exe")
					exeName, _ := os.Readlink(exePath)
					if exeName == "" {
						exeName = "unknown"
					}

					results = append(results, ProcessInfo{
						PID:        pid,
						Executable: exeName,
						Protocol:   protocol,
						Port:       port,
					})
					found = true
				}
			}
		}
	}

	if found && jsonOutput {
		jsonData, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(jsonData))
	} else if found {
		printResults(results, useColor)
	} else {
		fmt.Println("No process found using any of the specified ports")
	}
}

func FindPortWindows(ports []int, jsonOutput bool, useColor bool) {
	found := false
	var results []ProcessInfo

	cmd := exec.Command("netstat", "-ano")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create stdout pipe:", err)
		return
	}
	defer stdout.Close()

	if err := cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to start netstat:", err)
		return
	}

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		addrPort := fields[1]
		parts := strings.Split(addrPort, ":")
		if len(parts) < 2 {
			continue
		}
		portStr := parts[len(parts)-1]
		port, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}

		match := false
		for _, targetPort := range ports {
			if port == targetPort {
				match = true
				break
			}
		}
		if !match {
			continue
		}

		proto := strings.ToUpper(strings.Split(fields[0], "/")[0])
		if proto != "TCP" && proto != "UDP" {
			continue
		}

		var pid int
		if proto == "TCP" && len(fields) >= 5 && (fields[3] == "LISTENING" || strings.Contains(fields[3], "ESTABLISHED")) {
			pid, _ = strconv.Atoi(fields[4])
			found = true
		} else if proto == "UDP" && len(fields) >= 4 {
			pid, _ = strconv.Atoi(fields[3])
			found = true
		}

		if pid > 0 {
			exeName := getExeNameFromPIDWindows(pid)
			results = append(results, ProcessInfo{
				PID:        strconv.Itoa(pid),
				Executable: exeName,
				Protocol:   proto,
				Port:       port,
			})
		}
	}

	if found && jsonOutput {
		jsonData, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(jsonData))
	} else if found {
		printResults(results, useColor)
	} else {
		fmt.Println("No process found using any of the specified ports")
	}
}

func printResults(results []ProcessInfo, useColor bool) {
	if useColor {
		fmt.Printf("%s%-6s%s    %s%-6s%s    %s%-20s%s    %s%-5s%s\n",
			colorCyan, "Proto", colorReset,
			colorYellow, "PID", colorReset,
			colorGreen, "Executable", colorReset,
			colorBlue, "Port", colorReset)
		fmt.Println(strings.Repeat("-", 60))
	} else {
		fmt.Printf("%-6s    %-6s    %-20s    %-5s\n", "Proto", "PID", "Executable", "Port")
		fmt.Println(strings.Repeat("-", 60))
	}

	for _, res := range results {
		if useColor {
			fmt.Printf("%s%-6s%s    %s%-6s%s    %s%-20s%s    %s%-5d%s\n",
				colorCyan, res.Protocol, colorReset,
				colorYellow, res.PID, colorReset,
				colorGreen, res.Executable, colorReset,
				colorBlue, res.Port, colorReset)
		} else {
			fmt.Printf("%-6s    %-6s    %-20s    %-5d\n", res.Protocol, res.PID, res.Executable, res.Port)
		}
	}
}

func findPIDBySocketInode(inode string) string {
	pids, _ := filepath.Glob("/proc/[0-9]*")
	for _, pidDir := range pids {
		pid := filepath.Base(pidDir)

		fdDir := filepath.Join(pidDir, "fd")
		fds, err := os.ReadDir(fdDir)
		if err != nil {
			continue
		}

		for _, fd := range fds {
			fdPath := filepath.Join(fdDir, fd.Name())
			target, err := os.Readlink(fdPath)
			if err != nil {
				continue
			}

			if target == "socket:["+inode+"]" {
				return pid
			}
		}
	}
	return ""
}

func getExeNameFromPIDWindows(pid int) string {
	cmd := exec.Command("wmic", "process", "where", "processid="+strconv.Itoa(pid), "get", "executablepath")
	out, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "ExecutablePath") {
			return line
		}
	}
	return "unknown"
}

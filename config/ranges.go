package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func getConfigPath() string {
	exe, err := os.Executable()
	if err != nil {
		return "config/ip_ranges.txt"
	}
	return filepath.Join(filepath.Dir(exe), "config", "ip_ranges.txt")
}

func loadRanges(filename string) ([]string, error) {
	filePath := filepath.Join(getConfigDir(), filename)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open %s: %v", filePath, err)
	}
	defer file.Close()

	var ranges []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		ranges = append(ranges, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading %s: %v", filePath, err)
	}
	return ranges, nil
}

func getConfigDir() string {
	// Try current directory first
	if _, err := os.Stat("ipv4.txt"); err == nil {
		return "."
	}
	// Try config subdirectory
	if _, err := os.Stat("config/ipv4.txt"); err == nil {
		return "config"
	}
	// Fall back to executable directory
	exe, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(exe)
}

func LoadIPv4Ranges() ([]string, error) {
	return loadRanges("ipv4.txt")
}

func LoadIPv6Ranges() ([]string, error) {
	return loadRanges("ipv6.txt")
}

func LoadAllRanges() []string {
	exe, err := os.Executable()
	if err != nil {
		return nil
	}
	filePath := filepath.Join(filepath.Dir(exe), "config", "ip_ranges.txt")

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not open %s: %v\n", filePath, err)
		return nil
	}
	defer file.Close()

	var ranges []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		ranges = append(ranges, line)
	}
	return ranges
}

func GetIPv4Count(ranges []string) int64 {
	return getRangeCount(ranges)
}

func GetIPv6Count(ranges []string) int64 {
	return getRangeCount(ranges)
}

func getRangeCount(ranges []string) int64 {
	var total int64
	for _, r := range ranges {
		r = strings.TrimSpace(r)
		if r == "" {
			continue
		}
		parts := strings.Split(r, "/")
		if len(parts) != 2 {
			continue
		}
		var prefixLen int
		fmt.Sscanf(parts[1], "%d", &prefixLen)
		if strings.Contains(parts[0], ".") {
			total += 1 << uint(32-prefixLen)
		} else {
			total += 1 << uint(128-prefixLen)
		}
	}
	return total
}

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/surfrpt1/clean-ip-scanner/config"
	"github.com/surfrpt1/clean-ip-scanner/scanner"
	"github.com/surfrpt1/clean-ip-scanner/utils"
	"github.com/fatih/color"
)

const (
	defaultWorkers     = 10
	defaultTestCount   = 20
	defaultPingTimes   = 2
	defaultDownloadURL = "https://speed.cloudflare.com/__down?bytes=50000000"
	defaultUploadURL   = "https://speed.cloudflare.com/__up?bytes=25000000"
	defaultCPFile      = "./checkpoint.json"
)

func main() {
	var (
		scanMode    int
		workers     int
		testCount   int
		seed        int64
		save        bool
		resetCP     bool
		xrayMode    bool
		lightMode   bool
	)

	flag.IntVar(&scanMode, "mode", 0, "Scan mode: 0=normal, 1=download only, 2=upload only")
	flag.IntVar(&workers, "workers", defaultWorkers, "Number of concurrent workers")
	flag.IntVar(&testCount, "test-count", defaultTestCount, "Number of IPs to speed test")
	flag.IntVar(&testCount, "n", defaultTestCount, "Number of IPs to speed test (shorthand)")
	flag.Int64Var(&seed, "seed", 0, "Random seed (0 = current time)")
	flag.BoolVar(&save, "save", false, "Save results to file")
	flag.BoolVar(&resetCP, "reset", false, "Reset checkpoint and start fresh")
	flag.BoolVar(&xrayMode, "xray", false, "Use Xray mode (requires xray binary at ./xray/xray)")
	flag.BoolVar(&lightMode, "light", false, "Light mode: skip speed tests, only ping (recommended for iPad)")
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())
	scanner.OverridePingTimes(defaultPingTimes)

	utils.PrintDesigner()
	utils.PrintHeader()

	// Show mode info
	fmt.Println()
	if xrayMode {
		color.New(color.FgMagenta).Println("Mode: Xray (Protocol-aware)")
		if err := scanner.ValidateXrayConfig(); err != nil {
			color.New(color.FgRed).Printf("Config error: %v\n", err)
			color.New(color.FgYellow).Println("Please edit config/xray_config.txt (URL) or config/xray_config.json (JSON)")
			return
		}
		color.New(color.FgGreen).Println("Config validated OK")
	} else {
		color.New(color.FgCyan).Println("Mode: Normal (ICMP Ping)")
	}

	modeName := "Normal"
	switch scanMode {
	case 0:
		modeName = "Full (Ping + Download + Upload)"
	case 1:
		modeName = "Download only"
	case 2:
		modeName = "Upload only"
	}
	color.New(color.FgCyan).Printf("Test Mode: %s | Workers: %d | Test Count: %d\n\n", modeName, workers, testCount)

	// Checkpoint handling
	cp := scanner.LoadCheckpoint(defaultCPFile)
	if cp != nil && resetCP {
		os.Remove(defaultCPFile)
		cp = nil
		color.New(color.FgYellow).Println("Checkpoint reset.")
	}
	if cp != nil {
		color.New(color.FgGreen).Printf("Checkpoint found! Resuming from index %d\n", cp.ProgressIndex)
	}

	// Graceful shutdown
	stopCh := make(chan struct{})
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println()
		color.New(color.FgYellow).Println("Received interrupt, saving checkpoint...")
		close(stopCh)
	}()

	// Load IPs
	color.New(color.FgCyan).Println("Loading IP ranges...")
	ipv4, err := config.LoadIPv4Ranges()
	if err != nil {
		color.New(color.FgRed).Printf("Error loading IPv4 ranges: %v\n", err)
		return
	}
	ipv6, err := config.LoadIPv6Ranges()
	if err != nil {
		color.New(color.FgRed).Printf("Error loading IPv6 ranges: %v\n", err)
		return
	}
	ipv4Count := config.GetIPv4Count(ipv4)
	ipv6Count := config.GetIPv6Count(ipv6)
	totalCount := ipv4Count + ipv6Count

	color.New(color.FgGreen).Printf("Loaded %d IPv4 IPs, %d IPv6 IPs (Total: %d)\n",
		ipv4Count, ipv6Count, totalCount)

	var ips []scanner.CompactIP
	if seed != 0 {
		ips = scanner.GenerateIPsWithSeed(ipv4, ipv6, seed)
	} else {
		ips = scanner.GenerateIPsWithSeed(ipv4, ipv6, time.Now().UnixNano())
	}

	if cp != nil {
		startIdx := cp.ProgressIndex
		if startIdx >= len(ips) {
			color.New(color.FgRed).Println("Checkpoint index exceeds total IPs, starting fresh.")
			startIdx = 0
		}
		ips = ips[startIdx:]
		color.New(color.FgYellow).Printf("Skipped first %d IPs, testing remaining %d\n", startIdx, len(ips))
	}

	// Run scan
	var pingResults []scanner.PingResult
	if cp != nil && len(cp.PingResults) > 0 {
		pingResults = cp.GetPingResults()
		color.New(color.FgGreen).Printf("Loaded %d previous ping results from checkpoint\n", len(pingResults))
	} else {
		color.New(color.FgCyan).Printf("Starting scan of %d IPs...\n\n", len(ips))
		if xrayMode {
			pingResults = scanner.PingIPsViaXray(stopCh, ips, workers, cp, nil)
		} else {
			pingResults = scanner.PingIPsTCP(stopCh, ips, workers, cp, nil)
		}
	}

	// Speed test or light mode
	if lightMode {
		utils.PrintPingResults(pingResults, testCount)
		if save {
			utils.SavePingResults(pingResults, testCount, "results.txt")
			color.New(color.FgGreen).Println("Results saved to results.txt")
		}
	} else {
		switch scanMode {
		case 0:
			results := scanner.SpeedTest(stopCh, pingResults, defaultDownloadURL, defaultUploadURL, xrayMode, testCount)
			utils.PrintResults(results)
			if save {
				utils.SaveResults(results, "results.txt")
				utils.SaveSimpleResults(results, pingResults, "simple_results.txt")
				color.New(color.FgGreen).Println("Results saved to results.txt and simple_results.txt")
			}
		case 1:
			results := scanner.DownloadTest(stopCh, pingResults, defaultDownloadURL, xrayMode, testCount)
			utils.PrintResults(results)
			if save {
				utils.SaveResults(results, "results.txt")
				utils.SaveSimpleResults(results, pingResults, "simple_results.txt")
				color.New(color.FgGreen).Println("Results saved to results.txt and simple_results.txt")
			}
		case 2:
			results := scanner.UploadTest(stopCh, pingResults, defaultUploadURL, xrayMode, testCount)
			utils.PrintResults(results)
			if save {
				utils.SaveResults(results, "results.txt")
				utils.SaveSimpleResults(results, pingResults, "simple_results.txt")
				color.New(color.FgGreen).Println("Results saved to results.txt and simple_results.txt")
			}
		}
	}

	// Cleanup checkpoint
	if cp != nil {
		os.Remove(defaultCPFile)
	}
}

func init() {
	// Custom help
	flag.Usage = func() {
		utils.PrintHeader()
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  ./clean-ip-scanner [flags]")
		fmt.Println()
		fmt.Println("Flags:")
		fmt.Println("  -mode        Scan mode: 0=full (default), 1=download, 2=upload")
		fmt.Println("  -workers     Number of concurrent workers (default: 200)")
		fmt.Println("  -test-count  Number of IPs to speed test (default: 5)")
		fmt.Println("  -n           Shorthand for -test-count")
		fmt.Println("  -seed        Random seed for IP shuffling (0=current time)")
		fmt.Println("  -save        Save results to files")
		fmt.Println("  -reset       Reset checkpoint and start fresh")
		fmt.Println("  -xray        Use Xray mode (requires ./xray/xray binary)")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  ./clean-ip-scanner -save")
		fmt.Println("  ./clean-ip-scanner -mode 1 -n 10 -save")
		fmt.Println("  ./clean-ip-scanner -xray -save")
		fmt.Println("  ./clean-ip-scanner -reset -workers 100")
	}
}

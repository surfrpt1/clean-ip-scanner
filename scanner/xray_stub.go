//go:build !linux

package scanner

import "fmt"

func SetXrayOverridePort(p int)                              {}
func ValidateXrayConfig() error                              { return nil }
func PingIPsViaXray(stopCh <-chan struct{}, ips []CompactIP, workers int, cp *Checkpoint, existingResults []PingResult) []PingResult {
	fmt.Println("Xray mode only available on Linux")
	return nil
}

package scanner

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/fatih/color"
)

type SimpleBar struct {
	total     int
	current   int64
	found     int64
	startTime time.Time
	interval  int
	lastFound int
}

func newBar(total int, label, _ string) *SimpleBar {
	b := &SimpleBar{
		total:     total,
		startTime: time.Now(),
		interval:  25,
	}
	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Printf("Scanning %d IPs...\n\n", total)
	return b
}

func (b *SimpleBar) grow(num int, foundStr string) {
	newVal := atomic.AddInt64(&b.current, int64(num))
	if foundStr != "" {
		var f int64
		fmt.Sscanf(foundStr, "%d", &f)
		atomic.StoreInt64(&b.found, f)
	}
	current := int(newVal)
	found := int(atomic.LoadInt64(&b.found))

	if current%b.interval == 0 || current == b.total {
		elapsed := time.Since(b.startTime)
		percent := float64(current) / float64(b.total) * 100
		speed := float64(current) / elapsed.Seconds()
		remaining := 0.0
		if speed > 0 {
			remaining = float64(b.total-current) / speed
		}

		barWidth := 25
		filled := int(float64(barWidth) * float64(current) / float64(b.total))
		if filled > barWidth {
			filled = barWidth
		}
		if filled < 0 {
			filled = 0
		}
		bar := strings.Repeat("#", filled) + strings.Repeat("-", barWidth-filled)

		fmt.Printf("  [%s] %5.1f%%  %d/%d  |  Clean: %d  |  %.0f/sec  |  ETA: %ds\n",
			bar, percent, current, b.total, found, speed, int(remaining))
	}
}

func (b *SimpleBar) done() {
	elapsed := time.Since(b.startTime)
	green := color.New(color.FgGreen, color.Bold)
	fmt.Println()
	green.Printf("Done: %d tested, %d clean IPs found in %s\n\n",
		atomic.LoadInt64(&b.current), atomic.LoadInt64(&b.found),
		elapsed.Round(time.Second))
}

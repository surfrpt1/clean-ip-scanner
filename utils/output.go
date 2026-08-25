package utils

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/surfrpt1/clean-ip-scanner/scanner"
	"github.com/fatih/color"
)

const tableWidth = 66

var gradientStart = [3]int{57, 255, 20}
var gradientEnd = [3]int{255, 230, 0}

func rankColor(rank, total int) *color.Color {
	if total <= 1 {
		return color.RGB(gradientStart[0], gradientStart[1], gradientStart[2]).Add(color.Bold)
	}
	t := float64(rank) / float64(total-1)
	r := gradientStart[0] + int(float64(gradientEnd[0]-gradientStart[0])*t+0.5)
	g := gradientStart[1] + int(float64(gradientEnd[1]-gradientStart[1])*t+0.5)
	b := gradientStart[2] + int(float64(gradientEnd[2]-gradientStart[2])*t+0.5)
	return color.RGB(r, g, b).Add(color.Bold)
}

func PrintResults(results []scanner.IPResult) {
	fmt.Println()
	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println(strings.Repeat("=", tableWidth))
	cyan.Println(CenterText(tableWidth, "CLEAN IPs FOUND"))
	cyan.Println(strings.Repeat("=", tableWidth))
	fmt.Println()

	green := color.New(color.FgGreen, color.Bold)

	green.Printf("%-4s%-20s%-6s%-6s%-8s%-11s%-11s\n",
		"#", "IP", "S/R", "Loss", "Delay", "DL", "UL")
	cyan.Println(strings.Repeat("-", tableWidth))

	total := len(results)
	for i, r := range results {
		rank := fmt.Sprintf("%d.", i+1)
		sr := fmt.Sprintf("%d/%d", r.Sended, r.Received)
		loss := fmt.Sprintf("%.2f", r.LossRate)
		delay := fmt.Sprintf("%dms", r.Delay)
		downloadSpeed := fmt.Sprintf("%.2fMB/s", r.DownloadSpeed/1024/1024)
		uploadSpeed := fmt.Sprintf("%.2fMB/s", r.UploadSpeed/1024/1024)

		rankColor(i, total).Printf("%-4s%-20s%-6s%-6s%-8s%-11s%-11s\n",
			rank, r.IP.String(), sr, loss, delay, downloadSpeed, uploadSpeed)
	}

	cyan.Println(strings.Repeat("=", tableWidth))
}

func SaveResults(results []scanner.IPResult, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString("# Clean IPs\n")
	file.WriteString(fmt.Sprintf("# Generated at: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	file.WriteString(fmt.Sprintf("# Total IPs found: %d\n", len(results)))
	file.WriteString("#\n")
	file.WriteString("# Format: Rank | IP | Sent | Received | Loss | Avg Delay | Download Speed | Upload Speed\n")
	file.WriteString("#===========================================================================\n\n")

	for i, r := range results {
		line := fmt.Sprintf("%d. %s | Sent: %d | Recv: %d | Loss: %.2f | %dms | %.2f MB/s | %.2f MB/s\n",
			i+1,
			r.IP.String(),
			r.Sended,
			r.Received,
			r.LossRate,
			r.Delay,
			r.DownloadSpeed/1024/1024,
			r.UploadSpeed/1024/1024,
		)
		file.WriteString(line)
	}

	file.WriteString("\n# End of results\n")
	return nil
}

func SaveSimpleResults(topResults []scanner.IPResult, pingResults []scanner.PingResult, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, r := range topResults {
		file.WriteString(r.IP.String() + "\n")
	}

	file.WriteString("--------------------\n")

	for _, r := range pingResults {
		file.WriteString(r.IP.String() + "\n")
	}

	return nil
}

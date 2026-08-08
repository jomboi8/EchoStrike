package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"echostrike/internal/sender"
	"echostrike/pkg/syslog"

	"github.com/spf13/cobra"
)

var filePath string
var preserveTiming bool

var replayCmd = &cobra.Command{
	Use:   "replay",
	Short: "Replay logs from a file",
	Long:  `Read log lines from a file and send them to the syslog target.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize Sender
		var proto sender.Protocol = sender.UDP
		if protocol == "tcp" {
			proto = sender.TCP
		}
		if protocol == "tls" {
			proto = sender.TLS
		}

		fac, err := syslog.ParseFacility(strings.ToLower(facility))
		if err != nil {
			fmt.Printf("Error parsing facility: %v\n", err)
			os.Exit(1)
		}
		sev, err := syslog.ParseSeverity(strings.ToLower(severity))
		if err != nil {
			fmt.Printf("Error parsing severity: %v\n", err)
			os.Exit(1)
		}

		s, err := sender.NewSender(proto, host, port)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		defer s.Close()

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		count := 0
		fmt.Printf("Replaying logs from %s to %s:%d...\n", filePath, host, port)

		for scanner.Scan() {
			line := scanner.Text()

			if preserveTiming {
				time.Sleep(10 * time.Millisecond) // Mock "timing"
			}


			msg := syslog.NewMessage(line)
			msg.AppName = tag
			msg.Facility = fac
			msg.Severity = sev

			if err := s.Send(msg.String() + "\n"); err != nil {
				fmt.Printf("Error sending line %d: %v\n", count+1, err)
				continue
			}
			count++
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading file: %v\n", err)
		}

		fmt.Printf("Replay complete. Sent %d messages.\n", count)
	},
}

func init() {
	rootCmd.AddCommand(replayCmd)

	replayCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to log file")
	replayCmd.MarkFlagRequired("file")

	replayCmd.Flags().BoolVar(&preserveTiming, "preserve-timing", false, "Simulate original timing (mock)")

	replayCmd.Flags().StringVar(&host, "host", "127.0.0.1", "Target IP")
	replayCmd.Flags().IntVar(&port, "port", 514, "Target Port")
	replayCmd.Flags().StringVar(&protocol, "proto", "udp", "Protocol")
	replayCmd.Flags().StringVarP(&tag, "tag", "t", "replay", "Tag")
	replayCmd.Flags().StringVar(&facility, "facility", "local0", "Facility")
	replayCmd.Flags().StringVar(&severity, "severity", "info", "Severity")
}

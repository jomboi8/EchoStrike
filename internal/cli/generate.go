package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"echostrike/internal/generator"
	"echostrike/internal/ratelimiter"
	"echostrike/internal/sender"
	"echostrike/pkg/syslog"

	"github.com/spf13/cobra"
)

var (
	rate         int
	duration     time.Duration
	templateName string
	workerCount  int
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate high-volume syslog traffic from templates",
	Long: `Generate realistic syslog traffic using built-in templates.
Sends are fanned out across a bounded worker pool, each with its own
long-lived connection and throttled by a shared rate limiter so the
aggregate send rate matches --rate regardless of worker count.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Validate protocol
		var proto sender.Protocol
		switch strings.ToLower(protocol) {
		case "tcp":
			proto = sender.TCP
		case "udp":
			proto = sender.UDP
		case "tls":
			proto = sender.TLS
		default:
			fmt.Printf("Invalid protocol: %s\n", protocol)
			os.Exit(1)
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
		msgFormat := syslog.RFC3164
		if strings.ToLower(format) == "rfc5424" {
			msgFormat = syslog.RFC5424
		}

		gen := generator.New()
		if _, err := gen.Generate(templateName); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		workers := workerCount
		if workers <= 0 {
			workers = max(min(rate, 32), 1)
		}

		fmt.Printf("Starting generation: Template=%s Target=%s:%d Rate=%d/s Duration=%s Workers=%d\n",
			templateName, host, port, rate, duration, workers)

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()
		ctx, cancel := context.WithTimeout(ctx, duration)
		defer cancel()

		limiter := ratelimiter.New(rate)
		defer limiter.Stop()

		var sent atomic.Int64
		var failed atomic.Int64

		var wg sync.WaitGroup
		for range workers {
			s, err := sender.NewSender(proto, host, port)
			if err != nil {
				fmt.Printf("Error creating sender: %v\n", err)
				os.Exit(1)
			}

			wg.Add(1)
			go func(s *sender.Sender) {
				defer wg.Done()
				defer s.Close()

				for {
					if err := limiter.Wait(ctx); err != nil {
						// Context cancelled (duration elapsed or Ctrl+C).
						return
					}

					logMsg, err := gen.Generate(templateName)
					if err != nil {
						failed.Add(1)
						continue
					}

					msg := syslog.NewMessage(logMsg)
					msg.AppName = tag
					msg.Format = msgFormat
					msg.Facility = fac
					msg.Severity = sev

					if err := s.Send(msg.String() + "\n"); err != nil {
						failed.Add(1)
						continue
					}
					sent.Add(1)
				}
			}(s)
		}

		wg.Wait()

		fmt.Printf("\nCompleted. Sent %d messages", sent.Load())
		if n := failed.Load(); n > 0 {
			fmt.Printf(" (%d failed)", n)
		}
		fmt.Println(".")
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVar(&host, "host", "127.0.0.1", "Target IP/Hostname")
	generateCmd.Flags().IntVar(&port, "port", 514, "Target Port")
	generateCmd.Flags().StringVar(&protocol, "proto", "udp", "Protocol (tcp, udp, tls)")
	generateCmd.Flags().StringVarP(&templateName, "template", "T", "ssh-failed", "Log template name")
	generateCmd.Flags().IntVarP(&rate, "rate", "r", 1, "Aggregate logs per second across all workers")
	generateCmd.Flags().DurationVarP(&duration, "duration", "d", 10*time.Second, "Duration to run")
	generateCmd.Flags().IntVarP(&workerCount, "workers", "w", 0, "Concurrent sender workers (0 = auto, min(rate, 32))")

	generateCmd.Flags().StringVarP(&tag, "tag", "t", "echostrike", "Syslog tag/app-name")
	generateCmd.Flags().StringVar(&format, "format", "rfc3164", "Syslog format")
	generateCmd.Flags().StringVar(&facility, "facility", "local0", "Facility")
	generateCmd.Flags().StringVar(&severity, "severity", "info", "Severity")
}

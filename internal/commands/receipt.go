package commands

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"time"

	"github.com/google/subcommands"
	"github.com/swfrench/pushover-tool/internal/api/receipt"
)

// Receipt is a Command for interacting with the Receipt API.
type Receipt struct {
	receipt  string
	interval time.Duration
}

func (*Receipt) Name() string { return "receipt" }
func (*Receipt) Synopsis() string {
	return "Await message acknowledgement given a receipt."
}
func (*Receipt) Usage() string {
	return "receipt [-receipt <ID>] [-interval <duration>]\n"
}

func (m *Receipt) SetFlags(f *flag.FlagSet) {
	f.StringVar(&m.receipt, "receipt", "", "The receipt ID to wait upon. If empty / unset, reads from stdin.")
	f.DurationVar(&m.interval, "interval", 30*time.Second, "Interval between API calls while waiting for acknowledgement.")
}

func (m *Receipt) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	if m.receipt == "" {
		bs, _ := io.ReadAll(os.Stdin)
		m.receipt = string(bs)
	}
	client, err := receipt.NewClient(&receipt.Options{
		Token:    args[0].(string),
		Interval: m.interval,
	})
	if err != nil {
		log.Fatalf("Receipt Client initialization failed: %v", err)
	}
	if err := client.Wait(ctx, m.receipt); err != nil {
		log.Fatalf("Failed to verify message receipt: %v", err)
	}
	return subcommands.ExitSuccess
}

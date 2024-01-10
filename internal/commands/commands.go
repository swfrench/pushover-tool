package commands

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"time"

	"github.com/google/subcommands"
	"github.com/swfrench/pushover-tool/internal/api/message"
)

// Message is a Command for interacting with the Message API.
type Message struct {
	user      string
	message   string
	title     string
	emergency bool
	retry     time.Duration
	expire    time.Duration
}

func (*Message) Name() string     { return "message" }
func (*Message) Synopsis() string { return "Send a message to a user/group." }
func (*Message) Usage() string {
	return "message -user <key> [-message <content>] [-title <title>]\n"
}

func (m *Message) SetFlags(f *flag.FlagSet) {
	f.StringVar(&m.user, "user", "", "The user/group key to message (required).")
	f.StringVar(&m.message, "message", "", "The message to send. If empty / unset, reads from stdin.")
	f.StringVar(&m.title, "title", "pushover-tool", "The message title to send.")
	f.BoolVar(&m.emergency, "emergency", false, "Whether to send the message a emergency priority. If successfully send, the receipt ID will be printed to stdout.")
	f.DurationVar(&m.retry, "emergency_retry", 10*time.Minute, "Retry period for un-ACK'd emergency messages.")
	f.DurationVar(&m.expire, "emergency_expire", time.Hour, "Expiration age for un-ACK'd emergency messages.")
}

func (m *Message) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	if m.user == "" {
		log.Fatal("No recipient user specified")
	}
	if m.message == "" {
		bs, _ := io.ReadAll(os.Stdin)
		m.message = string(bs)
	}
	client, err := message.NewClient(&message.Options{Token: args[0].(string)})
	if err != nil {
		log.Fatalf("Message Client initialization failed: %v", err)
	}
	if err := client.Send(m.user, m.message, m.title, m.emergency, m.retry, m.expire); err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
	return subcommands.ExitSuccess
}

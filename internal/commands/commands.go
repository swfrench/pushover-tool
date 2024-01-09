package commands

import (
	"context"
	"flag"
	"io"
	"log"
	"os"

	"github.com/google/subcommands"
	"github.com/swfrench/pushover-tool/internal/api/message"
)

// TODO: Support variable message priority, printing a the receipt response
// parameter to stdout in the emergency priority case.

// Message is a Command for interacting with the Message API.
type Message struct {
	user    string
	message string
	title   string
}

func (*Message) Name() string     { return "message" }
func (*Message) Synopsis() string { return "Send a message to a user/group." }
func (*Message) Usage() string {
	return "message -user <key> [-message <content>] [-title <title>]\n"
}

func (m *Message) SetFlags(f *flag.FlagSet) {
	f.StringVar(&m.user, "user", "", "The user/group key to message.")
	f.StringVar(&m.message, "message", "", "The message to send. If empty / unset, reads from stdin.")
	f.StringVar(&m.title, "title", "", "The message title to send. If empty / unset, none is provided to the API.")
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
	title := &m.title
	if m.title == "" {
		title = nil
	}
	if err := client.Send(m.user, m.message, title); err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
	return subcommands.ExitSuccess
}

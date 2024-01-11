// pushover-tool is a tool for interacting with the Pushover API.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	"github.com/google/subcommands"
	"github.com/swfrench/pushover-tool/internal/commands"
)

var (
	tokenPath = flag.String("token_path", "", "Path to the Pushover API token file (i.e., app token). Note that environment expansion is applied to the path.")
	timeout   = flag.Duration("timeout", time.Hour, "The timeout applied to the operation underlying the selected subcommand.")
)

type tokenFile struct {
	Token string `json:"token"`
}

func mustReadToken() string {
	bs, err := os.ReadFile(*tokenPath)
	if err != nil {
		log.Fatalf("Failed to read token file (%q): %v", *tokenPath, err)
	}
	var tf tokenFile
	if err := json.Unmarshal(bs, &tf); err != nil {
		log.Fatalf("Failed to unmarshal token file (%q): %v", *tokenPath, err)
	}
	return tf.Token
}

type wrapper struct {
	subcommands.Command
}

func (w *wrapper) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	return w.Command.Execute(ctx, f, mustReadToken())
}

func withToken(c subcommands.Command) *wrapper {
	return &wrapper{c}
}

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(withToken(new(commands.Message)), "")
	subcommands.Register(withToken(new(commands.Receipt)), "")

	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	os.Exit(int(subcommands.Execute(ctx)))
}

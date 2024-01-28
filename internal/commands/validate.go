package commands

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/google/subcommands"
	"github.com/swfrench/pushover-tool/internal/api/validate"
)

// Validate is a Command for interacting with the Validate API.
type Validate struct {
	user string
}

func (*Validate) Name() string     { return "validate" }
func (*Validate) Synopsis() string { return "Validate a user/group key (printed to stdout)." }
func (*Validate) Usage() string    { return "validate -user <key>\n" }

func (v *Validate) SetFlags(f *flag.FlagSet) {
	f.StringVar(&v.user, "user", "", "The user/group key to validate (required).")
}

func (v *Validate) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	if v.user == "" {
		log.Fatal("No user/group key specified")
	}
	client, err := validate.NewClient(&validate.Options{Token: args[0].(string)})
	if err != nil {
		log.Fatalf("Validate Client initialization failed: %v", err)
	}
	status, err := client.Check(ctx, v.user)
	if err != nil {
		log.Fatalf("Failed to validate user/group key: %v", err)
	}
	fmt.Println(status.String())
	return subcommands.ExitSuccess
}

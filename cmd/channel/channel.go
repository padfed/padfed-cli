package channel

import (
	"github.com/spf13/cobra"

	"github.com/padfed/padfed-cli/cmd/channel/dump"
	"github.com/padfed/padfed-cli/cmd/channel/listen"
)

func Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "channel",
		Short:   "Interact with a channel",
		Aliases: []string{"ch"},
	}
	c.PersistentFlags().String("channel", "padfedchannel", "Channel name")
	c.AddCommand(dump.Command())
	c.AddCommand(listen.Command())
	return c
}

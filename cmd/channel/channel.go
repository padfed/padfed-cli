package channel

import (
	"github.com/padfed/padfed-cli/cmd/channel/dump"
	"github.com/padfed/padfed-cli/cmd/channel/listen"
	client "github.com/padfed/padfed-cli/fabric-client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	c := &cobra.Command{
		Use:               "channel",
		Short:             "Interact with a channel",
		Aliases:           []string{"ch"},
		PersistentPreRunE: InitChannel,
	}
	c.AddCommand(dump.Command())
	c.AddCommand(listen.Command())
	return c
}

func InitChannel(cmd *cobra.Command, args []string) error {
	cli := viper.Get("fabric-client").(*client.Client)
	ch, err := cli.Channel(viper.GetString("fabric.channel"))
	if err != nil {
		return err
	}
	viper.Set("channel", ch)
	return nil
}

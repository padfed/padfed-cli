package query

import (
	"encoding/json"
	"os"

	"github.com/padfed/padfed-cli/app/helpers"
	client "github.com/padfed/padfed-cli/fabric-client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	repeat uint64
)

func Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "query <function> [<arg1>]...",
		Short:   "Call a function in query mode",
		Aliases: []string{"q"},
		Args:    cobra.MinimumNArgs(1),
		RunE:    Run,
	}
	c.Flags().Uint64VarP(&repeat, "repeat", "r", 1, "Number of times to repeat the task")
	return c
}

func Run(cmd *cobra.Command, args []string) error {
	chn, err := cmd.Flags().GetString("channel")
	if err != nil {
		return err
	}
	ccn, err := cmd.Flags().GetString("chaincode")
	if err != nil {
		return err
	}
	cli := viper.Get("fabric-client").(*client.Client)
	ch, err := cli.Channel(chn)
	if err != nil {
		return err
	}
	cc := ch.Chaincode(ccn)
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	for r := uint64(0); r < repeat; r++ {
		res := cc.Query(args[0], helpers.Bytes(1, args))
		if enc.Encode(res) != nil {
			return err
		}
	}
	return nil
}

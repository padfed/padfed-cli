package execute

import (
	"encoding/json"

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
		Use:     "execute <function> [<arg1>]...",
		Short:   "Call a function in execute mode",
		Aliases: []string{"ex"},
		Args:    cobra.MinimumNArgs(1),
		RunE:    Run,
	}
	c.Flags().Uint64VarP(&repeat, "repeat", "r", 1, "Number of times to repeat the task")
	return c
}

func Run(cmd *cobra.Command, args []string) error {
	cc := viper.Get("chaincode").(*client.Chaincode)
	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetEscapeHTML(false)
	for r := uint64(0); r < repeat; r++ {
		res := cc.Execute(args[0], helpers.Bytes(1, args))
		if err := enc.Encode(res); err != nil {
			return err
		}
	}
	return nil
}

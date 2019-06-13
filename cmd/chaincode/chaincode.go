package chaincode

import (
	"github.com/spf13/cobra"

	"github.com/padfed/padfed-cli/cmd/chaincode/execute"
	"github.com/padfed/padfed-cli/cmd/chaincode/query"
)

func Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "chaincode",
		Short:   "Interact with a chaincode",
		Aliases: []string{"cc"},
	}
	c.PersistentFlags().String("channel", "padfedchannel", "Channel name")
	c.PersistentFlags().String("chaincode", "padfedcc", "Chaincode name")
	c.AddCommand(execute.Command())
	c.AddCommand(query.Command())
	return c
}

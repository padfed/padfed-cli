package chaincode

import (
	"github.com/padfed/padfed-cli/cmd/chaincode/execute"
	"github.com/padfed/padfed-cli/cmd/chaincode/query"
	"github.com/padfed/padfed-cli/cmd/channel"
	client "github.com/padfed/padfed-cli/fabric-client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	c := &cobra.Command{
		Use:               "chaincode",
		Short:             "Interact with a chaincode",
		Aliases:           []string{"cc"},
		PersistentPreRunE: InitChaincode,
	}
	c.AddCommand(execute.Command())
	c.AddCommand(query.Command())
	return c
}

func InitChaincode(cmd *cobra.Command, args []string) error {
	if err := channel.InitChannel(cmd, args); err != nil {
		return err
	}
	ch := viper.Get("channel").(*client.Channel)
	viper.Set("chaincode", ch.Chaincode(viper.GetString("fabric.chaincode")))
	return nil
}

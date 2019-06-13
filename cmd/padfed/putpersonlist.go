package padfed

import (
	"encoding/json"
	"os"

	"github.com/padfed/padfed-cli/app/helpers"
	"github.com/padfed/padfed-cli/cmd/chaincode"
	client "github.com/padfed/padfed-cli/fabric-client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func PutPersonList() *cobra.Command {
	c := &cobra.Command{
		Use:     "putpersonlist <person-list-json>...",
		Short:   "Put one or more person lists",
		Aliases: []string{"putpl"},
		RunE:    PutPersonListRun,
		PreRunE: chaincode.InitChaincode,
	}
	return c
}

func PutPersonListRun(cmd *cobra.Command, args []string) error {
	cc := viper.Get("chaincode").(*client.Chaincode)
	exit := 0
	enc := json.NewEncoder(os.Stdout)
	for _, arg := range helpers.Bytes(0, args) {
		res, err := padfedresponse(cc.Execute("PutPersonaList", [][]byte{arg}))
		if enc.Encode(res) != nil {
			return err
		}
		if res.Status != 200 {
			exit = 1
		}
	}
	if exit != 0 {
		os.Exit(exit)
	}
	return nil
}

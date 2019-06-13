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

func GetPerson() *cobra.Command {
	c := &cobra.Command{
		Use:     "getperson <person-id> ...",
		Short:   "Get a person by id",
		Aliases: []string{"getp"},
		Args:    cobra.MinimumNArgs(1),
		RunE:    GetPersonRun,
		PreRunE: chaincode.InitChaincode,
	}
	return c
}

func GetPersonRun(cmd *cobra.Command, args []string) error {
	cc := viper.Get("chaincode").(*client.Chaincode)
	exit := 0
	enc := json.NewEncoder(os.Stdout)
	for _, arg := range helpers.Bytes(0, args) {
		res, err := padfedresponse(cc.Query("GetPersona", [][]byte{arg}))
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

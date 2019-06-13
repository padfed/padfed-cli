package dump

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/padfed/padfed-cli/app/helpers"
	client "github.com/padfed/padfed-cli/fabric-client"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Print channel information",
		RunE:  Run,
	}
}

func Run(cmd *cobra.Command, args []string) error {
	chn, err := cmd.Flags().GetString("channel")
	helpers.ExitOnError(err)
	cli, err := client.New(client.User{
		Organization: viper.GetString("fabric.organization"),
		Certificate:  []byte(viper.GetString("fabric.user.certificate")),
		PrivateKey:   []byte(viper.GetString("fabric.user.privatekey")),
	})
	helpers.ExitOnError(err)
	ch, err := cli.Channel(chn)
	helpers.ExitOnError(err)
	l, err := ch.Ledger()
	helpers.ExitOnError(err)
	var i uint64 = 0
	for {
		block, err := l.Block(i)
		helpers.ExitOnError(err)
		io.Copy(os.Stdout, bytes.NewReader(helpers.MustMarshal(block)))
		fmt.Fprintln(os.Stdout)
		i++
	}
	return nil
}

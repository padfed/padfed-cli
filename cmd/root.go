package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/padfed/padfed-cli/cmd/padfed"

	"github.com/mitchellh/go-homedir"
	"github.com/padfed/padfed-cli/cmd/chaincode"
	"github.com/padfed/padfed-cli/cmd/channel"
	client "github.com/padfed/padfed-cli/fabric-client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile     string
	debug, verbose bool
)

func Command() *cobra.Command {
	root := &cobra.Command{
		Use:           "padfed-cli",
		Short:         "Padrón Federal Command Line Interface Tool",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	root.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Configuration file to use")
	root.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")
	root.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Output more data")

	root.PersistentFlags().String("channel", "padfedchannel", "Channel name")
	viper.BindPFlag("fabric.channel", root.PersistentFlags().Lookup("channel"))

	root.PersistentFlags().String("chaincode", "padfedcc", "Chaincode name")
	viper.BindPFlag("fabric.chaincode", root.PersistentFlags().Lookup("chaincode"))

	root.AddCommand(chaincode.Command())
	root.AddCommand(channel.Command())

	root.AddCommand(padfed.GetPerson())
	root.AddCommand(padfed.PutPerson())
	root.AddCommand(padfed.PutPersonList())
	return root
}

func Execute() {
	cobra.OnInitialize(LoadConfiguration)
	err := Command().Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}

func exitOnError(err error) {
	if err != nil {
		if viper.GetBool("debug") {
			fmt.Fprintf(os.Stderr, "ERROR: %+v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		}
		os.Exit(1)
	}
}

func LoadConfiguration() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	}
	home, err := homedir.Dir()
	exitOnError(err)
	viper.SetConfigName("padfed-cli")
	viper.AddConfigPath("/etc/padfed/")
	viper.AddConfigPath(filepath.Join(home, ".config", "padfed"))
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	exitOnError(err)
	viper.Set("debug", debug)
	viper.Set("verbose", verbose)
	cli, err := client.New(client.User{
		Organization: viper.GetString("fabric.organization"),
		Certificate:  []byte(viper.GetString("fabric.user.certificate")),
		PrivateKey:   []byte(viper.GetString("fabric.user.privatekey")),
	})
	exitOnError(err)
	viper.Set("fabric-client", cli)
}

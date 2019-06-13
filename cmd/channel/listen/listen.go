package listen

import (
	"io"
	"os"

	"github.com/apex/log"
	"github.com/dustin/go-humanize"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	"github.com/padfed/padfed-cli/app/helpers"
	client "github.com/padfed/padfed-cli/fabric-client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	outputFile   string
	outputAppend bool
	headerOnly   bool
	seek         string
	blockNumber  uint64
	pretty       bool
)

func Command() *cobra.Command {
	c := &cobra.Command{
		Use:   "listen",
		Short: "Listen for blocks and dump them",
		RunE:  Run,
	}
	c.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path")
	c.Flags().BoolVarP(&outputAppend, "append", "a", false, "Append to output file")
	c.Flags().BoolVar(&headerOnly, "header-only", false, "Only output block headers")
	c.Flags().StringVarP(&seek, "seek", "s", "end", "Seek mode (newest, oldest, from block)")
	c.Flags().Uint64VarP(&blockNumber, "block", "b", 0, "Block number")
	c.Flags().BoolVarP(&pretty, "pretty", "p", false, "Pretty print output")
	return c
}

func Run(cmd *cobra.Command, args []string) error {
	cli := viper.Get("fabric-client").(*client.Client)
	chn, err := cmd.Flags().GetString("channel")
	if err != nil {
		return err
	}
	ch, err := cli.Channel(chn)
	if err != nil {
		return err
	}
	var f io.Writer = os.Stdout
	if outputFile != "" {
		p := os.O_RDWR
		if !outputAppend {
			p = p | os.O_CREATE | os.O_TRUNC
		}
		f, err = os.OpenFile(outputFile, p, 0666)
		if err != nil {
			return err
		}
	}
	var blocks <-chan *common.Block
	switch seek {
	case "start":
		blocks, err = ch.Listen(nil, client.Start, 0)
	case "end":
		blocks, err = ch.Listen(nil, client.End, 0)
	case "block":
		blocks, err = ch.Listen(nil, client.Block, blockNumber)
	default:
		return errors.Errorf("invalid seek mode: %q", seek)
	}
	if err != nil {
		return err
	}
	for block := range blocks {
		log.Infof("block %d (%s)", block.Header.Number, humanize.Bytes(uint64(block.XXX_Size())))
		if headerOnly {
			err = helpers.Dump(f, block.Header, pretty)
		} else {
			err = helpers.Dump(f, block, pretty)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

package dump

import (
	"io"
	"os"
	"regexp"
	"strconv"

	"github.com/apex/log"
	"github.com/dustin/go-humanize"
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
	repeat       uint64
	blocks       []string
	pretty       bool
)

func Command() *cobra.Command {
	c := &cobra.Command{
		Use:   "dump",
		Short: "Dump ledger blocks",
		RunE:  Run,
	}
	c.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path")
	c.Flags().BoolVarP(&outputAppend, "append", "a", false, "Append to output file")
	c.Flags().BoolVar(&headerOnly, "header-only", false, "Only output block headers")
	c.Flags().Uint64VarP(&repeat, "repeat", "r", 1, "Number of times to repeat the task")
	c.Flags().StringSliceVarP(&blocks, "blocks", "b", nil, "Blocks to read (one or more ranges)")
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
	l, err := ch.Ledger()
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
	info, err := l.Info()
	if err != nil {
		return err
	}
	rngs, err := parseRanges(blocks, info.BCI.Height-1)
	if err != nil {
		return err
	}
	for r := uint64(0); r < repeat; r++ {
		for _, rng := range rngs {
			for num := rng.start; num <= rng.end; num++ {
				block, err := l.Block(num)
				if err != nil {
					return err
				}
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
		}
	}
	return nil
}

func parseRanges(blocks []string, height uint64) ([]rng, error) {
	rngs := []rng{}
	for _, block := range blocks {
		m := blocksRE.FindStringSubmatchIndex(block)
		if m == nil {
			return nil, errors.Errorf("invalid blocks expression: %q", block)
		}
		var err error
		var si uint64 = 0
		ss := string(blocksRE.ExpandString(nil, "$single$start", block, m))
		if ss != "" {
			si, err = strconv.ParseUint(ss, 10, 64)
			if err != nil {
				return nil, err
			}
		}
		var ei uint64 = height
		es := string(blocksRE.ExpandString(nil, "$single$end", block, m))
		if es != "" {
			ei, err = strconv.ParseUint(es, 10, 64)
			if err != nil {
				return nil, err
			}
		}
		rngs = append(rngs, rng{si, ei})
	}
	return rngs, nil

}

type rng struct {
	start, end uint64
}

var blocksRE = regexp.MustCompile(`^(?:(?P<single>\d+)|(?P<start>\d+)-(?P<end>\d+)?|-(?P<end>\d+))$`)

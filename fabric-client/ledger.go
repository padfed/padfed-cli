package client

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	"github.com/spf13/viper"
)

type Ledger struct {
	channel *Channel
	ledger  *ledger.Client
}

func (l *Ledger) Block(height uint64) (*common.Block, error) {
	opts := []ledger.RequestOption{}
	targets := viper.GetStringSlice("fabric.target")
	if len(targets) > 0 {
		opts = append(opts, ledger.WithTargetEndpoints(targets...))
	}
	return l.ledger.QueryBlock(height, opts...)
}

func (l *Ledger) Info() (*fab.BlockchainInfoResponse, error) {
	opts := []ledger.RequestOption{}
	targets := viper.GetStringSlice("fabric.target")
	if len(targets) > 0 {
		opts = append(opts, ledger.WithTargetEndpoints(targets...))
	}
	return l.ledger.QueryInfo(opts...)
}

package client

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/events/deliverclient/seek"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
)

type seekMode uint8

type Seek seekMode

const (
	Block Seek = iota
	Start
	End
)

type Channel struct {
	client   *Client
	channel  *channel.Client
	provider context.ChannelProvider
}

func (c *Channel) Chaincode(name string) *Chaincode {
	return &Chaincode{
		channel: c,
		name:    name,
	}
}

func (c *Channel) Ledger() (*Ledger, error) {
	l, err := ledger.New(c.provider)
	if err != nil {
		return nil, err
	}
	return &Ledger{
		channel: c,
		ledger:  l,
	}, nil
}

func (c *Channel) Listen(done <-chan struct{}, s Seek, block uint64) (<-chan *common.Block, error) {
	opts := []event.ClientOption{event.WithBlockEvents()}
	switch s {
	case Start:
		opts = append(opts, event.WithSeekType(seek.Oldest))
	case End:
		opts = append(opts, event.WithSeekType(seek.Newest))
	case Block:
		opts = append(opts, event.WithSeekType(seek.FromBlock), event.WithBlockNum(block))
	}
	e, err := event.New(c.provider, opts...)
	if err != nil {
		return nil, err
	}
	reg, in, err := e.RegisterBlockEvent()
	if err != nil {
		return nil, err
	}
	out := make(chan *common.Block)
	go func() {
		for {
			select {
			case ev := <-in:
				out <- ev.Block
			case <-done:
				e.Unregister(reg)
				close(out)
			}
		}
	}()
	return out, nil
}

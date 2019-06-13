package client

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	mspp "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

type Client struct {
	sdk      *fabsdk.FabricSDK
	identity *mspp.SigningIdentity
}

func (c *Client) Close() {
	c.sdk.Close()
}

func (c *Client) Channel(name string) (*Channel, error) {
	ctx := c.sdk.ChannelContext(name, fabsdk.WithIdentity(*c.identity))
	ch, err := channel.New(ctx)
	if err != nil {
		return nil, err
	}
	return &Channel{
		client:   c,
		channel:  ch,
		provider: ctx,
	}, nil
}

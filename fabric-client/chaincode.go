package client

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/multi"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/status"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Chaincode struct {
	channel *Channel
	name    string
}

type Response struct {
	Status      int32  `json:"status,omitempty"`
	Payload     []byte `json:"payload,omitempty"`
	Message     string `json:"message,omitempty"`
	Transaction string `json:"transaction,omitempty"`

	ErrorMessages []string          `json:"errormessages,omitempty"`
	ErrorStatus   *status.Status    `json:"errorstatus,omitempty"`
	Response      *channel.Response `json:"response,omitempty"`
}

func (c *Chaincode) Query(function string, arguments [][]byte) *Response {
	opts := []channel.RequestOption{}
	targets := viper.GetStringSlice("fabric.target")
	if len(targets) > 0 {
		opts = append(opts, channel.WithTargetEndpoints(targets...))
	}
	req := channel.Request{
		ChaincodeID: c.name,
		Fcn:         function,
		Args:        arguments,
	}
	res, err := c.channel.channel.Query(req, opts...)
	return processResponse(&res, err)
}

func (c *Chaincode) Execute(function string, arguments [][]byte) *Response {
	opts := []channel.RequestOption{}
	targets := viper.GetStringSlice("fabric.target")
	if len(targets) > 0 {
		opts = append(opts, channel.WithTargetEndpoints(targets...))
	}
	res, err := c.channel.channel.Execute(channel.Request{
		ChaincodeID: c.name,
		Fcn:         function,
		Args:        arguments,
	}, opts...)
	return processResponse(&res, err)
}

func processResponse(res *channel.Response, err error) *Response {
	var (
		ccstatus    int32  = res.ChaincodeStatus
		payload     []byte = res.Payload
		transaction string = string(res.TransactionID)
		message     string
		errmsgs     []string
		errstatus   *status.Status
	)
	if err != nil {
		err := errors.Cause(err)
		switch err := err.(type) {
		case multi.Errors:
			for _, err := range err {
				errmsgs = append(errmsgs, fmt.Sprintf("%v (%T)", err, err))
			}
		case *status.Status:
			message = err.Message
			errstatus = err
			errmsgs = append(errmsgs, fmt.Sprintf("%v (%T)", err, err))
			if len(err.Details) > 1 {
				payload = err.Details[1].([]byte)
			}
		default:
			errmsgs = append(errmsgs, fmt.Sprintf("%v (%T)", err, err))
		}
	}
	return &Response{
		Status:        ccstatus,
		Payload:       payload,
		Transaction:   transaction,
		Message:       message,
		Response:      res,
		ErrorMessages: errmsgs,
		ErrorStatus:   errstatus,
	}
}

package padfed

import (
	"encoding/json"

	client "github.com/padfed/padfed-cli/fabric-client"
)

type response struct {
	Status      int         `json:"status,omitempty"`
	Transaction string      `json:"transaction,omitempty"`
	Message     string      `json:"message,omitempty"`
	Payload     interface{} `json:"payload,omitempty"`
}

func padfedresponse(res *client.Response) (*response, error) {
	var payload interface{} = nil
	if len(res.Payload) > 0 && string(res.Payload[0:2]) != "{}" {
		if err := json.Unmarshal(res.Payload, &payload); err != nil {
			return nil, err
		}
	}
	r := &response{
		Status:      res.Status,
		Transaction: res.Transaction,
		Message:     res.Message,
		Payload:     payload,
	}
	return r, nil
}

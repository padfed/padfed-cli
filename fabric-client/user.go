package client

import (
	"io/ioutil"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	mspp "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type User struct {
	Organization string
	Certificate  []byte
	PrivateKey   []byte
}

func New(u User) (*Client, error) {
	con := viper.GetString("fabric.config")
	cfgy := map[interface{}]interface{}{}
	err := yaml.Unmarshal([]byte(con), &cfgy)
	if err != nil {
		return nil, err
	}
	if viper.GetBool("debug") {
		cfgy["client"].(map[interface{}]interface{})["logging"].(map[interface{}]interface{})["level"] = "debug"
	}
	bs, err := yaml.Marshal(cfgy)
	if err != nil {
		return nil, err
	}
	f, err := ioutil.TempFile("", "padfed-cli-config-*.yaml")
	if err != nil {
		return nil, err
	}
	defer os.Remove(f.Name())
	err = ioutil.WriteFile(f.Name(), bs, 0600)
	if err != nil {
		return nil, err
	}
	cfg := config.FromFile(f.Name())
	sdk, err := fabsdk.New(cfg)
	if err != nil {
		return nil, err
	}
	mspc, err := msp.New(sdk.Context(), msp.WithOrg(u.Organization))
	if err != nil {
		return nil, err
	}
	identity, err := mspc.CreateSigningIdentity(mspp.WithCert(u.Certificate), mspp.WithPrivateKey(u.PrivateKey))
	if err != nil {
		return nil, err
	}
	return &Client{
		sdk:      sdk,
		identity: &identity,
	}, nil
}

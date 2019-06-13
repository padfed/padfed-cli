module github.com/padfed/padfed-cli

go 1.12

require (
	github.com/apex/log v1.1.0
	github.com/cloudflare/cfssl v0.0.0-20190510060611-9c027c93ba9e // indirect
	github.com/dustin/go-humanize v1.0.0
	github.com/golang/mock v1.3.1 // indirect
	github.com/google/certificate-transparency-go v1.0.21 // indirect
	github.com/hyperledger/fabric-sdk-go v1.0.0-alpha5
	github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric v0.0.0-20190524192706-bfae339c63bf
	github.com/mitchellh/go-homedir v1.1.0
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v0.0.4
	github.com/spf13/viper v1.4.0
	golang.org/x/crypto v0.0.0-20190530122614-20be4c3c3ed5 // indirect
	golang.org/x/net v0.0.0-20190603091049-60506f45cf65 // indirect
	gopkg.in/yaml.v2 v2.2.2
)

replace github.com/hyperledger/fabric-sdk-go => github.com/lalloni/fabric-sdk-go v1.0.0-alpha5.0.20190613180105-4492670dd456

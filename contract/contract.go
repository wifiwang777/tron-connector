package contract

import (
	"github.com/wifiwang777/tron-connector/client"
	"github.com/wifiwang777/tron-connector/common"
)

type Contract struct {
	Address common.Address
	Client  *client.Client
}

package cidr

import (
	"fmt"
	"net"
)

type Network struct {
	*net.IPNet
	ip net.IP
}

func NewNetwork(s string) (*Network, error) {
	ip, n, err := net.ParseCIDR(s)
	if err != nil {
		return nil, fmt.Errorf("Error creating network: %v\n", err)
	}
	return &Network{n, ip}, nil
}

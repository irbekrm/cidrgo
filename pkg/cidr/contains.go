package cidr

import (
	"fmt"
	"net"
)

func (n *Network) ContainsIP(s string) (bool, error) {
	ip, err := parseIP(s)
	if err != nil {
		return false, err
	}
	return n.Contains(ip), nil
}

func (n *Network) ContainsSubnet(s string) (bool, error) {
	ip, subnet, err := parseSubnet(s)
	if err != nil {
		return false, err
	}
	sMaskSize, _ := subnet.Mask.Size()
	nMaskSize, _ := n.Mask.Size()
	if sMaskSize < nMaskSize {
		return false, nil
	}
	return n.Contains(ip), nil
}

func parseSubnet(s string) (net.IP, *net.IPNet, error) {
	ip, network, err := net.ParseCIDR(s)
	if err != nil {
		return nil, nil, fmt.Errorf("Error parsing subnet %s: %v\n", s, err)
	}
	return ip, network, nil
}

func parseIP(s string) (net.IP, error) {
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, fmt.Errorf("Error parsing IP %s\n", s)
	}
	return ip, nil
}

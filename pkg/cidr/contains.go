package cidr

import (
	"fmt"
	"net"
)

// ContainsIP returns whether an address is in range for the network
func (n *Network) ContainsIP(s string) (bool, error) {
	ip, err := parseIP(s)
	if err != nil {
		return false, err
	}
	return n.Contains(ip), nil
}

// ContainsIPAsHostAddress accepts an IP address and returns 3 boolean values indicating
// whether the IP address is in range for the network,
// is only used as network's own address, is only used as broadcast address and an error
func (n *Network) ContainsIPAsHostAddress(s string) (bool, bool, bool, error) {
	var inRange, onlyNetwork, onlyBroadcast bool
	ip, err := parseIP(s)

	if err != nil {
		return inRange, onlyNetwork, onlyBroadcast, err
	}
	inRange = n.Contains(ip)

	if n.has31Exception() || n.has32Exception() {
		return inRange, onlyNetwork, onlyBroadcast, nil
	}

	onlyNetwork = ip.Equal(n.IP)
	onlyBroadcast = n.isBroadcastAddress(ip)
	return inRange, onlyNetwork, onlyBroadcast, nil
}

// ContainsSubnet returns a boolean value indicating whether
// subnet is in network range and an error
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

func (n *Network) isBroadcastAddress(ip net.IP) bool {
	im := n.inverseMask()
	lastAddress := maskWithOR(im, n.IP)
	return ip.Equal(lastAddress)
}

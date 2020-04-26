package cidr

import (
	"fmt"
	"math"
	"net"
	"strings"
)

type Info struct {
	NetworkAddress         net.IP
	AvailableHostAddresses int
	AllAddresses           int
	Netmask                string
	FirstAddress           net.IP
	LastAddress            net.IP
}

func (n *Network) Describe() Info {
	all, available := n.addresses()
	nm := n.netmaskString()
	f := n.firstHostAddress()
	l := n.lastHostAddress()
	return Info{
		NetworkAddress:         n.IP,
		AllAddresses:           all,
		AvailableHostAddresses: available,
		Netmask:                nm,
		FirstAddress:           f,
		LastAddress:            l,
	}
}

func (n *Network) lastHostAddress() net.IP {
	im := n.inverseMask()
	last := maskWithOR(im, n.IP)
	if n.has32Exception() || n.has31Exception() {
		return last
	}
	last[len(last)-1]--
	return last
}

func (n *Network) firstHostAddress() net.IP {
	nip := n.IP.To4()
	if n.has31Exception() || n.has32Exception() {
		return nip
	}
	first := make(net.IP, len(nip))
	copy(first, nip)
	first[len(first)-1]++
	return first
}

func (n *Network) inverseMask() net.IPMask {
	m := n.Mask
	im := make(net.IPMask, len(m))
	for i, b := range m {
		im[i] = ^b
	}
	return im
}

func (n *Network) addresses() (int, int) {
	if n.has31Exception() {
		return 2, 2
	}
	if n.has32Exception() {
		return 1, 1
	}
	leadingBits, size := n.Mask.Size()
	lastBits := size - leadingBits
	a := math.Pow(2, float64(lastBits))
	return int(a), int(a) - 2
}

func (n *Network) netmaskString() string {
	m := n.Mask
	b := []byte(m)
	s := byteToString(b)
	return strings.Join(s, ".")
}

// Adapted from https://golang.org/src/net/ip.go?s=6471:6504#L238
func maskWithOR(mask net.IPMask, ip net.IP) net.IP {
	n := len(ip)
	if n != len(mask) {
		return nil
	}
	out := make(net.IP, n)
	for i := 0; i < n; i++ {
		out[i] = ip[i] | mask[i]
	}
	return out
}

func byteToString(b []byte) []string {
	s := make([]string, len(b))
	for i, v := range b {
		s[i] = fmt.Sprintf("%v", v)
	}
	return s
}

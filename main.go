package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"strings"
)

var (
	// Subcommands
	containsCommand = flag.NewFlagSet("contains", flag.ExitOnError)
	infoCommand     = flag.NewFlagSet("info", flag.ExitOnError)

	// Flag pointers
	ipPtr      = containsCommand.String("ip", "", "IP address. Will test if the network contains this IP. Exactly one of 'ip' or 'subnet' is required.")
	networkPtr = containsCommand.String("network", "", "Network in CIDR notation (Required)")
	subnetPtr  = containsCommand.String("subnet", "", "Subnet in CIDR notation. Will test if the network contains this subnet. Exactly one of 'ip' or 'subnet' is required.")

	networkPtrInfo = infoCommand.String("network", "", "Network in CIDR notation (Required)")
)

type networkInfo struct {
	networkAddress         string
	availableHostAddresses int
	allAddresses           int
	netmask                string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: subcommand required")
		os.Exit(1)
	}
	var outcome string
	switch os.Args[1] {
	case "contains":
		validateContains()
		outcome = contains()
	case "info":
		validateInfo()
		outcome = info()
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
	fmt.Println(outcome)
}

func validateContains() {
	containsCommand.Parse(os.Args[2:])
	if (*ipPtr == "" && *subnetPtr == "") || *networkPtr == "" {
		containsCommand.PrintDefaults()
		os.Exit(1)
	}
}

func validateInfo() {
	infoCommand.Parse(os.Args[2:])
	if *networkPtrInfo == "" {
		infoCommand.PrintDefaults()
		os.Exit(1)
	}
}

func info() string {
	_, network, err := net.ParseCIDR(*networkPtrInfo)
	if err != nil {
		log.Fatalf("Error parsing network CIDR: %v\n", err)
	}
	all, available := hosts(network)
	nm := netmask(network)
	i := &networkInfo{
		networkAddress:         network.IP.String(),
		allAddresses:           all,
		availableHostAddresses: available,
		netmask:                nm,
	}
	return fmt.Sprintf("Info:\nNetwork address: %s\nAll addresses: %d\nAvailable host addresses: %d\nNetmask: %s\n", i.networkAddress, i.allAddresses, i.availableHostAddresses, i.netmask)
}

func contains() string {
	var c bool
	_, network, err := net.ParseCIDR(*networkPtr)
	if err != nil {
		log.Fatalf("Error parsing network CIDR: %v\n", err)
	}
	if *ipPtr != "" {
		c, err = containsIP(*ipPtr, network)
	}
	if *subnetPtr != "" {
		c, err = containsSubnet(*subnetPtr, network)
	}
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprint(c)
}

func parseIP(s string) (net.IP, error) {
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, fmt.Errorf("Error parsing IP %s\n", s)
	}
	return ip, nil
}

func parseSubnet(s string) (net.IP, *net.IPNet, error) {
	ip, network, err := net.ParseCIDR(s)
	if err != nil {
		return nil, nil, fmt.Errorf("Error parsing subnet %s: %v\n", s, err)
	}
	return ip, network, nil
}

func containsIP(s string, network *net.IPNet) (bool, error) {
	ip, err := parseIP(s)
	if err != nil {
		return false, err
	}
	return network.Contains(ip), nil
}

func containsSubnet(s string, network *net.IPNet) (bool, error) {
	ip, subnet, err := parseSubnet(s)
	if err != nil {
		return false, err
	}
	sMaskSize, _ := subnet.Mask.Size()
	nMaskSize, _ := network.Mask.Size()
	if sMaskSize < nMaskSize {
		return false, nil
	}
	return network.Contains(ip), nil
}

func hosts(network *net.IPNet) (int, int) {
	if has31Exception(network) {
		return 2, 2
	}
	if has32Exception(network) {
		return 1, 1
	}
	leadingBits, size := network.Mask.Size()
	lastBits := size - leadingBits
	h := math.Pow(2, float64(lastBits))
	return int(h), int(h) - 2
}

func netmask(network *net.IPNet) string {
	m := network.Mask
	b := []byte(m)
	s := byteToString(b)
	return strings.Join(s, ".")
}

func byteToString(b []byte) []string {
	s := make([]string, len(b))
	for i, v := range b {
		s[i] = fmt.Sprintf("%v", v)
	}
	return s
}

func has31Exception(network *net.IPNet) bool {
	ones, bits := network.Mask.Size()
	return bits == 32 && ones == 31
}
func has32Exception(network *net.IPNet) bool {
	ones, bits := network.Mask.Size()
	return bits == 32 && ones == 32
}

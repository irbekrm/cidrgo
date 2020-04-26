package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/irbekrm/cidrgo/pkg/cidr"
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
	firstAddress           net.IP
	lastAddress            net.IP
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
	n, err := cidr.NewNetwork(*networkPtrInfo)
	if err != nil {
		log.Fatalf("Error parsing network CIDR: %v\n", err)
	}
	i := n.Describe()
	return fmt.Sprintf("Info:\nNetwork address: %s\nAll addresses: %d\nAvailable host addresses: %d\nNetmask: %s\nFirst host address: %v\nLast host address: %v\n", i.NetworkAddress, i.AllAddresses, i.AvailableHostAddresses, i.Netmask, i.FirstAddress, i.LastAddress)
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

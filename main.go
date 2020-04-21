package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

var (
	// Subcommands
	containsCommand = flag.NewFlagSet("contains", flag.ExitOnError)

	// containsCommand flag pointers
	ipPtr      = containsCommand.String("ip", "", "IP address. Will test if network contains this IP. Exactly one of ip or subnet is required")
	networkPtr = containsCommand.String("network", "", "Network in CIDR notation (Required)")
	subnetPtr  = containsCommand.String("subnet", "", "Subnet in CIDR notation. Will test if network contains this subnet. Exactly one of ip or subnet is required.")
)

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

func contains() string {
	var ip net.IP
	_, network, err := net.ParseCIDR(*networkPtr)
	if err != nil {
		log.Fatalf("Error parsing network CIDR: %v\n", err)
	}
	if *ipPtr != "" {
		ip, err = parseIP(*ipPtr)
	}
	if *subnetPtr != "" {
		ip, err = parseSubnet(*subnetPtr)
	}
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprint(network.Contains(ip))
}

func parseIP(s string) (net.IP, error) {
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, fmt.Errorf("Error parsing IP %s\n", s)
	}
	return ip, nil
}

func parseSubnet(s string) (net.IP, error) {
	ip, _, err := net.ParseCIDR(s)
	if err != nil {
		return nil, fmt.Errorf("Error parsing subnet %s: %v\n", s, err)
	}
	return ip, err
}

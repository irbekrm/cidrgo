# cidrgo

A CLI CIDR calculator

## Installation

1. [Install go](https://golang.org/doc/install) if you haven't already.

2.  Run ` go get -u github.com/irbekrm/cidrgo/...` 

## Commands

### info

##### Usage:
`$ cidrgo info -network NETWORK`

##### Description:

`info` command describes a CIDR range. 

Outputs:

`Network address` - the address of the network

`All addresses` - how many addresses are in the CIDR range

`Available host addresses` - how many addresses are in the CIDR range excluding network address and broadcast address. 

`Netmask` - netmask

`First host address` - first available host address in the range. (However, this is likely to differ on different IAAS providers as they reserve more addresses from a CIDR range)

`Last host address` - last available host address

##### Example:

```
$ cidrgo info -network 10.11.0.0/16

// Info:
// Network address: 10.11.0.0
// All addresses: 65536
// Available host addresses: 65534
// Netmask: 255.255.0.0
// First host address: 10.11.0.1
// Last host address: 10.11.255.254

```

### contains

##### Usage:
`$ cidrgo contains -network NETWORK -ip IP|-subnet SUBNET`

##### Description:

`contains` command can be used to determine if an IP address or a subnet is in a network. For an IP address `contains` will also determine whether the address is only used as network address or only used as broadcast address.

##### Examples:

```
$ cidrgo contains -network 10.0.0.0/14 -subnet 10.0.0.0/15

// true 


$ cidrgo contains -network 10.0.0.0/16 -ip 10.0.0.0

// Address in range for network: true
// Network address only: true
// Broadcast address only: false
```



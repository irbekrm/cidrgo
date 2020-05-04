package cidr

import (
	"testing"
)

func TestNetwork_ContainsIP(t *testing.T) {
	tests := []struct {
		name        string
		ip          string
		networkCIDR string
		want        bool
		wantErr     bool
	}{
		{
			name:        "Success- contains an IP",
			networkCIDR: "10.0.0.0/16",
			ip:          "10.0.0.0",
			want:        true,
		},
		{
			name:        "Success- does not contain an IP",
			networkCIDR: "10.13.14.15/24",
			ip:          "10.13.10.0",
			want:        false,
		},
		{
			name:        "Failure- CIDR range passed instead of an IP",
			networkCIDR: "10.0.0.0/16",
			ip:          "10.0.0.0/24",
			want:        false,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, err := NewNetwork(tt.networkCIDR)
			if err != nil {
				t.Errorf("Error creating network: %v\n", err)
			}
			got, err := n.ContainsIP(tt.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("Network.ContainsIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Network.ContainsIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetwork_ContainsIPAsHostAddress(t *testing.T) {
	tests := []struct {
		name            string
		networkCIDR     string
		ip              string
		isInRange       bool
		isNetwOnly      bool
		isBroadcastOnly bool
		wantErr         bool
	}{
		{
			name:        "Success: IP is in range and is not just broadcast/network address",
			networkCIDR: "10.0.0.0/16",
			ip:          "10.0.1.0",
			isInRange:   true,
		},
		{
			name:        "Success: IP is in range and is not just broadcast/network address for /32 CIDR",
			networkCIDR: "10.0.0.0/32",
			ip:          "10.0.0.0",
			isInRange:   true,
		},
		{
			name:        "Success: IP is in range and is not just broadcast/network address for /31 CIDR",
			networkCIDR: "10.0.0.0/31",
			ip:          "10.0.0.0",
			isInRange:   true,
		},
		{
			name:        "Success: IP is in range and is not just broadcast/network address for /31 CIDR",
			networkCIDR: "10.0.0.0/31",
			ip:          "10.0.0.1",
			isInRange:   true,
		},
		{
			name:            "Success: IP is in range, but is just a broadcast address",
			networkCIDR:     "10.0.0.0/16",
			ip:              "10.0.255.255",
			isInRange:       true,
			isBroadcastOnly: true,
		},
		{
			name:        "Success: IP is in range, but is just a network address",
			networkCIDR: "10.0.0.0/16",
			ip:          "10.0.0.0",
			isInRange:   true,
			isNetwOnly:  true,
		},
		{
			name:        "Success: IP is not in range (too high)",
			networkCIDR: "10.0.0.0/16",
			ip:          "10.1.0.0",
		},
		{
			name:        "Success: IP is not in range (too low)",
			networkCIDR: "10.1.0.0/16",
			ip:          "10.0.0.0",
		},
		{
			name:        "Failure: CIDR range passed instead of an IP address",
			networkCIDR: "10.0.0.0/16",
			ip:          "10.0.0.0/32",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, err := NewNetwork(tt.networkCIDR)
			if err != nil {
				t.Errorf("Error creating network: %v\n", err)
			}
			isInRange, isNetwOnly, isBroadcastOnly, err := n.ContainsIPAsHostAddress(tt.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("Network.ContainsIPAsHostAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if isInRange != tt.isInRange {
				t.Errorf("Network.ContainsIPAsHostAddress() isInRange = %v, want %v", isInRange, tt.isInRange)
			}
			if isNetwOnly != tt.isNetwOnly {
				t.Errorf("Network.ContainsIPAsHostAddress() isNetwOnly = %v, want %v", isNetwOnly, tt.isNetwOnly)
			}
			if isBroadcastOnly != tt.isBroadcastOnly {
				t.Errorf("Network.ContainsIPAsHostAddress() isBroadcastOnly = %v, want %v", isBroadcastOnly, tt.isBroadcastOnly)
			}
		})
	}
}

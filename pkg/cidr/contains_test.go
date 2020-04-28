package cidr

import (
	"testing"
)

func TestNetwork_ContainsIP(t *testing.T) {
	tests := []struct {
		name      string
		ip        string
		networkIP string
		want      bool
		wantErr   bool
	}{
		{
			name:      "Success- contains an IP",
			networkIP: "10.0.0.0/16",
			ip:        "10.0.0.0",
			want:      true,
		},
		{
			name:      "Success- does not contain an IP",
			networkIP: "10.13.14.15/24",
			ip:        "10.13.10.0",
			want:      false,
		},
		{
			name:      "Failure- CIDR range passed instead of an IP",
			networkIP: "10.0.0.0/16",
			ip:        "10.0.0.0/24",
			want:      false,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, err := NewNetwork(tt.networkIP)
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
package cidr

import (
	"net"
	"reflect"
	"testing"
)

func TestNetwork_Describe(t *testing.T) {
	tests := []struct {
		name     string
		netwCIDR string
		want     Info
	}{
		{
			name:     "Success",
			netwCIDR: "10.0.0.0/16",
			want: Info{
				NetworkAddress:         net.IP([]byte{10, 0, 0, 0}),
				AllAddresses:           65536,
				AvailableHostAddresses: 65534,
				Netmask:                "255.255.0.0",
				FirstAddress:           net.IP([]byte{10, 0, 0, 1}),
				LastAddress:            net.IP([]byte{10, 0, 255, 254}),
			},
		},
		{
			name:     "Success- /32 CIDR",
			netwCIDR: "10.0.0.0/32",
			want: Info{
				NetworkAddress:         net.IP([]byte{10, 0, 0, 0}),
				AllAddresses:           1,
				AvailableHostAddresses: 1,
				Netmask:                "255.255.255.255",
				FirstAddress:           net.IP([]byte{10, 0, 0, 0}),
				LastAddress:            net.IP([]byte{10, 0, 0, 0}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			netw, err := NewNetwork(tt.netwCIDR)
			if err != nil {
				t.Errorf("Error creating network: %v\n", err)
			}
			if got := netw.Describe(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Network.Describe() = %+v , want %+v", got, tt.want)
			}
		})
	}
}

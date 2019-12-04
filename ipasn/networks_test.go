// Copyright 2019 Freman/Fremnet (Shannon Wynter). All rights reserved.

package ipasn_test

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/freman/cymru/ipasn"
)

func TestNetworks(t *testing.T) {
	t.Parallel()

	testNets := ipasn.Networks{
		&net.IPNet{IP: net.IP{8, 8, 8, 0}, Mask: net.IPMask{255, 255, 255, 0}},
		&net.IPNet{IP: net.IP{1, 1, 1, 0}, Mask: net.IPMask{255, 255, 255, 0}},
	}

	privates := ipasn.DefaultPrivateNetworks()

	// Should be able to recursively define and call networks
	privates = append(privates, testNets)

	checks := []net.IP{net.IPv4(192, 168, 0, 4), net.IPv4(8, 8, 8, 8)}

	for _, ip := range checks {
		ip := ip
		t.Run(ip.String(), func(t *testing.T) {
			t.Parallel()
			require.True(t, privates.Contains(ip))
		})
	}
}

func TestNoNetworks(t *testing.T) {
	t.Parallel()

	nets := ipasn.NoPrivateNetworks()

	checks := []net.IP{net.IPv4(192, 168, 0, 4), net.IPv4(8, 8, 8, 8)}

	for _, ip := range checks {
		ip := ip
		t.Run(ip.String(), func(t *testing.T) {
			t.Parallel()
			require.False(t, nets.Contains(ip))
		})
	}
}

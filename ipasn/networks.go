// Copyright 2019 Freman/Fremnet (Shannon Wynter). All rights reserved.

package ipasn

import "net"

// NetworkFilter interface is based on *IP.Net->Contains magically permits
// the building of larger, or recursive networks.
type NetworkFilter interface {
	// Contains reports whether the network includes ip.
	Contains(ip net.IP) bool
}

// Networks is a list of NetworkFilters
type Networks []NetworkFilter

// Contains calls Contains on every NetworkFilter in the list
func (n Networks) Contains(ip net.IP) bool {
	for _, r := range n {
		if r.Contains(ip) {
			return true
		}
	}

	return false
}

// DefaultPrivateNetworks returns a lsit of commonly private ip ranges.
func DefaultPrivateNetworks() Networks {
	return Networks{
		&net.IPNet{IP: net.IP{0, 0, 0, 0}, Mask: net.IPMask{255, 0, 0, 0}},
		&net.IPNet{IP: net.IP{10, 0, 0, 0}, Mask: net.IPMask{255, 0, 0, 0}},
		&net.IPNet{IP: net.IP{100, 64, 0, 0}, Mask: net.IPMask{255, 192, 0, 0}},
		&net.IPNet{IP: net.IP{127, 0, 0, 0}, Mask: net.IPMask{255, 0, 0, 0}},
		&net.IPNet{IP: net.IP{172, 16, 0, 0}, Mask: net.IPMask{255, 240, 0, 0}},
		&net.IPNet{IP: net.IP{192, 0, 0, 0}, Mask: net.IPMask{255, 255, 255, 0}},
		&net.IPNet{IP: net.IP{192, 0, 2, 0}, Mask: net.IPMask{255, 255, 255, 0}},
		&net.IPNet{IP: net.IP{192, 88, 99, 0}, Mask: net.IPMask{255, 255, 255, 0}},
		&net.IPNet{IP: net.IP{192, 168, 0, 0}, Mask: net.IPMask{255, 255, 0, 0}},
		&net.IPNet{IP: net.IP{198, 18, 0, 0}, Mask: net.IPMask{255, 254, 0, 0}},
		&net.IPNet{IP: net.IP{198, 51, 100, 0}, Mask: net.IPMask{255, 255, 255, 0}},
		&net.IPNet{IP: net.IP{203, 0, 113, 0}, Mask: net.IPMask{255, 255, 255, 0}},
	}
}

type alwaysTheSameAnswer bool

func (a alwaysTheSameAnswer) Contains(_ net.IP) bool {
	return bool(a)
}

// NoPrivateNetworks disables the network filtering by always returning false
func NoPrivateNetworks() NetworkFilter {
	return alwaysTheSameAnswer(false)
}

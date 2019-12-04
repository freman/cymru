// Copyright 2019 Freman/Fremnet (Shannon Wynter). All rights reserved.

package ipasn_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/freman/cymru/ipasn"
)

type mockResolver func(ctx context.Context, name string) ([]string, error)

func (m mockResolver) LookupTXT(ctx context.Context, name string) ([]string, error) {
	return m(ctx, name)
}

//nolint:gochecknoglobals
var resolver = mockResolver(func(ctx context.Context, name string) ([]string, error) {
	switch name {
	case "31.108.90.216.origin.asn.cymru.com.":
		return []string{"23028 | 216.90.108.0/24 | US | arin | 1998-09-25"}, nil
	case "0.108.90.216.origin.asn.cymru.com.":
		return []string{"23028 | 216.90.108.0/24 | US | arin | 1998-09-25"}, nil
	case "8.6.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.2.0.0.b.0.6.8.4.1.0.0.2.origin6.asn.cymru.com.":
		return []string{"15169 | 2001:4860::/32 | US | arin | 2005-03-14"}, nil
	case "8.8.8.8.origin.asn.cymru.com.":
		return nil, nil
	case "31.108.90.216.peer.asn.cymru.com.":
		return []string{"701 1239 3549 3561 7132 | 216.90.108.0/24 | US | arin | 1998-09-25"}, nil
	case "8.8.8.8.peer.asn.cymru.com.":
		return nil, nil
	case "AS23028.asn.cymru.com.":
		return []string{"23028 | US | arin | 2002-01-04 | TEAM-CYMRU - Team Cymru Inc., US"}, nil
	case "AS1234.asn.cymru.com.":
		return []string{"1234 | EU | ripencc | 1993-09-01 | FORTUM-AS | Fortum, FI"}, nil
	case "AS911.asn.cymru.com.":
		return nil, nil
	}
	return nil, errors.New("what? " + name + " not found")
})

//nolint:gochecknoglobals
var testCases = struct {
	originCases []struct {
		input    net.IP
		expected ipasn.OriginInfo
		err      error
		str      string
	}
	peerCases []struct {
		input    net.IP
		expected ipasn.PeerInfo
		err      error
		str      string
	}
	asnCases []struct {
		input    int
		expected ipasn.ASNInfo
		err      error
		str      string
	}
}{
	originCases: []struct {
		input    net.IP
		expected ipasn.OriginInfo
		err      error
		str      string
	}{
		{
			net.IPv4(216, 90, 108, 31),
			ipasn.OriginInfo{
				ASN:       23028,
				Network:   &net.IPNet{IP: net.IP{216, 90, 108, 0}, Mask: net.IPMask{255, 255, 255, 0}},
				Country:   "US",
				Authority: "arin",
				Updated:   time.Unix(906681600, 0).UTC(),
			},
			nil,
			"23028 | 216.90.108.0/24 | US | arin | 1998-09-25",
		}, {
			net.IPv4(216, 90, 108, 0),
			ipasn.OriginInfo{
				ASN:       23028,
				Network:   &net.IPNet{IP: net.IP{216, 90, 108, 0}, Mask: net.IPMask{255, 255, 255, 0}},
				Country:   "US",
				Authority: "arin",
				Updated:   time.Unix(906681600, 0).UTC(),
			},
			nil,
			"23028 | 216.90.108.0/24 | US | arin | 1998-09-25",
		}, {
			net.IP{0x20, 0x1, 0x48, 0x60, 0xb0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x68},
			ipasn.OriginInfo{
				ASN: 15169,
				Network: &net.IPNet{
					IP:   net.IP{0x20, 0x1, 0x48, 0x60, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
					Mask: net.IPMask{0xff, 0xff, 0xff, 0xff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
				},
				Country:   "US",
				Authority: "arin",
				Updated:   time.Unix(1110758400, 0).UTC(),
			},
			nil,
			"15169 | 2001:4860::/32 | US | arin | 2005-03-14",
		}, {
			net.IPv4zero,
			ipasn.OriginInfo{},
			ipasn.ErrIPIsUnspecified,
			"",
		}, {
			net.IPv4(127, 0, 0, 1),
			ipasn.OriginInfo{},
			ipasn.ErrIPIsLoopback,
			"",
		}, {
			net.IPv4(224, 1, 2, 3),
			ipasn.OriginInfo{},
			ipasn.ErrIPIsMulticast,
			"",
		}, {
			net.IPv4(192, 168, 0, 1),
			ipasn.OriginInfo{},
			ipasn.ErrIPIsPrivate,
			"",
		}, {
			net.IPv4(1, 1, 1, 1),
			ipasn.OriginInfo{},
			errors.New("what? 1.1.1.1.origin.asn.cymru.com. not found"),
			"",
		}, {
			net.IPv4(8, 8, 8, 8),
			ipasn.OriginInfo{},
			ipasn.ErrNotFound,
			"",
		},
	},
	peerCases: []struct {
		input    net.IP
		expected ipasn.PeerInfo
		err      error
		str      string
	}{
		{
			net.IPv4(216, 90, 108, 31),
			ipasn.PeerInfo{
				ASNs:      []int{701, 1239, 3549, 3561, 7132},
				Network:   &net.IPNet{IP: net.IP{216, 90, 108, 0}, Mask: net.IPMask{255, 255, 255, 0}},
				Country:   "US",
				Authority: "arin",
				Updated:   time.Unix(906681600, 0).UTC(),
			},
			nil,
			"701 1239 3549 3561 7132 | 216.90.108.0/24 | US | arin | 1998-09-25",
		}, {
			net.IPv4zero,
			ipasn.PeerInfo{},
			ipasn.ErrIPIsUnspecified,
			"",
		}, {
			net.IPv4(127, 0, 0, 1),
			ipasn.PeerInfo{},
			ipasn.ErrIPIsLoopback,
			"",
		}, {
			net.IPv4(224, 1, 2, 3),
			ipasn.PeerInfo{},
			ipasn.ErrIPIsMulticast,
			"",
		}, {
			net.IPv4(192, 168, 0, 1),
			ipasn.PeerInfo{},
			ipasn.ErrIPIsPrivate,
			"",
		}, {
			net.IPv4(1, 1, 1, 1),
			ipasn.PeerInfo{},
			errors.New("what? 1.1.1.1.peer.asn.cymru.com. not found"),
			"",
		}, {
			net.IPv4(8, 8, 8, 8),
			ipasn.PeerInfo{},
			ipasn.ErrNotFound,
			"",
		},
	},
	asnCases: []struct {
		input    int
		expected ipasn.ASNInfo
		err      error
		str      string
	}{
		{
			23028,
			ipasn.ASNInfo{
				ASN:         23028,
				Country:     "US",
				Authority:   "arin",
				Updated:     time.Unix(1010102400, 0).UTC(),
				Description: "TEAM-CYMRU - Team Cymru Inc., US",
			},
			nil,
			"23028 | US | arin | 2002-01-04 | TEAM-CYMRU - Team Cymru Inc., US",
		}, {
			1234,
			ipasn.ASNInfo{
				ASN:         1234,
				Country:     "EU",
				Authority:   "ripencc",
				Updated:     time.Unix(746841600, 0).UTC(),
				Description: "FORTUM-AS | Fortum, FI",
			},
			nil,
			"1234 | EU | ripencc | 1993-09-01 | FORTUM-AS | Fortum, FI",
		}, {
			1111,
			ipasn.ASNInfo{},
			errors.New("what? AS1111.asn.cymru.com. not found"),
			"",
		}, {
			911,
			ipasn.ASNInfo{},
			ipasn.ErrNotFound,
			"",
		},
	},
}

func TestClientObject(t *testing.T) {
	t.Parallel()

	c := &ipasn.Client{Resolver: resolver}

	for i, test := range testCases.originCases {
		i, test := i, test
		t.Run(fmt.Sprintf("origin_%d", i), func(t *testing.T) {
			t.Parallel()
			got, err := c.Origin(context.TODO(), test.input)
			require.Equal(t, test.err, err)
			require.Equal(t, test.expected, got)
			require.Equal(t, test.str, got.String())
		})
	}

	for i, test := range testCases.peerCases {
		i, test := i, test
		t.Run(fmt.Sprintf("peer_%d", i), func(t *testing.T) {
			t.Parallel()
			got, err := c.Peer(context.TODO(), test.input)
			require.Equal(t, test.err, err)
			require.Equal(t, test.expected, got)
			require.Equal(t, test.str, got.String())
		})
	}

	for i, test := range testCases.asnCases {
		i, test := i, test
		t.Run(fmt.Sprintf("asn_%d", i), func(t *testing.T) {
			t.Parallel()
			got, err := c.ASN(context.TODO(), test.input)
			require.Equal(t, test.err, err)
			require.Equal(t, test.expected, got)
			require.Equal(t, test.str, got.String())
		})
	}
}

func TestDefaultClient(t *testing.T) {
	t.Parallel()

	ipasn.DefaultClient = &ipasn.Client{Resolver: resolver}

	for i, test := range testCases.originCases {
		i, test := i, test
		t.Run(fmt.Sprintf("origin_%d", i), func(t *testing.T) {
			t.Parallel()
			got, err := ipasn.Origin(context.TODO(), test.input)
			require.Equal(t, test.err, err)
			require.Equal(t, test.expected, got)
		})
	}

	for i, test := range testCases.peerCases {
		i, test := i, test
		t.Run(fmt.Sprintf("peer_%d", i), func(t *testing.T) {
			t.Parallel()
			got, err := ipasn.Peer(context.TODO(), test.input)
			require.Equal(t, test.err, err)
			require.Equal(t, test.expected, got)
		})
	}

	for i, test := range testCases.asnCases {
		i, test := i, test
		t.Run(fmt.Sprintf("asn_%d", i), func(t *testing.T) {
			t.Parallel()
			got, err := ipasn.ASN(context.TODO(), test.input)
			require.Equal(t, test.err, err)
			require.Equal(t, test.expected, got)
		})
	}
}

// TestOnlineDefaultResolver tests the online behavior with the default resolver.
// These tests are separate and isolated in case something in the online world changes
// and can be skipped by passing the short test flag
func TestOnlineDefaultResolver(t *testing.T) {
	if testing.Short() {
		t.Skip("Online tests not enabled in short tests")
		return
	}

	c := &ipasn.Client{}
	origin, err := c.Origin(context.TODO(), net.IPv4(216, 90, 108, 31))
	require.Equal(t, nil, err)
	require.Equal(t, ipasn.OriginInfo{
		ASN:       23028,
		Network:   &net.IPNet{IP: net.IP{216, 90, 108, 0}, Mask: net.IPMask{255, 255, 255, 0}},
		Country:   "US",
		Authority: "arin",
		Updated:   time.Unix(906681600, 0).UTC(),
	}, origin)

	peer, err := c.Peer(context.TODO(), net.IPv4(216, 90, 108, 31))
	require.Equal(t, nil, err)
	require.Equal(t, ipasn.PeerInfo{
		ASNs:      []int{3257, 23352},
		Network:   &net.IPNet{IP: net.IP{216, 90, 108, 0}, Mask: net.IPMask{255, 255, 255, 0}},
		Country:   "US",
		Authority: "arin",
		Updated:   time.Unix(906681600, 0).UTC(),
	}, peer)

	asn, err := c.ASN(context.TODO(), 23028)
	require.Equal(t, nil, err)
	require.Equal(t, ipasn.ASNInfo{
		ASN:         23028,
		Country:     "US",
		Authority:   "arin",
		Updated:     time.Unix(1010102400, 0).UTC(),
		Description: "TEAM-CYMRU - Team Cymru Inc., US",
	}, asn)
}

func ExampleClient_Origin() {
	client := new(ipasn.Client)

	origin, err := client.Origin(context.Background(), net.ParseIP("1.1.1.1"))
	if err != nil {
		panic(err)
	}

	fmt.Println(origin)

	// Output: 13335 | 1.1.1.0/24 | AU | apnic | 2011-08-11
}

func ExampleClient_Origin_iPv6() {
	client := new(ipasn.Client)

	origin, err := client.Origin(context.Background(), net.ParseIP("2606:4700:4700::1111"))
	if err != nil {
		panic(err)
	}

	fmt.Println(origin)

	// Output: 13335 | 2606:4700:4700::/48 | US | arin | 2011-11-01
}

func ExampleClient_Peer() {
	client := new(ipasn.Client)

	peer, err := client.Peer(context.Background(), net.ParseIP("1.1.1.1"))
	if err != nil {
		panic(err)
	}

	fmt.Println(peer)

	// Output: 174 1103 2381 2914 | 1.1.1.0/24 | AU | apnic | 2011-08-11
}

func ExampleClient_ASN() {
	client := new(ipasn.Client)

	asn, err := client.ASN(context.Background(), 13335)
	if err != nil {
		panic(err)
	}

	fmt.Println(asn)

	// Output: 13335 | US | arin | 2010-07-14 | CLOUDFLARENET - Cloudflare, Inc., US
}

// Copyright 2019 Freman/Fremnet (Shannon Wynter). All rights reserved.

package ipasn

import (
	"context"
	"net"
	"strconv"
	"strings"
	"time"
)

// Resolver permits the use of net/resolver or any other hand crafted resolver
// (eg something wrapped around miekg/dns) so long as it fulfils this interface.
type Resolver interface {
	LookupTXT(ctx context.Context, name string) ([]string, error)
}

// Client permits calling the Team Cymru IP-ASN mapping dns interface with relative
// ease and reliability.
//
// By default it will use net.DefaultResolver as the resolver but you can pass
// any resolver that implements the Resolver interface
// this will permit you to use the resolver of your choice (eg: miekg/dns) with little
// to no effort.
//
// Also by default this object blocks queries against private networks using a list
// of networks returned by DefaultPrivateNetworks() for your convenience you can
// configure this as NoPrivateNetworks()
//
// You can override either of these properties at any time.
//
// The original Team Cymru documentation specifies that a prefix (eg: 216.90.108)
// only send 108.90.216 but for simplicity, this package doesn't omit the leading
// 0 so it would send 0.108.90.216
type Client struct {
	Resolver        Resolver
	PrivateNetworks NetworkFilter
}

const dateFormat = `2006-01-02`

// Origin is used to map an IPv4 or IPv6 address or prefix to a corresponding
// BGP Origin ASN.
func (c *Client) Origin(ctx context.Context, ip net.IP) (o OriginInfo, err error) {
	if err := c.checkInputIP(ip); err != nil {
		return o, err
	}

	dat, err := c.lookupTXT(ctx, asLookupString(ip, "origin"))
	if err != nil {
		return o, err
	}

	if len(dat) == 5 {
		o.ASN, _ = strconv.Atoi(dat[0])
		_, o.Network, _ = net.ParseCIDR(dat[1])
		o.Country = dat[2]
		o.Authority = dat[3]
		o.Updated, _ = time.Parse(dateFormat, dat[4])
	}

	return o, nil
}

// Peer is used to map an IP address or prefix to the possible BGP peer ASNs that
// are one AS hop away from the BGP Origin ASN's prefix.
func (c *Client) Peer(ctx context.Context, ip net.IP) (p PeerInfo, err error) {
	if err := c.checkInputIP(ip); err != nil {
		return p, err
	}

	dat, err := c.lookupTXT(ctx, asLookupString(ip, "peer"))
	if err != nil {
		return p, err
	}

	if len(dat) == 5 {
		p.ASNs = parseASNList(dat[0])
		_, p.Network, _ = net.ParseCIDR(dat[1])
		p.Country = dat[2]
		p.Authority = dat[3]
		p.Updated, _ = time.Parse(dateFormat, dat[4])
	}

	return p, nil
}

// ASN is used to determine the AS description of a given BGP ASN.
// Notably this function returns the Description of the AS but not the network.
func (c *Client) ASN(ctx context.Context, asn int) (a ASNInfo, err error) {
	dat, err := c.lookupTXT(ctx, "AS"+strconv.Itoa(asn)+".asn.cymru.com.")
	if err != nil {
		return a, err
	}

	if len(dat) >= 5 {
		a.ASN, _ = strconv.Atoi(dat[0])
		a.Country = dat[1]
		a.Authority = dat[2]
		a.Updated, _ = time.Parse(dateFormat, dat[3])
		a.Description = strings.Join(dat[4:], " | ")
	}

	return a, nil
}

// lookupTXT checks there's a resolver (or makes one available) before
// simply forwarding the call to the lookupTXT
func (c *Client) lookupTXT(ctx context.Context, name string) ([]string, error) {
	if c.Resolver == nil {
		c.Resolver = net.DefaultResolver
	}

	vals, err := c.Resolver.LookupTXT(ctx, name)
	if err != nil {
		return nil, err
	}

	if len(vals) == 0 {
		return nil, ErrNotFound
	}

	return strings.Split(vals[0], " | "), nil
}

// isPrivateNetwork checks if the given ip falls in the list of private
// networks
//
// Optionally builds that list list from the default
func (c *Client) isPrivateNetwork(ip net.IP) bool {
	if c.PrivateNetworks == nil {
		c.PrivateNetworks = DefaultPrivateNetworks()
	}

	return c.PrivateNetworks.Contains(ip)
}

// checkInputIP performs basic sanity checking on the given IP to
// attempt to reduce the network traffic and lookups for things that
// realistically won't resolve.
func (c *Client) checkInputIP(ip net.IP) error {
	switch {
	case ip == nil || ip.IsUnspecified():
		return ErrIPIsUnspecified
	case ip.IsLoopback():
		return ErrIPIsLoopback
	case ip.IsMulticast():
		return ErrIPIsMulticast
	case c.isPrivateNetwork(ip):
		return ErrIPIsPrivate
	}

	return nil
}

// asLookupString does the mangling of the given ip to fit the required
// format of the dns request.
func asLookupString(ip net.IP, zone string) string {
	const hexDigit = "0123456789abcdef"

	zone6 := zone
	if zone == "origin" {
		zone6 += "6"
	}

	if ip.To4() != nil {
		end := len(ip) - 1

		return strings.Join([]string{
			strconv.Itoa(int(ip[end])),
			strconv.Itoa(int(ip[end-1])),
			strconv.Itoa(int(ip[end-2])),
			strconv.Itoa(int(ip[end-3])),
			zone,
			"asn.cymru.com.",
		}, ".")
	}

	// Must be IPv6
	buf := make([]byte, 0, len(ip)*4+len(zone6+".asn.cymru.com."))

	// Add it, in reverse, to the buffer
	for i := len(ip) - 1; i >= 0; i-- {
		buf = append(buf, hexDigit[ip[i]&0xF], '.', hexDigit[ip[i]>>4], '.')
	}

	buf = append(buf, zone6+".asn.cymru.com."...)

	return string(buf)
}

// parseASNList is cheap and nasty string to list of int parsing
func parseASNList(in string) []int {
	tmp := strings.Fields(in)
	r := make([]int, len(tmp))

	for i, t := range tmp {
		r[i], _ = strconv.Atoi(t)
	}

	return r
}

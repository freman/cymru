// Copyright 2019 Freman/Fremnet (Shannon Wynter). All rights reserved.

package ipasn

import (
	"context"
	"net"
)

// DefaultClient is a package level Cymru client object using net.DefaultResolver
// and DefaultPrivateNetworks()
//nolint:gochecknoglobals
var DefaultClient = &Client{}

// Origin is used to map an IPv4 or IPv6 address or prefix to a corresponding
// BGP Origin ASN.
func Origin(ctx context.Context, ip net.IP) (o OriginInfo, err error) {
	return DefaultClient.Origin(ctx, ip)
}

// Peer is used to map an IP address or prefix to the possible BGP peer ASNs that
// are one AS hop away from the BGP Origin ASN's prefix.
func Peer(ctx context.Context, ip net.IP) (p PeerInfo, err error) {
	return DefaultClient.Peer(ctx, ip)
}

// ASN is used to determine the AS description of a given BGP ASN.
// Notably this function returns the Description of the AS but not the network.
func ASN(ctx context.Context, asn int) (a ASNInfo, err error) {
	return DefaultClient.ASN(ctx, asn)
}

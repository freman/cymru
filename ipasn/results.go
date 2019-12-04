// Copyright 2019 Freman/Fremnet (Shannon Wynter). All rights reserved.

package ipasn

import (
	"fmt"
	"net"
	"time"
)

// OriginInfo is returned by Origin() and contains the BGP Origin ASN.
type OriginInfo struct {
	ASN       int
	Network   *net.IPNet
	Country   string
	Authority string
	Updated   time.Time
}

func (o OriginInfo) String() string {
	if o.ASN == 0 || o.Network == nil {
		return ""
	}

	return fmt.Sprintf("%d | %s | %s | %s | %s",
		o.ASN,
		o.Network.String(),
		o.Country,
		o.Authority,
		o.Updated.Format(dateFormat),
	)
}

// PeerInfo is returned by Peer() and contains BGP peer ASNs that are one
// AS hop away from the BGP Origin ASN's prefix.
type PeerInfo struct {
	ASNs      []int
	Network   *net.IPNet
	Country   string
	Authority string
	Updated   time.Time
}

func (p PeerInfo) String() string {
	if len(p.ASNs) == 0 || p.Network == nil {
		return ""
	}

	asns := fmt.Sprintf("%v", p.ASNs)

	return fmt.Sprintf("%s | %s | %s | %s | %s",
		asns[1:len(asns)-1],
		p.Network.String(),
		p.Country,
		p.Authority,
		p.Updated.Format(dateFormat),
	)
}

// ASNInfo is returned by ASN() and contains the AS description of a given
// BGP ASN.
type ASNInfo struct {
	ASN         int
	Country     string
	Authority   string
	Updated     time.Time
	Description string
}

func (a ASNInfo) String() string {
	if a.ASN == 0 {
		return ""
	}

	return fmt.Sprintf("%d | %s | %s | %s | %s",
		a.ASN,
		a.Country,
		a.Authority,
		a.Updated.Format(dateFormat),
		a.Description,
	)
}

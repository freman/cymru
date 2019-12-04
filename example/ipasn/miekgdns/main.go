// Copyright 2019 Freman/Fremnet (Shannon Wynter). All rights reserved.

package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/freman/cymru/ipasn"
	"github.com/miekg/dns"
)

type lookupTXT struct{}

func (l lookupTXT) LookupTXT(ctx context.Context, name string) ([]string, error) {
	r, _, err := new(dns.Client).ExchangeContext(ctx, new(dns.Msg).SetQuestion(name, dns.TypeTXT), "8.8.8.8:53")
	if err != nil {
		return nil, err
	}

	for _, answer := range r.Answer {
		if txt, isa := answer.(*dns.TXT); isa {
			return txt.Txt, nil
		}
	}

	return nil, err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please pass an IP address")
		os.Exit(1)
	}

	c := ipasn.Client{
		Resolver: &lookupTXT{},
	}

	origin, err := c.Origin(context.Background(), net.ParseIP(os.Args[1]))
	if err != nil {
		fmt.Println("Error looking up the origin:", err)
		os.Exit(1)
	}

	asn, err := c.ASN(context.Background(), origin.ASN)
	if err != nil {
		fmt.Println("Error looking up the description:", err)
		os.Exit(1)
	}

	fmt.Printf("ipset add badnetworks %s comment %q timeout 86400\n", origin.Network, asn.Description)
}

// Copyright 2019 Freman/Fremnet (Shannon Wynter). All rights reserved.

package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/freman/cymru/ipasn"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please pass an IP address")
		os.Exit(1)
	}

	origin, err := ipasn.Origin(context.Background(), net.ParseIP(os.Args[1]))
	if err != nil {
		fmt.Println("Error looking up the origin:", err)
		os.Exit(1)
	}

	asn, err := ipasn.ASN(context.Background(), origin.ASN)
	if err != nil {
		fmt.Println("Error looking up the description:", err)
		os.Exit(1)
	}

	fmt.Printf("ipset add badnetworks %s comment %q timeout 86400\n", origin.Network, asn.Description)
}

// Copyright 2019 Freman/Fremnet (Shannon Wynter). All rights reserved.

package ipasn

// Error is a string that will be returned by Cymru when things go bad
type Error string

func (s Error) Error() string {
	return string(s)
}

// Various errors that will be returned depending on how things go
const (
	ErrIPIsUnspecified Error = "IP is unspecified"
	ErrIPIsLoopback    Error = "IP is a loopback address"
	ErrIPIsMulticast   Error = "IP is a multicast address"
	ErrIPIsPrivate     Error = "IP is a private address"
	ErrNotFound        Error = "DNS result included no useful records"
)

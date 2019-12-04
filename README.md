# CYMRU

[Team Cymru](https://www.team-cymru.com/) related go libraries

## Badges

![](https://github.com/freman/cymru/workflows/Test/badge.svg) [![codecov](https://codecov.io/gh/freman/cymru/branch/master/graph/badge.svg)](https://codecov.io/gh/freman/cymru) [![Go Report Card](https://goreportcard.com/badge/github.com/freman/cymru)](https://goreportcard.com/report/github.com/freman/cymru) [![license](https://img.shields.io/github/license/freman/cymru.svg?maxAge=2592000)](LICENSE.md)

[![Documentation](https://godoc.org/github.com/freman/cymru?status.svg)](https://godoc.org/github.com/freman/cymru)

## Packages

### [**ipasn**](ipasn)

Interface for the [Team Cymru DNS IP-ASN mapping interface](https://www.team-cymru.com/IP-ASN-mapping.html#dns)

Make use of the DNS IP-ASN mapping interface to look up BGP prefixes and ASN's based on a given IP address.

eg:

```go
origin, err := ipasn.Origin(context.Background(), net.ParseIP("1.1.1.1"))
if err != nil {
    panic(err)
}

fmt.Println(origin)
```

Results in

```
13335 | 1.1.1.0/24 | AU | apnic | 2011-08-11
```




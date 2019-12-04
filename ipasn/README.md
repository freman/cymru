# IP-ASN

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
# dm-resolver

dm-resolver is a light way resolver library based on IP lookup and inspired by the default gRPC resolver.

## Usage

The library can be used as a resolver for the gRPC go client or can be used to keep up to date your HTTP connection pool up to date.

dm-resolver is also used [here](https://github.com/cperez08/h2-proxy) as a library in the h2-proxy to help to update the domain and refresh the http2 connection pool.

## Motivation

Trying to balance gRPC traffic properly in Kubernetes I came across a problem with the gRPC DNS resolver + Docker Alpine image + Kubernetes.

The problem laid in the balancing after changes in the initial mapping of the IPs against the domain. the gRPC client was actually balancing correctly the request, however, at the time Kubernetes scale up the pods or assign new IP for some reason to the pods the DNS resolver never updated the connection state hence stopped to balance correctly the request.

For instance:

- the initial lookup returned for the domain test.kubernetes 2 IPs (10.0.0.0 and 11.0.0.0)
- the gRPC client started to work correctly and balance (round-robin algorithm) the request properly
- after scaling up the new set of IPs was something like 3 IPs (10.0.0.0, 11.0.0.0 and 12.0.0.0)
- the gRPC client never sent a new request to the new IP 12.0.0.0


### Usage for gRPC

```go
package main

import (
    rsv "github.com/cperez08/dm-resolver/pkg/resolver"
    "google.golang.org/grpc"
    "google.golang.org/grpc/balancer/roundrobin"
    "google.golang.org/grpc/metadata"
    "google.golang.org/grpc/resolver"
)

var (
	scheme      = "my-scheme-name"
	host        = "service-address.com"
	port        = "50051"
	refreshRate = time.Duration(15) // will be parsed in secs
)

func main() { 
    address := fmt.Sprintf("%s:///%s", scheme, host+":"+port)
    conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }

    defer conn.Close()
    cli = mypkg.NewPkgServiceClient(conn)

    cli.Call(ctx, &myRequest{})

}

func init(){

    // scheme is the custom dns name
    // host is the ip or domain name
    // port number
    // true indicates if the resolver needs to watch for changes - only aplicable for domains
    // 50 indicates the value in seconds the resolver look for the changes in the domain.
    resolver.Register(rsv.NewDomainResolverBuilder(scheme, host, port, true, &refreshRate))
}

```

### Usage outside gRPC

```go
import "github.com/cperez08/dm-resolver/pkg/resolver"

func main(){

    listener := make(chan bool)
    // host that wants to be resolved (domain or IP)
    // port to connecto to
    // true indicates if the wants the resolver to watch for domain changes
    // time.Duration(50) in case the previous parameter is true, is the refresh rate (new domain lookup)
    refreshRate := time.Duration(50)
    // listener listen for changes in the IPs, in case there is no change in the initial set of ips nothing is triggered
    r = resolver.NewResolver(host, port, true, &refreshRate, listener)
    // StartResolver resolves the domain  the firstime and starts the domain watcher if enabled and if the address is not an IP
    r.StartResolver()

    // for knowing the current Addresses stored  by the resolver
    r.Addresses // return a list of string in the format host:port

    // for stopping the watcher
    r.Close()
}
```

Disclaimer: the issue commented above occurred on alpine >= 3.10.5 and <= 3.12 and was never tested with other Linux distros. In local environment the gRPC DNS resolver worked perfectly.
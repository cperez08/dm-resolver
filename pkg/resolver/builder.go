package resolver

import (
	"time"

	"google.golang.org/grpc/resolver"
)

// DomainResolverBuilder implements the Resolver.Builder interface
// the target field in the Dial gRPC function is in the way
// fmt.Sprintf("%s:///%s", scheme, address+":"+port)
// where the scheme is the name set in this constructor
type domainResolverBuilder struct {
	address     string
	port        string
	scheme      string
	needWatcher bool
	refreshRate *time.Duration
}

// NewDomainResolverBuilder creates a new instance for the DomainResolverBuilder
func NewDomainResolverBuilder(scheme, address, port string, needWatcher bool, refreshRate *time.Duration) *domainResolverBuilder {
	return &domainResolverBuilder{address, port, scheme, needWatcher, refreshRate}
}

// Build ...
func (b *domainResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := NewResolver(b.address, b.port, b.needWatcher, b.refreshRate, nil)
	r.target = target
	r.cc = cc
	r.updateState = true
	r.StartResolver()
	return r, nil
}

// Scheme ...
func (b *domainResolverBuilder) Scheme() string {
	return b.scheme
}

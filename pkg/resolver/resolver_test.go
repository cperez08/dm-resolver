package resolver

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/resolver"
)

var refreshRate = time.Duration(50)

func TestNewResolver(t *testing.T) {
	r := NewResolver("localhost", "8080", true, &refreshRate, nil)
	assert.NotNil(t, r)
	r = NewResolver("localhost", "8080", false, &refreshRate, nil)
	assert.NotNil(t, r)
	r = NewResolver("127.0.0.1", "8080", false, &refreshRate, nil)
	assert.NotNil(t, r)
	assert.Equal(t, 1, len(r.Addresses))
}

func TestStartResolver(t *testing.T) {
	r := NewResolver("127.0.0.1", "8080", false, &refreshRate, nil)
	r.StartResolver()
	assert.Equal(t, 1, len(r.Addresses))

	r = NewResolver("localhost", "8080", true, &refreshRate, nil)
	r.StartResolver()
	assert.True(t, len(r.Addresses) > 0)

	r = NewResolver("localhost", "8080", false, &refreshRate, nil)
	r.StartResolver()
	assert.True(t, len(r.Addresses) > 0)
}

func TestResolverFromBuilder(t *testing.T) {
	rb := NewDomainResolverBuilder("my-schema", "127.0.0.1", "8080", true, &refreshRate)
	rr, err := rb.Build(resolver.Target{Scheme: "test-schema", Endpoint: "127.0.0.1:8080"}, &TestResolver{}, resolver.BuildOptions{})
	assert.Nil(t, err)
	parsed := rr.(*DomainResolver)
	parsed.StartResolver()
	assert.Equal(t, 1, len(parsed.Addresses))

	rb = NewDomainResolverBuilder("my-schema", "localhost", "8080", true, &refreshRate)
	rr, err = rb.Build(resolver.Target{Scheme: "test-schema", Endpoint: "localhost:8080"}, &TestResolver{}, resolver.BuildOptions{})
	assert.Nil(t, err)
	parsed = rr.(*DomainResolver)
	parsed.StartResolver()
	assert.True(t, len(parsed.Addresses) > 0)
}

func TestResolveNow(t *testing.T) {
	r := NewResolver("localhost", "8080", true, &refreshRate, nil)
	r.ResolveNow(resolver.ResolveNowOptions{})
}

func TestCloseResolver(t *testing.T) {
	r := NewResolver("127.0.0.1", "8080", true, &refreshRate, nil)
	r.Close()

	r = NewResolver("localhost", "8080", true, &refreshRate, nil)
	r.StartResolver()
	r.Close()
}

func TestGetState(t *testing.T) {
	r := NewResolver("no-domain1234.com", "8080", false, &refreshRate, nil)
	r.StartResolver()
	state, isUpdated := r.getState()
	assert.Equal(t, 0, len(state.Addresses))
	assert.False(t, isUpdated)

	r.address = "localhost"
	state, isUpdated = r.getState()
	assert.True(t, len(state.Addresses) > 0)
	assert.True(t, isUpdated)

	r.address = "127.0.0.1"
	state, isUpdated = r.getState()
	assert.True(t, len(state.Addresses) > 0)
	assert.True(t, isUpdated)

	// without changes since start resolver already populated the addresses
	c := make(chan bool)
	r = NewResolver("localhost", "8080", false, &refreshRate, c)
	r.StartResolver()
	state, isUpdated = r.getState()
	assert.False(t, len(state.Addresses) > 0)
	assert.False(t, isUpdated)

	// increasing coverage by sending the update signal to the channel
	r.address = "localhost"
	r.Addresses = []string{"127.0.0.1"}

	go func() {
		r.getState()
	}()
	<-c
}

func TestWatchResolver(t *testing.T) {
	refresh := time.Duration(1)
	c1 := make(chan bool)
	c2 := make(chan bool)
	r := NewResolver("127.0.0.1", "8080", true, &refresh, c1)
	r.StartResolver()
	r2 := NewResolver("localhost", "8080", true, &refresh, c2)
	r2.StartResolver()

	assert.True(t, len(r2.Addresses) > 0)
	assert.True(t, len(r.Addresses) == 1)
}

func TestWatchResolverFromBuilder(t *testing.T) {
	refresh := time.Duration(1)
	rb := NewDomainResolverBuilder("my-schema", "localhost", "8080", true, &refresh)
	rr, err := rb.Build(resolver.Target{Scheme: "test-schema", Endpoint: "localhost:8080"}, &TestResolver{}, resolver.BuildOptions{})
	assert.Nil(t, err)
	parsed := rr.(*DomainResolver)

	<-parsed.ticker.C
	assert.True(t, len(parsed.Addresses) > 0)
}

package resolver

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
)

type TestResolver struct{}

func (r *TestResolver) UpdateState(resolver.State)              {}
func (r *TestResolver) ReportError(error)                       {}
func (r *TestResolver) NewAddress(addresses []resolver.Address) {}
func (r *TestResolver) NewServiceConfig(serviceConfig string)   {}
func (r *TestResolver) ParseServiceConfig(serviceConfigJSON string) *serviceconfig.ParseResult {
	return &serviceconfig.ParseResult{}
}

func TestNewDomainResolverBuilder(t *testing.T) {
	refreshPeriod := time.Duration(50)
	r := NewDomainResolverBuilder("test-schema", "127.0.0.1", "8080", false, &refreshPeriod)
	assert.NotNil(t, r)

	_, err := r.Build(resolver.Target{Scheme: "test-schema", Endpoint: "127.0.0.1:8080"}, &TestResolver{}, resolver.BuildOptions{})
	assert.Nil(t, err)
	assert.Equal(t, "test-schema", r.Scheme())
}

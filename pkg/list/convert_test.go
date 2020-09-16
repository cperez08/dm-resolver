package list

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/resolver"
)

func TestFromAddrToString(t *testing.T) {
	addr := []resolver.Address{{Addr: "my-domain.com"}}
	assert.Equal(t, []string{"my-domain.com"}, FromAddrToString(addr))

	addr = append(addr, resolver.Address{Addr: "localhost:8080"})
	assert.Equal(t, []string{"my-domain.com", "localhost:8080"}, FromAddrToString(addr))
	assert.Equal(t, []string{}, FromAddrToString([]resolver.Address{}))
}

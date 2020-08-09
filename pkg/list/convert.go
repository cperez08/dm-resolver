package list

import "google.golang.org/grpc/resolver"

// FromAddrToString converts a list of resolver address into
// a stirng list
func FromAddrToString(addrs []resolver.Address) []string {
	var rs []string
	for _, a := range addrs {
		rs = append(rs, a.Addr)
	}

	return rs
}

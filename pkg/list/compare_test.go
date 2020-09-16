package list

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompareListStr(t *testing.T) {
	base := []string{"1"}
	new := []string{}
	assert.Equal(t, true, CompareListStr(base, new))

	new = append(new, "1")
	assert.Equal(t, false, CompareListStr(base, new))

	new = append(new, "2")
	assert.Equal(t, true, CompareListStr(base, new))
	assert.Equal(t, false, CompareListStr([]string{}, []string{}))
	assert.Equal(t, true, CompareListStr([]string{"1"}, []string{"2"}))
}

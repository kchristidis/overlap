package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPoint(t *testing.T) {
	p := pointImpl{loc: 1}
	p.addTo(1)
	// See: https://npf.io/2017/08/lies/
	require.Equal(t, []int{1}, p.in)
}

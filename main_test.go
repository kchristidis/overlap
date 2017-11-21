package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPoint(t *testing.T) {
	p := pointImpl{location: 1}
	p.addTo("foo")
	// See: https://npf.io/2017/08/lies/
	require.Equal(t, []string{"foo"}, p.in)
}

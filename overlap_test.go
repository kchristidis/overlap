package overlap

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBadFile(t *testing.T) {
	f, err := ioutil.TempFile("", "bad.file")
	require.Nil(t, err, "could not create temp file")
	require.Nil(t, os.Remove(f.Name()), "could not remove temp file")
	_, err = Calculate(f.Name())
	require.Error(t, err, "expected error")
}

func TestPoint(t *testing.T) {
	p := pointImpl{loc: 1}
	p.addTo(1)
	// See: https://npf.io/2017/08/lies/
	require.Equal(t, []int{1}, p.in)
}

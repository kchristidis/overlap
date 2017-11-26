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

func TestBadLines(t *testing.T) {
	f, err := ioutil.TempFile("", "good.file")
	require.Nil(t, err, "could not create temp file")
	defer os.Remove(f.Name())
	lines := make([][]byte, 4)
	lines[0] = []byte("1\t0")
	lines[1] = []byte("foo\t1\t2")
	lines[2] = []byte("0\tfoo\t2")
	lines[3] = []byte("0\t1\tfoo")
	for i := 0; i < len(lines); i++ {
		_, err = f.WriteAt(lines[i], 0)
		require.Nil(t, err, "could not write to input file")
		_, err = Calculate(f.Name())
		require.Error(t, err, "expected error")
	}
}

func TestPoint(t *testing.T) {
	p := pointImpl{loc: 1}
	p.addTo(1)
	// See: https://npf.io/2017/08/lies/
	require.Equal(t, []int{1}, p.in)
}

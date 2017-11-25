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

func TestIncompleteLine(t *testing.T) {
	f, err := ioutil.TempFile("", "good.file")
	require.Nil(t, err, "could not create temp file")
	defer os.Remove(f.Name())
	badLine := []byte("1\t0") // 2 fields
	_, err = f.Write(badLine)
	require.Nil(t, err, "could not write to input file")
	_, err = Calculate(f.Name())
	require.Error(t, err, "expected error")
}

func TestBadLine(t *testing.T) {
	f, err := ioutil.TempFile("", "good.file")
	require.Nil(t, err, "could not create temp file")
	defer os.Remove(f.Name())
	badLine := []byte("foo\t1\t2") // 2 fields
	_, err = f.Write(badLine)
	require.Nil(t, err, "could not write to input file")
	_, err = Calculate(f.Name())
	require.Error(t, err, "expected error")
}

func TestPoint(t *testing.T) {
	p := pointImpl{loc: 1}
	p.addTo(1)
	// See: https://npf.io/2017/08/lies/
	require.Equal(t, []int{1}, p.in)
}

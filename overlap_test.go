package overlap

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPoint(t *testing.T) {
	p := pointImpl{loc: 0.5}
	p.addTo("foo")
	p.addTo("bar")
	// See: https://npf.io/2017/08/lies/
	require.Equal(t, len(p.belongsTo()), 2, "expected different segment count")
}

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
	lines := make([][]byte, 3)
	lines[0] = []byte("0,100")
	lines[1] = []byte("0,foo,200")
	lines[2] = []byte("0,100,foo")
	for i := range lines { // https://stackoverflow.com/a/39806983/2363529
		_, err = f.WriteAt(lines[i], 0)
		require.Nil(t, err, "could not write to input file")
		_, err = Calculate(f.Name())
		require.Error(t, err, "expected error")
	}
}

func TestCalculate(t *testing.T) {
	f, err := ioutil.TempFile("", "good.file")
	require.Nil(t, err, "could not create temp file")
	defer os.Remove(f.Name())
	// https://stackoverflow.com/a/39806983/2363529
	lines := []byte("foo,50.0,150.0" + "\n" + "bar,100,200.0")
	_, err = f.Write(lines)
	require.Nil(t, err, "could not write to input file")
	res, err := Calculate(f.Name())
	require.Nil(t, err, "expected no error")
	require.Equal(t, len(res), 2, "expected 1 overlap") // increment by 1 to account for the header
	require.Equal(t, res[1], []string{
		strconv.FormatFloat(50, 'f', -1, 64),      // overlapLength
		strconv.FormatFloat(100, 'f', -1, 64),     // overlapStart
		strconv.FormatFloat(150, 'f', -1, 64),     // overlapEnd
		strconv.Itoa(2),                           // segmentCount
		strings.Join([]string{"bar", "foo"}, ","), // segmentList
	}, "expected [100,150] overlap between the segments")
}

package overlap

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPoint(t *testing.T) {
	p := pointImpl{loc: 0.5}
	p.addTo(0)
	p.addTo(1)
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
	lines := make([][]byte, 4)
	lines[0] = []byte("0\t100")
	lines[1] = []byte("foo\t100\t200")
	lines[2] = []byte("0\tfoo\t200")
	lines[3] = []byte("0\t100\tfoo")
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
	lines := []byte("0\t50.0\t150.0" + "\n" + "1\t100\t200.0")
	_, err = f.Write(lines)
	require.Nil(t, err, "could not write to input file")
	res, err := Calculate(f.Name())
	require.Nil(t, err, "expected no error")
	require.Equal(t, len(res), 1, "expected 1 overlap")
	require.Equal(t, res[0], Result{
		OverlapLength: float64(50),
		OverlapStart:  float64(100),
		OverlapEnd:    float64(150),
		SegmentCount:  2,
		SegmentList:   []int{0, 1},
	}, "expected [100,150] overlap between the segments")
}

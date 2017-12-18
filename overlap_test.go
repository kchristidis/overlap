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
	_, err = Calculate(f.Name(), false)
	require.Error(t, err, "expected error")
}

func TestBadLines(t *testing.T) {
	f, err := ioutil.TempFile("", "good.file")
	require.Nil(t, err, "could not create temp file")
	defer os.Remove(f.Name())

	type testInput struct {
		line          []byte
		assumeHeaders bool
	}
	testInputs := []testInput{
		testInput{},
		testInput{line: []byte("foo,100"), assumeHeaders: false},
		testInput{line: []byte("foo,bar,200"), assumeHeaders: false},
		testInput{line: []byte("foo,100,bar"), assumeHeaders: false},
		testInput{line: []byte("foo,100,bar"), assumeHeaders: true},
	}
	for _, testInput := range testInputs { // https://stackoverflow.com/a/39806983/2363529
		_, err = f.WriteAt(testInput.line, 0)
		require.Nil(t, err, "could not write to input file")
		_, err = Calculate(f.Name(), testInput.assumeHeaders)
		require.Error(t, err, "expected error")
	}
}

func TestCalculate(t *testing.T) {
	f, err := ioutil.TempFile("", "good.file")
	require.Nil(t, err, "could not create temp file")
	defer os.Remove(f.Name())
	// https://stackoverflow.com/a/39806983/2363529
	lines := []byte("id,start,end" + "\n" +
		"foo,50.0,150.0" + "\n" +
		"bar,100,200.0" + "\n" +
		"baz,199,201")
	_, err = f.Write(lines)
	require.Nil(t, err, "could not write to input file")
	res, err := Calculate(f.Name(), true)
	require.Nil(t, err, "expected no error")
	require.Equal(t, len(res), 3, "expected 2 overlaps") // increment by 1 to account for the header
	require.Equal(t, res[1], []string{
		strconv.FormatFloat(50, 'f', -1, 64),      // overlapLength
		strconv.FormatFloat(100, 'f', -1, 64),     // overlapStart
		strconv.FormatFloat(150, 'f', -1, 64),     // overlapEnd
		strconv.Itoa(2),                           // segmentCount
		strings.Join([]string{"bar", "foo"}, ","), // segmentList
	})
	require.Equal(t, res[2], []string{
		strconv.FormatFloat(1, 'f', -1, 64),       // overlapLength
		strconv.FormatFloat(199, 'f', -1, 64),     // overlapStart
		strconv.FormatFloat(200, 'f', -1, 64),     // overlapEnd
		strconv.Itoa(2),                           // segmentCount
		strings.Join([]string{"bar", "baz"}, ","), // segmentList
	})
}

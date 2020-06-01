# overlap

[![DOI](https://zenodo.org/badge/111238120.svg)](https://zenodo.org/badge/latestdoi/111238120)
[![GoDoc](https://godoc.org/github.com/kchristidis/overlap?status.svg)](https://godoc.org/github.com/kchristidis/overlap)
[![Build Status](https://travis-ci.org/kchristidis/overlap.svg?branch=master)](https://travis-ci.org/kchristidis/overlap)

overlap is a library for identifying overlaps on a list of segments.

## Motivation

You are given a list of segments that looks like this:

![Candidate Segments](https://user-images.githubusercontent.com/14876848/34177937-62013f34-e4d3-11e7-9475-9a51b21095fe.png)

You are asked to identify any overlaps that are longer than `m`, or which have at least `n` segments overlapping.

This package allows you to do that, and identify, for instance, an overlap that is longer than 2 years and has at least 100 segments.

![Chosen Segments](https://user-images.githubusercontent.com/14876848/34177999-9ba191e4-e4d3-11e7-8f1e-e6caa01cd5ca.png)

## Installation

```bash
$ go get github.com/kchristidis/overlap
```

## Usage

```go
results, _ := overlap.Calculate("segments.csv", True)
f, _ := os.Create("overlaps.csv")
defer f.Close()
w := csv.NewWriter(f)
w.WriteAll(results) // calls Flush internally
```

You can also study [the sample binary](https://github.com/kchristidis/overlap/tree/master/cmd/overlap) provided in `cmd/overlap`, or read the package documentation in [GoDoc](http://godoc.org/github.com/kchristidis/overlap).

## Contributing

Contributions are welcome. Fork this library and submit a pull request.

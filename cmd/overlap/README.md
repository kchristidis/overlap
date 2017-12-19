# overlap

overlap identifies overlaps among segments. It is given a list of segments in a comma-separated values (CSV) file, and outputs a CSV file with the identified overlaps.

It is meant to demonstrate what a sample application built on [the package of the same name](https://github.com/kchristidis/overlap) looks like. 

## Installation

```bash
$ go get -u github.com/kchristidis/cmd/overlap
```

This will download and build the `overlap` utility, installing it in `$GOPATH/bin/overlap`.

## Usage

```bash
$ overlap [-headers] inputFile [outputFile]
```

Where:
* `-headers` (optional value, default value: false) is a flag that you set if your input CSV file has a header row describing its columns. If you do not set this flag, the input CSV file is assumed to _not_ include such a row.
* `inputFile` (required argument) is the reference to the file describing the segments.
* `outputFile` (optional argument) is the reference to the file that `overlap` will use to write the overlaps. This file does not need to exist before. If it does, it _will_ be overwritten. If this argument is ommitted, the results will be written to a file with the same name as the `inputFile` argument _and_ the `out_` prefix.

Here is an example of a valid invocation:

```bash
$ cd $GOPATH/src/github.com/kchristidis/cmd/overlap
$ overlap -headers testdata.csv result.csv
```

(`testdata.csv` contains data taken from [Dataport](https://dataport.cloud)'s repo.)

## Contributing

Contributions are welcome. Fork this library and submit a pull request.
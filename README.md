# pprofutils

pprofutils provides command line utilities for converting pprof files to Brendan Gregg's [folded text](https://github.com/brendangregg/FlameGraph#2-fold-stacks) format and vice versa.

## Install

pprofutils requires Go 1.16 and can be installed like this:

```
go install github.com/felixge/pprofutils/cmd/...
```

## Usage

Convert a pprof file to text:

```
pprof2text < ./test-fixtures/pprof.samples.cpu.001.pb.gz > example.txt
```

Convert a text file to pprof:

```
text2pprof < example.txt > example.pprof
```

## License

pprofutils is licensed under the MIT License.

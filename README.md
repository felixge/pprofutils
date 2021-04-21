[![documentation](http://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/felixge/pprofutils)

# pprofutils

pprofutils provides command line utilities for converting [pprof files](https://github.com/DataDog/go-profiler-notes/blob/main/pprof.md) to Brendan Gregg's [folded text](https://github.com/brendangregg/FlameGraph#2-fold-stacks) format and vice versa.

## Install

pprofutils requires Go 1.16 and can be installed like this:

```
go install github.com/felixge/pprofutils/cmd/...@latest
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

Warning: Converting from pprof to text is lossy. Only the first sample type will be converted, file names, lines, labels, and more will be dropped. Patches to make things less lossy would be welcome, but please open an issue first to discuss.

Create a delta profile that contains the difference `heap-b.pprof - heap-a.pprof`:

```
pprofdelta -o delta.pprof heap-a.pprof heap-b.pprof
```

## Tutorial: Generate a fake pprof profile

My primary use case for this tool is to quickly generate fake pprof profiles for creating educational content.

This can be done by simply creating a file called `profile.txt` with the following content:

```
main;foo 5
main;foo;bar 3
main;foobar 4
```

Then convert it to a pprof profile:

```
text2pprof < profile.txt > profile.pprof
```

And finally view it using pprof:

```
go tool pprof -http=:6060 profile.pprof
```

The resulting graphs should look like this:

![](./img/flamegraph.png)

![](./img/graph.png)

## Custom Extension: Multiple Sample Types

The `text2pprof` command supports a custom extension to the folded text format that allows users to specify multiple sample types.

This is done via a header that contains space separated `type/unit` sample types. The stack traces on the following lines must then contain one value for each sample type after the stack trace:

```
samples/count duration/nanoseconds
main;foo 5 50000000
main;foo;bar 3 30000000
main;foobar 4 40000000
```

The `pprof2text` command also supports outputting this format by passing the `-m` flag.

## License

pprofutils is licensed under the MIT License.

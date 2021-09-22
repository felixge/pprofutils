[![documentation](http://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/felixge/pprofutils)
[![ci test status](https://img.shields.io/github/workflow/status/felixge/pprofutils/Go?label=tests)](https://github.com/felixge/pprofutils/actions/workflows/go.yml?query=branch%3Amain)

# pprofutils

pprofutils is a swiss army knife for [pprof files](https://github.com/DataDog/go-profiler-notes/blob/main/pprof.md). You can use it as a command line utility or as a free web service.

- [**Install**](#install)
- [**Utilities**](#utilities): {{range $i, $util := .}}{{if $i}} · {{end}}[{{.Name}}](#{{.Name}}){{end}}
- [**Use Cases**](#use-cases): [Convert linux perf profiles to pprof](#convert-linux-perf-profiles-to-pprof)
- [**License**](#license)

## Install

pprofutils requires Go 1.16 and can be installed like this:

```
go install github.com/felixge/pprofutils/cmd/pprofutils@latest
```

Alternatively you can use it as a free web service hosted at https://pprof.to.

## Utilities

{{range $i := .}}### {{.Name}}

{{.LongHelp}}

#### Use {{.Name}} utility via cli

```
pprofutils {{.Name}} {{.ShortUsage}}{{if .Flags}}

FLAGS:{{range $name, $flag := .Flags}}
  -{{$name}}={{defaultval .Default}} {{.Usage}}{{end}}{{else}}{{end}}
```

#### Use {{.Name}} utility via web service

```
curl --data-binary @<input file> pprof.to/{{.Name}}{{queryflags .Flags}} > <output file>
```

{{examples .}}

{{end}}

## Use Cases

### Convert linux perf profiles to pprof

Convert a Linux `perf.data` profile to `pprof`, via Brendan Gregg's [`stackcollapse-perf.pl`](https://github.com/brendangregg/FlameGraph/blob/master/stackcollapse-perf.pl) script:

```bash
perf script | stackcollapse-perf.pl | pprofutils folded > perf.pprof
```

## License

pprofutils is licensed under the MIT License.

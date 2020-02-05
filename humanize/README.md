[tawesoft.co.uk/go/humanize](https://www.tawesoft.co.uk/go/humanize)
================================================================================

Lightweight human-readable numbers for Go.

The API is incomplete and may be subject to occasional breaking changes.

Download/update:

```
go get -u tawesoft.co.uk/go/humanize
```

Usage:

```
import "tawesoft.co.uk/go/humanize"
```

Docs:

[godoc.org/tawesoft.co.uk/go/humanize](https://godoc.org/tawesoft.co.uk/go/humanize)

Examples:

```
ParseBytes("1 KiB") => 1024
FormatBytesSI(nil, 1e6) => "1 MB"
FormatBytesIEC(nil, 1024*1024) => "1 MiB"
```

Versus dustin's [go-humanize](https://github.com/dustin/go-humanize)
--------------------------------------------------------------------------------

* `tawesoft.co.uk/go/humanize` parses about 5 times faster with fewer memory
allocations. Benchmark (YMMV):

```
BenchmarkTawesoftFormatBytes-4   	 3590067	       317 ns/op
BenchmarkDustinFormatBytes-4     	 2705889	       439 ns/op
BenchmarkTawesoftParseBytes-4    	55214542	        21.5 ns/op
BenchmarkDustinParseBytes-4      	10613521	       108 ns/op
```

* `github.com/dustin/go-humanize` is more complete, older, and has been tested
by more people

* `tawesoft.co.uk/go/humanize` handles fractional ammounts, such as "1.5 MB"

* `github.com/dustin/go-humanize` has a more stable API

* `tawesoft.co.uk/go/humanize` exposes lower-level components suitable for
constructing parsers and formatters of new numbers and unit types
